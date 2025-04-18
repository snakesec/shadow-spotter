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
Package http provides a 0 allocation wrapper around the fasthttp library.

The wrapper facilitates our redirect handling, our custom "target" and request building format,
and handling the various options we can attach to building a request.

Most structs in this package have a corresponding Acquire* and Release* function for using sync.Pool 'd objects.
This is to minimise allocations in your request hotloop. We recommend using Acquire* and Release* wherever possible
to ensure that unecessary allocations are avoided.

The Target type provides our wrapper around the full context needed to perform a http request. There are a few quirks
when using targets that one has to be weary of. These quirks are side-effects of the zero-locking and and minimal
allocation behaviour. Admittedly, this is very developer unfriendly, and changes to the API that make it more usable
while maintaining the performance are welcome.

 - We expect Targets to be instantiated then updated as required
 - Once Target.ParseHostHeader() has been called further modifications will not be respected
 - Target.HTTPClient() will cache the first set of options that are used. Future modifications are ignored
 - All operations on a target are thread safe to use



 */
package http