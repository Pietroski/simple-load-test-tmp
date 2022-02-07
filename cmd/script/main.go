package main

import (
	"YaloTest/internal/mock"
	"YaloTest/internal/model/messages"
	yalohttp "YaloTest/internal/transport/http"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	wg sync.WaitGroup
)

func main() {
	fmt.Println("Job began")

	fmt.Println("Loading files")
	// TODO: Load the data into memory - csv
	fileList := []messages.MessageRequest{}
	for i := 0; i < 100; i++ {
		fileList = append(fileList, mock.Data...)
	}
	fmt.Println("Finished loading data into memory")

	httpClient, err := yalohttp.NewClient()
	if err != nil {
		panic(err)
	}

	countSuccesses := int64(0)
	countFailures := int64(0)
	countPools := 0

	workerPool := make(chan int, 10)
	for _, item := range fileList {
		wg.Add(1)
		workerPool <- 1
		countPools++

		if countPools == 20 {
			countPools = 0
			<-time.After(1000 * time.Millisecond)
		}

		go func(i interface{}) {
			defer wg.Done()

			resp, err := httpClient.Post(&i)
			if err != nil {
				fmt.Println(resp.StatusCode)
				atomic.AddInt64(&countFailures, 1)
				<-workerPool
				return
			}

			atomic.AddInt64(&countSuccesses, 1)
			<-workerPool
		}(item)
	}
	fmt.Println("waiting until every goroutines graciously stops")
	wg.Wait()
	fmt.Println("all goroutines stopped")
	close(workerPool)
	fmt.Println("channel closed")

	fmt.Printf("successes: %v || failures: %v\n", countSuccesses, countFailures)

	fmt.Println("Job completed")
}
