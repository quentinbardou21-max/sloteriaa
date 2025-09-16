package main

import (
	"fmt"
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
func PoidsTotal(j *Joueur) int {
	total := 0
	for _, objet := range j.Inventaire {
		total += PoidsObjet(objet)
	}
	return total
}

func afficherInventaire(j *Joueur) {
	if estInventaireVide(j) {
		fmt.Println("🧳 Votre inventaire est vide.")
		return
	}
	fmt.Println("🧳 Inventaire :")
	for i, objet := range j.Inventaire {
		fmt.Printf("%d. %s\n", i+1, objet)
	}
}

func utiliserPotion(j *Joueur) {
	if retirerObjetParNom(j, "potion") {
		j.HP += 20
		if j.HP > j.HPMax {
			j.HP = j.HPMax
		}
		fmt.Printf("💖 Potion utilisée ! HP : %d/%d\n", j.HP, j.HPMax)
		return
	}
	fmt.Println("❌ Vous n'avez pas de potion !")
}

func retirerObjet(j *Joueur, index int) {
	if index < 0 || index >= len(j.Inventaire) {
		return
	}
	j.Inventaire = append(j.Inventaire[:index], j.Inventaire[index+1:]...)
}

// ajouterObjet ajoute un objet à l'inventaire
func ajouterObjet(j *Joueur, objet string) bool {
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
func retirerObjetParNom(j *Joueur, nom string) bool {
	for i, objet := range j.Inventaire {
		if strings.EqualFold(objet, nom) {
			retirerObjet(j, i)
			return true
		}
	}
	return false
}

func estInventaireVide(j *Joueur) bool {
	return len(j.Inventaire) == 0
}
