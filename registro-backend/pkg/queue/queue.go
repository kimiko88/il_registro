package queue

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Task represents a background job payload.
type Task struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskHandler func(ctx context.Context, task *Task) error

// Queue defines the interface for an asynchronous background job queue.
type Queue interface {
	Enqueue(ctx context.Context, taskType string, payload interface{}) (*Task, error)
	RegisterHandler(taskType string, handler TaskHandler)
	Start()
	Stop()
}

const defaultQueueKey = "registro:queue:tasks"

// RedisQueue implements an asynchronous distributed task queue using Redis lists.
type RedisQueue struct {
	rdb      *redis.Client
	handlers map[string]TaskHandler
	mu       sync.RWMutex
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		rdb:      client,
		handlers: make(map[string]TaskHandler),
		stopChan: make(chan struct{}),
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, taskType string, payload interface{}) (*Task, error) {
	var payloadBytes []byte
	var err error
	if b, ok := payload.([]byte); ok {
		payloadBytes = b
	} else {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	task := &Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   payloadBytes,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	if err := q.rdb.LPush(ctx, defaultQueueKey, data).Err(); err != nil {
		return nil, err
	}

	return task, nil
}

func (q *RedisQueue) RegisterHandler(taskType string, handler TaskHandler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[taskType] = handler
}

func (q *RedisQueue) Start() {
	q.wg.Add(1)
	go q.workerLoop()
}

func (q *RedisQueue) workerLoop() {
	defer q.wg.Done()
	for {
		select {
		case <-q.stopChan:
			return
		default:
			res, err := q.rdb.BRPop(context.Background(), 1*time.Second, defaultQueueKey).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				continue
			}

			if len(res) < 2 {
				continue
			}

			var task Task
			if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
				continue
			}

			q.mu.RLock()
			handler, exists := q.handlers[task.Type]
			q.mu.RUnlock()

			if exists && handler != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				_ = handler(ctx, &task)
				cancel()
			}
		}
	}
}

func (q *RedisQueue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
	_ = q.rdb.Close()
}

// MemoryQueue is an in-memory task queue used for testing or fallback.
type MemoryQueue struct {
	tasks    chan *Task
	handlers map[string]TaskHandler
	mu       sync.RWMutex
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		tasks:    make(chan *Task, 1000),
		handlers: make(map[string]TaskHandler),
		stopChan: make(chan struct{}),
	}
}

func (q *MemoryQueue) Enqueue(_ context.Context, taskType string, payload interface{}) (*Task, error) {
	var payloadBytes []byte
	var err error
	if b, ok := payload.([]byte); ok {
		payloadBytes = b
	} else {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	task := &Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   payloadBytes,
		CreatedAt: time.Now(),
	}

	select {
	case q.tasks <- task:
		return task, nil
	default:
		return nil, errors.New("queue is full")
	}
}

func (q *MemoryQueue) RegisterHandler(taskType string, handler TaskHandler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[taskType] = handler
}

func (q *MemoryQueue) Start() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-q.stopChan:
				return
			case task := <-q.tasks:
				q.mu.RLock()
				handler, exists := q.handlers[task.Type]
				q.mu.RUnlock()
				if exists && handler != nil {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					_ = handler(ctx, task)
					cancel()
				}
			}
		}
	}()
}

func (q *MemoryQueue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
}

// NewQueue creates a Queue backed by Redis if available, or in-memory fallback otherwise.
func NewQueue(redisURL string) Queue {
	if strings.TrimSpace(redisURL) != "" {
		opt, err := redis.ParseURL(redisURL)
		if err == nil {
			client := redis.NewClient(opt)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			pingErr := client.Ping(ctx).Err()
			if pingErr == nil {
				log.Println("queue: connected to Redis backend for distributed background job queue")
				return NewRedisQueue(client)
			}
			log.Printf("queue: Redis ping failed (%v), falling back to in-memory queue", pingErr)
			_ = client.Close()
		}
	}
	return NewMemoryQueue()
}
