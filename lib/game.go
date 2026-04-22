package gameloop

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

const FPS = 30

type Loop interface {
	Initialize() error
	Render() error
	Calculate(c keyboard.Key) error
}

func Start(gl Loop) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	if err := gl.Initialize(); err != nil {
		panic(err)
	}

	limiter := time.Tick((1000 * time.Millisecond) / FPS)
	for {
		<-limiter
		screen.Clear()
		fmt.Print("\033[H")
		go gl.Render()

		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		go gl.Calculate(key)
	}
}
