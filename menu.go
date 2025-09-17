package main

import (
	"fmt"
	"os"
	"time"

	"sloteriaa/internal/personnage"

	"runtime"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows"
)

func RunMenu() {
	enableANSIWindows()
	// Création d'un joueur test
	joueur := personnage.Personnage{
		Nom:        "Héros",
		Inventaire: []string{"Potion"},
		Argent:     50,
		PVActuels:  80,
		PVMax:      100,
	}

	afficherMenu(&joueur)
}

func afficherMenu(joueur *personnage.Personnage) {
	for {
		afficherTitre()
		options := []string{"Entrer dans le jeu", "Quitter"}
		selection := afficherMenuAvecFleches(options)
		clearMenuBody()

		switch selection {
		case 0:
			startGame(joueur)
		case 1:
			clearScreen()
			fmt.Println("Au revoir !")
			time.Sleep(1 * time.Second)
			os.Exit(0)
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

func startGame(joueur *personnage.Personnage) {
	for {
		afficherTitre()
		fmt.Printf("Bienvenue, %s ! PV: %d/%d - Argent: %d\n", joueur.Nom, joueur.PVActuels, joueur.PVMax, joueur.Argent)
		options := []string{"Afficher l'inventaire", "Utiliser une potion", "Retour au menu principal"}
		selection := afficherMenuAvecFleches(options)
		clearMenuBody()

		switch selection {
		case 0:
			afficherInventaire(joueur)
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			attendreEntree()
			clearMenuBody()
		case 1:
			utiliserPotion(joueur)
			fmt.Println("\nAppuyez sur Entrée pour continuer...")
			attendreEntree()
			clearMenuBody()
		case 2:
			return
		}
	}
}

<<<<<<< HEAD
// afficherMenuAvecFleches affiche une liste d'options navigable avec ↑/↓ et valide avec Entrée.
func afficherMenuAvecFleches(options []string) int {
	// Initialiser clavier (mode raw)
	if err := keyboard.Open(); err != nil {
		// fallback simple si impossible d'initialiser: retourner première option
		return 0
	}
	defer keyboard.Close()

	index := 0
	for {
		// Affiche les options avec un curseur
		for i, opt := range options {
			prefix := "  "
			if i == index {
				prefix = "> "
			}
			fmt.Printf("%s%s\n", prefix, opt)
		}
		// Lire touche
		char, key, err := keyboard.GetKey()
		if err != nil {
			return index
		}

		// Efface le bloc affiché
		for range options {
			fmt.Print("\033[A\033[2K")
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
			return index
		case keyboard.KeyEsc:
			return len(options) - 1 // ESC: retourne sur la dernière option (souvent "Retour")
		default:
			// permet aussi Enter via '\r' si nécessaire
			if char == '\r' || char == '\n' {
				return index
			}
		}
	}
}

func attendreEntree() {
	if err := keyboard.Open(); err != nil {
		// fallback: rien
		return
	}
	defer keyboard.Close()
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			return
		}
		if key == keyboard.KeyEnter {
			return
		}
	}
}

// enableANSIWindows active le support des séquences ANSI dans la console Windows
func enableANSIWindows() {
	if runtime.GOOS != "windows" {
		return
	}
	h := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return
	}
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(h, mode)
}
=======
// afficherInventaire et utiliserPotion sont définies dans inventaire.go
>>>>>>> 55990879bea5f86c431733b1adcdbc1699e30a0d
