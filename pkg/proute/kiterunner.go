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
	"strings"

	"gitlab.com/snake-security/shadowspotter/pkg/http"
	"github.com/hashicorp/go-multierror"
)

func APIsToKiterunnerRoutes(api []API, allowunsafe bool, forcemethod string, userquery string) ([]*http.Route, error) {
	var merr *multierror.Error
	ret := make([]*http.Route, 0)
	for _, v := range api {
		tmp, err := ToKiterunnerRoutes(v, allowunsafe, forcemethod, userquery)
		if err != nil {
			multierror.Append(merr, err)
		}
		ret = append(ret, tmp...)
		
	}
	return ret, merr.ErrorOrNil()
}

func ToKiterunnerRoutes(api API, allowunsafe bool, forcemethod string, userquery string) ([]*http.Route, error) {
	var merr *multierror.Error

	ret := make([]*http.Route, 0)
	for _, v := range api.Routes {
		// Skip these options since we don't actually care about the content here
		r, err := v.ToKiterunner(userquery, api.Headers(true)...)
		if err != nil {
			multierror.Append(merr, err)
		}
		switch string(r.Method) {
		case "HEAD", "OPTIONS", "CONNECT", "TRACE":
			// we're biased. skip these since they're noisy
			continue
		case "GET", "POST", "PUT", "DELETE", "PATCH":
			// these guys are alright
		}

		if !allowunsafe {

			switch string(r.Method) {
			case "PUT", "DELETE", "PATCH":
				continue
			}

		}

		if forcemethod != "" {

			if forcemethod != string(r.Method) {
				continue
			}

		}

		r.Source = api.ID
		ret = append(ret, r)
	}
	return ret, merr.ErrorOrNil()
}

func (r Route) ToKiterunner(customQuery string, extraHeaders ...KV) (*http.Route, error) {
	var err error

	method := strings.TrimSpace(strings.ToUpper(r.Method))
	switch method {
	case "HEAD", "OPTIONS", "CONNECT", "TRACE":
		// we're biased. skip these since they're noisy
	case "GET", "POST", "PUT", "DELETE", "PATCH":
		// these guys are alright
	default:
		method = "GET"
	}

	ret := &http.Route{
		Path:   []byte(r.path),
		Method: http.Method(method),
	}
	for _, h := range r.Headers(true) {
		ret.Headers = append(ret.Headers, http.Header{h.Key, h.Value})
	}

	for _, h := range extraHeaders {
		ret.Headers = append(ret.Headers, http.Header{h.Key, h.Value})
	}

	ct := ContentTypeFormEncoded
	if len(r.ContentType) > 0 {
		ct = r.ContentType[0]
		// Overwrite the form-data format with a proper boundary header
		if strings.Contains(string(ct), "form-data") {
			ct = "multipart/form-data; boundary=" + DefaultFormDataBoundary
		}
	}
	ret.Body = r.Body(true, ct)

	// only add content type if there is a body specified
	// or if its a non-get type
	if len(ret.Body) > 0 || (string(ret.Method) != string(http.GET)) {
		ret.Headers = append(ret.Headers, http.Header{"Content-Type", string(ct)})
	}

	tmp, err := r.Path(false)
	if err != nil {
		return ret, err
	}
	ret.Path = []byte( tmp )

	tmp, err = r.Query(false)

	if len(customQuery) > 0 {

		if strings.Contains(customQuery, "?") {
			customQuery = strings.Replace(customQuery, "?", "", 1)
		}

		queryStr := ""

		if len(tmp) > 0 {
			queryStr += tmp + "&" + customQuery
		} else {
			queryStr += customQuery
		}

		//r.Path = append(r.Path, []byte(pathStr)...)

		//ret.Query = []byte(tmp)
		ret.Query = []byte(queryStr)

	} else {

		ret.Query = []byte(tmp)

	}

	return ret, err
}
