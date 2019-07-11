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

func getFileCount(summary []LanguageSummary) int64 {
	var x int64

	for _, y := range summary {
		x += int64(len(y.Files))
	}

	return x
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
	var noFiles int64

	// The distributions are all buckets where things are grouped together
	// E.G. lineDistributionPerProject counts projects by the number of lines
	// so if the first project of 100 lines is added the map will look like 100:1
	// and if another is added of 100 lines it will be 100:2 and then if
	// we add a project of 50 lines it will be 100:2,50:1
	// And so on for each of the below
	// The exceptions are the map in maps, which add an additional spin by
	// counting per language so the above if the languages were C# and PHP would
	// be C# [100:1], PHP [100:1,50:1]
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

	filesPerProject := map[int64]int64{}      // Number of files in each project in buckets IE projects with 10 files or projects with 2
	projectsPerLanguage := map[string]int64{} // Number of projects which use a language
	filesPerLanguage := map[string]int64{}    // Number of files per language
	hasLicenceCount := map[string]int64{}     // Count of if a project has a licence file or not

	fileNamesCount := map[string]int64{}                     // Count of filenames
	fileNamesNoExtensionCount := map[string]int64{}          // Count of filenames without extensions
	fileNamesNoExtensionLowercaseCount := map[string]int64{} // Count of filenames tolower and no extensions
	complexityPerLanguage := map[string]int64{}              // Sum of complexity per language

	for file := range queue {
		summary, err := unmarshallContent(file.Content)

		// If no files we should exclude but count that
		cnt := getFileCount(summary)
		if cnt == 0 {
			noFiles++
		}

		if err == nil && cnt != 0 {
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
			getLicencePerProject(summary, hasLicenceCount)

			getComplexityPerLanguage(summary, complexityPerLanguage)
			getFilesPerLanguage(summary, filesPerLanguage)
			getProjectsPerLanguage(summary, projectsPerLanguage)

			// NB this might be too large and need to purge certain names over time
			getFileNamesCount(summary, fileNamesCount)
			getFileNamesNoExtensionCount(summary, fileNamesNoExtensionCount)
			getFileNamesNoExtensionLowercaseCount(summary, fileNamesNoExtensionLowercaseCount)
		}
	}

	endTime := makeTimestampSeconds()

	fmt.Println("NoFiles        ", noFiles)
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
	fmt.Println("TotalTime(s)   ", endTime-startTime)

	stats := fmt.Sprintf("NoFiles %d\nProjectCount %d\nLineCount %d\nCodeCount %d\nBlankCount %d\nCommentCount %d\nComplexityCount %d\nFileCount %d\nByteCount %d\nTotalTime(s) %d", noFiles, projectCount, lineCount, codeCount, blankCount, commentCount, complexityCount, fileCount, byteCount, endTime-startTime)
	_ = ioutil.WriteFile("./results/totalStats.txt", []byte(stats), 0600)

	v, _ := json.Marshal(lineDistributionPerProject)
	_ = ioutil.WriteFile("./results/lineDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(lineDistributionPerFile)
	_ = ioutil.WriteFile("./results/lineDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(lineDistributionPerLanguage)
	_ = ioutil.WriteFile("./results/lineDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(codeDistributionPerProject)
	_ = ioutil.WriteFile("./results/codeDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(codeDistributionPerFile)
	_ = ioutil.WriteFile("./results/codeDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(codeDistributionPerLanguage)
	_ = ioutil.WriteFile("./results/codeDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(commentDistributionPerProject)
	_ = ioutil.WriteFile("./results/commentDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(commentDistributionPerFile)
	_ = ioutil.WriteFile("./results/commentDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(commentDistributionPerLanguage)
	_ = ioutil.WriteFile("./results/commentDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(blankDistributionPerProject)
	_ = ioutil.WriteFile("./results/blankDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(blankDistributionPerFile)
	_ = ioutil.WriteFile("./results/blankDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(blankDistributionPerLanguage)
	_ = ioutil.WriteFile("./results/blankDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(complexityDistributionPerProject)
	_ = ioutil.WriteFile("./results/complexityDistributionPerProject.json", []byte(v), 0600)
	v, _ = json.Marshal(complexityDistributionPerFile)
	_ = ioutil.WriteFile("./results/complexityDistributionPerFile.json", []byte(v), 0600)
	v, _ = json.Marshal(complexityDistributionPerLanguage)
	_ = ioutil.WriteFile("./results/complexityDistributionPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(filesPerProject)
	_ = ioutil.WriteFile("./results/filesPerProject.json", []byte(v), 0600)

	v, _ = json.Marshal(projectsPerLanguage)
	_ = ioutil.WriteFile("./results/projectsPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(hasLicenceCount)
	_ = ioutil.WriteFile("./results/hasLicenceCount.json", []byte(v), 0600)

	v, _ = json.Marshal(fileNamesCount)
	_ = ioutil.WriteFile("./results/fileNamesCount.json", []byte(v), 0600)

	v, _ = json.Marshal(fileNamesNoExtensionCount)
	_ = ioutil.WriteFile("./results/fileNamesNoExtensionCount.json", []byte(v), 0600)

	v, _ = json.Marshal(fileNamesNoExtensionLowercaseCount)
	_ = ioutil.WriteFile("./results/fileNamesNoExtensionLowercaseCount.json", []byte(v), 0600)

	v, _ = json.Marshal(complexityPerLanguage)
	_ = ioutil.WriteFile("./results/complexityPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(filesPerLanguage)
	_ = ioutil.WriteFile("./results/filesPerLanguage.json", []byte(v), 0600)
}
