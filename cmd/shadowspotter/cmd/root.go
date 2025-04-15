package cmd

import (
	"fmt"
	"os"

	"gitlab.com/snake-security/shadowspotter/internal/art"
	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)


// These global variables can be configured with the corresponding lowercase flag
var (
	Verbose string // Verbose defines the logging level, either trace, debug, info, error, fatal
	Output  string // Output defines the output format, either pretty, text, json
	Quiet   bool // Quiet will hide the beautiful ascii art upon startup

	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Shadow-Spotter",
	Short: "Shadow-Spotter scan one or mulitple hosts",
	Long: "Shadow-Spotter is a new generation context-based web scanner designed for content discovery.\n\n\033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m",
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shadow-spotter.yaml)")

	rootCmd.PersistentFlags().StringVarP(&Verbose, "verbose", "v", "info", "level of logging verbosity. can be error,info,debug,trace")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "pretty", "output format. can be json,text,pretty")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "quiet mode. will mute unnecessary pretty text")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

func initLogging() {
	log.SetFormat(viper.GetString("output"))

	level := viper.GetString("verbose")
	if level != "" {
		if err := log.SetLevelString(level); err != nil {
			log.Fatal().Err(err).Msg("failed to initialize logging")
		}
	}
	log.Debug().Str("level", level).Str("format", viper.GetString("output")).Msg("custom log settings")

	if Output == "pretty" && !viper.GetBool("quiet") {
		art.WriteArtBytes(os.Stderr)
		fmt.Println("\t\t  \033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n\n")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kiterunner" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".shadow-spotter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
