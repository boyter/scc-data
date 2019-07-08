package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func getLineDistributionPerProject(summary []LanguageSummary, lineDistributionPerProject map[int64]int64) {
	for _, y := range summary {
		lineDistributionPerProject[y.Lines] = lineDistributionPerProject[y.Lines] + 1
	}
}

func getLineDistributionPerProjectPerLanguage(summary []LanguageSummary, lineDistributionPerProject map[string]map[int64]int64) {
	//for _, y := range summary {
	//	//lineDistributionPerProject[y.Lines] = lineDistributionPerProject[y.Lines] + 1
	//}
}

func getLineDistributionPerFile(summary []LanguageSummary, lineDistributionPerFile map[int64]int64) {
	for _, y := range summary {
		for _, x := range y.Files {
			lineDistributionPerFile[x.Lines] = lineDistributionPerFile[x.Lines] + 1
		}
	}
}

func main() {
	queue := make(chan File, 1000)
	_ = getFiles("./json/", queue)

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
	lineDistributionPerProjectPerLanguage := map[string]map[int64]int64{}
	//lineDistributionPerFilePerLanguage := map[string]map[int64]int64{}

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

			getLineDistributionPerProjectPerLanguage(summary, lineDistributionPerProjectPerLanguage)
		}
	}

	fmt.Println("ProjectCount   ", projectCount)
	fmt.Println("LineCount      ", lineCount)
	fmt.Println("CodeCount      ", codeCount)
	fmt.Println("BlankCount     ", blankCount)
	fmt.Println("CommentCount   ", commentCount)
	fmt.Println("ComplexityCount", complexityCount)
	fmt.Println("FileCount      ", fileCount)
	fmt.Println("ByteCount      ", byteCount)
	fmt.Println(lineDistributionPerProject)
	fmt.Println(lineDistributionPerFile)
}