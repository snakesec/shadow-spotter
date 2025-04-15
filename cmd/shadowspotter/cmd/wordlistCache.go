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
