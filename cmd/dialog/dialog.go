package dialog

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func Confirm(question string) bool {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	fmt.Println(question)

	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyCtrlC {
			return false
		}

		switch char {
		case 'y', 'Y':
			return true
		case 'n', 'N':
			return false
		}
	}
}
