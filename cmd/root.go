/*
Copyright © 2025 Rfluid
*/
package cmd

import (
	"os"

	cookie_flag "github.com/Rfluid/insta-tools/src/cookie/flag"
	log_flag "github.com/Rfluid/insta-tools/src/log/flag"
	output_flag "github.com/Rfluid/insta-tools/src/output/flag"
	thread_flag "github.com/Rfluid/insta-tools/src/thread/flag"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "insta-tools",
	Short: "A CLI tool to interact with Instagram's API",
	Long: `insta-tools is a command-line application designed to interact with the Instagram API, 
allowing users to perform various actions such as fetching followers and getting users.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.insta-tools.yaml)")
	// Implement latter: rootCmd.PersistentFlags().BoolVar(&uiMode, "ui", false, "Enable UI mode for enhanced user input")
	rootCmd.PersistentFlags().BoolVar(&log_flag.Logs, "logs", false, "Enable logs for better user experience")
	rootCmd.PersistentFlags().StringVar(&cookie_flag.Cookies, "cookies", "", "Set Instagram session cookies")
	rootCmd.PersistentFlags().StringVarP(&output_flag.OutputPath, "output", "o", "", "Set the output file path where results will be written")
	rootCmd.PersistentFlags().IntVar(&thread_flag.APIThreads, "threads", 4, "Number of threads to use in concurrent API calls")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
