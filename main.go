package main

import (
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// we'll just keep our constant variables for game play up here for clarity.
// Best practice would be to add to the Game struct and have a NewGame function return a new game before
// starting the game loop

const (
	sW           = 320
	sH           = 240
	paddleWidth  = 5.0
	paddleHeight = 50.0
	speed        = 4.0
	ballWidth    = 4.0
)

var p1 Paddle
var p2 Paddle
var b Ball

type Game struct{}

type Paddle struct {
	x      float32
	y      float32
	width  float32
	height float32
}
type Ball struct {
	x        float32
	y        float32
	width    float32 // It is a square, width = length
	vx       float32 // velocity x
	vy       float32 // velocity y
	v        float32 // velocity
	isInPlay bool    // false is the ready to serve the ball state. true is the ball is in play
}

// Sets the initial values for the player and ball entities.
// In a game with more than three entities you'd make a constructor function that would help you drop in all the defaults, etc.
func reset() {
	// Let's center the Y value, which is half the screen height - half the paddle height!
	p1 = Paddle{5.0, sH/2 - paddleHeight/2, paddleWidth, paddleHeight}
	p2 = Paddle{sW - 5.0 - paddleWidth, sH/2 - paddleHeight/2, paddleWidth, paddleHeight}
	b = Ball{sW/2 - ballWidth/2, sH/2 - ballWidth/2, ballWidth, 0, 0, 1, false}
}

func (p *Paddle) drawPaddle(screen *ebiten.Image) {
	// vector.FillRect(screen, 5, 20, 10, 50, color.White, false) // was this previously
	vector.FillRect(screen, p.x, p.y, p.width, p.height, color.White, false)
}

func (b *Ball) drawBall(screen *ebiten.Image) {
	vector.FillRect(screen, b.x, b.y, b.width, b.width, color.White, false)
}

// Emulates a coinflip and used to switch between + and - numbers in practice
func coinFlip() float32 {
	if rand.Float64() > .5 {
		return 1.0
	}
	return -1.0
}

// gives initial direction on serve
func (b *Ball) serveBall() {
	b.isInPlay = true //be explicit, rather than !b.isInPlay
	//set initial dx to dy
	b.vx = coinFlip()                      // uses a helper function to either start the ball going left or right, -1 or 1
	b.vy = rand.Float32() * 3 * coinFlip() // how diagonal will it be?
}

func (b *Ball) updateBall() {
	b.x = b.x + (b.v * b.vx)
}

func handleInput() {
	// PLAYER CONTROLS

	//Player 1 up
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p1.y = max(p1.y-speed, 0) //clamps to top of screen
	}
	//Player 1 down
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p1.y = min(p1.y+speed, sH-p1.height) // clamps to bottom of screen, taking the paddleHeight into account since the x,y is the TOP/left corner.
	}

	//Player 2 up
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p2.y = max(p2.y-speed, 0) //clamps to top of screen
	}
	//Player 2 down
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p2.y = min(p2.y+speed, sH-p2.height) // clamps to bottom of screen, taking the paddleHeight into account since the x,y is the TOP/left corner.
	}

	//Ball Start
	if b.isInPlay == false && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsMouseButtonPressed(ebiten.MouseButton0)) {
		b.serveBall()
	}

	// GAME OPTIONS
	// Fullscreen on/off -- easy!
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
}

func (g *Game) Update() error {
	handleInput()
	b.updateBall()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Pong!")
	p1.drawPaddle(screen)
	p2.drawPaddle(screen)
	b.drawBall(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return sW, sH
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GolangNYC Pong!")
	reset()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
