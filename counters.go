package main

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