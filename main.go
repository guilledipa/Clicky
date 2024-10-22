package main

import (
	"embed"
	"math/rand/v2"

	"github.com/eihigh/miniten"
)

//go:embed *.png
var fsys embed.FS

var (
	// Starting conditions
	scene = "title"
	score int
	x     = 200.0
	y     = 150.0
	vy    = 0.0 // Velocity of y
	g     = 0.1 // Gravity
	jump  = -4.0
	// Wall
	frames     int
	interval   = 144
	wallStartX = 640
	walls      = []wall{}
	wallWidth  = 20
	wallHeight = 360
	holeYMax   = 150
	holeHeight = 170
	// Gopher
	gopherWidth  = 60
	gopherHeight = 75

	isPrevClicked bool
	isJustClicked bool
)

type wall struct {
	wallX int
	holeY int
}

func main() {
	miniten.Run(draw)
}

func draw() {
	isJustClicked = miniten.IsClicked() && !isPrevClicked
	isPrevClicked = miniten.IsClicked()
	switch scene {
	case "title":
		drawTitle()
	case "game":
		drawGame()
	case "gameover":
		drawGameover()
	}
}

func drawTitle() {
	miniten.DrawImageFS(fsys, "bg.png", 0, 0)
	miniten.Println("Gopher Clicky game!")
	miniten.DrawImageFS(fsys, "gopher.png", int(x), int(y))
	if isJustClicked {
		scene = "game"
	}
}

func drawGame() {
	miniten.DrawImageFS(fsys, "bg.png", 0, 0)
	for i, wall := range walls {
		if wall.wallX < int(x) {
			score = i + 1
		}
	}
	miniten.Println("Score", score)
	if miniten.IsClicked() {
		vy = jump
	}
	vy += g
	y += vy
	miniten.DrawImageFS(fsys, "gopher.png", int(x), int(y))

	frames++
	if frames%interval == 0 {
		wall := wall{wallStartX, rand.N(holeYMax)}
		walls = append(walls, wall)
	}
	for i := range walls {
		walls[i].wallX -= 2
	}
	for _, wall := range walls {
		drawWalls(wall)
		// gopher
		aLeft := int(x)
		aTop := int(y)
		aRight := int(x) + gopherWidth
		aBottom := int(y) + gopherHeight
		// Colisiones
		// Top
		bLeft := wall.wallX
		bTop := wall.holeY - wallHeight
		bRight := wall.wallX + wallWidth
		bBottom := wall.holeY
		if hitTestRects(aLeft, aTop, aRight, aBottom, bLeft, bTop, bRight, bBottom) {
			scene = "gameover"
		}
		// Bot
		bLeft = wall.wallX
		bTop = wall.holeY + holeHeight
		bRight = wall.wallX + wallWidth
		bBottom = wall.holeY + holeHeight + wallHeight
		if hitTestRects(aLeft, aTop, aRight, aBottom, bLeft, bTop, bRight, bBottom) {
			scene = "gameover"
		}
		// Ceiling
		if y < 0 {
			scene = "gameover"
		}
		// Floor
		if 360 < y { // +float64(gopherHeight) .. Too brutal :P
			scene = "gameover"
		}
	}
}

func drawGameover() {
	miniten.DrawImageFS(fsys, "bg.png", 0, 0)
	miniten.DrawImageFS(fsys, "gopher.png", int(x), int(y))
	for _, wall := range walls {
		drawWalls(wall)
	}

	miniten.Println("Game Over")
	miniten.Println("Score", score)
	if isJustClicked {
		scene = "title"
		// Reset initial game conditions
		x = 200.0
		y = 150.0
		vy = 0.0
		frames = 0
		walls = []wall{}
		score = 0
	}
}

func drawWalls(wall wall) {
	miniten.DrawImageFS(fsys, "wall.png", wall.wallX, wall.holeY-wallHeight)
	miniten.DrawImageFS(fsys, "wall.png", wall.wallX, wall.holeY+holeHeight)
}

func hitTestRects(aLeft, aTop, aRight, aBottom, bLeft, bTop, bRight, bBottom int) bool {
	return aLeft < bRight &&
		bLeft < aRight &&
		aTop < bBottom &&
		bTop < aBottom
}
