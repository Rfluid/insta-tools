package following_service

import (
	"fmt"
	"sync"
	"time"

	log_service "github.com/Rfluid/insta-tools/src/log/service"
	"github.com/pterm/pterm"
)

type fetchResult struct {
	Following []map[string]interface{}
	NextMaxID string
	Err       error
}

// GetAll retrieves *all* following concurrently using a manager–worker pattern.
func GetAll(
	userID string,
	cookies map[string]string,
	count int,
	initialMaxID string,
	threads int,
	sleepTime int,
) ([]map[string]interface{}, error) {
	// ------------------------------------------------------------------------
	// Data structures
	// ------------------------------------------------------------------------
	var (
		allFollowing []map[string]interface{}
		dataMu       sync.Mutex // protects allFollowing
		globalErr    error
		errMu        sync.Mutex
	)

	// Channel of tasks, where each task is "fetch the next page for this maxID"
	taskChan := make(chan string)

	// Channel of results, each worker sends back a FetchResult
	resultsChan := make(chan fetchResult)

	var wg sync.WaitGroup // WaitGroup for workers

	// ------------------------------------------------------------------------
	// Worker pool
	// ------------------------------------------------------------------------
	worker := func() {
		defer wg.Done()

		for maxID := range taskChan {
			log_service.LogConditionally(
				pterm.DefaultLogger.Info,
				fmt.Sprintf("Fetching following for maxID: %s", maxID),
			)
			result, err := Get(userID, cookies, count, maxID)
			if err != nil {
				// Send error back
				resultsChan <- fetchResult{
					Err: fmt.Errorf("fetch error for maxID=%s: %w", maxID, err),
				}
				continue
			}

			batch, ok := result["users"].([]interface{})
			if !ok {
				// Invalid response format
				resultsChan <- fetchResult{
					Err: fmt.Errorf("invalid response format; missing 'users' array for maxID=%s", maxID),
				}
				continue
			}

			// Convert []interface{} → []map[string]interface{}
			var batchFollowing []map[string]interface{}
			for _, item := range batch {
				if fm, ok := item.(map[string]interface{}); ok {
					batchFollowing = append(batchFollowing, fm)
				}
			}

			nextMaxID, _ := result["next_max_id"].(string)

			// Optional rate limiting
			time.Sleep(time.Duration(sleepTime) * time.Second)

			// Send success result
			resultsChan <- fetchResult{
				Following: batchFollowing,
				NextMaxID: nextMaxID,
				Err:       nil,
			}
		}
	}

	// Spin up N workers
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker()
	}

	// ------------------------------------------------------------------------
	// Manager goroutine
	// ------------------------------------------------------------------------
	// The manager sends tasks (maxIDs) to workers, *and* collects results.
	// We track how many tasks are "in flight" so we know when to stop.
	// Because tasks can generate new tasks (i.e. nextMaxID), we dynamically
	// feed them back into the worker pool.
	managerWg := sync.WaitGroup{}
	managerWg.Add(1)

	go func() {
		defer managerWg.Done()

		// Start by feeding the initial maxID
		inFlight := 1
		taskChan <- initialMaxID

		// Keep reading results until inFlight == 0
		for inFlight > 0 {
			res, ok := <-resultsChan
			if !ok {
				// If resultsChan is closed unexpectedly, we break
				break
			}

			// Decrement in-flight count for the completed task
			inFlight--

			if res.Err != nil {
				// Record error (if you only want the first error, store it once)
				errMu.Lock()
				if globalErr == nil {
					globalErr = res.Err
				}
				errMu.Unlock()
				// We can keep going, or break early, up to you.
				continue
			}

			// Append the returned following
			dataMu.Lock()
			allFollowing = append(allFollowing, res.Following...)
			dataMu.Unlock()

			// If there's a nextMaxID, enqueue a new task
			if res.NextMaxID != "" {
				inFlight++
				taskChan <- res.NextMaxID
			}
		}

		// No more tasks will be generated -> close the worker channel
		close(taskChan)
	}()

	// ------------------------------------------------------------------------
	// Wait for workers to finish
	// ------------------------------------------------------------------------
	wg.Wait()

	// When all workers are done reading from taskChan, they exit
	// so no one else will write to resultsChan. Now we can close resultsChan.
	close(resultsChan)

	// The manager goroutine might still be waiting in the for-loop above,
	// but it will break out once resultsChan is closed. Wait for manager to finish.
	managerWg.Wait()

	return allFollowing, globalErr
}
