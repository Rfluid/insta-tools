/*
Copyright © 2025 Rfluid
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	cookie_service "github.com/Rfluid/insta-tools/src/cookie/service"
	log_service "github.com/Rfluid/insta-tools/src/log/service"
	output_service "github.com/Rfluid/insta-tools/src/output/service"
	user_service "github.com/Rfluid/insta-tools/src/user/service"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user [username]",
	Short: "Retrieve the user of an Instagram account",
	Long: `This command fetches the user of a given Instagram username.

Example:
  insta-tools user zuck --cookies "<your_cookies>"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the username from the command arguments
		username := args[0]

		// Parse cookies
		cookies := cookie_service.ParseCookies()

		// Fetch user profile info
		log_service.LogConditionally(pterm.DefaultLogger.Info, fmt.Sprintf("Fetching user for %s", username))

		data, err := user_service.Get(username, cookies)
		if err != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Error fetching user ID: %s", err))
			os.Exit(1)
		}

		// Convert map[string]interface{} to JSON
		resultJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Failed to convert data to JSON: %s", err))
			os.Exit(1)
		}

		// Print or save output
		output_service.PrintConditionally(string(resultJSON))
		if err := output_service.WriteConditionally(string(resultJSON)); err != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Error writing output: %s", err))
			os.Exit(1)
		}

		log_service.LogConditionally(pterm.DefaultLogger.Info, "User retrieval process completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
