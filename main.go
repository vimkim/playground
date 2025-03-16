package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	tickRate            = 200 * time.Millisecond
	newBlockProbability = 0.1
)

type FallingBlock struct {
	x, y int
}

type Game struct {
	playerX, playerY int
	blocks           []FallingBlock
	score            int
	width, height    int
}

func NewGame(width, height int) *Game {
	return &Game{
		playerX: width / 2,
		playerY: height - 1,
		blocks:  []FallingBlock{},
		score:   0,
		width:   width,
		height:  height,
	}
}

// Update the game state: spawn new blocks, move blocks down, and update score.
func (g *Game) Update() {
	// For each column in the playable area, randomly spawn a new block.
	for x := 0; x < g.width; x++ {
		if rand.Float64() < newBlockProbability {
			g.blocks = append(g.blocks, FallingBlock{x: x, y: 0})
		}
	}

	// Move blocks down and remove those that have moved off-screen.
	newBlocks := []FallingBlock{}
	for _, b := range g.blocks {
		b.y++
		if b.y < g.height {
			newBlocks = append(newBlocks, b)
		}
	}
	g.blocks = newBlocks

	// Increase score as time passes.
	g.score++
}

// CheckCollision returns true if any block occupies the same position as the player.
func (g *Game) CheckCollision() bool {
	for _, b := range g.blocks {
		if b.x == g.playerX && b.y == g.playerY {
			return true
		}
	}
	return false
}

// drawBorder draws a simple border around the outer area.
func drawBorder(s tcell.Screen, width, height int) {
	// Top and bottom borders.
	for x := 0; x < width; x++ {
		s.SetContent(x, 0, '-', nil, tcell.StyleDefault)
		s.SetContent(x, height-1, '-', nil, tcell.StyleDefault)
	}
	// Left and right borders.
	for y := 0; y < height; y++ {
		s.SetContent(0, y, '|', nil, tcell.StyleDefault)
		s.SetContent(width-1, y, '|', nil, tcell.StyleDefault)
	}
	// Corners.
	s.SetContent(0, 0, '+', nil, tcell.StyleDefault)
	s.SetContent(width-1, 0, '+', nil, tcell.StyleDefault)
	s.SetContent(0, height-1, '+', nil, tcell.StyleDefault)
	s.SetContent(width-1, height-1, '+', nil, tcell.StyleDefault)
}

// drawGame renders the game state into the terminal.
func drawGame(s tcell.Screen, game *Game, outerWidth, outerHeight int) {
	s.Clear()
	drawBorder(s, outerWidth, outerHeight)

	// The playable area is inside the border (offset by 1).
	offsetX, offsetY := 1, 1

	// Draw each falling block.
	for _, b := range game.blocks {
		x := b.x + offsetX
		y := b.y + offsetY
		if x < outerWidth-1 && y < outerHeight-1 {
			s.SetContent(x, y, '#', nil, tcell.StyleDefault)
		}
	}

	// Draw the player with a highlighted style.
	playerX := game.playerX + offsetX
	playerY := game.playerY + offsetY
	playerStyle := tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack)
	if playerX < outerWidth-1 && playerY < outerHeight-1 {
		s.SetContent(playerX, playerY, '@', nil, playerStyle)
	}

	// Draw score on the top border.
	scoreStr := fmt.Sprintf(" Score: %d ", game.score)
	for i, r := range scoreStr {
		if 2+i < outerWidth-1 {
			s.SetContent(2+i, 0, r, nil, tcell.StyleDefault)
		}
	}

	s.Show()
}

// putCentered displays a centered message.
func putCentered(s tcell.Screen, msg string) {
	width, height := s.Size()
	startX := (width - len(msg)) / 2
	startY := height / 2
	for i, r := range msg {
		s.SetContent(startX+i, startY, r, nil, tcell.StyleDefault)
	}
	s.Show()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize the tcell screen.
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}
	defer s.Fini()

	// Clear the screen and determine its size.
	s.Clear()
	outerWidth, outerHeight := s.Size()
	// Define the playable area inside a 1-character border.
	playableWidth := outerWidth - 2
	playableHeight := outerHeight - 2

	game := NewGame(playableWidth, playableHeight)

	// Create a ticker for game updates.
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	// Create a channel to receive keyboard events.
	evChan := make(chan tcell.Event)
	go func() {
		for {
			evChan <- s.PollEvent()
		}
	}()

	gameOver := false
Loop:
	for {
		select {
		case ev := <-evChan:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyLeft:
					if game.playerX > 0 {
						game.playerX--
					}
				case tcell.KeyRight:
					if game.playerX < game.width-1 {
						game.playerX++
					}
				case tcell.KeyEscape:
					break Loop
				default:
					r := ev.Rune()
					if r == 'q' || r == 'Q' {
						break Loop
					}
				}
			}
		case <-ticker.C:
			game.Update()
			if game.CheckCollision() {
				gameOver = true
				break Loop
			}
			drawGame(s, game, outerWidth, outerHeight)
		}
	}

	// Game over: clear the screen and display a final score message.
	s.Clear()
	msg := fmt.Sprintf("Game Over! Final Score: %d. Press any key to exit.", game.score)
	putCentered(s, msg)
	// Wait for any key press before exiting.
	s.PollEvent()
}
