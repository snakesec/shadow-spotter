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
The errors package provides a custom error type and utilities used when performing recursive
analysis of kitebuilder and proute apis.

The error type will provide context about what depth and what specific component
of the API yielded the error, and the printing utilities assist in graphically
representing the errors with corresponding nested depth.

Usage

	import errors2 "gitlab.com/snake-security/shadowspotter/pkg/errors"

	...

	if err := inAPI.EncodeStringSlice(output); err != nil {
		var merr *multierror.Error
		if errors.As(err, &merr) {
			for _, v := range merr.Errors {
				errors2.PrintError(v, 0)
			}
		} else {
			return fmt.Errorf("converting to txt output error: %w", err)
		}
	}

 */
package errors
