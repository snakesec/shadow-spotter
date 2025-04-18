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
	"strings"
)

type HeaderField struct {
	Key Field
	Value Field
}

type FieldType int
const (
	String FieldType = iota
	UUID
	Int
	Date
	Timestamp
	Format
)

type Field struct {
	Key string `toml:"key" json:"key" mapstructure:"key"`
	Type FieldType `toml:"type" json:"type" mapstructure:"type"`
}

func (f *Field) Bytes() []byte {
	if len(f.Key) == 0 {
		return nil
	}
	return []byte(f.Key)
}

// StringToField will break down a URL path into fields
// This is a convenience function to convert /foo/bar into []Field{{"foo"}, {"bar"}}
func StringToFields(in string) (ret []Field) {
	for _, v := range strings.Split(in, "/") {
		ret = append(ret, Field{Key: v})
	}
	return ret
}

// Write the string to the specified writer
func (f *Field) AppendBytes(dst []byte) []byte {
	return append(dst, f.Key...)
}