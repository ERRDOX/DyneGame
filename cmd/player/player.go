package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/eiannone/keyboard"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(conn)
	defer conn.Close()

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key != 0 {
			// Handling special keys
			var keyString string
			switch key {
			case keyboard.KeyArrowUp:
				keyString = "38"
			case keyboard.KeyArrowDown:
				keyString = "40"
			case keyboard.KeyArrowLeft:
				keyString = "37"
			case keyboard.KeyArrowRight:
				keyString = "39"
			default:
				keyString = fmt.Sprintf("%v", key)
			}
			conn.Write([]byte(keyString))
		} else if char != 0 {
			conn.Write([]byte(string(char)))
		}
		data, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		fmt.Print("server message: received  data:", data, "    ")

		fmt.Printf("You pressed: %q, key code: %v\n", char, key)
		if key == keyboard.KeyEsc {
			break
		}
	}

	fmt.Println("Program exited")
}
