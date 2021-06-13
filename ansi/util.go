package ansi

func arrayMax(lines []string) int {
	var max int
	for _, line := range lines {
		if len(line) > max {
			max = len(line)
		}
	}

	return max
}
