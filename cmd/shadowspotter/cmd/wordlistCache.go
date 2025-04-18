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

package cmd

import (
	"gitlab.com/snake-security/shadowspotter/internal/wordlist"
	"gitlab.com/snake-security/shadowspotter/pkg/context"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/spf13/cobra"
)

var (
	noCache bool = false
)

// wordlistSaveCmd represents the wordlistCache command
var wordlistSaveCmd = &cobra.Command{
	Use:   "save [wordlists ...]",
	Short: "save the wordlists specified (full filename or alias)",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `save will download the wordlists specified to ~/.cache/shadowspotter/wordlists

you can use the alias or the full filename listed in [Shadow-Spotter wordlist list]
`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := wordlist.Save(context.Context(), noCache, args...); err != nil {
			log.Fatal().Err(err).Msg("failed to list wordlists")
		}
	},
}

func init() {
	wordlistCmd.AddCommand(wordlistSaveCmd)
	wordlistSaveCmd.Flags().BoolVar(&noCache, "no-cache", noCache, "delete the local files matching the names and pull fresh files")
}
