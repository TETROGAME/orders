package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func Run(ctx context.Context, n int, jobs []func(context.Context) error) error {
	jobs_count := len(jobs)
	if n <= 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}
	var wg sync.WaitGroup
	errs := make([]error, jobs_count)
	ch := make(chan struct{}, n)
	for i, job := range jobs {
		if ctx.Err() != nil {
			wg.Wait()
			return ctx.Err()
		}
		select {
		case ch <- struct{}{}:
			wg.Go(func() {
				defer func() { <-ch }()
				err := job(ctx)
				errs[i] = err
			})
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		}

	}
	wg.Wait()
	return errors.Join(ctx.Err(), errors.Join(errs...))
}

func Work(idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Printf("Finished job %v\n", idx)
}

func main() {
	n := 20
	jobs_count := 10
	jobs := make([]func(context.Context) error, 0, jobs_count)
	for i := range 10 {
		jobs = append(jobs,
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				fmt.Printf("Job %v done\n", i)
				return nil
			})
	}

	err := Run(context.Background(), n, jobs)
	fmt.Printf("Error: %v", err)
}
