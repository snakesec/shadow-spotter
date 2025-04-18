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
Package testServer provides fasthttp testing server that's configured to respond
at very high RPS (capable of supporting up to 100k RPS).

This server is used when testing and benchmarking kiterunner to ensure that performance
does not degrade between feature enhancements and that expected behaviour does not change

The server is used for testing, and should not be used in a production environment.
We provide no support on how to use the testServer.
 */
package main
