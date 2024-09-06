//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

var ws js.Value // Global variable to hold the WebSocket object

func main() {
	// Create a channel to block the main function from exiting
	c := make(chan struct{}, 0)

	// Register a callback function for keyboard input
	js.Global().Set("sendKeyPress", js.FuncOf(keyPress))
	js.Global().Set("sendKeyRelease", js.FuncOf(keyRelease)) // Register key release function

	// Open a WebSocket connection
	connectWebSocket()

	// Block forever to keep the Go runtime alive
	<-c
}

// connectWebSocket sets up the WebSocket connection
func connectWebSocket() {
	ws = js.Global().Get("WebSocket").New("ws://localhost:8080/ws")

	// Define event handlers for the WebSocket
	ws.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("WebSocket connection opened")
		return nil
	}))

	ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].Get("data").String()
		fmt.Println("Received message:", data)
		return nil
	}))

	ws.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("WebSocket error")
		return nil
	}))

	ws.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("WebSocket connection closed")
		return nil
	}))
}

// keyPress is a callback function that handles keypress events
func keyPress(this js.Value, p []js.Value) interface{} {
	key := p[0].String() // Directly get the key name
	fmt.Println("Key pressed:", key)

	// Send the key to WebSocket server
	if ws.Truthy() {
		ws.Call("send", key)
	} else {
		fmt.Println("WebSocket is not connected")
	}

	return nil
}

// keyRelease is a callback function that handles key release events
func keyRelease(this js.Value, p []js.Value) interface{} {
	key := p[0].String()
	fmt.Println("Key released:", key)

	// Send the key release to WebSocket server
	if ws.Truthy() {
		ws.Call("send", "release:"+key) // Example of sending key release event
	} else {
		fmt.Println("WebSocket is not connected")
	}

	return nil
}
