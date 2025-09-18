package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func enterAltScreen() { fmt.Print("\033[?1049h\033[H") }
func exitAltScreen()  { fmt.Print("\033[?1049l") }
func hideCursor()     { fmt.Print("\033[?25l") }
func showCursor()     { fmt.Print("\033[?25h") }

func clearHome()      { fmt.Print("\033[H") }
func clearScreenAll() { fmt.Print("\033[H\033[J") }

// selectWithArrows renders a header and options, navigable via ↑/↓, Enter validates, ESC cancels
func selectWithArrows(header string, options []string) (int, bool) {
	if err := keyboard.Open(); err != nil {
		// Fallback: always pick first option
		return 0, false
	}
	defer keyboard.Close()

	index := 0
	for {
		clearHome()
		clearScreenAll()
		fmt.Println()
		if header != "" {
			fmt.Println(header)
			fmt.Println()
		}
		for i, opt := range options {
			prefix := "  "
			if i == index {
				prefix = "> "
			}
			fmt.Printf("%s%s\n", prefix, opt)
		}
		fmt.Println()
		fmt.Println("                    Contrôles: ↑/↓ naviguer | Entrée valider | Q sortir/retour")
		fmt.Println()
		ch, key, err := keyboard.GetKey()
		if err != nil {
			return index, false
		}
		switch key {
		case keyboard.KeyArrowUp:
			if index > 0 {
				index--
			} else {
				index = len(options) - 1
			}
		case keyboard.KeyArrowDown:
			if index < len(options)-1 {
				index++
			} else {
				index = 0
			}
		case keyboard.KeyEnter:
			return index, false
		case keyboard.KeyEsc:
			return index, true
		}
		if ch == 'q' || ch == 'Q' {
			return index, true
		}
	}
}
