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

package kiterunner

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"gitlab.com/snake-security/shadowspotter/pkg/http"
)

type ProgressBar interface {
	Incr(n int64)
	AddTotal(n int64)
}

type NullProgressBar struct {
	total int64
	hits  int64
}

func (n *NullProgressBar) Incr(v int64) {
	atomic.AddInt64(&n.hits, v)
}

func (n *NullProgressBar) AddTotal(v int64) {
	atomic.AddInt64(&n.total, v)
}

var _ ProgressBar = &NullProgressBar{}

type Config struct {
	MaxParallelHosts     int           `toml:"max_parallel_hosts" json:"max_parallel_hosts" mapstructure:"max_parallel_hosts"`
	MaxConnPerHost       int           `toml:"max_conn_per_host" json:"max_conn_per_host" mapstructure:"max_conn_per_host"`
	WildcardDetection    bool          `json:"wildcard_detection"`
	Delay                time.Duration `toml:"delay_ms" json:"delay_ms" mapstructure:"delay_ms"`
	HTTP                 http.Config   `toml:"http" json:"http" mapstructure:"http"`
	QuarantineThreshold  int64
	PreflightCheckRoutes []*http.Route // these are the routes use to calculate the baseline. If the slice is empty, no baselines will be created so requests will match on the status codes
	ProgressBar          ProgressBar
	RequestValidators    []RequestValidator
}

func NewDefaultConfig() *Config {
	return &Config{
		MaxParallelHosts: 50,
		MaxConnPerHost:   5,
		Delay:            1 * time.Duration(0),
		// we have no default status codes, we rely on our wildcard detection
		ProgressBar:          &NullProgressBar{},
		PreflightCheckRoutes: append([]*http.Route{}, PreflightCheckRoutes...),
		RequestValidators: []RequestValidator{
			&KnownBadSitesValidator{},
			&WildcardResponseValidator{},
		},
	}
}

type ErrBadConfig struct {
	fields []string
}

func (e *ErrBadConfig) Error() string {
	return fmt.Sprintf("config has invalid values in: %v", strings.Join(e.fields, ", "))
}

func (c *Config) Validate() error {
	badFields := make([]string, 0)
	if c.MaxConnPerHost < 1 {
		badFields = append(badFields, "MaxConnPerHost")
	}
	if c.MaxParallelHosts < 1 {
		badFields = append(badFields, "MaxParallelHosts")
	}
	if len(badFields) != 0 {
		return &ErrBadConfig{fields: badFields}
	}

	if c.ProgressBar == nil {
		c.ProgressBar = &NullProgressBar{}
	}

	actualValidators := make([]RequestValidator, 0)
	for _, v := range c.RequestValidators {
		if v != nil {
			actualValidators = append(actualValidators, v)
		}
	}
	c.RequestValidators = actualValidators

	return nil
}

type ConfigOption func(*Config)

func MaxTimeout(n time.Duration) ConfigOption {
	return func(c *Config) {
		c.HTTP.Timeout = n
	}
}

func Delay(n time.Duration) ConfigOption {
	return func(c *Config) {
		c.Delay = n
	}
}

func MaxRedirects(n int) ConfigOption {
	return func(c *Config) {
		c.HTTP.MaxRedirects = n
	}
}

func MaxConnPerHost(v int) ConfigOption {
	return func(c *Config) {
		c.MaxConnPerHost = v
	}
}

func MaxParallelHosts(v int) ConfigOption {
	return func(c *Config) {
		c.MaxParallelHosts = v
	}
}

func ReadBody(v bool) ConfigOption {
	return func(c *Config) {
		c.HTTP.ReadBody = v
	}
}

func ReadHeaders(v bool) ConfigOption {
	return func(c *Config) {
		c.HTTP.ReadHeaders = v
	}
}

func BlacklistDomains(in []string) ConfigOption {
	return func(o *Config) {
		o.HTTP.BlacklistRedirects = append(o.HTTP.BlacklistRedirects, in...)
	}
}

func WildcardDetection(enabled bool) ConfigOption {
	return func(o *Config) {
		o.WildcardDetection = enabled
	}
}

func AddRequestFilter(f RequestValidator) ConfigOption {
	return func(o *Config) {
		if f != nil {
			o.RequestValidators = append(o.RequestValidators, f)
		}
	}
}

// SkipPreflight will zero out the preflight check routes
func SkipPreflight(enabled bool) ConfigOption {
	return func(o *Config) {
		if enabled {
			o.PreflightCheckRoutes = o.PreflightCheckRoutes[:0]
		}
	}
}

func AddProgressBar(p ProgressBar) ConfigOption {
	return func(o *Config) {
		o.ProgressBar = p
	}
}

func TargetQuarantineThreshold(n int64) ConfigOption {
	return func(o *Config) {
		o.QuarantineThreshold = n
	}
}

func SetPreflightCheckRoutes(r []*http.Route) ConfigOption {
	return func(o *Config) {
		o.PreflightCheckRoutes = append(o.PreflightCheckRoutes[:0], r...)
	}
}


func HTTPExtraHeaders(h []http.Header) ConfigOption {
	return func(o *Config) {
		o.HTTP.ExtraHeaders = append(o.HTTP.ExtraHeaders, h...)
	}
}