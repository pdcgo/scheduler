package scheduler

import (
	"context"

	"github.com/pdcgo/common_conf/pdc_common"
)

type SchedulerConfig struct{}

type Scheduler struct {
	ctx        context.Context
	limitGuard chan int
	store      TaskStore
}

func (s *Scheduler) SpawnTask(name string, desc string, handler Taskhandler) (cancelTask func()) {
	s.limitGuard <- 1

	ctx, cancelTask := context.WithCancel(s.ctx)
	task := Task{
		Ctx: ctx,
		State: &TaskState{
			Store:       s.store,
			ID:          CreateTaskID(),
			Name:        name,
			Description: desc,
			Status:      TaskPending,
		},
		TaskFunc: handler,
	}

	s.store.TaskAdd(task.State)

	go func() {
		defer func() {
			<-s.limitGuard
		}()

		err := task.Run()
		if err != nil {
			pdc_common.ReportError(err)
		}

	}()

	return cancelTask
}

func (s *Scheduler) Summary() (*TaskSummary, error) {
	return s.store.TaskSummary()
}

func NewScheduler(ctx context.Context, store TaskStore, limit int) *Scheduler {
	guard := make(chan int, limit)
	sched := Scheduler{
		ctx:        ctx,
		limitGuard: guard,
		store:      store,
	}
	return &sched
}
