package main

import (
	"math/rand"
	"time"

	"github.com/rthornton128/goncurses"
)

type Obstacle struct {
	x, y int
}

func main() {
	// Initialize ncurses
	stdscr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	defer goncurses.End()

	// Set terminal modes
	goncurses.Echo(false)
	goncurses.CBreak(true)
	stdscr.Keypad(true)
	stdscr.Timeout(0) // non-blocking input

	// Get screen dimensions
	maxY, maxX := stdscr.MaxYX()

	// Player settings (placed at the bottom center)
	playerX := maxX / 2
	playerY := maxY - 1

	var obstacles []Obstacle
	rand.Seed(time.Now().UnixNano())

	// Game timing
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	frameCount := 0
	gameOver := false

	// Main game loop
	for !gameOver {
		// Non-blocking input
		ch := stdscr.GetChar()
		if ch != 0 {
			switch ch {
			case 'a', goncurses.KEY_LEFT:
				if playerX > 0 {
					playerX--
				}
			case 'd', goncurses.KEY_RIGHT:
				if playerX < maxX-1 {
					playerX++
				}
			case 'q': // quit game
				gameOver = true
			}
		}

		// Every 10 frames, add a new obstacle at a random x position at the top.
		if frameCount%10 == 0 {
			newObstacle := Obstacle{
				x: rand.Intn(maxX),
				y: 0,
			}
			obstacles = append(obstacles, newObstacle)
		}

		// Update obstacles: move down and check for collision with the player.
		newObstacles := []Obstacle{}
		for _, obs := range obstacles {
			obs.y++
			if obs.y == playerY && obs.x == playerX {
				gameOver = true
			}
			if obs.y < maxY {
				newObstacles = append(newObstacles, obs)
			}
		}
		obstacles = newObstacles

		// Draw the frame
		stdscr.Clear()
		for _, obs := range obstacles {
			stdscr.MoveAddChar(obs.y, obs.x, 'X')
		}
		stdscr.MoveAddChar(playerY, playerX, 'A')
		stdscr.Refresh()

		// Wait for the next tick
		<-ticker.C
		frameCount++
	}

	// Display Game Over message
	stdscr.Clear()
	msg := "Game Over! Press any key..."
	stdscr.MovePrint(maxY/2, (maxX-len(msg))/2, msg)
	stdscr.Refresh()
	stdscr.GetChar() // Wait for a key press before exiting.
}
