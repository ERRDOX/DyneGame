// go:build exclude
package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("Could not initialize SDL: %s\n", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Key Press/Release Example", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Could not create window: %s\n", err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		log.Fatalf("Could not get window surface: %s\n", err)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.KeyboardEvent:
				if t.State == sdl.PRESSED {
					fmt.Printf("Key Pressed: %s\n", sdl.GetKeyName(sdl.Keycode(t.Keysym.Sym)))
				} else if t.State == sdl.RELEASED {
					fmt.Printf("Key Released: %s\n", sdl.GetKeyName(sdl.Keycode(t.Keysym.Sym)))
				}
			}
		}

		// Update the window surface
		surface.FillRect(nil, sdl.MapRGB(surface.Format, 0, 0, 0))
		window.UpdateSurface()

		// Small delay to avoid high CPU usage
		sdl.Delay(16)
	}

	fmt.Println("Program exited")
}
