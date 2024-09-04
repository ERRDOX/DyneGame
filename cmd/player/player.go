package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/eiannone/keyboard"
)

func main() {
	port := flag.String("port", "8080", "Port to connect to")
	flag.Parse()
	conn, err := net.Dial("tcp", "localhost:"+*port)
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
				keyString = "ArrowUp"
			case keyboard.KeyArrowDown:
				keyString = "ArrowDown"
			case keyboard.KeyArrowLeft:
				keyString = "ArrowLeft"
			case keyboard.KeyArrowRight:
				keyString = "ArrowRight"
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
