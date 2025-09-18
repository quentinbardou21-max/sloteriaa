package main

import (
	"fmt"
	"sloteriaa/internal/personnage"
	"sloteriaa/struct/objet"
	"strings"

	"github.com/eiannone/keyboard"
)

// Limite de poids totale autorisée dans l'inventaire
const PoidsMaxInventaire = 50

// poidsConnus mappe les objets connus à leur poids (par défaut 1 si inconnu)
var poidsConnus = map[string]int{
	"potion": 1,
}

// PoidsObjet retourne le poids d'un objet (insensible à la casse)
func PoidsObjet(nom string) int {
	if p, ok := poidsConnus[strings.ToLower(nom)]; ok {
		return p
	}
	return 1
}

// estMateriau vérifie si un objet est un matériau de craft
func estMateriau(nom string) bool {
	materiaux := []string{
		"cuir", "cuir renforcé", "fer", "fer renforcé",
		"pierre", "pierre précieuse", "gemme", "cristal",
		"bois", "bois dur", "écorce", "résine",
		"os", "os ancien", "dent", "griffe",
		"plume", "plume rare", "écailles", "écailles de dragon",
		"soie", "soie d'araignée", "fil", "corde",
		"poudre", "poudre magique", "essence", "élixir",
		"minerai", "minerai rare", "cristal de mana", "gemme de pouvoir",
	}

	nomLower := strings.ToLower(nom)
	for _, mat := range materiaux {
		if strings.Contains(nomLower, strings.ToLower(mat)) {
			return true
		}
	}
	return false
}

// PoidsTotal calcule le poids total actuel de l'inventaire
func PoidsTotal(j *personnage.Personnage) int {
	total := 0
	for _, objet := range j.Inventaire {
		total += PoidsObjet(objet)
	}
	return total
}

func afficherInventaire(j *personnage.Personnage) {
	fmt.Println("🧳 Inventaire :")
	fmt.Printf("Or: %d\n", j.Argent)
	if estInventaireVide(j) {
		fmt.Println("Votre inventaire est vide.")
		return
	}

	noms, counts := compterItems(j)

	// Séparer les matériaux et les équipements
	var materiaux []string
	var equipements []string

	for _, item := range noms {
		if estMateriau(item) {
			materiaux = append(materiaux, item)
		} else {
			equipements = append(equipements, item)
		}
	}

	// Afficher la section Matériaux
	if len(materiaux) > 0 {
		fmt.Println("\n📦 MATÉRIAUX :")
		for i, item := range materiaux {
			label := item
			// Marquer les objets droppés
			if strings.Contains(item, "[DROPPÉ]") {
				label = strings.Replace(label, "[DROPPÉ]", "🎁", 1)
			}
			if counts[item] > 1 {
				fmt.Printf("  %d. %s x%d\n", i+1, label, counts[item])
			} else {
				fmt.Printf("  %d. %s\n", i+1, label)
			}
		}
	}

	// Afficher la section Équipements
	if len(equipements) > 0 {
		fmt.Println("\n⚔️ ÉQUIPEMENTS :")
		for i, item := range equipements {
			suffix := ""
			if estArmeEquipee(j, item) || estArmureEquipee(j, item) {
				suffix = "  [Équipé]"
			}
			label := item

			// Afficher les descriptions des potions
			switch strings.ToLower(item) {
			case "potion":
				label += " (+20 PV)"
			case "potion majeure":
				label += " (+50 PV)"
			case "potion force":
				label += " (+2 Force, 3 combats)"
			case "potion agilite":
				label += " (+2 Agilité, 3 combats)"
			case "potion endurance":
				label += " (+2 Endurance, 3 combats)"
			case "antidote":
				label += " (Guérit statuts)"
			case "elixir vie":
				label += " (+100 PV)"
			}

			// Afficher les stats des armes
			if arme, ok := trouverArmeParNom(item); ok && arme.Nom != "" {
				label += fmt.Sprintf(" (ATK %d)", arme.EffetAttaque)
			}

			// Afficher les stats des armures
			if armure, ok := trouverArmureParNom(item); ok && armure.Nom != "" {
				label += fmt.Sprintf(" (DEF %d)", armure.EffetDefense)
			}

			// Marquer les objets droppés
			if strings.Contains(item, "[DROPPÉ]") {
				label = strings.Replace(label, "[DROPPÉ]", "🎁", 1)
			}
			if counts[item] > 1 {
				fmt.Printf("  %d. %s x%d%s\n", i+1, label, counts[item], suffix)
			} else {
				fmt.Printf("  %d. %s%s\n", i+1, label, suffix)
			}
		}
	}
}

// afficherInventaireInteractif permet de naviguer avec ↑/↓ et d'utiliser l'objet sélectionné avec Entrée
func afficherInventaireInteractif(j *personnage.Personnage) {
	// Si l'inventaire est vide, afficher et attendre une entrée
	if estInventaireVide(j) {
		afficherInventaire(j)
		fmt.Println("Appuyez sur Entrée pour continuer...")
		attendreEntree()
		return
	}

	if err := keyboard.Open(); err != nil {
		// fallback: simple affichage
		afficherInventaire(j)
		return
	}
	defer keyboard.Close()

	index := 0
	for {
		// Render gold and list (grouped) with cursor
		noms, counts := compterItems(j)
		fmt.Printf("Or: %d\n", j.Argent)
		for i, item := range noms {
			prefix := "  "
			if i == index {
				prefix = "> "
			}
			suffix := ""
			if estArmeEquipee(j, item) || estArmureEquipee(j, item) {
				suffix = "  [Équipé]"
			}
			label := item
			// Afficher les descriptions des potions
			switch strings.ToLower(item) {
			case "potion":
				label += " (+20 PV)"
			case "potion majeure":
				label += " (+50 PV)"
			case "potion force":
				label += " (+2 Force, 3 combats)"
			case "potion agilite":
				label += " (+2 Agilité, 3 combats)"
			case "potion endurance":
				label += " (+2 Endurance, 3 combats)"
			case "antidote":
				label += " (Guérit statuts)"
			case "elixir vie":
				label += " (+100 PV)"
			}

			// Afficher les stats des armes
			if arme, ok := trouverArmeParNom(item); ok && arme.Nom != "" {
				label += fmt.Sprintf(" (ATK %d)", arme.EffetAttaque)
			}

			// Afficher les stats des armures
			if armure, ok := trouverArmureParNom(item); ok && armure.Nom != "" {
				label += fmt.Sprintf(" (DEF %d)", armure.EffetDefense)
			}

			// Marquer les objets droppés
			if strings.Contains(item, "[DROPPÉ]") {
				label = strings.Replace(label, "[DROPPÉ]", "🎁", 1)
			}
			if counts[item] > 1 {
				fmt.Printf("%s%s x%d%s\n", prefix, label, counts[item], suffix)
			} else {
				fmt.Printf("%s%s%s\n", prefix, label, suffix)
			}
		}

		// Input
		char, key, err := keyboard.GetKey()
		if err != nil {
			return
		}

		// Clear rendered lines (gold line + items)
		for i := 0; i < len(noms)+1; i++ {
			fmt.Print("\033[A\033[2K")
		}

		switch key {
		case keyboard.KeyArrowUp:
			if index > 0 {
				index--
			} else {
				index = len(noms) - 1
			}
		case keyboard.KeyArrowDown:
			if index < len(noms)-1 {
				index++
			} else {
				index = 0
			}
		case keyboard.KeyEnter:
			// Use selected item by name
			if index >= 0 && index < len(noms) {
				_ = utiliserObjetNom(j, noms[index])
				// clamp after potential change
				noms2, _ := compterItems(j)
				if index >= len(noms2) && len(noms2) > 0 {
					index = len(noms2) - 1
				}
			}
			// brief feedback line
			fmt.Println("(Objet utilisé. Appuyez sur Entrée pour continuer / ESC pour quitter)")
			// wait for key then clear the line
			_, k2, _ := keyboard.GetKey()
			fmt.Print("\033[A\033[2K")
			if k2 == keyboard.KeyEsc {
				return
			}
		case keyboard.KeyEsc:
			return
		default:
			if char == '\r' || char == '\n' {
				// treat as Enter
				if index >= 0 {
					_ = utiliserObjetNom(j, noms[index])
					noms2, _ := compterItems(j)
					if index >= len(noms2) && len(noms2) > 0 {
						index = len(noms2) - 1
					}
				}
			}
		}
		if char == 'q' || char == 'Q' {
			return
		}
	}
}

func afficherInventaireInteractifOld(j *personnage.Personnage) {
	// Afficher l'inventaire
	afficherInventaire(j)

	// Si l'inventaire est vide, attendre une entrée et retourner
	if estInventaireVide(j) {
		fmt.Println("Appuyez sur Entrée pour continuer...")
		attendreEntree()
		return
	}

	if err := keyboard.Open(); err != nil {
		// fallback: simple affichage
		afficherInventaire(j)
		return
	}
	defer keyboard.Close()

	index := 0
	for {
		// Render gold and list (grouped) with cursor
		noms, counts := compterItems(j)
		fmt.Printf("Or: %d\n", j.Argent)
		for i, item := range noms {
			prefix := "  "
			if i == index {
				prefix = "> "
			}
			suffix := ""
			if estArmeEquipee(j, item) || estArmureEquipee(j, item) {
				suffix = "  [Équipé]"
			}
			label := item
			// Afficher les descriptions des potions
			switch strings.ToLower(item) {
			case "potion":
				label += " (+20 PV)"
			case "potion majeure":
				label += " (+50 PV)"
			case "potion force":
				label += " (+2 Force, 3 combats)"
			case "potion agilite":
				label += " (+2 Agilité, 3 combats)"
			case "potion endurance":
				label += " (+2 Endurance, 3 combats)"
			case "antidote":
				label += " (Guérit statuts)"
			case "elixir vie":
				label += " (+100 PV)"
			}

			// Afficher les stats des armes
			if arme, ok := trouverArmeParNom(item); ok && arme.Nom != "" {
				label += fmt.Sprintf(" (ATK %d)", arme.EffetAttaque)
			}

			// Afficher les stats des armures
			if armure, ok := trouverArmureParNom(item); ok && armure.Nom != "" {
				label += fmt.Sprintf(" (DEF %d)", armure.EffetDefense)
			}

			// Marquer les objets droppés
			if strings.Contains(item, "[DROPPÉ]") {
				label = strings.Replace(label, "[DROPPÉ]", "🎁", 1)
			}
			if counts[item] > 1 {
				fmt.Printf("%s%s x%d%s\n", prefix, label, counts[item], suffix)
			} else {
				fmt.Printf("%s%s%s\n", prefix, label, suffix)
			}
		}

		// Input
		char, key, err := keyboard.GetKey()
		if err != nil {
			return
		}

		// Clear rendered lines (gold line + items)
		for i := 0; i < len(noms)+1; i++ {
			fmt.Print("\033[A\033[2K")
		}

		switch key {
		case keyboard.KeyArrowUp:
			if index > 0 {
				index--
			} else {
				index = len(noms) - 1
			}
		case keyboard.KeyArrowDown:
			if index < len(noms)-1 {
				index++
			} else {
				index = 0
			}
		case keyboard.KeyEnter:
			// Use selected item by name
			if index >= 0 && index < len(noms) {
				_ = utiliserObjetNom(j, noms[index])
				// clamp after potential change
				noms2, _ := compterItems(j)
				if index >= len(noms2) && len(noms2) > 0 {
					index = len(noms2) - 1
				}
			}
			// brief feedback line
			fmt.Println("(Objet utilisé. Appuyez sur Entrée pour continuer / ESC pour quitter)")
			// wait for key then clear the line
			_, k2, _ := keyboard.GetKey()
			fmt.Print("\033[A\033[2K")
			if k2 == keyboard.KeyEsc {
				return
			}
		case keyboard.KeyEsc:
			return
		default:
			if char == '\r' || char == '\n' {
				// treat as Enter
				if index >= 0 {
					_ = utiliserObjetNom(j, noms[index])
					noms2, _ := compterItems(j)
					if index >= len(noms2) && len(noms2) > 0 {
						index = len(noms2) - 1
					}
				}
			}
		}
		if char == 'q' || char == 'Q' {
			return
		}
	}
}

func utiliserPotion(j *personnage.Personnage) {
	if retirerObjetParNom(j, "potion") {
		j.PVActuels += 20
		if j.PVActuels > j.PVMax {
			j.PVActuels = j.PVMax
		}
		fmt.Printf("💖 Potion utilisée ! PV : %d/%d\n", j.PVActuels, j.PVMax)
		return
	}
	fmt.Println("❌ Vous n'avez pas de potion !")
}

func retirerObjet(j *personnage.Personnage, index int) {
	if index < 0 || index >= len(j.Inventaire) {
		return
	}
	j.Inventaire = append(j.Inventaire[:index], j.Inventaire[index+1:]...)
}

// ajouterObjet ajoute un objet à l'inventaire
func ajouterObjet(j *personnage.Personnage, objet string) bool {
	poidsActuel := PoidsTotal(j)
	poidsAjout := PoidsObjet(objet)
	if poidsActuel+poidsAjout > PoidsMaxInventaire {
		fmt.Printf("❌ Trop lourd: %s (poids %d). Poids actuel %d/%d.\n", objet, poidsAjout, poidsActuel, PoidsMaxInventaire)
		return false
	}
	j.Inventaire = append(j.Inventaire, objet)
	return true
}

// retirerObjetParNom retire le premier objet correspondant (insensible à la casse)
// et retourne true si un objet a été retiré
func retirerObjetParNom(j *personnage.Personnage, nom string) bool {
	for i := range j.Inventaire {
		if strings.EqualFold(j.Inventaire[i], nom) {
			retirerObjet(j, i)
			return true
		}
	}
	return false
}

func estInventaireVide(j *personnage.Personnage) bool {
	return len(j.Inventaire) == 0
}

// utiliserObjetNom permet d'utiliser un objet par son nom (insensible à la casse).
// - Potion: soigner et consommer
// - Arme: équiper (met à jour p.Attaque), ne consomme pas
// - Armure: afficher/equiper visuellement (ne consomme pas)
func utiliserObjetNom(j *personnage.Personnage, nom string) bool {
	// Gérer toutes les potions
	switch strings.ToLower(nom) {
	case "potion":
		utiliserPotion(j)
		return true
	case "potion majeure":
		utiliserPotionMajeure(j)
		return true
	case "potion force":
		utiliserPotionForce(j)
		return true
	case "potion agilite":
		utiliserPotionAgilite(j)
		return true
	case "potion endurance":
		utiliserPotionEndurance(j)
		return true
	case "antidote":
		utiliserAntidote(j)
		return true
	case "elixir vie":
		utiliserElixirVie(j)
		return true
	}

	// Tente une correspondance avec les armes connues via clés et noms affichés
	if arme, ok := trouverArmeParNom(nom); ok {
		// Toggle: si déjà équipée, on déséquipe
		if strings.EqualFold(j.Attaque, arme.Nom) {
			j.Attaque = ""
			fmt.Printf("🔪 Arme déséquipée: %s\n", arme.Nom)
		} else {
			j.Attaque = arme.Nom
			fmt.Printf("🔪 Arme équipée: %s (Attaque %d)\n", arme.Nom, arme.EffetAttaque)
			objet.AfficherArme(arme)
		}
		return true
	}

	// Tente une correspondance avec les armures connues via clés et noms affichés
	if arm, ok := trouverArmureParNom(nom); ok {
		if arm.Nom == "" {
			fmt.Println("❌ Armure inconnue.")
			return false
		}
		if j.ArmuresEquipees == nil {
			j.ArmuresEquipees = make(map[string]bool)
		}
		// Toggle equip/desequip
		if j.ArmuresEquipees[arm.Nom] {
			delete(j.ArmuresEquipees, arm.Nom)
			fmt.Printf("🛡️ Armure déséquipée: %s\n", arm.Nom)
		} else {
			j.ArmuresEquipees[arm.Nom] = true
			fmt.Printf("🛡️ Armure équipée: %s (DEF %d)\n", arm.Nom, arm.EffetDefense)
		}
		objet.AfficherArmure(arm)
		return true
	}

	fmt.Println("❌ Objet inconnu/impropre à l'utilisation.")
	return false
}

// estArmeEquipee indique si le texte d'un item correspond à l'arme actuellement équipée
func estArmeEquipee(j *personnage.Personnage, item string) bool {
	if j.Attaque == "" {
		return false
	}
	// correspondance via clés et noms affichés
	if arme, ok := trouverArmeParNom(item); ok {
		return strings.EqualFold(j.Attaque, arme.Nom)
	}
	return false
}

// estArmureEquipee indique si le texte d'un item correspond à une armure équipée dans un slot
func estArmureEquipee(j *personnage.Personnage, item string) bool {
	if j.ArmuresEquipees == nil {
		return false
	}
	// normaliser par nom d'affichage (objet.CreerArmure renvoie .Nom)
	if arm, ok := trouverArmureParNom(item); ok {
		return j.ArmuresEquipees[arm.Nom]
	}
	return false
}

// compterItems regroupe l'inventaire par nom et renvoie l'ordre et les quantités
func compterItems(j *personnage.Personnage) ([]string, map[string]int) {
	counts := make(map[string]int)
	order := []string{}
	for _, it := range j.Inventaire {
		if _, ok := counts[it]; !ok {
			order = append(order, it)
		}
		counts[it]++
	}
	return order, counts
}

// utiliserObjetSelection utilise l'objet à l'index (1-based pour l'affichage) si possible
func utiliserObjetSelection(j *personnage.Personnage, indexAffiche int) bool {
	index := indexAffiche - 1
	if index < 0 || index >= len(j.Inventaire) {
		fmt.Println("❌ Index invalide.")
		return false
	}
	nom := j.Inventaire[index]
	ok := utiliserObjetNom(j, nom)
	// Consommation uniquement pour potion (déjà gérée par utiliserPotion via retirerObjetParNom)
	return ok
}

// --- Helpers de correspondance objets ---

func trouverArmeParNom(nom string) (objet.Arme, bool) {
	// Liste des clés d'armes supportées par objet.CreerArme
	cles := []string{
		"EpeeRouillee", "EpeeFer", "EpeeMagique", "EpeeCourte",
		"Hache", "HacheDeCombat", "HacheDeBataille",
		"ArcBois", "ArcLong", "ArcElfe",
	}
	needle := strings.ToLower(strings.TrimSpace(nom))
	for _, cle := range cles {
		a := objet.CreerArme(cle)
		if strings.EqualFold(needle, cle) || strings.EqualFold(needle, a.Nom) {
			return a, true
		}
	}
	return objet.Arme{}, false
}

func trouverArmureParNom(nom string) (objet.Armure, bool) {
	cles := []string{
		// Casques
		"CasqueCuir", "CasqueCuirRenforce", "CasqueFer", "CasqueFerRenforce",
		// Plastrons
		"PlastronCuir", "PlastronCuirRenforce", "PlastronFer", "PlastronFerRenforce",
		// Pantalons
		"PantalonCuir", "PantalonCuirRenforce", "PantalonFer", "PantalonFerRenforce",
		// Chaussures
		"BottesCuir", "BottesCuirRenforce", "BottesFer", "BottesFerRenforce",
	}
	needle := strings.ToLower(strings.TrimSpace(nom))
	for _, cle := range cles {
		ar := objet.CreerArmure(cle)
		if strings.EqualFold(needle, cle) || strings.EqualFold(needle, ar.Nom) {
			return ar, true
		}
	}
	return objet.Armure{}, false
}

// Nouvelles fonctions de potions

func utiliserPotionMajeure(j *personnage.Personnage) {
	if retirerObjetParNom(j, "potion majeure") {
		j.PVActuels += 50
		if j.PVActuels > j.PVMax {
			j.PVActuels = j.PVMax
		}
		fmt.Printf("💖 Potion majeure utilisée ! PV : %d/%d\n", j.PVActuels, j.PVMax)
		return
	}
	fmt.Println("❌ Vous n'avez pas de potion majeure !")
}

func utiliserPotionForce(j *personnage.Personnage) {
	if retirerObjetParNom(j, "potion force") {
		j.BuffForce += 2
		j.BuffCombats = 3
		fmt.Printf("💪 Potion de force utilisée ! Force +2 pour 3 combats\n")
		return
	}
	fmt.Println("❌ Vous n'avez pas de potion de force !")
}

func utiliserPotionAgilite(j *personnage.Personnage) {
	if retirerObjetParNom(j, "potion agilite") {
		j.BuffAgilite += 2
		j.BuffCombats = 3
		fmt.Printf("🏃 Potion d'agilité utilisée ! Agilité +2 pour 3 combats\n")
		return
	}
	fmt.Println("❌ Vous n'avez pas de potion d'agilité !")
}

func utiliserPotionEndurance(j *personnage.Personnage) {
	if retirerObjetParNom(j, "potion endurance") {
		j.BuffEndurance += 2
		j.BuffCombats = 3
		// Augmenter temporairement les PV max
		oldPVMax := j.PVMax
		j.PVMax += 20     // +2 Endurance = +20 PV
		j.PVActuels += 20 // Bonus immédiat de PV
		fmt.Printf("❤️ Potion d'endurance utilisée ! Endurance +2, PV Max %d→%d pour 3 combats\n", oldPVMax, j.PVMax)
		return
	}
	fmt.Println("❌ Vous n'avez pas de potion d'endurance !")
}

func utiliserAntidote(j *personnage.Personnage) {
	if retirerObjetParNom(j, "antidote") {
		fmt.Printf("🧪 Antidote utilisé ! Tous les statuts négatifs sont guéris\n")
		// Note: La logique de guérison des statuts sera gérée dans le combat
		return
	}
	fmt.Println("❌ Vous n'avez pas d'antidote !")
}

func utiliserElixirVie(j *personnage.Personnage) {
	if retirerObjetParNom(j, "elixir vie") {
		j.PVActuels += 100
		if j.PVActuels > j.PVMax {
			j.PVActuels = j.PVMax
		}
		fmt.Printf("✨ Élixir de vie utilisé ! PV : %d/%d\n", j.PVActuels, j.PVMax)
		return
	}
	fmt.Println("❌ Vous n'avez pas d'élixir de vie !")
}
