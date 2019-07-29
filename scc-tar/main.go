package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func getFileKeysS3(output chan string) {
	svc, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2")},
	)

	s3client := s3.New(svc)

	count := 0

	err = s3client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String("sloccloccode"),
		Prefix: aws.String(""),
	}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, value := range page.Contents {
			count++
			output <- *value.Key

			//if count >= 5000 {
			//	return false
			//}
		}

		return true
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	close(output)
}

// Takes in a channel of S3 keys and goes and gets em for processing
func getFilesS3(input chan string, output chan File) {
	for key := range input {
		data, err := clientReadS3File("sloccloccode", key)

		if err == nil {
			output <- File{
				Filename: key,
				Content:  data,
			}
		} else {
			// If we get an error then back off for a while
			fmt.Println(err.Error())
			time.Sleep(10 * time.Second)
		}
	}
}

// Read a file from s3 into memory
func clientReadS3File(bucket string, key string) ([]byte, error) {
	svc, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2")},
	)

	s3client := s3.New(svc)

	results, err := s3client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	defer results.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, results.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type File struct {
	Filename string
	Content  []byte
}

// Adds a file into the tar we are writing out
func addFile(tw *tar.Writer, file File) error {

	// now lets create the header as needed for this file within the tarball
	header := new(tar.Header)
	header.Name = file.Filename
	header.Size = int64(len(file.Content))

	// write the header to the tarball archive
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	// copy the file data to the tarball
	if _, err := io.Copy(tw, bytes.NewReader(file.Content)); err != nil {
		return err
	}

	return nil
}

// Download all of the files from S3 and stuff them into a very large tar file
// so we can download it easily and process
func main() {
	// this is for processing for real
	keys := make(chan string, 11000000) // large enough to hold everything
	queue := make(chan File, 1000)
	go getFileKeysS3(keys)

	var wg sync.WaitGroup

	// Spawn off goroutines to fetch from s3
	go func() {
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				getFilesS3(keys, queue)
				wg.Done()
			}()
		}
		wg.Wait()
		close(queue)
	}()

	// set up the output file
	file, err := os.Create("output.tar.gz")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// set up the gzip writer
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	count := 0
	for f := range queue {
		count++
		if count%100 == 0 {
			fmt.Println(f.Filename, count)
		}

		err = addFile(tw, f)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
