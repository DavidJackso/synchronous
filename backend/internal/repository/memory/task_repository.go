package memory

import (
	"fmt"
	"sync"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type TaskRepository struct {
	tasks map[string]*entity.Task
	mu    sync.RWMutex
}

func NewTaskRepository() interfaces.TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*entity.Task),
	}
}

func (r *TaskRepository) Create(task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) GetByID(id string) (*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	return task, nil
}

func (r *TaskRepository) GetBySessionID(sessionID string) ([]*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []*entity.Task
	for _, task := range r.tasks {
		if task.SessionID == sessionID {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) Update(task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("task with ID %s not found", id)
	}

	delete(r.tasks, id)
	return nil
}
