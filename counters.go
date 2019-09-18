package main

import (
	"strings"
)

func getLineDistributionPerProject(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		x[y.Lines] = x[y.Lines] + 1
	}
}

func getLineDistributionPerLanguage(summary []LanguageSummary, x map[string]map[int64]int64) {
	for _, y := range summary {
		_, ok := x[y.Name]

		if ok {
			m := x[y.Name]
			m[y.Lines] = m[y.Lines] + 1
			x[y.Name] = m
		} else {
			m := map[int64]int64{}
			m[y.Lines] = 1
			x[y.Name] = m
		}
	}
}

func getLineDistributionPerFile(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		for _, z := range y.Files {
			x[z.Lines] = x[z.Lines] + 1
		}
	}
}

////////////////////////////////////

func getCodeDistributionPerProject(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		x[y.Code] = x[y.Code] + 1
	}
}

func getCodeDistributionPerLanguage(summary []LanguageSummary, x map[string]map[int64]int64) {
	for _, y := range summary {
		_, ok := x[y.Name]

		if ok {
			m := x[y.Name]
			m[y.Code] = m[y.Code] + 1
			x[y.Name] = m
		} else {
			m := map[int64]int64{}
			m[y.Code] = 1
			x[y.Name] = m
		}
	}
}

func getCodeDistributionPerFile(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		for _, z := range y.Files {
			x[z.Code] = x[z.Code] + 1
		}
	}
}

///////////////////////////////////

func getCommentDistributionPerProject(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		x[y.Comment] = x[y.Comment] + 1
	}
}

func getCommentDistributionPerLanguage(summary []LanguageSummary, x map[string]map[int64]int64) {
	for _, y := range summary {
		_, ok := x[y.Name]

		if ok {
			m := x[y.Name]
			m[y.Comment] = m[y.Comment] + 1
			x[y.Name] = m
		} else {
			m := map[int64]int64{}
			m[y.Comment] = 1
			x[y.Name] = m
		}
	}
}

func getCommentDistributionPerFile(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		for _, z := range y.Files {
			x[z.Comment] = x[z.Comment] + 1
		}
	}
}

///////////////////////////////////

func getComplexityDistributionPerProject(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		x[y.Complexity] = x[y.Complexity] + 1
	}
}

func getComplexityDistributionPerLanguage(summary []LanguageSummary, x map[string]map[int64]int64) {
	for _, y := range summary {
		_, ok := x[y.Name]

		if ok {
			m := x[y.Name]
			m[y.Complexity] = m[y.Complexity] + 1
			x[y.Name] = m
		} else {
			m := map[int64]int64{}
			m[y.Complexity] = 1
			x[y.Name] = m
		}
	}
}

func getComplexityDistributionPerFile(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		for _, z := range y.Files {
			x[z.Complexity] = x[z.Complexity] + 1
		}
	}
}

///////////////////////////////////

func getBlankDistributionPerProject(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		x[y.Blank] = x[y.Blank] + 1
	}
}

func getBlankDistributionPerLanguage(summary []LanguageSummary, x map[string]map[int64]int64) {
	for _, y := range summary {
		_, ok := x[y.Name]

		if ok {
			m := x[y.Name]
			m[y.Blank] = m[y.Blank] + 1
			x[y.Name] = m
		} else {
			m := map[int64]int64{}
			m[y.Blank] = 1
			x[y.Name] = m
		}
	}
}

func getBlankDistributionPerFile(summary []LanguageSummary, x map[int64]int64) {
	for _, y := range summary {
		for _, z := range y.Files {
			x[z.Blank] = x[z.Blank] + 1
		}
	}
}

////////////////////////

func getFilesPerProject(summary []LanguageSummary, x map[int64]int64) {
	var c int64
	for _, y := range summary {
		c += int64(len(y.Files))
	}

	x[c] = x[c] + 1
}

////////////////////////

func getComplexityPerLanguage(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		x[y.Name] = x[y.Name] + y.Complexity
	}
}

/////////////////////////

func getProjectsPerLanguage(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		x[y.Name] = x[y.Name] + 1
	}
}

/////////////////////////

// Gets the count of files per language name
func getFilesPerLanguage(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		x[y.Name] = x[y.Name] + int64(len(y.Files))
	}
}

/////////////////////////

func getLicencePerProject(summary []LanguageSummary, x map[string]int64) {
	hasLicence := false

	for _, y := range summary {
		if y.Name == "License" {
			hasLicence = true
		}
	}

	if hasLicence {
		x["Yes"] = x["Yes"] + 1
	} else {
		x["No"] = x["No"] + 1
	}
}

/////////////////////////

func getFileNamesCount(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		for _, z := range y.Files {
			x[z.Filename] = x[z.Filename] + 1
		}
	}
}

func getFileNamesNoExtensionCount(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		for _, z := range y.Files {

			l := strings.IndexAny(z.Filename, ".")

			if l == -1 || l == 0 {
				x[z.Filename] = x[z.Filename] + 1
			} else {
				x[z.Filename[:l]] = x[z.Filename[:l]] + 1
			}
		}
	}
}

func getFileNamesNoExtensionLowercaseCount(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		for _, z := range y.Files {

			n := strings.ToLower(z.Filename)
			l := strings.IndexAny(n, ".")

			if l == -1 {
				x[n] = x[n] + 1
			} else {
				x[n[:l]] = x[n[:l]] + 1
			}
		}
	}
}

//////////////

func getMostComplex(summary []LanguageSummary, filename string, x Largest) Largest {
	for _, y := range summary {
		for _, z := range y.Files {
			if z.Complexity > x.Value {
				x = Largest{
					Value:    z.Complexity,
					Name:     filename,
					Location: z.Location,
					Filename: z.Filename,
					Blank:    z.Blank,
					Comment:  z.Comment,
					Code:     z.Code,
					Lines:    z.Lines,
					Bytes:    z.Bytes,
				}
			}
		}
	}

	return x
}

func getMostComplexPerLanguage(summary []LanguageSummary, filename string, x map[string]Largest) {
	for _, y := range summary {
		for _, z := range y.Files {

			if z.Complexity > x[z.Language].Value {
				x[z.Language] = Largest{
					Value:    z.Complexity,
					Name:     filename,
					Location: z.Location,
					Filename: z.Filename,
					Blank:    z.Blank,
					Comment:  z.Comment,
					Code:     z.Code,
					Lines:    z.Lines,
					Bytes:    z.Bytes,
				}
			}
		}
	}
}

//////////////

func getMostComplexWeighted(summary []LanguageSummary, filename string, x Largest) Largest {
	for _, y := range summary {
		for _, z := range y.Files {
			var t int64

			if z.Lines != 0 {
				t = int64(z.Complexity / z.Lines)
			}

			if t > x.Value {
				x = Largest{
					Value:    t,
					Name:     filename,
					Location: z.Location,
					Filename: z.Filename,
					Blank:    z.Blank,
					Comment:  z.Comment,
					Code:     z.Code,
					Lines:    z.Lines,
					Bytes:    z.Bytes,
				}
			}
		}
	}

	return x
}

func getMostComplexWeightedPerLanguage(summary []LanguageSummary, filename string, x map[string]Largest) {
	for _, y := range summary {
		for _, z := range y.Files {
			var t int64

			if z.Lines != 0 {
				t = int64(z.Complexity / z.Lines)
			}

			if t > x[z.Language].Value {
				x[z.Language] = Largest{
					Value:    t,
					Name:     filename,
					Location: z.Location,
					Filename: z.Filename,
					Blank:    z.Blank,
					Comment:  z.Comment,
					Code:     z.Code,
					Lines:    z.Lines,
					Bytes:    z.Bytes,
				}
			}
		}
	}
}

//////////////

func getLargest(summary []LanguageSummary, filename string, x Largest) Largest {
	for _, y := range summary {
		for _, z := range y.Files {
			if z.Bytes > x.Value {
				x = Largest{
					Value:      z.Bytes,
					Name:       filename,
					Location:   z.Location,
					Filename:   z.Filename,
					Blank:      z.Blank,
					Comment:    z.Comment,
					Code:       z.Code,
					Lines:      z.Lines,
					Bytes:      z.Bytes,
					Complexity: z.Complexity,
				}
			}
		}
	}

	return x
}

func getLargestPerLanguage(summary []LanguageSummary, filename string, x map[string]Largest) {
	for _, y := range summary {
		for _, z := range y.Files {

			if z.Bytes > x[z.Language].Value {
				x[z.Language] = Largest{
					Value:      z.Bytes,
					Name:       filename,
					Location:   z.Location,
					Filename:   z.Filename,
					Blank:      z.Blank,
					Comment:    z.Comment,
					Code:       z.Code,
					Lines:      z.Lines,
					Bytes:      z.Bytes,
					Complexity: z.Complexity,
				}
			}
		}
	}
}

//////////////

func getMostCommented(summary []LanguageSummary, filename string, x Largest) Largest {
	for _, y := range summary {
		for _, z := range y.Files {
			if z.Comment > x.Value {
				x = Largest{
					Value:      z.Comment,
					Name:       filename,
					Location:   z.Location,
					Filename:   z.Filename,
					Blank:      z.Blank,
					Comment:    z.Comment,
					Code:       z.Code,
					Lines:      z.Lines,
					Bytes:      z.Bytes,
					Complexity: z.Complexity,
				}
			}
		}
	}

	return x
}

func getMostCommentedPerLanguage(summary []LanguageSummary, filename string, x map[string]Largest) {
	for _, y := range summary {
		for _, z := range y.Files {

			if z.Comment > x[z.Language].Value {
				x[z.Language] = Largest{
					Value:      z.Comment,
					Name:       filename,
					Location:   z.Location,
					Filename:   z.Filename,
					Blank:      z.Blank,
					Comment:    z.Comment,
					Code:       z.Code,
					Lines:      z.Lines,
					Bytes:      z.Bytes,
					Complexity: z.Complexity,
				}
			}
		}
	}
}

///////////////

func getMostLines(summary []LanguageSummary, filename string, x Largest) Largest {
	for _, y := range summary {
		for _, z := range y.Files {
			if z.Lines > x.Value {
				x = Largest{
					Value:      z.Lines,
					Name:       filename,
					Location:   z.Location,
					Filename:   z.Filename,
					Blank:      z.Blank,
					Comment:    z.Comment,
					Code:       z.Code,
					Lines:      z.Lines,
					Bytes:      z.Bytes,
					Complexity: z.Complexity,
				}
			}
		}
	}

	return x
}

func getMostLinesPerLanguage(summary []LanguageSummary, filename string, x map[string]Largest) {
	for _, y := range summary {
		for _, z := range y.Files {

			if z.Lines > x[z.Language].Value {
				x[z.Language] = Largest{
					Value:      z.Lines,
					Name:       filename,
					Location:   z.Location,
					Filename:   z.Filename,
					Blank:      z.Blank,
					Comment:    z.Comment,
					Code:       z.Code,
					Lines:      z.Lines,
					Bytes:      z.Bytes,
					Complexity: z.Complexity,
				}
			}
		}
	}
}

///////////////

func getYmlOrYaml(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		if y.Name == "YAML" {
			for _, z := range y.Files {

				if z.Extension == "yml" || strings.HasSuffix(z.Extension, ".yml") {
					x["yml"] = x["yml"] + 1
				} else if z.Extension == "yaml" || strings.HasSuffix(z.Extension, ".yaml") {
					x["yaml"] = x["yaml"] + 1
				} else {
					x["other"] = x["other"] + 1
				}
			}
		}
	}
}

func getPurity(summary []LanguageSummary, x map[int64]int64) {
	x[int64(len(summary))] = x[int64(len(summary))] + 1
}

func getPurityByLanguage(summary []LanguageSummary, x map[string]map[int64]int64) {
	for _, y := range summary {
		_, ok := x[y.Name]

		if ok {
			m := x[y.Name]
			m[int64(len(summary))] = m[int64(len(summary))] + 1
			x[y.Name] = m
		} else {
			m := map[int64]int64{}
			m[int64(len(summary))] = 1
			x[y.Name] = m
		}
	}
}

func getFactoryCount(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {
		if y.Name == "Java" {
			for _, z := range y.Files {
				x["count"] = x["count"] + 1

				n := strings.ToLower(z.Filename)

				if strings.Contains(n, "factoryfactoryfactory") {
					x["factoryfactoryfactory"] = x["factoryfactoryfactory"] + 1
				} else if strings.Contains(n, "factoryfactory") {
					x["factoryfactory"] = x["factoryfactory"] + 1
				} else if strings.Contains(n, "factory") {
					x["factory"] = x["factory"] + 1
				}
			}
		}
	}
}

func getCursingByLanguage(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {

		for _, z := range y.Files {
			if containsCurse(z.Filename) != "" {
				x[y.Name] = x[y.Name] + 1
			}
		}
	}
}

func getCursingByWord(summary []LanguageSummary, x map[string]int64) {
	for _, y := range summary {

		for _, z := range y.Files {
			c := containsCurse(z.Filename)
			if c != "" {
				x[c] = x[c] + 1
			}
		}
	}
}

func containsCurse(name string) string {
	l := strings.ToLower(name)

	for _, c := range curseWords {
		if strings.HasPrefix(l, c+".") {
			return c
		}
	}

	return ""
}

func getGitIgnore(summary []LanguageSummary, x map[int64]int64) {
	hasGitignore := false
	for _, y := range summary {

		if y.Name == "gitignore" {
			hasGitignore = true
			x[int64(len(y.Files))] = x[int64(len(y.Files))] + 1
		}
	}

	if !hasGitignore {
		x[0] = x[0] + 1
	}
}

func getHasCoffeeScriptAndTypeScript(summary []LanguageSummary, x map[string]int64) {
	hasCoffeeScript := false
	hasTypeScript := false

	for _, y := range summary {

		if y.Name == "CoffeeScript" {
			hasCoffeeScript = true
		}

		if y.Name == "TypeScript" {
			hasTypeScript = true
		}
	}

	if hasTypeScript && hasCoffeeScript {
		x["Both"] = x["Both"] + 1
	} else {
		x["Nope"] = x["Nope"] + 1
	}
}

func getHasTypeScriptExclusively(summary []LanguageSummary, x map[string]int64) {
	hasJavaScript := false
	hasTypeScript := false

	for _, y := range summary {
		if y.Name == "JavaScript" {
			hasJavaScript = true
		}

		if y.Name == "TypeScript" {
			hasTypeScript = true
		}
	}

	if hasTypeScript && !hasJavaScript {
		x["Exclusive"] = x["Exclusive"] + 1
	} else {
		x["Both"] = x["Both"] + 1
	}
}

func getUpperLowerOrMixedCase(summary []LanguageSummary, x map[string]int64) {
	hasUpper := false
	hasLower := false

	for _, y := range summary {
		for _, t := range y.Files {
			for _, x := range t.Filename {
				z := string(x)

				if strings.ToUpper(z) == z {
					hasUpper = true
				}

				if strings.ToLower(z) == z {
					hasLower = true
				}
			}
		}
	}

	if hasUpper && hasLower {
		x["Both"] = x["Both"] + 1
		return
	}

	if hasUpper {
		x["Upper"] = x["Upper"] + 1
	}

	if hasLower {
		x["Lower"] = x["Lower"] + 1
	}
}

func getUpperLowerOrMixedCaseIgnoreExt(summary []LanguageSummary, x map[string]int64) {
	hasUpper := false
	hasLower := false

	for _, y := range summary {
		for _, t := range y.Files {

			name := t.Filename
			index := strings.Index(t.Filename, ".")

			if index != -1 {
				name = t.Filename[:index]
			}

			for _, x := range name {
				z := string(x)

				if strings.ToUpper(z) == z {
					hasUpper = true
				}

				if strings.ToLower(z) == z {
					hasLower = true
				}
			}
		}
	}

	if hasUpper && hasLower {
		x["Both"] = x["Both"] + 1
		return
	}

	if hasUpper {
		x["Upper"] = x["Upper"] + 1
	}

	if hasLower {
		x["Lower"] = x["Lower"] + 1
	}
}
