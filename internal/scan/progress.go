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

package scan

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/vbauerster/mpb/v6"
)

type ProgressBar struct {
	Pb       *mpb.Progress
	Requests *progressbar.ProgressBar
}

func NewProgress(max int64) *ProgressBar {
	pb := mpb.New(
		mpb.WithOutput(os.Stderr),
	)
	requestb := progressbar.NewOptions64(max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(5),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetVisibility(true),
		progressbar.OptionSpinnerType(14),
		// progressbar.OptionFullWidth(),
	)
	return &ProgressBar{
		Pb:       pb,
		Requests: requestb,
	}
}

func (b *ProgressBar) Incr(n int64) {
	b.Requests.Add64(n)
}

func (b *ProgressBar) AddTotal(n int64) {
	b.Requests.ChangeMax64(b.Requests.GetMax64() + n)
}

