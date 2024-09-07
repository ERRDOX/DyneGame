// go:build exclude

//go:build windows
// +build windows

package main

/*
#cgo CFLAGS: -Wall
#cgo LDFLAGS: -lgdi32 -luser32

#include <windows.h>

void listenKeys() {
    while (1) {
        for (int key = 8; key <= 190; key++) {
            SHORT state = GetAsyncKeyState(key);
            if (state & 0x8000) {
                printf("Key Pressed: %d\n", key);
            } else if (state & 1) {
                printf("Key Released: %d\n", key);
            }
        }
        Sleep(50);
    }
}
*/
import "C"

func main() {
	println("Starting key listener...")
	C.listenKeys() // Calls the C function to listen to key events
}
