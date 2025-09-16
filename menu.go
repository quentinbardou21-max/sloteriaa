package main

import (
	"fmt"
	"os"
	"time"
)

type Joueur struct {
	Nom        string
	Inventaire []string
	Or         int
	HP         int
	HPMax      int
}

func main() {
	// Création d'un joueur test
	joueur := Joueur{
		Nom:        "Héros",
		Inventaire: []string{},
		Or:         50,
		HP:         100,
		HPMax:      100,
	}

	afficherMenu(&joueur)
}

func afficherMenu(joueur *Joueur) {
	for {
		afficherTitre()
		fmt.Println("1. Entrer dans le jeu")
		fmt.Println("2. Quitter")
		fmt.Print("\nVotre choix : ")

		var choix int
		fmt.Scanln(&choix)

		clearMenuBody()

		switch choix {
		case 1:
			startGame()
		case 2:
			clearScreen()
			fmt.Println("Au revoir !")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		default:
			fmt.Println("Choix invalide, appuyez sur Entrée pour réessayer...")
			fmt.Scanln()
			clearMenuBody()
		}
	}
}

func afficherTitre() {
	clearScreen()
	fmt.Println(`
███████ ██       ██████  ████████ ███████ ██████  ██  █████  
██      ██      ██    ██    ██    ██      ██   ██ ██ ██   ██ 
███████ ██      ██    ██    ██    █████   ██████  ██ ███████ 
     ██ ██      ██    ██    ██    ██      ██   ██ ██ ██   ██ 
███████ ███████  ██████     ██    ███████ ██   ██ ██ ██   ██ 
`)
}

func clearMenuBody() {
	for i := 0; i < 20; i++ { // adapte selon la taille du menu/jeu
		fmt.Print("\033[A\033[2K")
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func startGame() {
	afficherTitre()
	fmt.Println("🎮 Le jeu démarre...")
	fmt.Println("\nAppuyez sur Entrée pour revenir au menu.")
	fmt.Scanln()
	clearMenuBody()
}
