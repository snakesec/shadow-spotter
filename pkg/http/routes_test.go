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

import "testing"

func Test_getDepth(t *testing.T) {
	type args struct {
		path  string
		depth int64
	}
	tests := []struct {
		name    string
		args    args
		wantRet string
	}{
		{"simple", args{"/foo/bar", 2}, "/foo/bar"},
		{"simple shorter", args{"/foo/bar/baz", 2}, "/foo/bar"},
		{"simple shorter again", args{"/foo", 2}, "/foo"},
		{"simple shorter root", args{"/", 2}, "/"},
		{"simple no prefix", args{"foo/bar", 2}, "/foo/bar"},
		{"shorter", args{"foo", 2}, "/foo"},
		{"longer", args{"foo/bar/baz", 2}, "/foo/bar"},
		{"longer", args{"foo/bar/baz", 0}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := getDepth(tt.args.path, tt.args.depth); gotRet != tt.wantRet {
				t.Errorf("getDepth() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
