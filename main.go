package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Grid      [][]*Node
	nodeStart *Node
	nodeEnd   *Node
}

var game = &Game{}

const (
	MAP_HEIGHT  = 25
	MAP_WIDTH   = 32
	NODE_SIZE   = 12
	NODE_BORDER = 3
)

func init() {
	game.Grid = make([][]*Node, MAP_WIDTH)

	for x := 0; x < MAP_WIDTH; x++ {

		game.Grid[x] = make([]*Node, MAP_HEIGHT)

		for y := 0; y < MAP_HEIGHT; y++ {
			node := &Node{Obstacle: false, Visited: false, X: x, Y: y, Parent: nil}
			game.Grid[x][y] = node
		}
	}
	log.Println("grid built.")

	// add neighbours to each node
	for x := 0; x < MAP_WIDTH; x++ {
		for y := 0; y < MAP_HEIGHT; y++ {
			node := game.Grid[x][y]
			if y > 0 {
				node.Neighbours = append(node.Neighbours, game.Grid[x][y-1])
			}
			if y < MAP_HEIGHT-1 {
				node.Neighbours = append(node.Neighbours, game.Grid[x][y+1])
			}

			if x > 0 {
				node.Neighbours = append(node.Neighbours, game.Grid[x-1][y])
			}
			if x < MAP_WIDTH-1 {
				node.Neighbours = append(node.Neighbours, game.Grid[x+1][y])
			}

			//diagonal neighbours
			if x > 0 && y > 0 {
				node.Neighbours = append(node.Neighbours, game.Grid[x-1][y-1])
			}
			if y < MAP_HEIGHT-1 && x > 0 {
				node.Neighbours = append(node.Neighbours, game.Grid[x-1][y+1])
			}
			if y > 0 && x < MAP_WIDTH-1 {
				node.Neighbours = append(node.Neighbours, game.Grid[x+1][y-1])
			}
			if y < MAP_HEIGHT-1 && x < MAP_WIDTH-1 {
				node.Neighbours = append(node.Neighbours, game.Grid[x+1][y+1])
			}
		}
	}

	log.Println("neighbours added")

	// init start and end node
	game.nodeStart = game.Grid[5][5]
	game.nodeEnd = game.Grid[22][18]

	SolveAStar(game.nodeStart, game.nodeEnd)
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		tileX := x / NODE_SIZE
		tileY := y / NODE_SIZE

		if tileX >= 0 && tileX < MAP_WIDTH && tileY >= 0 && tileY < MAP_HEIGHT {
			// mouse click while holding the shift key positions the start node
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				g.nodeStart = g.Grid[tileX][tileY]
			} else if ebiten.IsKeyPressed(ebiten.KeyControlLeft) {
				// mouse click and ctrl-left positions the end node
				g.nodeEnd = g.Grid[tileX][tileY]
			} else {
				// just a click positions an obstacle
				g.Grid[tileX][tileY].Obstacle = !g.Grid[tileX][tileY].Obstacle
			}

		}
		if result := SolveAStar(g.nodeStart, g.nodeEnd); !result {
			log.Printf("aborted! there is no solution after %d tires", MAX_TRIES)
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for x := 0; x < MAP_WIDTH; x++ {
		for y := 0; y < MAP_HEIGHT; y++ {
			node := game.Grid[x][y]

			for _, n := range node.Neighbours {
				vector.StrokeLine(
					screen,
					float32(x*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
					float32(y*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
					float32(n.X*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
					float32(n.Y*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
					1.5,
					color.RGBA{50, 150, 240, 255},
					false,
				)
			}
		}
	}

	for x := 0; x < MAP_WIDTH; x++ {
		for y := 0; y < MAP_HEIGHT; y++ {
			node := game.Grid[x][y]

			if node == g.nodeStart {
				vector.DrawFilledRect(
					screen,
					float32(x)*NODE_SIZE+NODE_BORDER,
					float32(y)*NODE_SIZE+NODE_BORDER,
					NODE_SIZE-NODE_BORDER,
					NODE_SIZE-NODE_BORDER,
					color.RGBA{10, 200, 10, 255},
					false,
				)
			} else if node == g.nodeEnd {
				vector.DrawFilledRect(
					screen,
					float32(x)*NODE_SIZE+NODE_BORDER,
					float32(y)*NODE_SIZE+NODE_BORDER,
					NODE_SIZE-NODE_BORDER,
					NODE_SIZE-NODE_BORDER,
					color.RGBA{200, 10, 10, 255},
					false,
				)
			} else if node.Obstacle {
				vector.DrawFilledRect(
					screen,
					float32(x)*NODE_SIZE+NODE_BORDER,
					float32(y)*NODE_SIZE+NODE_BORDER,
					NODE_SIZE-NODE_BORDER,
					NODE_SIZE-NODE_BORDER,
					color.RGBA{150, 150, 200, 255},
					false,
				)
			} else {
				if node.Visited {
					vector.DrawFilledRect(
						screen,
						float32(x)*NODE_SIZE+NODE_BORDER,
						float32(y)*NODE_SIZE+NODE_BORDER,
						NODE_SIZE-NODE_BORDER,
						NODE_SIZE-NODE_BORDER,
						color.RGBA{50, 50, 50, 180},
						false,
					)
				} else {
					vector.DrawFilledRect(
						screen,
						float32(x)*NODE_SIZE+NODE_BORDER,
						float32(y)*NODE_SIZE+NODE_BORDER,
						NODE_SIZE-NODE_BORDER,
						NODE_SIZE-NODE_BORDER,
						color.RGBA{50, 50, 200, 255},
						false,
					)
				}
			}
		}
	}
	// draw the found path
	if g.nodeEnd != nil {
		p := g.nodeEnd

		for p.Parent != nil {
			vector.StrokeLine(
				screen,
				float32(p.X*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
				float32(p.Y*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
				float32(p.Parent.X*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
				float32(p.Parent.Y*NODE_SIZE+(NODE_SIZE+NODE_BORDER)/2),
				1.5,
				color.RGBA{250, 250, 10, 255},
				false,
			)
			p = p.Parent
		}
	}

	ebitenutil.DebugPrintAt(screen, "place obstacle: left mouse click", 1, 1)
	ebitenutil.DebugPrintAt(screen, "position start node: left mouse click + shift", 1, 11)
	ebitenutil.DebugPrintAt(screen, "position end node: left mouse click + ctrl", 1, 21)
}

// }} Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 300
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Pathfinding")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
