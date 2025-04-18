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
	"context"
	"fmt"

	"gitlab.com/snake-security/shadowspotter/pkg/http"
)

// Run will perform the same operation as RunAsync. This wraps the channels with allocated structs and returns the results
// This is safe to call from concurrent threads and will use separate worker pools for each call.
// your callback will be invoked on each result received so you can asynchronously process the results if you wish
// All the results will still be returned by the []*Result slice. Modifying the result in the callback is considered
// undefined behaviour
func (e *Engine) Run(ctx context.Context, input []*http.Target) ([]*Result, error) {
	return e.RunCallback(ctx, input)
}

// RunCallback will run the scan against the provided input, calling the provided callbacks on each result.
// the callbacks can be used to log the error in realtime, or perform other processes. You should not
// modify or use the Target, or route from the Result as this may have unintended side effects
func (e *Engine) RunCallback(ctx context.Context, input []*http.Target, cb ...func(r *Result, c *Config)) ([]*Result, error) {
	tx, rx, err := e.RunAsync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to run scan: %w", err)
	}

	res := make([]*Result, 0)

	// send our input to the tx channel in a separate thread
	go func() {
		for _, v := range input {
			select {
			case <-ctx.Done():
				// log.Trace().Err(ctx.Err()).Str("goroutine", "tx worker").Msg("context cancellation received")
				return
			case tx <- v:
				// log.Trace().Str("goroutine", "tx worker").Str("target", v.String()).Msg("message sent")
			}
		}
		// log.Trace().Str("goroutine", "tx worker").Msg("closing tx")
		close(tx)
	}()

	// collect results in main thread
	for v := range rx {
		res = append(res, v)
		for _, c := range cb {
			c(v, e.config)
		}
	}

	return res, nil
}

// RunCallbackNoResult will run the scan against the provided input, calling the provided callbacks on each result.
// the callbacks can be used to log the error in realtime, or perform other processes. You should not
// modify or use the Target, or route from the Result as this may have unintended side effects
// This function does not return the results as they are released immediately after all callbacks are called.
// It is unsafe to use the result after your callbacks return
// Use this when you don't require using the result after the callback, e.g. writing to disk/printing to output
func (e *Engine) RunCallbackNoResult(ctx context.Context, input []*http.Target, cb ...func(r *Result, c *Config)) (error) {
	tx, rx, err := e.RunAsync(ctx)
	if err != nil {
		return fmt.Errorf("failed to run scan: %w", err)
	}

	// send our input to the tx channel in a separate thread
	go func() {
		for _, v := range input {
			select {
			case <-ctx.Done():
				// log.Trace().Err(ctx.Err()).Str("goroutine", "tx worker").Msg("context cancellation received")
				return
			case tx <- v:
				// log.Trace().Str("goroutine", "tx worker").Str("target", v.String()).Msg("message sent")
			}
		}
		// log.Trace().Str("goroutine", "tx worker").Msg("closing tx")
		close(tx)
	}()

	// collect results in main thread
	for v := range rx {
		for _, c := range cb {
			c(v, e.config)
		}
		v.Release()
	}

	return nil
}
