package cmd

import (
	"github.com/spf13/cobra"
)

// kidebuilderCmd represents the kitebuilder command
var wordlistCmd = &cobra.Command{
	Use:   "wordlist",
	Short: "look at your cached wordlists and remote wordlists",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `used to help manage your wordlists in your .cache/kiterunner`,
}

func init() {
	rootCmd.AddCommand(wordlistCmd)
}
