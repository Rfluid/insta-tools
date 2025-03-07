/*
Copyright © 2025 Rfluid
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	cookie_service "github.com/Rfluid/insta-tools/src/cookie/service"
	followers_flag "github.com/Rfluid/insta-tools/src/followers/flag"
	followers_service "github.com/Rfluid/insta-tools/src/followers/service"
	log_service "github.com/Rfluid/insta-tools/src/log/service"
	output_service "github.com/Rfluid/insta-tools/src/output/service"
	thread_flag "github.com/Rfluid/insta-tools/src/thread/flag"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// followersCmd represents the followers command
var followersCmd = &cobra.Command{
	Use:   "followers [userID] [count] [maxID]",
	Short: "Retrieve a list of Instagram followers",
	Long: `This command fetches followers from Instagram using the API.

It requires:
1. A valid userID.
2. A batch count (number of followers per request).
3. An optional maxID to paginate requests.`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse arguments
		userID := args[0]
		count, err := strconv.Atoi(args[1])
		if err != nil {
			pterm.DefaultLogger.Error("Invalid count argument. Must be an integer.")
			os.Exit(1)
		}
		maxID := ""
		if len(args) == 3 {
			maxID = args[2]
		}

		// Parse cookies
		cookies := cookie_service.ParseCookies()

		// Check if retrieving all followers
		if followers_flag.RetrieveAll {
			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Fetching ALL followers for userID: %s with count: %d and initial maxID: %s", userID, count, maxID),
			)

			// Fetch all followers using pagination
			followers, reqErr := followers_service.GetAll(userID, cookies, count, maxID, thread_flag.APIThreads, followers_flag.SleepTime)
			if reqErr != nil {
				pterm.DefaultLogger.Error(fmt.Sprintf("Error fetching all followers: %s\nOnly partial results available", reqErr))
			}

			// Convert to JSON
			resultJSON, err := json.MarshalIndent(followers, "", "  ")
			if err != nil {
				pterm.DefaultLogger.Error(fmt.Sprintf("Failed to convert data to JSON: %s", err))
				os.Exit(1)
			}

			// Print or save output
			output_service.PrintConditionally(string(resultJSON))
			output_service.WriteConditionally(string(resultJSON))
			if reqErr != nil {
				os.Exit(1)
			}

			return
		}

		// Fetch a single batch of followers
		log_service.LogConditionally(
			pterm.DefaultLogger.Info,
			fmt.Sprintf("Fetching followers for userID: %s with count: %d and maxID: %s", userID, count, maxID),
		)

		data, reqErr := followers_service.Get(userID, cookies, count, maxID)
		if reqErr != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Error fetching followers: %s", err))
		}

		// Convert map[string]interface{} to JSON
		resultJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Failed to convert data to JSON: %s", err))
			os.Exit(1)
		}

		output_service.PrintConditionally(string(resultJSON))
		output_service.WriteConditionally(string(resultJSON))
		if reqErr != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(followersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// followersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// followersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	followersCmd.Flags().BoolVarP(&followers_flag.RetrieveAll, "all", "a", false, "Retrieve all followers using pagination")
	followersCmd.Flags().IntVar(&followers_flag.SleepTime, "sleep", 0, "Seconds to wait between API requests when using --all")
}
