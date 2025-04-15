package cmd

import (
	"gitlab.com/snake-security/shadowspotter/internal/wordlist"
	"gitlab.com/snake-security/shadowspotter/pkg/context"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/spf13/cobra"
)

var (
)

// wordlistListCmd represents the wordlistList command
var wordlistListCmd = &cobra.Command{
	Use:   "list",
	Short: "list the wordlists cached and available",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `list will show all the remote wordlists and all the local wordlists
with the corresponding abbreviations for use
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt, err := wordlist.FormatFromString(Output)
		if err != nil {
			log.Fatal().Err(err).Msg("invalid format")
		}
		if err := wordlist.List(context.Context(), wordlist.OutputFormat(fmt)); err != nil {
			log.Fatal().Err(err).Msg("failed to list wordlists")
		}
	},
}

func init() {
	wordlistCmd.AddCommand(wordlistListCmd)
}
