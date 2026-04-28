package utils

import (
	"sort"
)

// MergeStringMaps merges multiple maps into one. Later maps overwrite earlier ones
// for duplicate keys.
func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// SortedKeys returns the keys of a map in sorted order.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// InvertMap swaps keys and values in a map. If there are duplicate values,
// only one will be preserved (non-deterministic).
func InvertMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[v] = k
	}
	return result
}

// FilterMap returns a new map containing only entries where the predicate returns true.
func FilterMap(m map[string]string, predicate func(key, value string) bool) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}
