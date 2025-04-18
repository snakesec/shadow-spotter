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

package http

import (
	"fmt"
	"strconv"
	"strings"
)

type Range struct {
	Min int
	Max int
}

func (r Range) String() string {
	if r.Min == r.Max && r.Max == 0 {
		return ""
	}
	if r.Min == r.Max {
		return strconv.Itoa(r.Min)
	}
	return fmt.Sprintf("%d-%d", r.Min, r.Max)
}

// RangeFromString will return a range from a string like 5-10
func RangeFromString(in string) (ret Range, err error) {
	if !strings.Contains(in, "-") {
		// treat it as a single value
		ret.Min, err = strconv.Atoi(in)
		if err != nil {
			return ret, fmt.Errorf("unable to parse range: %w", err)
		}
		ret.Max = ret.Min
		return ret, nil
	}
	v := strings.SplitN(in, "-", -1)
	if len(v) != 2 {
		return ret, fmt.Errorf("unexpected format for range")
	}

	ret.Min, err = strconv.Atoi(v[0])
	if err != nil {
		return ret, fmt.Errorf("unable to parse range min: %w", err)
	}

	ret.Max, err = strconv.Atoi(v[1])
	if err != nil {
		return ret, fmt.Errorf("unable to parse range max: %w", err)
	}

	if ret.Min > ret.Max {
		return ret, fmt.Errorf("invalid range. min is not lower than max")
	}
	return ret, nil
}