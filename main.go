package main

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	Filename        string
	Name            string
	Content         []byte
	LanguageSummary LanguageSummary
}

type Largest struct {
	Name       string
	Location   string
	Filename   string
	Value      int64
	Bytes      int64
	Lines      int64
	Code       int64
	Comment    int64
	Blank      int64
	Complexity int64
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
					Filename: fileInfo.Name(),
					Name:     fileInfo.Name(),
					Content:  content,
				}
			}
		}
	}

	close(output)
	return nil
}

func getFilesTar(output chan File) {
	file, err := os.Open("./output.tar")
	if err != nil {
		fmt.Println("error: There is a problem with os.Open:" + err.Error())
	}

	tr := tar.NewReader(file)
	for {
		hdr, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		bs, _ := ioutil.ReadAll(tr)

		output <- File{
			Name:    hdr.Name,
			Content: bs,
		}
	}

	close(output)
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

	// below is for testing locally
	queue := make(chan File, 1000)
	//go getFiles("./json/", queue)
	go getFilesTar(queue)

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

	sourceCount := map[string]int64{} // Count of each source github/bitbucket/gitlab

	ymlOrYaml := map[string]int64{} // yaml or yml extension?

	mostComplex := Largest{}                               // Holds details of the most complex file
	mostComplexPerLanguage := map[string]Largest{}         // Most complex of each file type
	mostComplexWeighted := Largest{}                       // Holds details of the most complex file weighted by lines NB useless because it only picks up minified files
	mostComplexWeightedPerLanguage := map[string]Largest{} // Most complex of each file type weighted by lines
	largest := Largest{}                                   // Holds details of the largest file in bytes
	largestPerLanguage := map[string]Largest{}             // largest file per language
	longest := Largest{}                                   // Holds details of the longest file in lines
	longestPerLanguage := map[string]Largest{}             // longest file per language
	mostCommented := Largest{}                             // Holds details of the most commented file in lines
	mostCommentedPerLanguage := map[string]Largest{}       // most commented file per language

	pureProjects := map[int64]int64{}                      // Purity of projects, IE how many languages are used in a project by count
	pureProjectsByLanguage := map[string]map[int64]int64{} // Purity by language

	javaFactory := map[string]int64{} // Count of factoryfactoryfactory factoryfactory and factory

	cursingByLanguage := map[string]int64{} // Cursing names by language
	cursingByWord := map[string]int64{} // Cursing names by most common curse word

	multipleGitIgnore := map[int64]int64{} // See how many projects use none, single or multiple gitignore files

	hasCoffeeScriptAndTypescript := map[string]int64{} // Count of projects with both languages
	hasTypeScriptExclusively := map[string]int64{} // Count of projects with just typescript

	upperLowerOrMixedCase := map[string]int64{} // Count if we have upper/lower or mixed case in the name
	upperLowerOrMixedCaseIgnoreExt := map[string]int64{} // Count if we have upper/lower or mixed case in the name ignoring ext

	count := 0
	for file := range queue {
		count++
		if count%100 == 0 {
			fmt.Println("Processing", file.Name, count)
		}

		summary, err := unmarshallContent(file.Content)

		if strings.HasPrefix(file.Name, "bitbucket") {
			sourceCount["bitbucket"] = sourceCount["bitbucket"] + 1
		}
		if strings.HasPrefix(file.Name, "github") {
			sourceCount["github"] = sourceCount["github"] + 1
		}
		if strings.HasPrefix(file.Name, "gitlab") {
			sourceCount["gitlab"] = sourceCount["gitlab"] + 1
		}

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

			// Purge counts as mentioned above to avoid blowing memory budget
			if count%100 == 0 {
				fileNamesCount = cullCountMap(fileNamesCount)
				fileNamesNoExtensionCount = cullCountMap(fileNamesNoExtensionCount)
				fileNamesNoExtensionLowercaseCount = cullCountMap(fileNamesNoExtensionLowercaseCount)
			}

			mostComplex = getMostComplex(summary, file.Filename, mostComplex)
			getMostComplexPerLanguage(summary, file.Filename, mostComplexPerLanguage)

			mostComplexWeighted = getMostComplexWeighted(summary, file.Filename, mostComplexWeighted)
			getMostComplexWeightedPerLanguage(summary, file.Filename, mostComplexWeightedPerLanguage)

			largest = getLargest(summary, file.Filename, largest)
			getLargestPerLanguage(summary, file.Filename, largestPerLanguage)
			longest = getMostLines(summary, file.Filename, longest)
			getMostLinesPerLanguage(summary, file.Filename, longestPerLanguage)

			mostCommented = getMostCommented(summary, file.Filename, mostCommented)
			getMostCommentedPerLanguage(summary, file.Filename, mostCommentedPerLanguage)

			getYmlOrYaml(summary, ymlOrYaml)
			getPurity(summary, pureProjects)
			getPurityByLanguage(summary, pureProjectsByLanguage)
			getFactoryCount(summary, javaFactory)
			getCursingByLanguage(summary, cursingByLanguage)
			getCursingByWord(summary, cursingByWord)
			getGitIgnore(summary, multipleGitIgnore)

			getHasCoffeeScriptAndTypeScript(summary, hasCoffeeScriptAndTypescript)
			getHasTypeScriptExclusively(summary, hasTypeScriptExclusively)

			getUpperLowerOrMixedCase(summary, upperLowerOrMixedCase)
			getUpperLowerOrMixedCaseIgnoreExt(summary, upperLowerOrMixedCaseIgnoreExt)
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

	v, _ = json.Marshal(sourceCount)
	_ = ioutil.WriteFile("./results/sourceCount.json", []byte(v), 0600)

	v, _ = json.Marshal(mostComplex)
	_ = ioutil.WriteFile("./results/mostComplex.json", []byte(v), 0600)
	v, _ = json.Marshal(mostComplexPerLanguage)
	_ = ioutil.WriteFile("./results/mostComplexPerLanguage.json", []byte(v), 0600)
	v, _ = json.Marshal(mostComplexWeighted)
	_ = ioutil.WriteFile("./results/mostComplexWeighted.json", []byte(v), 0600)
	v, _ = json.Marshal(mostComplexWeightedPerLanguage)
	_ = ioutil.WriteFile("./results/mostComplexWeightedPerLanguage.json", []byte(v), 0600)
	v, _ = json.Marshal(largest)
	_ = ioutil.WriteFile("./results/largest.json", []byte(v), 0600)
	v, _ = json.Marshal(largestPerLanguage)
	_ = ioutil.WriteFile("./results/largestPerLanguage.json", []byte(v), 0600)
	v, _ = json.Marshal(mostCommented)
	_ = ioutil.WriteFile("./results/mostCommented.json", []byte(v), 0600)
	v, _ = json.Marshal(mostCommentedPerLanguage)
	_ = ioutil.WriteFile("./results/mostCommentedPerLanguage.json", []byte(v), 0600)
	v, _ = json.Marshal(longest)
	_ = ioutil.WriteFile("./results/longest.json", []byte(v), 0600)
	v, _ = json.Marshal(longestPerLanguage)
	_ = ioutil.WriteFile("./results/longestPerLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(ymlOrYaml)
	_ = ioutil.WriteFile("./results/ymlOrYaml.json", []byte(v), 0600)

	v, _ = json.Marshal(pureProjects)
	_ = ioutil.WriteFile("./results/pureProjects.json", []byte(v), 0600)
	v, _ = json.Marshal(pureProjectsByLanguage)
	_ = ioutil.WriteFile("./results/pureProjectsByLanguage.json", []byte(v), 0600)

	v, _ = json.Marshal(javaFactory)
	_ = ioutil.WriteFile("./results/javaFactory.json", []byte(v), 0600)

	v, _ = json.Marshal(cursingByLanguage)
	_ = ioutil.WriteFile("./results/cursingByLanguage.json", []byte(v), 0600)
	v, _ = json.Marshal(cursingByWord)
	_ = ioutil.WriteFile("./results/cursingByWord.json", []byte(v), 0600)

	v, _ = json.Marshal(multipleGitIgnore)
	_ = ioutil.WriteFile("./results/multipleGitIgnore.json", []byte(v), 0600)

	v, _ = json.Marshal(hasCoffeeScriptAndTypescript)
	_ = ioutil.WriteFile("./results/hasCoffeeScriptAndTypescript.json", []byte(v), 0600)

	v, _ = json.Marshal(hasTypeScriptExclusively)
	_ = ioutil.WriteFile("./results/hasTypeScriptExclusively.json", []byte(v), 0600)

	v, _ = json.Marshal(upperLowerOrMixedCase)
	_ = ioutil.WriteFile("./results/upperLowerOrMixedCase.json", []byte(v), 0600)
	v, _ = json.Marshal(upperLowerOrMixedCaseIgnoreExt)
	_ = ioutil.WriteFile("./results/upperLowerOrMixedCaseIgnoreExt.json", []byte(v), 0600)
}
