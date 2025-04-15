package cmd

import (
	"github.com/spf13/cobra"
)

// kidebuilderCmd represents the kitebuilder command
var kidebuilderCmd = &cobra.Command{
	Use:   "kb",
	Short: "manipulate the kitebuilder schema",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `manipuate the kitebuilder schema in various ways`,
}

func init() {
	rootCmd.AddCommand(kidebuilderCmd)
}
