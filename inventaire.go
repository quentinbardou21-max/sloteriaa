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

// PoidsTotal calcule le poids total actuel de l'inventaire
func PoidsTotal(j *personnage.Personnage) int {
	total := 0
	for _, objet := range j.Inventaire {
		total += PoidsObjet(objet)
	}
	return total
}

func afficherInventaire(j *personnage.Personnage) {
	if estInventaireVide(j) {
		fmt.Println("🧳 Votre inventaire est vide.")
		return
	}
	fmt.Println("🧳 Inventaire :")
	for i, item := range j.Inventaire {
		suffix := ""
		if estArmeEquipee(j, item) || estArmureEquipee(j, item) {
			suffix = "  [Équipé]"
		}
		fmt.Printf("%d. %s%s\n", i+1, item, suffix)
	}
}

// afficherInventaireInteractif permet de naviguer avec ↑/↓ et d'utiliser l'objet sélectionné avec Entrée
func afficherInventaireInteractif(j *personnage.Personnage) {
	if estInventaireVide(j) {
		fmt.Println("🧳 Votre inventaire est vide.")
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
		// Render list with cursor
		for i, item := range j.Inventaire {
			prefix := "  "
			if i == index {
				prefix = "> "
			}
			suffix := ""
			if estArmeEquipee(j, item) || estArmureEquipee(j, item) {
				suffix = "  [Équipé]"
			}
			fmt.Printf("%s%s%s\n", prefix, item, suffix)
		}

		// Input
		char, key, err := keyboard.GetKey()
		if err != nil {
			return
		}

		// Clear rendered lines
		for range j.Inventaire {
			fmt.Print("\033[A\033[2K")
		}

		switch key {
		case keyboard.KeyArrowUp:
			if index > 0 {
				index--
			} else {
				index = len(j.Inventaire) - 1
			}
		case keyboard.KeyArrowDown:
			if index < len(j.Inventaire)-1 {
				index++
			} else {
				index = 0
			}
		case keyboard.KeyEnter:
			// Use selected item
			if utiliserObjetSelection(j, index+1) {
				// If potion consumed, shrink list and clamp index
				if index >= len(j.Inventaire) && len(j.Inventaire) > 0 {
					index = len(j.Inventaire) - 1
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
				if utiliserObjetSelection(j, index+1) {
					if index >= len(j.Inventaire) && len(j.Inventaire) > 0 {
						index = len(j.Inventaire) - 1
					}
				}
			}
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
	for i, objet := range j.Inventaire {
		if strings.EqualFold(objet, nom) {
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
	if strings.EqualFold(nom, "potion") {
		utiliserPotion(j)
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
		// Toggle par type: équipe/déséquipe le slot correspondant
		switch arm.Type {
		case objet.TypeBouclier:
			if strings.EqualFold(j.Bouclier, arm.Nom) {
				j.Bouclier = ""
				fmt.Printf("🛡️ Bouclier déséquipé: %s\n", arm.Nom)
			} else {
				j.Bouclier = arm.Nom
				fmt.Printf("🛡️ Bouclier équipé: %s (Défense %d)\n", arm.Nom, arm.EffetDefense)
				objet.AfficherArmure(arm)
			}
		case objet.TypeCasque:
			if strings.EqualFold(j.Casque, arm.Nom) {
				j.Casque = ""
				fmt.Printf("🛡️ Casque déséquipé: %s\n", arm.Nom)
			} else {
				j.Casque = arm.Nom
				fmt.Printf("🛡️ Casque équipé: %s (Défense %d)\n", arm.Nom, arm.EffetDefense)
				objet.AfficherArmure(arm)
			}
		case objet.TypePlastron:
			if strings.EqualFold(j.Plastron, arm.Nom) {
				j.Plastron = ""
				fmt.Printf("🛡️ Plastron déséquipé: %s\n", arm.Nom)
			} else {
				j.Plastron = arm.Nom
				fmt.Printf("🛡️ Plastron équipé: %s (Défense %d)\n", arm.Nom, arm.EffetDefense)
				objet.AfficherArmure(arm)
			}
		case objet.TypePantalon:
			if strings.EqualFold(j.Pantalon, arm.Nom) {
				j.Pantalon = ""
				fmt.Printf("🛡️ Pantalon déséquipé: %s\n", arm.Nom)
			} else {
				j.Pantalon = arm.Nom
				fmt.Printf("🛡️ Pantalon équipé: %s (Défense %d)\n", arm.Nom, arm.EffetDefense)
				objet.AfficherArmure(arm)
			}
		case objet.TypeChaussure:
			if strings.EqualFold(j.Chaussures, arm.Nom) {
				j.Chaussures = ""
				fmt.Printf("🛡️ Chaussures déséquipées: %s\n", arm.Nom)
			} else {
				j.Chaussures = arm.Nom
				fmt.Printf("🛡️ Chaussures équipées: %s (Défense %d)\n", arm.Nom, arm.EffetDefense)
				objet.AfficherArmure(arm)
			}
		}
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
	if arm, ok := trouverArmureParNom(item); ok {
		switch arm.Type {
		case objet.TypeBouclier:
			return strings.EqualFold(j.Bouclier, arm.Nom)
		case objet.TypeCasque:
			return strings.EqualFold(j.Casque, arm.Nom)
		case objet.TypePlastron:
			return strings.EqualFold(j.Plastron, arm.Nom)
		case objet.TypePantalon:
			return strings.EqualFold(j.Pantalon, arm.Nom)
		case objet.TypeChaussure:
			return strings.EqualFold(j.Chaussures, arm.Nom)
		}
	}
	return false
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
		// Boucliers
		"BouclierBois", "BouclierFer",
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
