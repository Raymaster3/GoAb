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
  ###### These failed requests are because the server has assigned only 768 workers, so ab is giving us 32 errors by exception and an extra 32 errors by discrepancies in the length of the message.
##### Adding -k:
* TPS: 16400
* Test time: 61 ms
* Time per request: 48ms
