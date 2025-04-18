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

package kiterunner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWildcardResponses_AddWildcard(t *testing.T) {
	type args struct {
		wr WildcardResponse
	}
	tests := []struct {
		name     string
		w        WildcardResponses
		args     args
		expected WildcardResponses
	}{
		{"nil", nil, args{WildcardResponse{DefaultWordCount: 1}}, WildcardResponses{{DefaultWordCount: 1}}},
		{"simple", WildcardResponses{}, args{WildcardResponse{DefaultWordCount: 1}}, WildcardResponses{{DefaultWordCount: 1}}},
		{"simple + 1", WildcardResponses{{DefaultWordCount: 2}}, args{WildcardResponse{DefaultWordCount: 1}}, WildcardResponses{{DefaultWordCount: 2}, {DefaultWordCount: 1}}},
		{"dedupe", WildcardResponses{{DefaultWordCount: 1}}, args{WildcardResponse{DefaultWordCount: 1}}, WildcardResponses{{DefaultWordCount: 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w, _ = tt.w.UniqueAdd(tt.args.wr)
			assert.ElementsMatch(t, tt.w, tt.expected)
		})
	}
}
