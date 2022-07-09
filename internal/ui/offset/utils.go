package offset

func countNewlines(v string) int {
	var n int
	for _, r := range v {
		if r == '\n' {
			n++
		}
	}

	return n
}
