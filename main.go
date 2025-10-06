package main

import (
	"context"
	"fmt"
	"time"
)

// leakyWorker simulates a task that ignores the context cancellation signal.
// This goroutine will continue running (and logging) indefinitely, even after
// the parent context is cancelled, leading to a goroutine leak.
func leakyWorker(ctx context.Context) {
	fmt.Printf("Leaky Worker started. It will never exit gracefully.\n")

	// In a real application, this loop might be processing data, listening on a
	// network connection, or performing database queries. Since it never checks
	// 'ctx.Done()', it cannot be stopped by context cancellation.
	for {
		time.Sleep(500 * time.Millisecond)
		// The worker is still running and wasting resources
		fmt.Printf("Leaky Worker: Doing work...\n")
	}
}

// safeWorker simulates a task that checks the context cancellation signal.
// This goroutine uses the ctx.Done() channel to exit cleanly.
func safeWorker(ctx context.Context) {
	fmt.Printf("Safe Worker started. It checks ctx.Done().\n")

	// This is where a real worker would typically perform its task, possibly
	// using a `select` statement to handle both work and cancellation.
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // Always stop timers/tickers when done

	for {
		select {
		case <-ticker.C:
			// Simulates doing some periodic work
			fmt.Print("Safe Worker Doing work...\n")

		case <-ctx.Done():
			// **CRITICAL:** The context was cancelled. Clean up and exit the goroutine.
			fmt.Print("Safe Worker received cancellation signal from ctx.Done(). Exiting now.\n")
			// Print the reason for cancellation
			fmt.Printf("Cancellation reason: %v\n", ctx.Err())
			return // Exit the goroutine cleanly
		}
	}
}

func main() {
	fmt.Println("Starting Context Demonstration...")
	fmt.Println("---------------------------------")

	// 1. Create a parent context and a cancellable child context.
	// `cancel` is the function we call to signal cancellation to the context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called if the function exits prematurely

	// 2. Start both workers
	go leakyWorker(context.Background()) // Note: The leaky worker often ignores the context entirely or uses a long-lived one. Here we use a different context to emphasize the leak.
	go safeWorker(ctx)

	// 3. Let the workers run for a short time
	fmt.Println("\nAllowing workers to run for 1.5 seconds...")
	time.Sleep(1500 * time.Millisecond)

	// 4. Cancel the context
	fmt.Println("\n>>> Calling cancel() to signal the context to stop <<<")
	cancel() // This closes the ctx.Done() channel

	// 5. Wait to see the effect
	// The Safe Worker will receive the signal and exit.
	// The Leaky Worker will ignore the signal and continue running.
	fmt.Println("Waiting 2 seconds for workers to respond to cancellation...")
	time.Sleep(2000 * time.Millisecond)

	fmt.Println("\n---------------------------------")
	fmt.Println("Demonstration complete.")
	fmt.Println("Safe Worker has exited gracefully.")
	fmt.Println("Leaky Worker is still running in the background (goroutine leak).")
}
