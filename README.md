# GoAb
To use goab you first need to have [golang installed](https://go.dev/doc/install) in your computer.

## How to use
To execute goab just simply type in the terminal “go run goab.go [OPTIONS] [URL]”. 

The following are the available options:
* -n: sets the total number of requests, by default is 1
* -c: sets the concurrency level, by default is 1
* -k: reuse existing connections
## How it works
GoAb works by using an initial loop creating exactly c (concurrency level) calls. Each of this concurrent calls will make a new call on finnish if there are still calls to be made. This ensures that there will be exactly c calls being executed at any given time during the execution of the test. 
## Results and comparatives
In order to confirm that goab is working as expected a local nginx server was used. We will be testing and comparing the results obtained with goab and ab on this server. 
### AB
##### With -n 100 -c 20:
* TPS: 18900
* Test time: 5 ms
* Time per request: 1ms
##### Adding -k:
* TPS: 27480
* Test time: 4 ms
* Time per request: 0.73ms
##### With -n 1000 -c 200:
* TPS: 18600
* Test time: 10 ms
* Time per request: 10ms
##### Adding -k:
* TPS: 33700
* Test time: 54 ms
* Time per request: 5.9 ms
##### With -n 1000 -c 800:
* TPS: 17900
* Test time: 56 ms
* Time per request: 44.5 ms
* Failed: 64
  ###### These failed requests are because the server has assigned only 768 workers, so ab is giving us 32 errors by exception and an extra 32 errors by discrepancies in the length of the message. Also, because we are executing the test in the same machine as the server is running we can see that there is a correlation between the concurrency level and the request time. This happens because as we increase the amount of concurrent calls of our program, we are limiting the resources available to serve our page, and this makes the server take more time per request.
##### Adding -k:
* TPS: 16400
* Test time: 61 ms
* Time per request: 48ms

### GOAB
##### With -n 100 -c 20:
* TPS: 3660
* Test time: 27 ms
* Time per request: 1.1ms
##### Adding -k:
* TPS: 4200
* Test time: 23 ms
* Time per request: 1.4ms
##### With -n 1000 -c 200:
* TPS: 5950
* Test time: 160 ms
* Time per request: 56 ms
##### Adding -k:
* TPS: 7220
* Test time: 138 ms
* Time per request: 36 ms  
##### With -n 1000 -c 800:
* TPS: 5400
* Test time: 184 ms
* Time per request: 141 ms
  ###### In this case there are no errors because the threads are getting blocked by synchronization and they are not really simultaneous.
##### Adding -k:
* TPS: 8550
* Test time: 116 ms
* Time per request: 64 ms 
## Observations
As we can see the implementation is far from perfect, in fact, the numbers are not even similar in some cases. I believe this is caused by my implementation of parallelism and data sharing. This would explain why there are no errors when making more concurrent calls than the server can handle, and why the times are so big. But i can confidently say that the keep-alive functionality is working properly and is opening new connections only in the first batch of calls.  
Its also worth noting that performing the same tests in a remote server does give us a more natural and similar result. For example using facebook.com as the url and -n 1000 -c 800, we get 175 transfers per second for GOAB and 200 transfers per second in the case of AB.
