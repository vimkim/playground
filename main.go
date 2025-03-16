package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	screenWidth  = 30
	screenHeight = 20
	playerChar   = "A"
	obstacleChar = "X"
	emptyChar    = " "
)

type Obstacle struct {
	x, y int
}

func clearScreen() {
	// ANSI escape sequence to clear screen and move cursor to top-left.
	fmt.Print("\033[2J\033[H")
}

func drawScreen(playerX int, obstacles []Obstacle) {
	// Create a 2D grid representing the screen.
	screen := make([][]string, screenHeight)
	for i := 0; i < screenHeight; i++ {
		screen[i] = make([]string, screenWidth)
		for j := 0; j < screenWidth; j++ {
			screen[i][j] = emptyChar
		}
	}

	// Place obstacles
	for _, obs := range obstacles {
		if obs.y >= 0 && obs.y < screenHeight && obs.x >= 0 && obs.x < screenWidth {
			screen[obs.y][obs.x] = obstacleChar
		}
	}

	// Place player on the bottom row.
	if playerX >= 0 && playerX < screenWidth {
		screen[screenHeight-1][playerX] = playerChar
	}

	// Draw the screen.
	for _, row := range screen {
		for _, ch := range row {
			fmt.Print(ch)
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize player position.
	playerX := screenWidth / 2

	// Create a slice to hold obstacles.
	var obstacles []Obstacle

	// Set up a ticker for game updates.
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Channel to capture interrupts (Ctrl+C) to exit gracefully.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Initialize keyboard input.
	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	// Game loop
	gameOver := false
	frameCount := 0

	clearScreen()
	go func() {
		// Listen for keyboard events.
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Println("Error reading key:", err)
				return
			}
			if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
				gameOver = true
				return
			}

			// Move left or right
			if char == 'a' || key == keyboard.KeyArrowLeft {
				if playerX > 0 {
					playerX--
				}
			} else if char == 'd' || key == keyboard.KeyArrowRight {
				if playerX < screenWidth-1 {
					playerX++
				}
			}
		}
	}()

	for !gameOver {
		select {
		case <-sigChan:
			gameOver = true
		case <-ticker.C:
			// Every few frames, add a new obstacle at a random column at the top.
			if frameCount%10 == 0 {
				newObs := Obstacle{
					x: rand.Intn(screenWidth),
					y: 0,
				}
				obstacles = append(obstacles, newObs)
			}

			// Move obstacles down.
			newObstacles := []Obstacle{}
			for _, obs := range obstacles {
				obs.y++
				// Check collision: if an obstacle reaches the bottom row and collides with the player.
				if obs.y == screenHeight-1 && obs.x == playerX {
					gameOver = true
				}
				// Keep obstacles that are still visible.
				if obs.y < screenHeight {
					newObstacles = append(newObstacles, obs)
				}
			}
			obstacles = newObstacles

			// Redraw screen.
			clearScreen()
			drawScreen(playerX, obstacles)

			frameCount++
		}
	}

	// Game Over screen
	clearScreen()
	fmt.Println("Game Over!")
	fmt.Println("Press any key to exit...")
	keyboard.GetKey() // wait for a key press before exit.
}
