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
#### With -n 100 -c 20
