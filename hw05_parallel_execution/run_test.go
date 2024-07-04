package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		tasks = append(tasks, func() error {
			err := fmt.Errorf("error from task %d", 1)
			time.Sleep(time.Millisecond * 2)
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
		tasks = append(tasks, func() error {
			err := fmt.Errorf("error from task %d", 2)
			time.Sleep(time.Millisecond * 2)
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
		tasks = append(tasks, func() error {
			err := fmt.Errorf("error from task %d", 3)
			time.Sleep(time.Millisecond * 5)
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
		tasks = append(tasks, func() error {
			err := fmt.Errorf("error from task %d", 4)
			time.Sleep(time.Millisecond * 3)
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
		tasks = append(tasks, func() error {
			err := fmt.Errorf("error from task %d", 5)
			time.Sleep(time.Millisecond * 1)
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})

		workersCount := 2
		maxErrorsCount := 2
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		tasks = append(tasks, func() error {
			fmt.Println("Task 1; Duration 2")
			time.Sleep(time.Millisecond * 2)
			atomic.AddInt32(&runTasksCount, 1)
			fmt.Println("Task 1 ended")
			return nil
		})
		tasks = append(tasks, func() error {
			fmt.Println("Task 2; Duration 2")
			time.Sleep(time.Millisecond * 2)
			atomic.AddInt32(&runTasksCount, 1)
			fmt.Println("Task 2 ended")
			return nil
		})
		tasks = append(tasks, func() error {
			fmt.Println("Task 3; Duration 5")
			time.Sleep(time.Millisecond * 5)
			atomic.AddInt32(&runTasksCount, 1)
			fmt.Println("Task 3 ended")
			return nil
		})
		tasks = append(tasks, func() error {
			fmt.Println("Task 4; Duration 3")
			time.Sleep(time.Millisecond * 3)
			atomic.AddInt32(&runTasksCount, 1)
			fmt.Println("Task 4 ended")
			return nil
		})
		tasks = append(tasks, func() error {
			fmt.Println("Task 5; Duration 1")
			time.Sleep(time.Millisecond * 1)
			atomic.AddInt32(&runTasksCount, 1)
			fmt.Println("Task 5 ended")
			return nil
		})

		workersCount := 2
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}
