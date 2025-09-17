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

func RunMenu() {
	// Création d'un joueur test
	joueur := Joueur{
		Nom:        "Héros",
		Inventaire: []string{"Potion"},
		Or:         50,
		HP:         80,
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
			startGame(joueur)
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

func startGame(joueur *Joueur) {
	for {
		afficherTitre()
		fmt.Printf("Bienvenue, %s ! HP: %d/%d - Or: %d\n", joueur.Nom, joueur.HP, joueur.HPMax, joueur.Or)
		fmt.Println("\n1. Afficher l'inventaire")
		fmt.Println("2. Utiliser une potion")
		fmt.Println("3. Retour au menu principal")
		fmt.Print("\nVotre choix : ")

		var choix int
		fmt.Scanln(&choix)

		clearMenuBody()

		switch choix {
		case 1:
			afficherInventaire(joueur)
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
			clearMenuBody()
		case 2:
			utiliserPotion(joueur)
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			fmt.Scanln()
			clearMenuBody()
		case 3:
			return
		default:
			fmt.Println("Choix invalide, appuyez sur Entrée pour réessayer...")
			fmt.Scanln()
			clearMenuBody()
		}
	}
}

// afficherInventaire et utiliserPotion sont définies dans inventaire.go
