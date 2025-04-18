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

import "math/rand"

var (
	ASCIINum              = "0123456789"
	ASCIIHex              = "0123456789abcdefABCDEF"
	ASCIIPrintableNoSpace = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&\\'()*+,-./:;<=>?@[\\\\]^_`{|}~"
	ASCIISpecia           = "!\"#$%&\\'()*+,-./:;<=>?@[\\\\]^_`{|}~"
	ASCIIAlpha            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ASCIIAlphaNum         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultCharSet        = ASCIINum
)

func RandomString(rng *rand.Rand, charset string, length int) string {
	if len(charset) == 0 {
		charset = ASCIINum
	}
	b := make([]rune, length)
	for i := range b {
		if rng != nil {
			b[i] = rune(charset[rng.Intn(len(charset))])
		} else {
			b[i] = rune(charset[rand.Intn(len(charset))])
		}
	}
	return string(b)
}
