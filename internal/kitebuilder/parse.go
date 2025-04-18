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
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	errors2 "gitlab.com/snake-security/shadowspotter/pkg/errors"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"gitlab.com/snake-security/shadowspotter/pkg/kitebuilder"
	"gitlab.com/snake-security/shadowspotter/pkg/proute"
	"github.com/hashicorp/go-multierror"
)

type ScanOptions struct {
	Debug bool
}

func NewDefaultScanOptions() *ScanOptions {
	return &ScanOptions{}
}

func Debug(enabled bool) ScanOption {
	return func(o *ScanOptions) {
		o.Debug = enabled
	}
}

type ScanOption func(o *ScanOptions)

func ScanStdin(ctx context.Context, opts ...ScanOption) error {
	return DebugPrintReader(ctx, os.Stdin, opts...)
}

func ScanFile(ctx context.Context, filename string, opts ...ScanOption) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer jsonFile.Close()
	return DebugPrintReader(ctx, jsonFile, opts...)
}

func fixOutputFilename(filename string) string {
	if !strings.HasSuffix(filename, ".kite") {
		log.Info().Str("filename", filename).Msg(".kite extension added to filename")
		return fmt.Sprintf("%s.kite", filename)
	}
	return filename
}

func CompileFile(ctx context.Context, input string, outputFile string, opts ...ScanOption) error {
	jsonFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer jsonFile.Close()
	return Compile(ctx, jsonFile, outputFile, opts...)
}

func Compile(ctx context.Context, r io.Reader, outputFile string, opts ...ScanOption) error {
	// add the .kite extension if it doesnt exist
	outputFile = fixOutputFilename(outputFile)

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var merr *multierror.Error
	api, err := kitebuilder.SlowLoadJSONBytes(data)
	if errors.As(err, &merr) {
		for _, v := range merr.Errors {
			errors2.PrintError(v, 0)
		}
	} else if err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	apis, err := proute.FromKitebuilderAPIs(api)
	if errors.As(err, &merr) {
		log.Error().Msg("errors while parsing apis")
		for _, v := range merr.Errors {
			errors2.PrintError(v, 0)
		}
	} else if err != nil {
		return fmt.Errorf("failed to parse api: %w", err)
	}

	if err := proute.APIS(apis).EncodeProtoFile(outputFile); err != nil {
		return fmt.Errorf("failed to encode apis: %w", err)
	}

	return nil
}

func DebugPrintReader(ctx context.Context, r io.Reader, opts ...ScanOption) error {
	options := NewDefaultScanOptions()
	for _, o := range opts {
		o(options)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if options.Debug {
		return DebugPrintBytes(data)
	} else {
		return PrintBytes(data)
	}
}

func DebugPrintBytes(data []byte) error {
	log.Debug().Msg("debug printing")

	var merr *multierror.Error
	api, err := kitebuilder.SlowLoadJSONBytes(data)
	if errors.As(err, &merr) {
		for _, v := range merr.Errors {
			errors2.PrintError(v, 0)
		}
	} else if err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	// kitebuilder.PrintAPIs(api)
	// for _, v := range api {
	//	pr := proute.FromKitebuilderAPI(v)
	//	pr.DebugPrint()
	// }
	for _, v := range api {
		tmp, err := proute.FromKitebuilderAPI(v)
		if errors.As(err, &merr) {
			log.Error().Str("id", v.ID).Msg("errors while parsing api")
			for _, v := range merr.Errors {
				errors2.PrintError(v, 1)
			}
		} else if err != nil {
			return fmt.Errorf("failed to parse api: %w", err)
		}

		wcr, err := proute.ToKiterunnerRoutes(tmp, true, "", "")
		if errors.As(err, &merr) {
			log.Error().Str("id", v.ID).Msg("errors while building routes")
			for _, v := range merr.Errors {
				errors2.PrintError(v, 1)
			}
		} else if err != nil {
			return fmt.Errorf("failed to parse api: %w", err)
		}
		for _, v := range wcr {
			_ = v
			// log.Log().Object("route", v)
		}
	}
	return nil
}

func PrintBytes(data []byte) error {
	api, err := kitebuilder.LoadJSONBytes(data)
	if err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	kitebuilder.PrintAPIs(api)
	return nil
}
