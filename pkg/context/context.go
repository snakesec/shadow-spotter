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

package context

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gitlab.com/snake-security/shadowspotter/pkg/log"
)

var (
	ctx            context.Context
	cancel         context.CancelFunc
	ctxInitialized sync.Once
)

// AddInterruptCancellation will add an interrupt handler that will catch the first SIGTERM and cancel the context
// upon second SIGTERM, the program will exit immediately
// This wrapping allows for graceful shutdown of the application
func AddInterruptCancellation(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		interrupts := 0
		for {
			select {
			case <-c:
				interrupts++
				if interrupts > 1 {
					log.Info().Msg("Received multiple interrupt signals. Exiting")
					os.Exit(1)
				}
				log.Info().Msg("Received interrupt signal")
				cancel()
			case <-ctx.Done():
			}
		}
	}()
}

// InitContext will initialize the global context used to catch interrupts. This is automatically called
// by Context and Cancel
func InitContext() {
	ctxInitialized.Do(func() {
		ctx, cancel = context.WithCancel(context.Background())
		AddInterruptCancellation(ctx, cancel)
	})
}

// Context will initialize the global context and attach the interrupt handler that will cancel the context
// upon SIGTERM. This is safe to call from multiple goroutines and will always return the same context
func Context() context.Context {
	InitContext()
	return ctx
}

// Cancel will cancel the global context. Calling this multiple times is the equivalent of cancelling
// the same context multiple times
func Cancel() {
	InitContext()
	cancel()
}

