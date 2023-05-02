package scheduler

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type TaskSummary struct {
	All       int32
	Pending   int32
	Process   int32
	Failed    int32
	Cancel    int32
	Completed int32
}

type TaskStore interface {
	TaskUpdate(state *TaskState)
	TaskAdd(state *TaskState)
	TaskSummary() (*TaskSummary, error)
}

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskProcess   TaskStatus = "process"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
	TaskCancel    TaskStatus = "cancel"
)

func CreateTaskID() string {
	id := uuid.New()
	data := strings.Split(id.String(), "-")
	return data[len(data)-1]
}

type TaskState struct {
	Store        TaskStore
	ID           string
	Name         string
	Description  string
	ErrorMessage string
	Status       TaskStatus
	Progress     float32
}

func (s *TaskState) UpdateProgress(progres float32) {
	s.Progress = progres
	s.Store.TaskUpdate(s)
}

func (s *TaskState) UpdateStatus(status TaskStatus) {
	s.Status = status
	s.Store.TaskUpdate(s)
}

type Taskhandler func(state *TaskState) error

type Task struct {
	Ctx      context.Context
	State    *TaskState
	TaskFunc Taskhandler
}

func (task *Task) Run() error {
	var err error
	task.State.UpdateStatus(TaskProcess)
	defer func() {
		if err != nil {
			task.State.ErrorMessage = err.Error()
			task.State.UpdateStatus(TaskFailed)
		} else {
			task.State.UpdateStatus(TaskCompleted)
		}
	}()

	select {
	case <-task.Ctx.Done():
		task.State.UpdateStatus(TaskCancel)
		return err
	default:
		err = task.TaskFunc(task.State)
	}

	return err

}
