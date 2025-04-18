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

package kitebuilder

import (
	"fmt"
	"strings"

	"gitlab.com/snake-security/shadowspotter/pkg/log"
)

func PrintAPIs(apis []API) {
	headerTypes := make(map[string]int)
	routes := 0
	routeDistribution := make(map[int]int)
	parameterTypes := make(map[string]int)
	for _, api := range apis {
		url := api.URL
		if url == "" {
			url = "<no-url>"
		}

		securityOpts := make([]string, 0)
		for _, v := range api.SecurityDefinitions {
			securityOpts = append(securityOpts, fmt.Sprintf("{%s:%s(%s)}", v.In, v.Name, v.Type))
		}
		log.Info().Msgf("%s [%s] %s", url, api.ID, strings.Join(securityOpts, " "))

		apiRoutes := 0

		for path, ops := range api.Paths {
			for method, op := range ops {
				params := make([]string, 0)
				for _, p := range op.Parameters {
					params = append(params, fmt.Sprintf("{%s:%s(%s)}", p.In, p.Name, p.Type))
					parameterTypes[p.Type] += 1
					if p.Schema != nil && !p.Schema.IsZero() {
						if p.Schema.Type == "" {
							// v, _:= json.Marshal(p.Schema)
							for _, v := range p.Schema.Properties {
								if v.Type == "string" {
									headerTypes[strings.ToLower(v.Format)] += 1
								}
							}
							for _, v := range p.Schema.AllOf {
								if v.Type == "string" {
									headerTypes[strings.ToLower(v.Format)] += 1
								}
							}
						}
					}
				}
				_, _ = path, method
				log.Info().Msgf("\t%s %s %s", method, path, strings.Join(params, " "))
				routes += 1
				apiRoutes += 1
			}
		}
		routeDistribution[apiRoutes] += 1
	}
	log.Info().Interface("v", headerTypes).Msg("schema types")
	log.Info().Interface("v", parameterTypes).Msg("parameter types")
	log.Info().Interface("v", routeDistribution).Msg("route api distribution")
	log.Info().
		Int("apis", len(apis)).
		Int("routes", routes).
		Msg("analysis complete")
}
