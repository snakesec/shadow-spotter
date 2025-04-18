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

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <input> <output>",
	Short: "convert an input file format into the specified output file format",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `convert an input file format into the specified output file format

this will determine the conversion based on the extensions of the input and the output
we support the following filetypes: txt, json, kite
You can convert any of the following into the corresponding types 

-d Debug mode will attempt to convert the schema with error handling
-v=debug Debug verbosity will print out the errors for the schema`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		output := args[1]

		if err := kitebuilder.ConvertFiles(context.Context(), input, output); err != nil {
			log.Fatal().Err(err).Msg("failed to convert files")
		}
		log.Info().Msg("conversion complete")
	},
}

func init() {
	kidebuilderCmd.AddCommand(convertCmd)
	convertCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug the parsing")
}
