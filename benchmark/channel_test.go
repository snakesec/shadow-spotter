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

package benchmark

import (
	"context"
	"sync"
	"testing"
)

func BenchmarkWaitGroup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
		}()
		wg.Wait()
	}
}

func BenchmarkChannel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		done := make(chan bool)
		go func() {
			done <- true
		}()
		<-done
	}
}

func BenchmarkSelectChannel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := 0; n < b.N; n++ {
		done := make(chan bool)
		go func() {
			select { case <-ctx.Done():
			case done <- true:
			}
		}()
		<-done
	}
}

func BenchmarkDoubleSelectChannel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := 0; n < b.N; n++ {
		done := make(chan bool)
		go func() {
			select {
			case <-ctx.Done():
			case done <- true:
			}
		}()
		select {
		case <-ctx.Done():
		case <-done:
		}
	}
}
