package main

import (
	"fmt"
	"runtime"
	"sloteriaa/internal/personnage"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("Note: Si les couleurs s'affichent mal, essayez dans Windows Terminal/PowerShell récent.")
	}

	fmt.Println("--- Création d'un personnage par l'utilisateur ---")
	p := personnage.CreationPersonnage()
	personnage.AfficherInfos(p)
}
