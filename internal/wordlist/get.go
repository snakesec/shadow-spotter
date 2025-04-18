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

package wordlist

import (
	"context"
	"fmt"
)

// Get will retrieve all the names from remote or local sources
func Get(ctx context.Context, names ...string) (ret []WordlistMetadata, err error) {
	if err := CreateLocalCache(); err != nil {
		return nil, fmt.Errorf("failed to create local cache dir: %w", err)
	}

	// don't download existing things in the cache
	local, err := GetLocalDirListing()
	if err != nil {
		return nil, fmt.Errorf("failed to get local dir listing: %w", err)
	}

	got := make(map[string]WordlistMetadata)
	for _, v := range local {
		got[v.Filename] = v
		got[v.Shortname] = v
	}

	missing := make([]string, 0)
	for _, v := range names {
		if w, ok := got[v]; ok {
			ret = append(ret, w)
			continue
		}
		missing = append(missing, v)
	}

	if len(missing) == 0 {
		return ret, nil
	}

	// Check against the remote list of what we can fetch
	remote, err := GetRemoteWordlists()
	if err != nil {
		return ret, fmt.Errorf("failed to get remote wordlists")
	}

	remotegot := make(map[string]WordlistMetadata)
	for _, v := range remote {
		remotegot[v.Filename] = v
		remotegot[v.Shortname] = v
	}

	nonExist := make([]string, 0)
	for _, v := range missing {
		if vv, ok := remotegot[v]; ok {
			ret = append(ret, vv)
			continue
		}
		nonExist = append(nonExist, v)
	}
	if len(nonExist) != 0 {
		return ret, fmt.Errorf("invalid names: %v", nonExist)
	}

	return ret, err
}
