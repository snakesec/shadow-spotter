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

package proute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromStringSlice(t *testing.T) {
	type args struct {
		in     []string
		source string
		opts   []PRouteOption
	}
	tests := []struct {
		name    string
		args    args
		want    API
		wantErr bool
	}{
		{"simple", args{[]string{"/foo"}, "sometext.file", nil}, API{URL: "sometext.file", Routes: []Route{{Method: "GET", TemplatePath: "/foo"}}}, false},
		{"simple add slash", args{[]string{"foo"}, "sometext.file", nil}, API{URL: "sometext.file", Routes: []Route{{Method: "GET", TemplatePath: "/foo"}}}, false},
		{"simple two paths", args{[]string{"/foo", "/bar"}, "sometext.file", nil}, API{URL: "sometext.file", Routes: []Route{{Method: "GET", TemplatePath: "/foo"},{Method: "GET", TemplatePath: "/bar"}}}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromStringSlice(tt.args.in, tt.args.source, tt.args.opts...)
			if !tt.wantErr {
				assert.Nil(t, err)
			}

			tt.want.ID = got.ID
			assert.Equal(t, tt.want, got)
		})
	}
}
