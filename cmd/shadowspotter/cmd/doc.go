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

/*
Package cmd provides all the commands for the kiterunner binary.

The commands are separated by file, with the prefix for each command corresponding to the parent command, i.e.
kitebuilderCompile.go corresponds to kb compile

there are a few global CLI flags that can be used to configure how kiterunner will operate. These are defined
by the globally exposed variables

The server can be started up across multiple ports to provide the simulation of scanning
multiple hosts

Usage

	go run ./cmd/testServer -p 14000-14500

 */
package cmd
