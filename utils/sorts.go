package utils

func CountTechstacks(techstacks []Techstack) (result []TechstackCount) {
	counts := make(map[string]map[string]int)

	for _, tech := range techstacks {
		if _, ok := counts[tech.Name]; !ok {
			counts[tech.Name] = make(map[string]int)
		}
		counts[tech.Name][tech.Category]++
	}

	for name, categories := range counts {
		for category, count := range categories {
			result = append(result, TechstackCount{name, category, count})
		}
	}

	return
}

func FindMaxCountPerCategory(techstacks []TechstackCount) (result []TechstackCount) {
	maxCounts := make(map[string]TechstackCount)

	for _, tech := range techstacks {
		if maxTech, ok := maxCounts[tech.Category]; !ok || tech.Count > maxTech.Count {
			maxCounts[tech.Category] = tech
		}
	}

	for _, tech := range maxCounts {
		result = append(result, tech)
	}

	return result
}
