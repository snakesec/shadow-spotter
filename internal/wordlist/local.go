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
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gitlab.com/snake-security/shadowspotter/pkg/log"
	humanize "github.com/dustin/go-humanize"
)


func GetLocalDirPanic() (string) {
	ret, err := GetLocalDir()
	if err != nil {
		panic(err)
	}
	return ret
}

func GetLocalDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	dir := usr.HomeDir
	newpath := filepath.Join(dir, ".cache", "shadowspotter", "wordlists")
	return newpath, err
}

func CreateLocalCache() error {
	newpath, err := GetLocalDir()
	if err != nil {
		return fmt.Errorf("failed to get local dir: %w", err)
	}
	log.Debug().Str("path", newpath).Msg("creating directory")
	if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	return nil
}

func GetLocalDirListing() ([]WordlistMetadata, error) {
	dir, err := GetLocalDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get local dir: %w", err)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	ret := make([]WordlistMetadata, 0)
	for _, v := range files {
		r := WordlistMetadata{
			Shortname: nameConvert(v.Name()),
			Filename: v.Name(),
			Cached: true,
			FileSize: humanize.Bytes(uint64(v.Size())),
			Source: "local",
		}
		ret = append(ret, r)
	}

	return ret, nil
}