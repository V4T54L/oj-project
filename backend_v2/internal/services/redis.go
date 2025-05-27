package services

import (
	"context"
	"encoding/json"
	"online-judge/internal/models"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(addr string) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisService{client: rdb}
}

func (r *RedisService) Close() {
	r.client.Close()
}

func (r *RedisService) StartResultWorker(
	ctx context.Context, handleResultFunc func(*models.ExecuteCodeResponse), wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// log.Println("🛠️ Result worker started...")
		for {
			select {
			case <-ctx.Done():
				// log.Println("🛑 Result worker shutting down...")
				return
			default:
				res, err := r.client.BLPop(ctx, 5*time.Second, "results_queue").Result()
				if err != nil {
					if err == redis.Nil {
						continue // no result yet
					}
					if ctx.Err() != nil {
						// log.Println("Context canceled, exiting worker...")
						return
					}
					// log.Printf("Redis BLPOP error: %v", err)
					time.Sleep(1 * time.Second)
					continue
				}

				var result models.ExecuteCodeResponse
				if err := json.Unmarshal([]byte(res[1]), &result); err != nil {
					// log.Printf("Invalid result JSON: %v", err)
					continue
				}

				handleResultFunc(&result)
			}
		}
	}()
}

func (r *RedisService) ExecuteCode(ctx context.Context, language string, payload models.ExecuteCodePayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return r.client.RPush(ctx, language, data).Err()
}
