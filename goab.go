package main

import (
	/*"io/ioutil"
	"log"
	"net/http"*/
	"fmt"
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

func makeCall() {
	atomic.AddInt32(&callsDone, 1)
	resp, err := http.Get(serverLink)
	if err != nil {
		//log.Fatalln(err)
		atomic.AddInt32(&failedRequests, 1)
	}
	_ = resp

	if callsLeft > 0 {
		/*callsDone := numberOfCalls - callsLeft
		if callsDone%(numberOfCalls/10) == 0 {
			fmt.Println("Completed", callsDone, "requests")
		}*/
		atomic.AddInt32(&callsLeft, -1)
		// Make another call
		makeCall()
	} else {
		wg.Done()
	}

	//We Read the response body on the line below.
	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)

	}*/
	//Convert the body to type string
	//sb := string(body)
	//log.Printf(sb)
}

func printConfig() {
	fmt.Println("Total calls:", numberOfCalls)
	fmt.Println("Concurrency level:", concurrentCalls)
	fmt.Println("Keep-Alive:", keepAlive)
	fmt.Println("Server URL:", serverLink)
	fmt.Println("")
}
func printResults() {
	fmt.Println("Failed calls:", failedRequests)
	fmt.Println("Failed calls %:", float64(float64(failedRequests)/float64(numberOfCalls)*100.0), "%")
}
func main() {
	inputConfig()
	printConfig()

	start := time.Now()

	/*for i := 0; i < numberOfCalls; i += concurrentCalls {
		for j := 0; j < concurrentCalls; j++ {
			wg.Add(1)
			// We make the concurrent calls
			go func() {
				makeCall()
				wg.Done()
			}()
		}
		// We wait for the calls to end before dispatching the next wave
		wg.Wait()
	}*/

	for j := 0; j < concurrentCalls; j++ {
		wg.Add(1)
		callsLeft--
		// We make the concurrent calls
		go makeCall()
	}
	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("Test time:", elapsed)
	fmt.Println("TPS(#/sec):", float64(numberOfCalls)/float64(elapsed.Seconds()))
	printResults()
	/*resp, err := http.Get("https://google.com/")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)*/
}
