package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusError     TaskStatus = "error" // Если понадобится отслеживать ошибки

	timeMin = 3
	timeMax = 5
)

type TaskI interface {
	// Do Выполняет задачу
	Do()
	// Info возвращает информацию о задаче в TaskInfo
	Info() TaskInfo
}

type TaskStatus string

type TaskInfo struct {
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	Duration  string     `json:"duration"`
}

type TaskID struct {
	ID string `json:"id"`
}

// Task - реализация задачи, просто "спит" от 3 до 5 минут.
type Task struct {
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt time.Time  `json:"finished_at"`
	mu         sync.Mutex
}

func (t *Task) Do() {
	t.mu.Lock()
	t.StartedAt = time.Now()
	t.Status = StatusRunning
	t.mu.Unlock()

	num := rand.Intn(timeMax-timeMin+1) + timeMin
	time.Sleep(time.Minute * time.Duration(num))

	t.mu.Lock()
	t.FinishedAt = time.Now()
	t.Status = StatusCompleted
	t.mu.Unlock()
}

func (t *Task) Info() TaskInfo {
	t.mu.Lock()
	status := t.Status
	created := t.CreatedAt
	started := t.StartedAt
	finished := t.FinishedAt
	t.mu.Unlock()

	var duration time.Duration
	switch status {
	case StatusRunning:
		duration = time.Since(started)
	case StatusPending:
		duration = 0
	default:
		duration = finished.Sub(started)
	}
	durationStr := fmt.Sprintf("%d sec", int(duration.Seconds()))

	return TaskInfo{
		Status:    status,
		CreatedAt: created,
		Duration:  durationStr,
	}
}
