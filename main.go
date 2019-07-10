package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

type FileJob struct {
	Language           string
	PossibleLanguages  []string
	Filename           string
	Extension          string
	Location           string
	Content            []byte
	Bytes              int64
	Lines              int64
	Code               int64
	Comment            int64
	Blank              int64
	Complexity         int64
	WeightedComplexity float64
	Hash               []byte
	Binary             bool
}

type LanguageSummary struct {
	Name               string
	Bytes              int64
	Lines              int64
	Code               int64
	Comment            int64
	Blank              int64
	Complexity         int64
	Count              int64
	WeightedComplexity float64
	Files              []FileJob
}

type File struct {
	Name            string
	Content         []byte
	LanguageSummary LanguageSummary
}

// Gets from disk, but could be S3 just as easily
func getFiles(directory string, output chan File) error {
	all, err := ioutil.ReadDir(directory)

	if err != nil {
		return errors.New("unable to read directory: " + directory)
	}

	for _, fileInfo := range all {
		if !fileInfo.IsDir() {
			content, err := ioutil.ReadFile(filepath.Join(directory, fileInfo.Name()))

			if err == nil {
				output <- File{
					Name:    fileInfo.Name(),
					Content: content,
				}
			}
		}
	}

	close(output)
	return nil
}

func unmarshallContent(content []byte) ([]LanguageSummary, error) {
	var summary []LanguageSummary
	if err := json.Unmarshal(content, &summary); err != nil {
		return summary, err
	}

	return summary, nil
}

func getLineCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		x += y.Lines
	}

	return x
}

func getCodeCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		x += y.Code
	}

	return x
}

func getBlankCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		x += y.Blank
	}

	return x
}

func getComplexityCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		x += y.Complexity
	}

	return x
}

func getCommentCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		x += y.Comment
	}

	return x
}

func getByteCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		for _, z := range y.Files {
			x += z.Bytes
		}
	}

	return x
}

func makeTimestampSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func main() {
	startTime := makeTimestampSeconds()
	queue := make(chan File, 1000)
	go getFiles("./json/", queue)

	var projectCount int64
	var lineCount int64
	var codeCount int64
	var blankCount int64
	var commentCount int64
	var complexityCount int64
	var fileCount int64
	var byteCount int64

	lineDistributionPerProject := map[int64]int64{}
	lineDistributionPerFile := map[int64]int64{}
	lineDistributionPerLanguage := map[string]map[int64]int64{}

	codeDistributionPerProject := map[int64]int64{}
	codeDistributionPerFile := map[int64]int64{}
	codeDistributionPerLanguage := map[string]map[int64]int64{}

	commentDistributionPerProject := map[int64]int64{}
	commentDistributionPerFile := map[int64]int64{}
	commentDistributionPerLanguage := map[string]map[int64]int64{}

	blankDistributionPerProject := map[int64]int64{}
	blankDistributionPerFile := map[int64]int64{}
	blankDistributionPerLanguage := map[string]map[int64]int64{}

	complexityDistributionPerProject := map[int64]int64{}
	complexityDistributionPerFile := map[int64]int64{}
	complexityDistributionPerLanguage := map[string]map[int64]int64{}

	filesPerProject := map[int64]int64{}


	for file := range queue {
		summary, err := unmarshallContent(file.Content)

		if err == nil {
			projectCount++
			lineCount += getLineCount(summary)
			codeCount += getCodeCount(summary)
			blankCount += getBlankCount(summary)
			commentCount += getCommentCount(summary)
			complexityCount += getComplexityCount(summary)
			fileCount += int64(len(summary))
			byteCount += getByteCount(summary)

			getLineDistributionPerProject(summary, lineDistributionPerProject)
			getLineDistributionPerFile(summary, lineDistributionPerFile)
			getLineDistributionPerLanguage(summary, lineDistributionPerLanguage)

			getCodeDistributionPerProject(summary, codeDistributionPerProject)
			getCodeDistributionPerFile(summary, codeDistributionPerFile)
			getCodeDistributionPerLanguage(summary, codeDistributionPerLanguage)

			getCommentDistributionPerProject(summary, commentDistributionPerProject)
			getCommentDistributionPerFile(summary, commentDistributionPerFile)
			getCommentDistributionPerLanguage(summary, commentDistributionPerLanguage)

			getBlankDistributionPerProject(summary, blankDistributionPerProject)
			getBlankDistributionPerFile(summary, blankDistributionPerFile)
			getBlankDistributionPerLanguage(summary, blankDistributionPerLanguage)

			getComplexityDistributionPerProject(summary, complexityDistributionPerProject)
			getComplexityDistributionPerFile(summary, complexityDistributionPerFile)
			getComplexityDistributionPerLanguage(summary, complexityDistributionPerLanguage)

			getFilesPerProject(summary, filesPerProject)
		}
	}

	endTime := makeTimestampSeconds()

	fmt.Println("ProjectCount   ", projectCount)
	fmt.Println("LineCount      ", lineCount)
	fmt.Println("CodeCount      ", codeCount)
	fmt.Println("BlankCount     ", blankCount)
	fmt.Println("CommentCount   ", commentCount)
	fmt.Println("ComplexityCount", complexityCount)
	fmt.Println("FileCount      ", fileCount)
	fmt.Println("ByteCount      ", byteCount)
	fmt.Println("StartTime      ", startTime)
	fmt.Println("EndTime        ", endTime)
	fmt.Println("TotalTime(s)   ", endTime - startTime)

	stats := fmt.Sprintf("ProjectCount %d\nLineCount %d\nCodeCount %d\nBlankCount %d\nCommentCount %d\nComplexityCount %d\nFileCount %d\nByteCount %d\nTotalTime(s) %d", projectCount, lineCount, codeCount, blankCount, commentCount, complexityCount, fileCount, byteCount, endTime - startTime)
	_ = ioutil.WriteFile("totalStats.txt", []byte(stats), 0600)

	v, _ := json.Marshal(lineDistributionPerProject)
	_ = ioutil.WriteFile("lineDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(lineDistributionPerFile)
	_ = ioutil.WriteFile("lineDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(lineDistributionPerLanguage)
	_ = ioutil.WriteFile("lineDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(codeDistributionPerProject)
	_ = ioutil.WriteFile("codeDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(codeDistributionPerFile)
	_ = ioutil.WriteFile("codeDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(codeDistributionPerLanguage)
	_ = ioutil.WriteFile("codeDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(commentDistributionPerProject)
	_ = ioutil.WriteFile("commentDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(commentDistributionPerFile)
	_ = ioutil.WriteFile("commentDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(commentDistributionPerLanguage)
	_ = ioutil.WriteFile("commentDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(blankDistributionPerProject)
	_ = ioutil.WriteFile("blankDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(blankDistributionPerFile)
	_ = ioutil.WriteFile("blankDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(blankDistributionPerLanguage)
	_ = ioutil.WriteFile("blankDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(complexityDistributionPerProject)
	_ = ioutil.WriteFile("complexityDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(complexityDistributionPerFile)
	_ = ioutil.WriteFile("complexityDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(complexityDistributionPerLanguage)
	_ = ioutil.WriteFile("complexityDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(filesPerProject)
	_ = ioutil.WriteFile("filesPerProject.json", []byte(v), 0600)
}
