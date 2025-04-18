/*
Shadow-Spotter Next Gen Content Discovery
Copyright (C) 2024  Weidsom Nascimento - SNAKE Security

Based on kiterunner from AssetNote

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package convert

// IntMapToSlice will return all the keys for the given intmap
func IntMapToSlice(m map[int]interface{}) (ret []int) {
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

// IntSliceToMap will return a map with keys created from the slice
func IntSliceToMap(v []int) map[int]interface{} {
	ret := make(map[int]interface{})
	for _, vv := range v {
		ret[vv] = struct{}{}
	}
	return ret
}

// StringMapToSlice will return all the keys for the given string map
func StringMapToSlice(m map[string]interface{}) (ret []string) {
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

// UniqueStrings will remove duplicates preserving order of the input
func UniqueStrings(in []string) (out []string) {
	seen := make(map[string]interface{})
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return
}
