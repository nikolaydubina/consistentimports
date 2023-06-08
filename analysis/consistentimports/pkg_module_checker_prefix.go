package consistentimports

import "strings"

// PrefixPkgModuleChecker decides that packages are in same module by looking how many prefix segments in path they match
type PrefixPkgModuleChecker struct {
	NumSegments uint
}

func (s PrefixPkgModuleChecker) IsSameModule(packagePathA, packagePathB string) bool {
	partsA := strings.Split(packagePathA, "/")
	partsB := strings.Split(packagePathB, "/")

	var countSame uint = 0
	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		if partsA[i] != partsB[i] {
			break
		}
		countSame++
	}

	return countSame >= s.NumSegments
}
