package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	for _, link := range links {
		checkLink(link)
	}
}

func checkLink(link string) {
	_, err := http.Get(link) // Blocking call. When this runs, Main Go routine can do nothing else!
	if err != nil {
		fmt.Println(link, "might be down!")
		return
	}

	fmt.Println(link, "is up!")
}

// Notes:
// 1. How our code is being executed right now:
//    With our slice of links -> Take first link from slice -> Make request -> GET http://google.com -> Wait for a response, log it
//							  -> Take next link, make request -> GET http://facebook.com -> Wait for a response, log it
//							  -> Repeat for others
//
//    Basically it's in series, every single time we're making a request, We sit around and wait for the response to come back before making the next
//    So our aim to to make requests in parallel
//
// 2. When we launch a Go program, i.e. when we compile it and execute it, we automatically create one Go routine. You can think of a routine as being
//    something that exists inside of our running program or our process. This go routine takes every line of code inside of our program and executes them one by one.
//    Actual compiled form of our code might look a little bit differently than what we have.
//
// 3. Syntax of a Go routine:
//    go checkLink(link)
//    go - creates a new thread go routine
//    checkLink(link) - the function the newly created thread would run
//
// 4. What happens when we spawn multiple Go routines inside our program?
//                       One CPU Core
//                             |
//                             v
//                       Go Scheduler                      <- Scheduler runs one routine until it finishes or makes a blocking call (like an HTTP request)
//           / \              / \              / \
//            |                |                |
//            V                V                V
//       Go Routine       Go Routine       Go Routine
//
//   The most important thing to understand here is that even though we are launching multiple routines, only one is being executed or running at any given time.
//   So the purpose of this Go scheduler is to monitor the code that is running inside of each of these Go Routines. As soon as the scheduler detects that one routine
//   has finished running all of the code inside of it, so essentially all the code inside of a given function or when the scheduler it detects that a function has made a
//   blocking call like the HTTP request that we are making then it says okay you know what? You Go routine right here, you thing that just finished or has some blocking
//   code that is being executed. You're done for right now. We are going to pause you. And instead we're going to start executing this other Go routine. So essentially
//   even though we are spawning multiple Go routines, they are not actually being executed truly at the same time. Whenever we have one CPU, so this one CPU is only running
//   the code inside of one Go routine at a time and we rely upon this go scheduler to decide which Go routine is being executed. So in the blink of an eye like we might run
//   this routine right here for a fraction, then a fraction of a second and then jump over to that and then jump back over to this one. Thus, the scheduler is working very
//   quickly behind the scenes to handle all these different routines as best as it can and cycle through them very very quickly.
//
// 5. What happens when we have multiple CPU cores on our local machine?
//    By default, Go tries to only use ONE core, but we can easily change that
//    One CPU Core       One CPU Core      One CPU Core
//       |  / \             |  / \           |  / \
//       v   |              v   |            v   |
//                       Go Scheduler                      <- Scheduler runs one thread on each "logical" core
//           / \              / \              / \
//            |                |                |
//            V                V                V
//       Go Routine       Go Routine       Go Routine
//
//   when we have multiple CPU cores, each one can run one single Go routine at a time.
//   And so the Go scheduler might say oh okay we've got three separate routines and we have three separate CPU cores.
//   So rather than monitoring each routine and attempting to run only one at a time, the scheduler will instead assign one routine
//   to this core, another one to the second core and the last one to the third core. So soon as we have multiple CPU cores then
//   we're talking about running multiple chunks of code truly concurrently.
//
// 6. Concurrency - we can have multiple threads executing code. If one thread blocks, another one is picked up and worked on
//                         One Core
//                             |
//                             v
//                    Pick one Go routine
//           / \              / \              / \
//            |                |                |
//            V                V                V
//       Go Routine       Go Routine       Go Routine
//
//   So when we say something is concurrent we are simply saying that our program has the ability to run  different things kind of at the same time
//   but not really at the same time because we have one core. We're only picking one Go routine. So all we're saying with concurrency is that we can
//   kind of schedule work to be done throughout each other. We don't necessarily have to wait for one Go routine to finish before going onto the next one.
//
// 7. Parallelism = multiple threads executed at the exact same time, like nanosecond. Requires multiple CPUs
//                         One Core                                                One Core
//                             |                                                       |
//                             v                                                       v
//             Pick     one      Go     routine                        Pick     one      Go     routine
//           / \              / \              / \                   / \              / \              / \
//            |                |                |                     |                |                |
//            V                V                V                     V                V                V
//       Go Routine       Go Routine       Go Routine            Go Routine       Go Routine       Go Routine
//
// 8. Bug we're going to see as soon as we implement Go routines:
//    Our Running Program
//    Main routine - created when we launched the program
//    Child Go routine  ---\
//    Child Go routine  -------> Child routines created by 'go' keyword
//    Child Go routine  ---/
//
//    Child routines are not quite given the same level of respect, I guess for lack of a better term, we'll say respect as the main routine is.
