package arch

import (
	"strconv"
	"strings"
)

// compareVersions compares two version strings in the same way as arch repo‑add.
// It returns 1 if v1 > v2, -1 if v1 < v2, and 0 if equal.
func compareVersions(v1, v2 string) int {
	epoch1, ver1 := splitEpoch(v1)
	epoch2, ver2 := splitEpoch(v2)
	if epoch1 != epoch2 {
		if epoch1 > epoch2 {
			return 1
		}
		return -1
	}
	return vercmp(ver1, ver2)
}

// splitEpoch splits a version string into its epoch and the rest.
// If no epoch is present, it returns 0 as the epoch.
func splitEpoch(version string) (int, string) {
	i := strings.IndexByte(version, ':') //nolint:modernize // more performant
	if i == -1 {
		return 0, version
	}
	epoch, _ := strconv.Atoi(version[:i])
	return epoch, version[i+1:]
}

// vercmp compares two version strings (without epoch) segment by segment.
//
//nolint:gocognit
func vercmp(a, b string) int {
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		segA, isNumA := getVersionSegment(a, i)
		segB, isNumB := getVersionSegment(b, j)

		switch {
		case isNumA && isNumB:
			// Trim leading zeros
			segATrim := strings.TrimLeft(segA, "0")
			segBTrim := strings.TrimLeft(segB, "0")
			if segATrim == "" {
				segATrim = "0"
			}
			if segBTrim == "" {
				segBTrim = "0"
			}
			// Compare numeric values by length first, then lexicographically.
			if len(segATrim) != len(segBTrim) {
				if len(segATrim) > len(segBTrim) {
					return 1
				}
				return -1
			}
			if segATrim != segBTrim {
				if segATrim > segBTrim {
					return 1
				}
				return -1
			}
		case isNumA != isNumB:
			// Numeric segments are considered lower than non-numeric.
			if isNumA {
				return -1
			}
			return 1
		default:
			// Compare non-numeric segments character by character.
			minLen := min(len(segA), len(segB))
			for k := range minLen {
				oa := versionCharOrder(segA[k])
				ob := versionCharOrder(segB[k])
				if oa != ob {
					if oa > ob {
						return 1
					}
					return -1
				}
			}
			if len(segA) != len(segB) {
				if len(segA) > len(segB) {
					return 1
				}
				return -1
			}
		}

		i += len(segA)
		j += len(segB)
	}
	return 0
}

// getVersionSegment returns the next contiguous segment from s starting at index i.
// It returns the segment and a boolean indicating whether the segment is numeric.
func getVersionSegment(s string, i int) (string, bool) {
	if i >= len(s) {
		return "", false
	}
	// Determine whether the segment is numeric.
	isNum := (s[i] >= '0' && s[i] <= '9')
	start := i
	for i < len(s) {
		if isNum {
			if s[i] < '0' || s[i] > '9' {
				break
			}
		} else {
			if s[i] >= '0' && s[i] <= '9' {
				break
			}
		}
		i++
	}
	return s[start:i], isNum
}

// versionCharOrder returns an integer that represents the sort order for a given character.
// Tilde (~) sorts lower than any other character.
func versionCharOrder(c byte) int {
	if c == '~' {
		return -1
	}
	return int(c)
}
