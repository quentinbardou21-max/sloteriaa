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
	hideCursor()       // Cache le curseur au début du jeu
	defer showCursor() // Réaffiche le curseur à la fin du jeu

	// Création d'un joueur test
	joueur := personnage.Personnage{
		Nom:        "Héros",
		Inventaire: []string{"Épée rouillée"},
		Argent:     100,
		PVActuels:  80,
		PVMax:      100,
		Attaque:    "Épée",
	}

	afficherMenu(&joueur)
}

func afficherMenu(joueur *personnage.Personnage) {
	for {
		afficherTitre()
		options := []string{"Continuer", "Nouvelle partie", "Quitter"}
		selection := afficherMenuAvecFleches(options)
		clearMenuBody()

		switch selection {
		case 0: // Continuer
			if gs, err := LoadGame(); err == nil {
				StartGameFromSave(gs)
			} else {
				fmt.Println("Aucune sauvegarde trouvée.")
				time.Sleep(1 * time.Second)
			}
		case 1: // Nouvelle partie
			_ = DeleteSave() // ignore l'erreur si pas de sauvegarde
			StartGameNew()
		case 2: // Quitter
			clearScreen()
			showCursor() // Réaffiche le curseur avant de quitter
			fmt.Println("Au revoir !")
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}
}

func afficherTitre() {
	clearScreen()
	fmt.Print(`
███████ ██       ██████  ████████ ███████ ██████  ██  █████  
██      ██      ██    ██    ██    ██      ██   ██ ██ ██   ██ 
███████ ██      ██    ██    ██    █████   ██████  ██ ███████ 
     ██ ██      ██    ██    ██    ██      ██   ██ ██ ██   ██ 
███████ ███████  ██████     ██    ███████ ██   ██ ██ ██   ██ 
`)
	fmt.Println("Appuyez sur ↑/↓ pour choisir, Entrée pour valider, Q pour quitter.")
}

func clearMenuBody() {
	for i := 0; i < 15; i++ { // adapte selon la taille du menu/jeu
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
		options := []string{"Inventaire (↑/↓ + Entrée)", "Retour au menu principal"}
		selection := afficherMenuAvecFleches(options)
		clearMenuBody()

		switch selection {
		case 0:
			afficherInventaireInteractif(joueur)
		case 1:
			return
		}
	}
}

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
		ch, key, err := keyboard.GetKey()
		if err != nil {
			return
		}
		if key == keyboard.KeyEnter || key == keyboard.KeyEsc || ch == 'q' || ch == 'Q' {
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
