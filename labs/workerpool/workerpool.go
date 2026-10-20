package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func Run(ctx context.Context, n int, jobs []func(context.Context) error) error {
	var wg sync.WaitGroup
	errs := make([]error, len(jobs))
	ch := make(chan struct{}, n)
	for i, job := range jobs {
		ch <- struct{}{}
		wg.Go(func() {
			defer func() { <-ch }()
			err := job(ctx)
			errs[i] = err
		})
	}
	wg.Wait()
	return errors.Join(errs...)
}
func Work(idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Printf("Finished job %v\n", idx)

}
func main() {
	// var wg sync.WaitGroup
	n := 5
	jobs := make([]func(context.Context) error, 0, n)
	for i := range 10 {
		jobs = append(jobs,
			func(ctx context.Context) error {
				time.Sleep(time.Second)
				fmt.Printf("Job %v done\n", i)
				return nil
			})
	}

	err := Run(context.Background(), n, jobs)
	fmt.Printf("%T\n", err)
	fmt.Printf("Error: %v", err)

}
