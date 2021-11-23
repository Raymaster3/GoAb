# GoAb
To use goab you first need to have golang installed in your computer.

## How to use
To execute goab just simply type in the terminal “go run goab.go [OPTIONS] [URL]”. 

The following are the available options:
* -n: sets the total number of requests, by default is 1
* -c: sets the concurrency level, by default is 1
* -k: reuse existing connections
## How it works
GoAb works by using an initial loop creating exactly c (concurrency level) calls. Each of this concurrent calls will make a new call on finnish if there are still calls to be made. This ensures that there will be exactly c calls being executed at any given time during the execution of the test. 
