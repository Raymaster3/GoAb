package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Configuration
var numberOfCalls int = 1
var concurrentCalls int = 1
var keepAlive bool = false
var serverLink string

// Metrics
var failedRequests int32 = 0
var totalTimeOfRequests int64 = 0

// Temporals
var callsLeft int32
var callsDone int32 = 0

// Synchronization
var wg sync.WaitGroup

// Getting the program arguments and setting the variables accordingly
func inputConfig() {
	if len(os.Args) > 1 { // We are changing the default config
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-n":
				var err error
				numberOfCalls, err = strconv.Atoi(os.Args[i+1])
				if err == nil {
					if numberOfCalls <= 0 {
						numberOfCalls = 1
					}
					i++
				} else {
					fmt.Println("Error")
				}
			case "-c":
				var err error
				concurrentCalls, err = strconv.Atoi(os.Args[i+1])
				if err == nil {
					if concurrentCalls <= 0 {
						concurrentCalls = 1
					}
					i++
				} else {
					fmt.Println("Error")
				}
			case "-k":
				keepAlive = true
			default:
				serverLink = os.Args[i]
			}
		}
		callsLeft = int32(numberOfCalls)
	}
}

func makeCall(t *http.Transport, m sync.Mutex) {

	client := &http.Client{
		Transport: t, // Shared transport by default
		Timeout:   10 * time.Second,
	}

	// Using the same transport for all the connections allow us to use keep-alive
	// To not use the functionality we need to use the default one for each call
	if !keepAlive {
		client.Transport = http.DefaultTransport
	}

	startTime := time.Now()
	res, err1 := client.Get(serverLink)

	m.Lock()
	if callsDone <= int32(concurrentCalls) {
		totalTimeOfRequests += time.Since(startTime).Milliseconds()
	}
	m.Unlock()

	if err1 != nil {
		atomic.AddInt32(&failedRequests, 1)
	} else {
		if res.StatusCode < 200 || res.StatusCode > 299 {
			fmt.Println(res.StatusCode)
			atomic.AddInt32(&failedRequests, 1)
		}
	}

	atomic.AddInt32(&callsDone, 1)

	// If the body of the message is read to completion and then closed the next call
	// may reuse the existing connection
	if keepAlive && err1 == nil {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}

	if callsLeft > 0 {
		atomic.AddInt32(&callsLeft, -1)
		// Make another call
		makeCall(t, m)
	} else {
		wg.Done()
	}
}

func printConfig() {
	fmt.Println("Total calls:", numberOfCalls)
	fmt.Println("Concurrency level:", concurrentCalls)
	fmt.Println("Keep-Alive:", keepAlive)
	fmt.Println("Server URL:", serverLink)
	fmt.Println("")
}
func printResults(elapsedTime int64) {
	fmt.Println("Failed calls:", failedRequests)
	fmt.Println("Failed calls %:", float64(float64(failedRequests)/float64(numberOfCalls)*100.0), "%")
	fmt.Println("")
	fmt.Println("Time per request:", float64(totalTimeOfRequests)/float64(concurrentCalls), "(mean [ms])")
	fmt.Println("Time per request:", float64(elapsedTime)/float64(numberOfCalls), "(mean across all concurrent requests [ms])")
}
func main() {
	inputConfig()
	printConfig()

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 200
	t.MaxConnsPerHost = 200
	t.MaxIdleConnsPerHost = 200

	m := sync.Mutex{}

	start := time.Now()

	for j := 0; j < concurrentCalls; j++ {
		wg.Add(1)
		callsLeft--
		// We make the concurrent calls
		go makeCall(t, m)
	}
	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("Test time:", elapsed)
	fmt.Println("TPS(#/sec):", float64(numberOfCalls)/float64(elapsed.Seconds()))

	printResults(elapsed.Milliseconds())
}
