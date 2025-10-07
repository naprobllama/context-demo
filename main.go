package main

import (
	"context"
	"fmt"
	"time"
)

// leakyWorker simulates a task that ignores the context cancellation signal.
// This goroutine will continue running (and logging) indefinitely, even after
// the parent context is cancelled, leading to a goroutine leak.
func leakyCauldron(ctx context.Context) {
	fmt.Printf("Entering the Leaky Cauldron. It will never exit gracefully.\n")

	// This worker ignores the context, leading to a leak.
	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Leaky Cauldron Doing work...\n")
	}
}

// hogwarts simulates a task that checks the context cancellation signal.
// It now uses context.Cause() to report the specific reason for cancellation.
func hogwarts(ctx context.Context) {
	fmt.Printf("Entering Hogwarts. It will check if ctx.Done().\n")

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // Always stop timers/tickers when done

	for {
		select {
		case <-ticker.C:
			// Simulates doing some periodic work
			fmt.Print("Hogwarts Doing work...\n")

		case <-ctx.Done():
			// **CRITICAL:** The context was cancelled.
			fmt.Print("Hogwart's received cancellation signal from ctx.Done(). Exiting now.\n")

			// ctx.Err() will now contain the basic cancellation error (e.g., context canceled)
			fmt.Printf("Cancellation error (ctx.Err()): %v\n", ctx.Err())

			// Use context.Cause() to retrieve the specific error passed during the cancel call.
			cause := context.Cause(ctx)
			fmt.Printf("Cancellation cause (context.Cause()): %v\n", cause)

			return // Exit the goroutine cleanly
		}
	}
}

func main() {
	fmt.Print("\n\nStarting Context Demonstration with Cancel Cause...\n\n")
	fmt.Println("---------------------------------------------------")

	ctx, cancel := context.WithCancelCause(context.Background())

	// Use defer to call cancel with a nil cause for standard function exit cleanup.
	defer cancel(nil)

	// Start both workers
	go leakyCauldron(context.Background())
	go hogwarts(ctx)

	// Let the workers run for a short time
	fmt.Println("\nAllowing workers to run for 1.5 seconds...")
	time.Sleep(1500 * time.Millisecond)

	// Cancel the context, providing a specific cause.
	causeError := fmt.Errorf("Voldemort is here: all tasks stopped")
	fmt.Printf("\n>>> Calling cancel(cause) with cause: '%v' <<<\n", causeError)
	cancel(causeError) // Pass the cause error here

	// 5. Wait to see the effect
	fmt.Print("Waiting 2 seconds for workers to respond to cancellation...\n\n\n")
	time.Sleep(2000 * time.Millisecond)

	fmt.Print("\n\n---------------------------------------------------\n")
	fmt.Print("Demonstration complete. \n\n")
	fmt.Println("Hogwarts has shutdown gracefully, reporting the 'voldemort is here' cause.")
	fmt.Println("Leaky Cauldron is still running (goroutine leak).")
}
