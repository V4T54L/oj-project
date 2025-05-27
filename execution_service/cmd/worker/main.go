package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("👷 Worker service starting...")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Graceful shutdown on SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Use env var in real setup
	})
	defer rdb.Close()

	// Start worker
	startWorker(ctx, rdb, &wg)

	// Wait for signal
	<-sigs
	// log.Println("🔻 Shutdown signal received.")
	cancel()

	// Wait for worker to finish
	wg.Wait()
	// log.Println("✅ Worker exited cleanly.")
}

func startWorker(ctx context.Context, rdb *redis.Client, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// log.Println("🛠️ Worker started...")

		for {
			select {
			case <-ctx.Done():
				// log.Println("🛑 Worker context canceled. Exiting...")
				return
			default:
				res, err := rdb.BLPop(ctx, 5*time.Second, "python").Result()
				if err != nil {
					if err == redis.Nil {
						continue
					}
					if ctx.Err() != nil {
						log.Println("Context canceled during BLPOP")
						return
					}
					log.Printf("BLPOP error: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}

				var task ExecuteCodePayload
				if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
					log.Printf("Invalid task JSON: %v", err)
					continue
				}

				// log.Printf("🔧 Processing task ID %d: %s", task.ID, task.Name)

				// Simulate work
				time.Sleep(2 * time.Second)

				testResults := make([]TestCaseResult, 0, len(task.TestCases))
				for _, v := range task.TestCases {
					testResults = append(testResults, TestCaseResult{
						ID:             v.ID,
						Input:          v.Input,
						Output:         v.ExpectedOutput,
						ExpectedOutput: v.ExpectedOutput,
						RuntimeMS:      30,
						MemoryKB:       4,
						Status:         "Accepted",
					})
				}

				result := ExecuteCodeResponse{
					ID:              task.ID,
					Status:          "Accepted",
					TestCaseResults: testResults,
					ExecutionType:   "Submission",
				}

				data, _ := json.Marshal(result)
				if err := rdb.RPush(ctx, "results_queue", data).Err(); err != nil {
					log.Printf("❌ Failed to push result: %v", err)
				} else {
					log.Printf("✅ Pushed result for task ID %d", task.ID)
				}
			}
		}
	}()
}

type ProblemTestCase struct {
	ID             int    `json:"id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type TestCaseResult struct {
	ID             int    `json:"id"`
	Input          string `json:"input"`
	Output         string `json:"output"`
	ExpectedOutput string `json:"expected_output"`
	RuntimeMS      int    `json:"runtime_ms"`
	MemoryKB       int    `json:"memory_kb"`
	Status         string `json:"status"` // TLE, MLE, Acccepted, Wrong Answer, $Error.message
}

type ExecuteCodePayload struct {
	ID             int               `json:"id"`
	LanguageID     int               `json:"language_id"`
	Code           string            `json:"code"`
	TestCases      []ProblemTestCase `json:"test_cases"`
	RuntimeLimitMS int               `json:"runtime_limit_ms"`
	MemoryLimitKB  int               `json:"memory_limit_kb"`
	ExecutionType  string            `json:"execution_type"` // Run, Submit, Validation
}

type ExecuteCodeResponse struct {
	ID              int              `json:"id"`
	Status          string           `json:"status"` // TLE, MLE, Acccepted, Wrong Answer, $Error.message
	TestCaseResults []TestCaseResult `json:"test_case_results"`
	ExecutionType   string           `json:"execution_type"` // Run, Submit, Validation
}
