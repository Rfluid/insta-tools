package following_service

import (
	"fmt"
	"sync"
	"time"

	log_service "github.com/Rfluid/insta-tools/src/log/service"
	"github.com/pterm/pterm"
)

func GetAll(userID string, cookies map[string]string, count int, maxID string, threads int, sleepTime int) ([]map[string]interface{}, error) {
	var following []map[string]interface{}
	hasMore := true
	retrievedCount := 0

	// Mutex for safely updating shared data
	var mutex sync.Mutex

	// Channel for managing API requests
	taskChan := make(chan string, threads)
	var wg sync.WaitGroup

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		"Starting to retrieve following...",
	)

	// Worker function for concurrent API calls
	worker := func() {
		for currentMaxID := range taskChan {
			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Fetching following with maxID: %s...", currentMaxID),
			)

			// Make request
			result, err := Get(userID, cookies, count, currentMaxID)
			if err != nil {
				log_service.LogConditionally(
					pterm.DefaultLogger.Error,
					fmt.Sprintf("Error fetching following: %v", err),
				)
				continue
			}

			// Extract following from response
			batchFollowing, ok := result["users"].([]interface{})
			if !ok {
				log_service.LogConditionally(
					pterm.DefaultLogger.Error,
					"Invalid response format",
				)
				continue
			}

			// Convert following to map[string]interface{}
			var batchResults []map[string]interface{}
			for _, f := range batchFollowing {
				if followerMap, ok := f.(map[string]interface{}); ok {
					batchResults = append(batchResults, followerMap)
				}
			}

			// Update maxID for pagination
			nextMaxID, _ := result["next_max_id"].(string)
			if nextMaxID == "" {
				hasMore = false
			}

			// Safely update shared data
			mutex.Lock()
			following = append(following, batchResults...)
			retrievedCount += len(batchResults)
			mutex.Unlock()

			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Retrieved %d following in this batch. Total so far: %d\n", len(batchResults), retrievedCount),
			)

			// Add a delay to avoid hitting API rate limits
			time.Sleep(time.Duration(sleepTime) * time.Second)

			// If there are more following, queue the next request
			if hasMore {
				taskChan <- nextMaxID
			}
		}
		wg.Done()
	}

	// Start worker threads
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker()
	}

	// Add initial request with provided maxID
	taskChan <- maxID

	// Close channel when all tasks are completed
	go func() {
		wg.Wait()
		close(taskChan)
	}()

	wg.Wait()

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		fmt.Sprintf("Retrieval complete. Retrieved %d following.\n", len(following)),
	)
	return following, nil
}
