package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/pdcgo/scheduler"
)

type DBtaskStore struct{}

func (d *DBtaskStore) TaskUpdate(state *scheduler.TaskState) {
	log.Println(state.ID, state.Status)
}
func (d *DBtaskStore) TaskAdd(state *scheduler.TaskState) {}
func (d *DBtaskStore) TaskSummary() (*scheduler.TaskSummary, error) {
	sum := scheduler.TaskSummary{}
	return &sum, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 30)
		cancel()
	}()

	store := DBtaskStore{}
	sched := scheduler.NewScheduler(ctx, &store, 4)

	for i := 0; i < 500; i++ {
		sec := rand.Intn(5)

		cancelTask := sched.SpawnTask("test", "", func(state *scheduler.TaskState) error {

			time.Sleep(time.Second * time.Duration(sec))
			return nil
		})

		if sec == 4 {
			cancelTask()
		}
	}
}
