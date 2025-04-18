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

	"gitlab.com/snake-security/shadowspotter/pkg/kitebuilder"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFromKitebuilderAPI(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want API
	}{
		{name: "simple.json", args: args{"./testdata/simple.json"},
			want: API{
				ID:  "0cc39f78f36dbe7fe7ea94c0f2687d269d728f96",
				URL: "projectplay.xyz",
				HeaderCrumbs: []Crumb{
					RandomStringCrumb{
						Name:    "X-Developer-Key",
						Charset: ASCIIHex,
						Length:  32,
					},
				},
				Routes: []Route{
					{
						TemplatePath: "/onigokko/player",
						Method:       "post",
						HeaderCrumbs: []Crumb{
							StaticCrumb{
								K: "Token",
								V: "example-token",
							},
							StaticCrumb{
								K: "Token",
								V: "example-token",
							},
						},
						QueryCrumbs: []Crumb{
							StaticCrumb{
								K: "QueryParam",
								V: "example-query",
							},
						},
						BodyCrumbs: []Crumb{
							ObjectCrumb{
								Name: "player",
								Elements: []Crumb{
									FloatCrumb{
										Name:  "id",
										Fixed: true,
										Val:   5,
									},
									StaticCrumb{
										K: "name",
										V: "Nathan Reline",
									},
								},
							},
						},
						ContentType: []ContentType{"application/json"},
					},
				},
			},
		},
		{name: "content-type-json", args: args{"./testdata/content-type-param.json"},
			want: API{
				ID:  "0cc39f78f36dbe7fe7ea94c0f2687d269d728f96",
				URL: "projectplay.xyz",
				HeaderCrumbs: []Crumb{
					RandomStringCrumb{
						Name:    "X-Developer-Key",
						Charset: ASCIIHex,
						Length:  32,
					},
				},
				Routes: []Route{
					{
						TemplatePath: "/onigokko/player",
						Method:       "post",
						HeaderCrumbs: []Crumb{
							StaticCrumb{
								K: "Token",
								V: "example-token",
							},
							StaticCrumb{
								K: "Token",
								V: "example-token",
							},
						},
						BodyCrumbs: []Crumb{
							ObjectCrumb{
								Name: "player",
								Elements: []Crumb{
									FloatCrumb{
										Name:  "id",
										Fixed: true,
										Val:   5,
									},
									StaticCrumb{
										K: "name",
										V: "Nathan Reline",
									},
								},
							},
						},
						ContentType: []ContentType{"application/json"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSchema, err := kitebuilder.LoadJSONFile(tt.args.filename)
			assert.Nil(t, err)
			for _, v := range gotSchema {
				got, err := FromKitebuilderAPI(v)
				assert.Nil(t, err)
				spew.Dump(err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAPI_ToKitebuilderAPI_Reflection(t *testing.T) {
	tests := []struct {
		name string
		in   API
	}{
		{"simple",
			API{
				ID:  "0cc39f78f36dbe7fe7ea94c0f2687d269d728f96",
				URL: "projectplay.xyz",
				HeaderCrumbs: []Crumb{
					RandomStringCrumb{
						Name:    "X-Developer-Key",
						Charset: ASCIIHex,
						Length:  32,
					},
				},
				Routes: []Route{
					{
						TemplatePath: "/onigokko/player",
						Method:       "post",
						HeaderCrumbs: []Crumb{
							StaticCrumb{
								K: "Token",
								V: "example-token",
							},
							StaticCrumb{
								K: "Token",
								V: "example-token",
							},
						},
						BodyCrumbs: []Crumb{
							ObjectCrumb{
								Name: "player",
								Elements: []Crumb{
									FloatCrumb{
										Name:  "id",
										Fixed: true,
										Val:   5,
									},
									StaticCrumb{
										K: "name",
										V: "Nathan Reline",
									},
								},
							},
						},
						ContentType: []ContentType{"application/json"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.ToKitebuilderAPI()
			assert.Nil(t, err)

			back, err := FromKitebuilderAPI(got)
			assert.Nil(t, err)

			assert.Equal(t, tt.in, back)
		})
	}
}
