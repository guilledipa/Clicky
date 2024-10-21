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
	frames     = 0
	interval   = 144
	wallStartX = 640
	wallXs     = []int{}
	wallWidth  = 20
	wallHeight = 360
	holeYs     = []int{}
	holeYMax   = 150
	holeHeight = 170
	// Gopher
	gopherWidth  = 60
	gopherHeight = 75

	isPrevClicked = false
	isJustClicked = false
)

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
	miniten.Println("クリックしてスタート")
	miniten.DrawImageFS(fsys, "gopher.png", int(x), int(y))
	if isJustClicked {
		scene = "game"
	}
}

func drawGame() {
	miniten.DrawImageFS(fsys, "bg.png", 0, 0)
	for i, wallX := range wallXs {
		if wallX < int(x) {
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
		wallXs = append(wallXs, wallStartX)
		holeYs = append(holeYs, rand.N(holeYMax))
	}
	for i := range wallXs {
		wallXs[i] -= 2
	}
	for i := range wallXs {
		wallX := wallXs[i]
		holeY := holeYs[i]
		miniten.DrawImageFS(fsys, "wall.png", wallX, holeY-wallHeight)
		miniten.DrawImageFS(fsys, "wall.png", wallX, holeY+holeHeight)

		// gopher
		aLeft := int(x)
		aTop := int(y)
		aRight := int(x) + gopherWidth
		aBottom := int(y) + gopherHeight
		// Colisiones
		bLeft := wallX
		bTop := holeY - wallHeight
		bRight := wallX + wallWidth
		bBottom := holeY
		if aLeft < bRight &&
			bLeft < aRight &&
			aTop < bBottom &&
			bTop < aBottom {
			scene = "gameover"
		}
		bLeft = wallX
		bTop = holeY + holeHeight
		bRight = wallX + wallWidth
		bBottom = holeY + holeHeight + wallHeight
		if aLeft < bRight &&
			bLeft < aRight &&
			aTop < bBottom &&
			bTop < aBottom {
			scene = "gameover"
		}
		// Ceiling
		if y < 0 {
			scene = "gameover"
		}
		if 360 < y { // +float64(gopherHeight) .. Too brutal :P
			scene = "gameover"
		}
	}
}

func drawGameover() {
	miniten.DrawImageFS(fsys, "bg.png", 0, 0)
	miniten.DrawImageFS(fsys, "gopher.png", int(x), int(y))
	for i := range wallXs {
		wallX := wallXs[i]
		holeY := holeYs[i]
		miniten.DrawImageFS(fsys, "wall.png", wallX, holeY-wallHeight)
		miniten.DrawImageFS(fsys, "wall.png", wallX, holeY+holeHeight)
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
		wallXs = []int{}
		holeYs = []int{}
		score = 0
	}
}
