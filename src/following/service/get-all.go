package following_service

import (
	"fmt"
	"sync"
	"time"

	log_service "github.com/Rfluid/insta-tools/src/log/service"
	"github.com/pterm/pterm"
)

// GetAll retrieves all following concurrently with pagination support.
func GetAll(
	userID string,
	cookies map[string]string,
	count int,
	maxID string,
	threads int,
	sleepTime int,
) ([]map[string]interface{}, error) {
	var (
		following      []map[string]interface{}
		retrievedCount int

		// Global mutex for shared data
		dataMu sync.Mutex

		// WaitGroup to know when all workers have finished
		wg sync.WaitGroup

		// Shared error tracking
		globalErr error
		errMu     sync.Mutex

		// A flag to indicate whether there are more pages to fetch
		hasMore      = true
		hasMoreMutex sync.Mutex
	)

	// Channel to queue up "maxID" tasks.
	taskChan := make(chan string, threads)
	// Channel to collect errors from workers.
	errCh := make(chan error, threads)

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		"Starting to retrieve following...",
	)

	// -------------------------------------------------------
	// Worker function
	// -------------------------------------------------------
	worker := func() {
		defer wg.Done()

		for currentMaxID := range taskChan {
			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Fetching following with maxID: %s...", currentMaxID),
			)

			// 1. Make the request
			result, err := Get(userID, cookies, count, currentMaxID)
			if err != nil {
				errCh <- fmt.Errorf("error fetching following with maxID %s: %w", currentMaxID, err)
				continue
			}

			// 2. Extract the "users" field
			batchFollowing, ok := result["users"].([]interface{})
			if !ok {
				errCh <- fmt.Errorf("invalid response format; missing 'users' array for maxID %s", currentMaxID)
				continue
			}

			// Convert the raw slice of interfaces into []map[string]interface{}
			var batchResults []map[string]interface{}
			for _, item := range batchFollowing {
				if fm, ok := item.(map[string]interface{}); ok {
					batchResults = append(batchResults, fm)
				}
			}

			// 3. Safely update global slice & count
			dataMu.Lock()
			following = append(following, batchResults...)
			retrievedCount += len(batchResults)
			dataMu.Unlock()

			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Retrieved %d following in this batch. Total so far: %d", len(batchResults), retrievedCount),
			)

			// 4. Find the next maxID for pagination
			nextMaxID, _ := result["next_max_id"].(string)

			// 5. Respect the API rate limits by sleeping
			time.Sleep(time.Duration(sleepTime) * time.Second)

			// 6. Decide if we have more pages to fetch
			if nextMaxID == "" {
				// If there's no next page, signal no more
				hasMoreMutex.Lock()
				hasMore = false
				hasMoreMutex.Unlock()
			} else {
				// If we still have more pages, queue nextMaxID
				hasMoreMutex.Lock()
				if hasMore {
					taskChan <- nextMaxID
				}
				hasMoreMutex.Unlock()
			}
		}
	}

	// -------------------------------------------------------
	// Spin up N worker goroutines
	// -------------------------------------------------------
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker()
	}

	// -------------------------------------------------------
	// Error collector goroutine
	// -------------------------------------------------------
	var errWg sync.WaitGroup
	errWg.Add(1)
	go func() {
		defer errWg.Done()
		for e := range errCh {
			if e != nil {
				errMu.Lock()
				// Record the first non-nil error, or extend to track them all if desired.
				if globalErr == nil {
					globalErr = e
				}
				errMu.Unlock()
			}
		}
	}()

	// -------------------------------------------------------
	// Feed the initial maxID into the pipeline
	// -------------------------------------------------------
	taskChan <- maxID

	// Wait for all workers to finish
	wg.Wait()

	// Close channels now that all workers are done
	close(taskChan)
	close(errCh)

	errWg.Wait()

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		fmt.Sprintf("Retrieval complete. Retrieved %d following.", len(following)),
	)

	return following, globalErr
}
