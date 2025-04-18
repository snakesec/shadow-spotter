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
	"gitlab.com/snake-security/shadowspotter/internal/kitebuilder"
	"gitlab.com/snake-security/shadowspotter/pkg/context"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/spf13/cobra"
)

var (
	debug = false
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse <filename>",
	Short: "parse an kitebuilder schema and print out the prettified data",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `parse an kitebuilder schema and print out the prettified data.
this can also be configured to compile the schema into a compact/compressed
format (which is actually a captnproto serialized blob)

passing - as filename will read from stdin. This will read all of stdin
before processing and will not parse the input as streaming input

-d Debug mode will attempt to parse the schema with error handling
-v=debug Debug verbosity will print out the errors for the schema`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		if filename == "-" {
			if err := kitebuilder.ScanStdin(context.Context(), kitebuilder.Debug(debug)); err != nil {
				log.Fatal().Err(err).Msg("failed to read from stdin")
			}
		} else {
			if err := kitebuilder.ScanFile(context.Context(), filename, kitebuilder.Debug(debug)); err != nil {
				log.Fatal().Err(err).Msg("failed to read from stdin")
			}
		}
	},
}

func init() {
	kidebuilderCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	 parseCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug the parsing")
}
