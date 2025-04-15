package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These global variables are injected at build time to provide the
// version command
var (
	Version = "v0.0.1"
	Commit  = "commit"
	Date    = "today"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of the binary you're running",
	Long: "\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n" + `this shows you the version of the binary that is running`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s - %s\n", Version, Commit)
		fmt.Printf("Built on %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
