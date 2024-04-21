package maputils

import "strings"

func RecursiveGet(key string, m map[string]any) (any, bool) {
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		v, found := m[keys[0]]
		return v, found
	}

	v, found := m[keys[0]].(map[string]any)
	if !found {
		return nil, false
	}

	return RecursiveGet(strings.Join(keys[1:], "."), v)
}
