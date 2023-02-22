
/*
888888 888888 88   88 8888P 8888P 888888 88""Yb 
  88   88__   88   88   dP    dP  88__   88__dP 
  88   88""   Y8   8P  dP    dP   88""   88"Yb  
  88   88     'YbodP' d8888 d8888 888888 88  Yb

by n0nexist.github.io / www.n0nexist.gq
a web fuzzer written in golang
*/

package main

import (
	"fmt"
	"os"
	"bufio"
	"sync"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"strconv"
)

// global variables
var (

	// colors
	reset  = "\033[0m"
	cyan   = "\033[36m"
	blue   = "\033[34m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	magenta = "\033[35m"

	// other global variables
	myurl = ""
	totalrequests = 0
	donerequests = 0
	percentage = 0.0
)

// stick common status codes togheter with colors
func replaceStatusCode(statusCode int) string {
	switch statusCode {
		case 500:
			return "\033[33m500\033[0m (Internal server error)" // Yellow
		case 503:
			return "\033[33m503\033[0m (Service unavailable)" // Yellow
		case 200:
			return "\033[32m200\033[0m (Ok)" // Green
		case 301:
			return "\033[35m301\033[0m (Redirect)" // Magenta
		case 302:
			return "\033[35m302\033[0m (Temp. redirect)" // Magenta
		// Add other status codes here with their respective colors
		default:
			return strconv.Itoa(statusCode)
	}
}

// function that does the request
func doRequest(url string, payload string) {
	// Increment done requests
	donerequests++

	// Replace "tFUZZER" in the URL with the payload
	url = strings.Replace(url, "tFUZZER", payload, -1)

	// Make a GET request to the URL
	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the status code, response word count, response size, and response time
	var statuscode = resp.StatusCode;
	var words = len(strings.Split(string(body), " "));
	var respsize = len(body);

	if statuscode != 404 {
		fmt.Printf("%s[%s%s%s] %s| %s %s=>%s %d words (%s%d%s bytes) %s=>%s %s%v%s\n", reset,blue,payload,reset,reset,replaceStatusCode(statuscode),yellow,reset,words,green,respsize,reset,yellow,reset,green,elapsed,reset);
	}
	
}

// check if the file exists
func doesFileExist(fileName string) bool{
	_ , error := os.Stat(fileName)
  
	if os.IsNotExist(error) {
	  return false
	} else {
	  return true
	}
}

// shows help & quit
func showhelp(){
	fmt.Printf("%sERROR %s=>%s ./tFuzzer %s(%surl.com/query?=%stFUZZER%s) (%swordlist%s) (%sthreads%s)%s\n",red,yellow,reset,red,reset,yellow,red,reset,red,reset,red,reset)
	os.Exit(0)
}


// function called by a thread; it makes the request
func worker(mycurrentline chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range mycurrentline {
		percentage = float64(donerequests) / float64(totalrequests) * 100.0
		fmt.Printf("[%s%d %srequests done out of %s%d%s total, %s%.2f%s%%%s] \r",blue,donerequests,reset,blue,totalrequests,reset,blue,percentage,cyan,reset)
		doRequest(myurl, line)
	}
}

// get the number of lines in a file
func getLines(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		showhelp()
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		showhelp()
	}

	return lineCount
}

func main() {

	// Check system arguments
	if len(os.Args) != 4 {
		showhelp()
	}

	myurl = os.Args[1]
	mywordlist := os.Args[2]
	mythreads,isAnumber := strconv.Atoi(os.Args[3])

	// check validity of system arguments
	if !strings.Contains(myurl,"tFUZZER") || !doesFileExist(mywordlist) || isAnumber != nil {
		showhelp()
	}

	fmt.Print(cyan)

	// tFuzzer starts from here
	fmt.Println(`
888888 888888 88   88 8888P 8888P 888888 88""Yb 
  88   88__   88   88   dP    dP  88__   88__dP 
  88   88""   Y8   8P  dP    dP   88""   88"Yb  
  88   88     'YbodP' d8888 d8888 888888 88  Yb`)

	fmt.Println(blue)

	fmt.Println("[ get request fuzzer ]\n[ by github.com/n0nexist ]\n",reset)

	fmt.Printf("%sFUZZING %s=>%s %s\n%sWORDLIST %s=>%s %s\n%sTHREADS %s=>%s %d%s\n\n",cyan,yellow,blue,myurl,cyan,yellow,blue,mywordlist,cyan,yellow,blue,mythreads,reset)

	totalrequests = getLines(mywordlist)

	// Open the file containing the lines to process
	file, err := os.Open(mywordlist)
	if err != nil {
		showhelp()
	}
	defer file.Close()

	// Create a scanner to read the lines from the file
	scanner := bufio.NewScanner(file)

	// Create a channel to pass lines to the worker goroutines
	mycurrentline := make(chan string)

	// Create a wait group to wait for all worker goroutines to finish
	var wg sync.WaitGroup

	// Start the worker goroutines
	for i := 0; i < mythreads; i++ {
		wg.Add(1)
		go worker(mycurrentline, &wg)
	}

	// Read lines from the file and pass them to the worker goroutines
	for scanner.Scan() {
		line := scanner.Text()
		mycurrentline <- line
	}

	// Close the channel to signal the worker goroutines to exit
	close(mycurrentline)

	// Wait for all worker goroutines to finish
	wg.Wait()

	// All done
	fmt.Printf("%s+%s done fuzzing\n",cyan,reset)
}
