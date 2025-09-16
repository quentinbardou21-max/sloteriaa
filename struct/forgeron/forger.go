package forgeron

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-tty"

	"github.com/quent/sloteriaa/struct/objet"
)

// Définition des matériaux utilisables à la forge
type Materiau string

const (
	Fer            Materiau = "Lingot de fer"
	Bois           Materiau = "Planche de bois"
	Cuir           Materiau = "Cuir tanné"
	EssenceMagique Materiau = "Essence magique"
)

// Coût d'une recette par matériau
type Cout map[Materiau]int

// Recette d'artisanat d'une arme humaine (non-monstre)
type Recette struct {
	CleArme    string
	NomAffiche string
	Cout       Cout
}

// Inventaire des matériaux du joueur
type InventaireMateriaux map[Materiau]int

func (inv InventaireMateriaux) AAssez(c Cout) bool {
	for m, q := range c {
		if inv[m] < q {
			return false
		}
	}
	return true
}

func (inv InventaireMateriaux) Debiter(c Cout) {
	for m, q := range c {
		inv[m] -= q
		if inv[m] < 0 {
			inv[m] = 0
		}
	}
}

// Catalogue des recettes pour armes humaines (clés compatibles avec objet.CreerArme)
func RecettesArmesHumaines() []Recette {
	return []Recette{
		{CleArme: "EpeeRouillee", NomAffiche: "Épée rouillée", Cout: Cout{Fer: 1}},
		{CleArme: "EpeeCourte", NomAffiche: "Épée courte", Cout: Cout{Fer: 2, Cuir: 1}},
		{CleArme: "EpeeFer", NomAffiche: "Épée en fer", Cout: Cout{Fer: 4, Cuir: 1}},
		{CleArme: "Hache", NomAffiche: "Hache lourde", Cout: Cout{Fer: 5, Bois: 2}},
		{CleArme: "HacheDeCombat", NomAffiche: "Hache de combat", Cout: Cout{Fer: 4, Bois: 2, Cuir: 1}},
		{CleArme: "HacheDeBataille", NomAffiche: "Hache de bataille", Cout: Cout{Fer: 7, Bois: 3}},
		{CleArme: "ArcBois", NomAffiche: "Arc en bois", Cout: Cout{Bois: 4, Cuir: 1}},
		{CleArme: "ArcLong", NomAffiche: "Arc long", Cout: Cout{Bois: 6, Cuir: 2}},
		{CleArme: "ArcElfe", NomAffiche: "Arc elfique", Cout: Cout{Bois: 5, Cuir: 2, EssenceMagique: 1}},
		{CleArme: "EpeeMagique", NomAffiche: "Épée magique", Cout: Cout{Fer: 5, EssenceMagique: 2}},
	}
}

// Interface console: point d'entrée de la forge
func RunForge(inv InventaireMateriaux) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== Forge ===")
		afficherInventaire(inv)
		fmt.Print("Voulez-vous parler au forgeron ? (o/n): ")
		choix, _ := reader.ReadString('\n')
		choix = strings.TrimSpace(strings.ToLower(choix))
		fmt.Println()

		if choix == "n" || choix == "non" {
			fmt.Println("Vous quittez la forge.")
			return
		}
		if choix != "o" && choix != "oui" {
			fmt.Println("Réponse invalide. Réessayez.")
			continue
		}

		afficherCatalogue()
		fmt.Print("Entrez le numéro de l'arme à forger (ou vide pour annuler): ")
		entree, _ := reader.ReadString('\n')
		entree = strings.TrimSpace(entree)
		if entree == "" {
			fmt.Println("Retour au forgeron...")
			continue
		}
		idx, err := strconv.Atoi(entree)
		if err != nil {
			fmt.Println("Entrée invalide.")
			continue
		}

		recettes := RecettesArmesHumaines()
		if idx < 1 || idx > len(recettes) {
			fmt.Println("Numéro hors liste.")
			continue
		}

		recette := recettes[idx-1]
		if !inv.AAssez(recette.Cout) {
			fmt.Println("Vous n'avez pas assez de matériaux pour cette arme.")
			fmt.Println("Coût requis:")
			afficherCout(recette.Cout)
			fmt.Println()
			continue
		}

		inv.Debiter(recette.Cout)
		arme := objet.CreerArme(recette.CleArme)
		fmt.Println("Fabrication réussie ! Voici votre arme :")
		objet.AfficherArme(arme)
	}
}

func afficherInventaire(inv InventaireMateriaux) {
	if len(inv) == 0 {
		fmt.Println("Inventaire matériaux: (vide)")
		return
	}
	fmt.Println("Inventaire matériaux:")
	for m, q := range inv {
		fmt.Printf("  - %s x%d\n", m, q)
	}
	fmt.Println()
}

func afficherCatalogue() {
	recettes := RecettesArmesHumaines()
	fmt.Println("--- Catalogue d'armes humaines ---")
	for i, r := range recettes {
		fmt.Printf("%d) %s\n", i+1, r.NomAffiche)
		arme := objet.CreerArme(r.CleArme)
		fmt.Printf("   -> Attaque: %d | Poids: %d\n", arme.EffetAttaque, arme.Poids)
		fmt.Print("   Coût: ")
		afficherCoutInline(r.Cout)
		fmt.Println()
	}
	fmt.Println()
}

func afficherCout(c Cout) {
	for m, q := range c {
		fmt.Printf("  - %s x%d\n", m, q)
	}
}

func afficherCoutInline(c Cout) {
	items := make([]string, 0, len(c))
	for m, q := range c {
		items = append(items, fmt.Sprintf("%s x%d", m, q))
	}
	fmt.Print(strings.Join(items, ", "))
}

// -------- TUI améliorée (navigation clavier) --------

// RunForgeTUI propose une interface plus visuelle avec sélection et panneau de détails
func RunForgeTUI(inv InventaireMateriaux) {
	recettes := RecettesArmesHumaines()
	if len(recettes) == 0 {
		fmt.Println("Aucune recette disponible.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	selection := 0

	for {
		renderForgeTUI(recettes, inv, selection)
		fmt.Print("Commandes: [z] haut, [s] bas, [c] craft, [q] quitter > ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))

		switch line {
		case "z":
			if selection > 0 {
				selection--
			}
		case "s":
			if selection < len(recettes)-1 {
				selection++
			}
		case "c":
			choisie := recettes[selection]
			if !inv.AAssez(choisie.Cout) {
				fmt.Println("\nMatériaux insuffisants pour fabriquer cette arme. Appuyez Entrée...")
				reader.ReadString('\n')
				continue
			}
			inv.Debiter(choisie.Cout)
			arme := objet.CreerArme(choisie.CleArme)
			fmt.Println("\nFabrication réussie !")
			objet.AfficherArme(arme)
			fmt.Println("Appuyez Entrée pour revenir à la forge...")
			reader.ReadString('\n')
		case "q":
			fmt.Println("Fermeture de la forge.")
			return
		}
	}
}

func renderForgeTUI(recettes []Recette, inv InventaireMateriaux, selection int) {
	clearScreenTUI()
	fmt.Println("=== Forge (TUI) ===")
	afficherInventaire(inv)

	// Liste avec sélection
	fmt.Println("--- Armes humaines ---")
	for i, r := range recettes {
		prefix := "  "
		if i == selection {
			prefix = "> "
		}
		fmt.Printf("%s%d) %s\n", prefix, i+1, r.NomAffiche)
	}

	// Panneau de détails pour l'arme sélectionnée
	choisie := recettes[selection]
	arme := objet.CreerArme(choisie.CleArme)
	fmt.Println("\n--- Détails ---")
	fmt.Printf("Nom: %s\n", arme.Nom)
	fmt.Printf("Description: %s\n", arme.Description)
	fmt.Printf("Attaque: %d | Poids: %d\n", arme.EffetAttaque, arme.Poids)
	fmt.Print("Coût: ")
	afficherCoutInline(choisie.Cout)
	if inv.AAssez(choisie.Cout) {
		fmt.Println("  [Disponible]")
	} else {
		fmt.Println("  [Matériaux insuffisants]")
	}
	fmt.Println()
}

func clearScreenTUI() {
	// Efface l'écran et replace le curseur (ANSI)
	fmt.Print("\033[H\033[2J")
}

// -------- TUI temps réel avec flèches (sans Entrée) --------

// RunForgeInteractive lit les touches en temps réel (flèches, c, q) via go-tty
func RunForgeInteractive(inv InventaireMateriaux) {
	recettes := RecettesArmesHumaines()
	if len(recettes) == 0 {
		fmt.Println("Aucune recette disponible.")
		return
	}

	t, err := tty.Open()
	if err != nil {
		fmt.Println("Impossible d'initialiser le TTY:", err)
		return
	}
	defer t.Close()

	selection := 0
	for {
		renderForgeTUI(recettes, inv, selection)
		key := readKey(t)
		switch key {
		case "up":
			if selection > 0 {
				selection--
			}
		case "down":
			if selection < len(recettes)-1 {
				selection++
			}
		case "c":
			choisie := recettes[selection]
			if !inv.AAssez(choisie.Cout) {
				messagePause(t, "Matériaux insuffisants pour fabriquer cette arme. (tapez une touche)")
				continue
			}
			inv.Debiter(choisie.Cout)
			arme := objet.CreerArme(choisie.CleArme)
			clearScreenTUI()
			fmt.Println("Fabrication réussie ! Voici votre arme :")
			objet.AfficherArme(arme)
			messagePause(t, "Retour à la forge (tapez une touche)")
		case "q":
			clearScreenTUI()
			fmt.Println("Fermeture de la forge.")
			return
		}
	}
}

// readKey convertit les séquences de touches en identifiants simples
func readKey(t *tty.TTY) string {
	r, err := t.ReadRune()
	if err != nil {
		return ""
	}
	if r == 0x1b { // ESC
		// attente séquence CSI: ESC [ A/B/C/D
		r2, _ := t.ReadRune()
		r3, _ := t.ReadRune()
		if r2 == '[' {
			switch r3 {
			case 'A':
				return "up"
			case 'B':
				return "down"
			case 'C':
				return "right"
			case 'D':
				return "left"
			}
		}
		return ""
	}
	// lettres de commande
	s := strings.ToLower(string(r))
	if s == "c" || s == "q" {
		return s
	}
	return ""
}

func messagePause(t *tty.TTY, msg string) {
	fmt.Println(msg)
	_, _ = t.ReadRune()
}
