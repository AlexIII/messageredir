package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a buffered channel with a capacity of 1
	ch := make(chan int)

	// Simulate sending data without waiting for a receiver
	go func() {
		for i := 0; i < 100; i++ {
			select {
			case ch <- i: // Attempt to send value to the channel
				fmt.Println("Sent:", i)
			default: // If the channel is full, discard the value
				fmt.Println("Discarded:", i)
			}
			time.Sleep(100 * time.Millisecond) // Simulate some work
		}
	}()

	// Simulate a receiver with some delay
	time.Sleep(1300 * time.Millisecond)
	fmt.Println("Receiving values:")

	// Receive values from the channel
	for i := 0; i < 3; i++ {
		val := <-ch // Wait for a value to be received
		fmt.Println("Received:", val)
	}
}
