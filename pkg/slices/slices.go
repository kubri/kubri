package slices

func Filter[S ~[]E, E any](s S, filter func(E) bool) S {
	var i int
	for _, e := range s {
		if filter(e) {
			s[i] = e
			i++
		}
	}

	// Overwrite truncated elements to allow pointer values to be garbage collected.
	var zero E
	for j := i; j < len(s); j++ {
		s[j] = zero
	}

	return s[:i]
}
