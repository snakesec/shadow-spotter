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
Package context provides utilities wrapping the native go/context package
for catching and handling multiple interrupts.

The main use-case is to to attach an interrupt signal handler to the context.
This can be used from your CLI applications to ensure a graceful shutdown of the scanning
and to clean up any resources mid-flight

	import "gitlab.com/snake-security/shadowspotter/pkg/context"

	...

	if err := scan.ScanDomainOrFile(context.Context(), domain, opts...); err != nil {
		log.Fatal().Err(err).Msg("failed to scan domain")
	}
 */
package context
