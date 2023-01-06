package source

import "golang.org/x/mod/semver"

type ByVersion []*Release

func (vs ByVersion) Len() int      { return len(vs) }
func (vs ByVersion) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
func (vs ByVersion) Less(i, j int) bool {
	cmp := semver.Compare(vs[i].Version, vs[j].Version)
	if cmp != 0 {
		return cmp > 0
	}
	return vs[i].Version > vs[j].Version
}
