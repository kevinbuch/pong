package pong

import (
	"log"
	"os"
	"time"
)

type Pong struct {
	Player1 Player
	Player2 Player
	Ball    Ball
	Ui      Ui
}

type Window interface {
	ColorOn(color int)
	ColorOff(color int)
	AttributeOn(attr int)
	AttributeOff(attr int)
	Print(y int, x int, text string)
	Move(y int, x int)
	MaxRowAndColumn() (y int, x int)
	CurrentRowAndColumn() (y int, x int)
	Delete()
}

type Ui interface {
	GetChar() rune
	Erase()
	Refresh()
	Draw(Window)
	NewWindow(height int, width int) Window
	MaxRowAndColumn() (y int, x int)
}

func listen(ui Ui, player1 Player, player2 Player) {
	for {
		switch ui.GetChar() {
		case player1.UpKey:
			player1.Paddle.Up()
		case player1.DownKey:
			player1.Paddle.Down()
		case player2.UpKey:
			player2.Paddle.Up()
		case player2.DownKey:
			player2.Paddle.Down()
		case 0:
			player1.Paddle.Stopped()
			player2.Paddle.Stopped()
		}
	}
}

func (game Pong) Start() {
	go game.Player1.Paddle.Update()
	go game.Player2.Paddle.Update()
	go game.Ball.Update(&game.Player1.Paddle, &game.Player2.Paddle)
	go listen(game.Ui, game.Player1, game.Player2)

	ticks := time.NewTicker(time.Second / 16)

	for range ticks.C {
		game.Ui.Erase()
		game.Ui.Draw(game.Ball.Draw())
		game.Ui.Draw(game.Player1.Paddle.Draw())
		game.Ui.Draw(game.Player2.Paddle.Draw())
		game.Ui.Refresh()
	}
}

func NewGame() Pong {
	out, _ := os.Create("pong.log")
	log.SetOutput(out)

	ui := NewUi()

	paddle1 := NewPaddle(Left, ui)
	player1 := NewPlayer('f', 'd', paddle1)

	paddle2 := NewPaddle(Right, ui)
	player2 := NewPlayer('k', 'j', paddle2)

	ball := NewBall(ui)

	return Pong{
		Player1: player1,
		Player2: player2,
		Ball:    ball,
		Ui:      ui,
	}
}
