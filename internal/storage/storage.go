package storage

import (
	"errors"
	"sync"
	"time"

	"workmate/internal/models"

	"github.com/google/uuid"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrTaskRunning  = errors.New("task is running")
)

// StorageI - репозиторий
type StorageI interface {
	Add() string
	GetInfo(id string) (models.TaskInfo, error)
	GetTasks() []models.TaskID
	Delete(id string) error
}

// Storage - in-memory реализация репозитория
type Storage struct {
	mu    sync.RWMutex
	tasks map[string]models.TaskI
}

func New() *Storage {
	return &Storage{
		tasks: make(map[string]models.TaskI),
	}
}

// Add создает новую задачу и запускает её выполнение в отдельной горутине.
//
// В текущей реализации models.Task Worker Pool избыточен,
// т.к. задачи i/o bound и не обращаются ко внешним сервисам. В дальнейшем можно прикрутить Worker Pool над Storage при необходимости.
func (s *Storage) Add() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()
	task := &models.Task{
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	s.tasks[id] = task

	go func(t *models.Task) {
		time.Sleep(time.Second * 20) // чтобы увидеть Pending статус, т.к задача сразу начинает выполнятся
		t.Do()
	}(task)

	return id
}

// GetInfo возвращает информацию по задаче.
func (s *Storage) GetInfo(id string) (models.TaskInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.TaskInfo{}, ErrTaskNotFound
	}

	return task.Info(), nil
}

// GetTasks возвращает список всех задач.
func (s *Storage) GetTasks() []models.TaskID {
	var res []models.TaskID
	for id := range s.tasks {
		res = append(res, models.TaskID{ID: id})
	}
	return res
}

// Delete Удаляет задачу из хранилища, при этом нельзя удалить задачу со статусом "running" и "pending".
// В дальнейшем можно реализовать отмену запущенных задач флагом или контекстом.
func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}

	info := task.Info()
	if info.Status == models.StatusRunning || info.Status == models.StatusPending {
		return ErrTaskRunning
	}

	delete(s.tasks, id)
	return nil
}
