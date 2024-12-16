// go:build exclude

package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// Define input_event structure (from linux/input.h)
type input_event struct {
	Time  [2]uint32 // Time of event
	Type  uint16    // Event type (EV_KEY for keyboard events)
	Code  uint16    // Event code (key code)
	Value int32     // Key value (1=key press, 0=key release)
}

func main() {
	// Open the keyboard input event file (you may need to change the event number)
	file, err := os.Open("/dev/input/event0")
	if err != nil {
		fmt.Println("Error opening input device:", err)
		return
	}
	defer file.Close()

	// Read input events in an infinite loop
	for {
		var event input_event
		// Read the binary input_event from the device file
		err := binary.Read(file, binary.LittleEndian, &event)
		if err != nil {
			fmt.Println("Error reading input event:", err)
			return
		}

		// Check if it's a key event
		if event.Type == 1 {
			if event.Value == 1 {
				fmt.Printf("Key Pressed: %d\n", event.Code)
			} else if event.Value == 0 {
				fmt.Printf("Key Released: %d\n", event.Code)
			}
		}
		time.Sleep(50 * time.Millisecond) // Delay for CPU efficiency
	}
}
