package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func readFile(filename string, output chan string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		output <- text
	}

	close(output)
}

func processPath(s string) string {
	s = strings.ToLower(s)
	split := strings.Split(s, "/")

	if len(split) != 5 {
		return ""
	}

	sp := []string{}

	for _, s := range split {
		sp = append(sp, cleanString(s))
	}

	filename := strings.Replace(sp[2], ".com", "", -1)
	filename = strings.Replace(filename, ".org", "", -1)
	filename += "." + sp[3]
	filename += "." + strings.Replace(sp[4], ".git", "", -1) + ".json"

	return filename
}

func cleanString(s string) (string) {
	reg, err := regexp.Compile("[^a-z0-9-._]+")
	if err != nil {
		log.Fatal(err)
	}

	processedString := reg.ReplaceAllString(s, "")

	return processedString
}

func process(id int, s string) {
	fmt.Println("processing", s)

	// Clean target just to be sure
	cmdArgs := []string{
		"-rf",
		"/tmp/scc-tmp-path-" + strconv.Itoa(id),
	}

	cmd := exec.Command("rm", cmdArgs...)
	err := cmd.Run()

	if err != nil {
		fmt.Println("rm start", err.Error())
		return
	}

	// Run git clone against the target
	// 180 seconds seems enough as the kernel itself takes about 60 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, "git", "clone", "--depth=1", s + ".git", "/tmp/scc-tmp-path-" + strconv.Itoa(id))

	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	resp, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("git clone timed out")
		return
	}

	if err != nil {
		fmt.Println("git clone non-zero exit code", string(resp))
		return
	}

	// Run scc against what we just cloned
	fileName := processPath(s)

	if fileName == "" {
		return
	}

	cmdArgs = []string{
		"-f",
		"json",
		"-o",
		"/tmp/" + fileName,
		"/tmp/scc-tmp-path-" + strconv.Itoa(id),
	}

	cmd = exec.Command("scc", cmdArgs...)
	err = cmd.Run()

	if err != nil {
		fmt.Println("scc", err.Error())
	}

	err = uploadS3File("sloccloccode", fileName, "/tmp/"+fileName)
	if err != nil {
		fmt.Println("s3 upload", err.Error())
	}
	fmt.Println("uploaded now cleaning up")

	// Cleanup
	cmdArgs = []string{
		"-rf",
		"/tmp/" + fileName,
	}

	cmd = exec.Command("rm", cmdArgs...)
	err = cmd.Run()

	if err != nil {
		fmt.Println("rm cleanup filename", err.Error())
		return
	}

	cmdArgs = []string{
		"-rf",
		"/tmp/scc-tmp-path-" + strconv.Itoa(id),
	}

	cmd = exec.Command("rm", cmdArgs...)
	err = cmd.Run()

	if err != nil {
		fmt.Println("rm cleanup", err.Error())
		return
	}
}

func uploadS3File(bucket string, key string, filePath string) error {
	svc, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2")},
	)

	if err != nil {
		return err
	}

	s3client := s3.New(svc)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	_, err = s3client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

func main() {
	queue := make(chan string, 1000)
	go readFile("urls.txt", queue)

	var wg sync.WaitGroup

	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(id int) {
			fmt.Println("start goroutine " + strconv.Itoa(id))
			for s := range queue {
				process(id, s)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}