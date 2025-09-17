package main

import (
	"fmt"
<<<<<<< Updated upstream
	"sloteriaa/internal/personnage"
	"strings"
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
=======
	"strings"
)

func afficherInventaire(j *Joueur) {
>>>>>>> Stashed changes
	if estInventaireVide(j) {
		fmt.Println("🧳 Votre inventaire est vide.")
		return
	}
	fmt.Println("🧳 Inventaire :")
	for i, objet := range j.Inventaire {
		fmt.Printf("%d. %s\n", i+1, objet)
	}
}

<<<<<<< Updated upstream
func utiliserPotion(j *personnage.Personnage) {
	if retirerObjetParNom(j, "potion") {
		j.PVActuels += 20
		if j.PVActuels > j.PVMax {
			j.PVActuels = j.PVMax
		}
		fmt.Printf("💖 Potion utilisée ! PV : %d/%d\n", j.PVActuels, j.PVMax)
		return
=======
func utiliserPotion(j *Joueur) {
	for i, objet := range j.Inventaire {
		if strings.ToLower(objet) == "potion" {
			j.HP += 20
			if j.HP > j.HPMax {
				j.HP = j.HPMax
			}
			fmt.Printf("💖 Potion utilisée ! HP : %d/%d\n", j.HP, j.HPMax)
			retirerObjet(j, i) // retire la potion de l'inventaire
			return
		}
>>>>>>> Stashed changes
	}
	fmt.Println("❌ Vous n'avez pas de potion !")
}

<<<<<<< Updated upstream
func retirerObjet(j *personnage.Personnage, index int) {
=======
func retirerObjet(j *Joueur, index int) {
>>>>>>> Stashed changes
	if index < 0 || index >= len(j.Inventaire) {
		return
	}
	j.Inventaire = append(j.Inventaire[:index], j.Inventaire[index+1:]...)
}

<<<<<<< Updated upstream
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
=======
func estInventaireVide(j *Joueur) bool {
>>>>>>> Stashed changes
	return len(j.Inventaire) == 0
}
