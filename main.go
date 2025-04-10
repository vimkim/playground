package main

import (
	"fmt"
	"strings"
	"time"
)

type Job struct {
	ID    int
	Value int
}

func worker(id int, jobs chan Job, results chan<- Job) {
	for job := range jobs {
		// fmt.Printf("Goroutine ID: %d, Job ID: %d, Job Value: %d\n", id, job.ID, job.Value)
		fmt.Printf("%s%d:%d:%d\n", strings.Repeat("        ", id), id, job.ID, job.Value)
		time.Sleep(1 * time.Second)
		if job.Value >= 10 {
			results <- job
		} else {
			job.Value++
			jobs <- job
		}
	}
}

func main() {
	arr := []int{3, 5, 11, 10, 7, 2, 8} // Input array
	jobs := make(chan Job)
	results := make(chan Job)

	// Spawn 4 goroutines
	for i := 1; i <= 4; i++ {
		go worker(i, jobs, results)
	}

	go func() {
		// Convert array elements to Job structs and send to jobs channel
		for i, val := range arr {
			jobs <- Job{ID: i + 1, Value: val}
		}
	}()

	resultArr := make([]Job, len(arr))
	for range arr {
		result := <-results
		fmt.Printf("Job ID: %d, Final Value: %d\n", result.ID, result.Value)
		resultArr = append(resultArr, result)
	}
	for i, val := range resultArr {
		fmt.Printf("Result %d: ID: %d, Value: %d\n", i+1, val.ID, val.Value)
	}

	close(jobs)
	close(results)
}
