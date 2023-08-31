package utils

func UniqueValues(slice []int) []int {
	frequency := map[int]int{}
	var result []int

	for _, v := range slice {
		frequency[v]++
	}

	for _, item := range slice {
		if frequency[item] == 1 {
			result = append(result, item)
		}
	}

	return result
}
