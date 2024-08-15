package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded     = errors.New("errors limit exceeded")
	ErrInvalidGoroutinesNumber = errors.New("invalid number of goroutines")
	ErrEmptyTasksList          = errors.New("empty tasks list")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrEmptyTasksList
	}

	if n <= 0 {
		return ErrInvalidGoroutinesNumber
	}

	errCh := make(chan error)
	taskCh := make(chan Task, len(tasks))
	signalCh := make(chan struct{})

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(errCh, signalCh, taskCh)
		}()
	}

	go func() {
		defer close(errCh)
		wg.Wait()
	}()

	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	var count int
	for {
		_, ok := <-errCh
		if ok {
			count++
			if count == m {
				close(signalCh)
			}
		} else {
			break
		}
	}

	if count >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(errCh chan error, signalCh <-chan struct{}, taskCh <-chan Task) {
	for task := range taskCh {
		select {
		case <-signalCh:
			return
		default:
			if err := task(); err != nil {
				errCh <- err
			}
		}
	}
}
