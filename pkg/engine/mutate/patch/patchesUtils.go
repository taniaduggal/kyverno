package patch

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kyverno/kyverno/pkg/utils/wildcard"
	"github.com/mattbaird/jsonpatch"
)

func ConvertPatches(in ...jsonpatch.JsonPatchOperation) [][]byte {
	var out [][]byte
	for _, patch := range in {
		if patch, err := patch.MarshalJSON(); err == nil {
			out = append(out, patch)
		}
	}
	return out
}

func generatePatches(src, dst []byte) ([]jsonpatch.JsonPatchOperation, error) {
	if pp, err := jsonpatch.CreatePatch(src, dst); err != nil {
		return nil, err
	} else {
		return FilterAndSortPatches(pp), err
	}
}

// FilterAndSortPatches
// 1. filters out patches with the certain paths
// 2. sorts the removal patches(with same path) by the key of index
// in descending order. The sort is required as when removing multiple
// elements from an array, the elements must be removed in descending
// order to preserve each index value.
func FilterAndSortPatches(originalPatches []jsonpatch.JsonPatchOperation) []jsonpatch.JsonPatchOperation {
	patches := filterInvalidPatches(originalPatches)

	result := make([]jsonpatch.JsonPatchOperation, len(patches))
	index := getIndexToBeReversed(patches)

	if len(index) == 0 {
		return patches
	}

	start := 0
	for _, idx := range index {
		end := idx[0]
		copy(result[start:end], patches[start:end])
		reversedPatches := reverse(patches, idx)
		copy(result[end:], reversedPatches)
		start = idx[1] + 1
	}
	copy(result[start:], patches[start:])
	return result
}

func getIndexToBeReversed(patches []jsonpatch.JsonPatchOperation) [][]int {
	removePaths := make([]string, len(patches))

	for i, patch := range patches {
		if patch.Operation == "remove" {
			removePaths[i] = patch.Path
		}
	}
	return getRemoveInterval(removePaths)
}

func getRemoveInterval(removePaths []string) [][]int {
	// get paths end with '/number'
	regex := regexp.MustCompile(`\/\d+$`)
	for i, path := range removePaths {
		if !regex.Match([]byte(path)) {
			removePaths[i] = ""
		}
	}

	res := [][]int{}
	for i := 0; i < len(removePaths); {
		if removePaths[i] != "" {
			baseDir := filepath.Dir(removePaths[i])
			j := i + 1
			for ; j < len(removePaths); j++ {
				curDir := filepath.Dir(removePaths[j])
				if baseDir != curDir {
					break
				}
			}
			if i != j-1 {
				res = append(res, []int{i, j - 1})
			}
			i = j
		} else {
			i++
		}
	}

	return res
}

func reverse(patches []jsonpatch.JsonPatchOperation, interval []int) []jsonpatch.JsonPatchOperation {
	res := make([]jsonpatch.JsonPatchOperation, interval[1]-interval[0]+1)
	j := 0
	for i := interval[1]; i >= interval[0]; i-- {
		res[j] = patches[i]
		j++
	}
	return res
}

// filterInvalidPatches filters out patch with the following path:
// - not */metadata/name, */metadata/namespace, */metadata/labels, */metadata/annotations
// - /status
func filterInvalidPatches(patches []jsonpatch.JsonPatchOperation) []jsonpatch.JsonPatchOperation {
	res := []jsonpatch.JsonPatchOperation{}
	for _, patch := range patches {
		if ignorePatch(patch.Path) {
			continue
		}

		res = append(res, patch)
	}
	return res
}

func ignorePatch(path string) bool {
	if wildcard.Match("/spec/triggers/*/metadata/*", path) {
		return false
	}

	if wildcard.Match("*/metadata", path) {
		return false
	}

	if strings.Contains(path, "/metadata") {
		if !strings.Contains(path, "/metadata/name") &&
			!strings.Contains(path, "/metadata/namespace") &&
			!strings.Contains(path, "/metadata/annotations") &&
			!strings.Contains(path, "/metadata/labels") &&
			!strings.Contains(path, "/metadata/ownerReferences") &&
			!strings.Contains(path, "/metadata/generateName") &&
			!strings.Contains(path, "/metadata/finalizers") {
			return true
		}
	}

	return false
}
