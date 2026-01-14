package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"strconv"

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
	paddleHeight = 30.0
	speed        = 4.0
	ballWidth    = 4.0
	maxSpeed     = 5.0
	deadZone     = .3
)

var p1 Paddle
var p2 Paddle
var b Ball
var player1 Player
var player2 Player

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

type Player struct {
	score int
}

// Sets the initial values for the player and ball entities.
// In a game with more than three entities you'd make a constructor function that would help you drop in all the defaults, etc.
func reset() {
	// Let's center the Y value, which is half the screen height - half the paddle height!
	p1 = Paddle{5.0, sH/2 - paddleHeight/2, paddleWidth, paddleHeight}
	p2 = Paddle{sW - 5.0 - paddleWidth, sH/2 - paddleHeight/2, paddleWidth, paddleHeight}
	b = Ball{sW/2 - ballWidth/2, sH/2 - ballWidth/2, ballWidth, 0, 0, 1, false}
}

func resetPlayers() {
	player1 = Player{0}
	player2 = Player{0}
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
	//set initial velocity x and velocity y
	b.vx = coinFlip()                      // uses a helper function to either start the ball going left or right, -1 or 1
	b.vy = rand.Float32() * 3 * coinFlip() // how diagonal will it be? The '*3' adds eccentricity to the angle which makes for a more dynamic game
}

// a classic Axis-Aligned Bounding-Box Collission check that check if 2 rectangles intersect.
func aabb(ax, ay, aw, ah, bx, by, bw, bh float32) bool {
	return ax < bx+bw &&
		ax+aw > bx &&
		ay < by+bh &&
		ay+ah > by
}

func (b *Ball) updateBall() {
	//check for top || bottom screen colissions to cause a bounce. Note that we invert the velocity of the y to change the direction.
	if b.y <= 0 || b.y >= sH-b.width {
		b.vy = -b.vy
	}
	//check left -> if off screen player 2 scores a point and reset the ball and paddles to centered positions
	if b.x <= 0 {
		player2.score += 1
		reset()
		fmt.Println("Player 1: ", player1.score, "Player 2: ", player2.score)
	}
	//check right
	if b.x >= sW {
		player1.score += 1
		reset()
		fmt.Println("Player 1: ", player1.score, "Player 2: ", player2.score)
	}
	//check for paddle collision. Ball shouldn't know about Paddle, but this is a small game, so no point to abstract.
	//Paddle 1
	if aabb(b.x, b.y, b.width, b.width, p1.x, p1.y, p1.width, p1.height) && b.vx < 0 {
		b.vx = -b.vx
		b.vy += rand.Float32() / 3 * coinFlip()
		b.v = min(maxSpeed, b.v+.5)
	}

	//Paddle 2
	if aabb(b.x, b.y, b.width, b.width, p2.x, p2.y, p2.width, p2.height) && b.vx > 0 {
		b.vx = -b.vx
		b.vy += rand.Float32() / 3 * coinFlip()
		b.v = min(maxSpeed, b.v+.5)
	}

	//at last, move the ball!
	b.x = b.x + (b.v * b.vx)
	b.y = b.y + (b.v * b.vy)
}

func handleInput() {
	// PLAYER CONTROLS

	// Player 1 up
	// Note the 'deadZone' constant is to eliminate jitter (i,e, subtly pressed, false positves, etc on the directional pad/control)
	if ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.StandardGamepadAxisValue(0, ebiten.StandardGamepadAxisLeftStickVertical) < -deadZone {
		p1.y = max(p1.y-speed, 0) //clamps to top of screen
	}
	// Player 1 down
	if ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.StandardGamepadAxisValue(0, ebiten.StandardGamepadAxisLeftStickVertical) > deadZone {
		p1.y = min(p1.y+speed, sH-p1.height) // clamps to bottom of screen, taking the paddleHeight into account since the x,y is the TOP/left corner.
	}

	// Player 2 up
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.StandardGamepadAxisValue(1, ebiten.StandardGamepadAxisLeftStickVertical) < -deadZone {
		p2.y = max(p2.y-speed, 0) //clamps to top of screen
	}

	// Player 2 down
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) ||
		ebiten.StandardGamepadAxisValue(1, ebiten.StandardGamepadAxisLeftStickVertical) > deadZone {
		p2.y = min(p2.y+speed, sH-p2.height) // clamps to bottom of screen, taking the paddleHeight into account since the x,y is the TOP/left corner.
	}

	// Ball Start with Space, Enter, Mouse, or GAMEPAD: Y,B,A
	if b.isInPlay == false && (ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyEnter) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButton0) ||
		ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) ||
		ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) ||
		ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton2) ||
		ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton0) ||
		ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton1) ||
		ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton2)) {
		b.serveBall()
	}

	// GAME OPTIONS
	// Fullscreen on/off -- easy!
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// RESET On gamepad controller button 8 is the Select Button
	if inpututil.IsKeyJustPressed(ebiten.KeyR) ||
		ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton8) ||
		ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton8) {
		reset()
		resetPlayers()
	}
}

func (g *Game) Update() error {
	handleInput()
	b.updateBall()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//Show player score
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(player1.score), 40, 10)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(player2.score), sW-40, 10)

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
	resetPlayers()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
