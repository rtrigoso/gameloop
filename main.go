package main

import (
	"fmt"
	"sync"

	"github.com/eiannone/keyboard"
	gameloop "renfoc.us/game_loop/lib"
)

var (
	direction string
	mu        sync.RWMutex
)

type Data struct{}

func (data Data) Initialize() error {
	mu.Lock()
	defer mu.Unlock()
	direction = "starting..."

	return nil
}

func (data Data) Render() error {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println(direction)

	return nil
}

func (data Data) Calculate(key keyboard.Key) error {
	mu.Lock()
	defer mu.Unlock()
	switch key {
	case keyboard.KeyArrowDown:
		direction = "down"
	case keyboard.KeyArrowUp:
		direction = "up"
	case keyboard.KeyArrowLeft:
		direction = "left"
	case keyboard.KeyArrowRight:
		direction = "right"
	case keyboard.KeySpace:
		direction = "space"
	}

	return nil
}

func main() {
	data := Data{}
	gameloop.Start(data)
}
