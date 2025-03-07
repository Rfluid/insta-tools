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
	following_flag "github.com/Rfluid/insta-tools/src/following/flag"
	following_service "github.com/Rfluid/insta-tools/src/following/service"
	log_service "github.com/Rfluid/insta-tools/src/log/service"
	output_service "github.com/Rfluid/insta-tools/src/output/service"
	thread_flag "github.com/Rfluid/insta-tools/src/thread/flag"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// followingCmd represents the following command
var followingCmd = &cobra.Command{
	Use:   "following [userID] [count] [maxID]",
	Short: "Retrieve a list of users the target user is following",
	Long: `This command fetches the list of users that the given user is following.

Arguments:
1. A valid userID.
2. A batch count (number of followings per request).
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

		// Check if retrieving all following
		if following_flag.RetrieveAll {
			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Fetching ALL following for userID: %s with count: %d and initial maxID: %s", userID, count, maxID),
			)

			// Fetch all following using pagination
			following, reqErr := following_service.GetAll(userID, cookies, count, maxID, thread_flag.APIThreads, following_flag.SleepTime)
			if reqErr != nil {
				pterm.DefaultLogger.Error(fmt.Sprintf("Error fetching all following: %s. Only partial results available", reqErr))
			}

			// Convert to JSON
			resultJSON, err := json.MarshalIndent(following, "", "  ")
			if err != nil {
				pterm.DefaultLogger.Error(fmt.Sprintf("Failed to convert data to JSON: %s", err))
				os.Exit(1)
			}

			// Print or save output
			output_service.PrintConditionally(string(resultJSON))
			output_service.PrintConditionally(string(resultJSON))
			if err := output_service.WriteConditionally(string(resultJSON)); err != nil {
				pterm.DefaultLogger.Error(fmt.Sprintf("Error writing output: %s", err))
				os.Exit(1)
			}
			if reqErr != nil {
				os.Exit(1)
			}

			return
		}

		// Fetch a single batch of following
		log_service.LogConditionally(
			pterm.DefaultLogger.Info,
			fmt.Sprintf("Fetching following for userID: %s with count: %d and maxID: %s", userID, count, maxID),
		)

		data, reqErr := following_service.Get(userID, cookies, count, maxID)
		if reqErr != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Error fetching following: %s", reqErr))
		}

		// Convert map[string]interface{} to JSON
		resultJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Failed to convert data to JSON: %s", err))
			os.Exit(1)
		}

		output_service.PrintConditionally(string(resultJSON))
		output_service.PrintConditionally(string(resultJSON))
		if err := output_service.WriteConditionally(string(resultJSON)); err != nil {
			pterm.DefaultLogger.Error(fmt.Sprintf("Error writing output: %s", err))
			os.Exit(1)
		}
		if reqErr != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(followingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// followingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// followingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	followingCmd.Flags().BoolVarP(&following_flag.RetrieveAll, "all", "a", false, "Retrieve all followings using pagination")
	followingCmd.Flags().IntVar(&following_flag.SleepTime, "sleep", 1, "Seconds to wait between API requests when using --all")
}
