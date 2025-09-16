package main

import (
	"fmt"
	"runtime"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Note: Si les couleurs s'affichent mal, essayez dans Windows Terminal/PowerShell récent.")
	}

	RunMenu()
}
