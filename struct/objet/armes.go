package objet

import "fmt"

type TypeObjet string

const TypeArme TypeObjet = "Arme"

type Arme struct {
	Nom          string
	Description  string
	Type         TypeObjet
	EffetAttaque int
	Poids        int
}

// Arme uniquement utilisable par les monstres, avec des stats uniques
type ArmeMonstre struct {
	Nom          string
	Description  string
	Type         TypeObjet
	EffetAttaque int
	Poids        int
	Instabilite  int // Variance potentielle des dégâts (0-10)
	Sauvagerie   int // Brutalité/saignement potentiel (0-10)
}

// Création d'une arme selon son nom
func CreerArme(nom string) Arme {
	switch nom {
	// Épées
	case "EpeeRouillee":
		return Arme{"Épée rouillée", "Vieille épée peu puissante", TypeArme, 15, 5}
	case "EpeeFer":
		return Arme{"Épée en fer", "Épée solide et fiable", TypeArme, 35, 8}
	case "EpeeMagique":
		return Arme{"Épée magique", "Épée enchantée par la magie ancienne", TypeArme, 60, 6}
	case "EpeeCourte":
		return Arme{"Épée courte", "Épée rapide et maniable", TypeArme, 30, 4}

	// Haches
	case "Hache":
		return Arme{"Hache lourde", "Hache massive et lourde", TypeArme, 40, 10}
	case "HacheDeCombat":
		return Arme{"Hache de combat", "Hache équilibrée pour le combat", TypeArme, 35, 8}
	case "HacheDeBataille":
		return Arme{"Hache de bataille", "Hache puissante à deux mains", TypeArme, 50, 12}

	// Arcs
	case "ArcBois":
		return Arme{"Arc en bois", "Arc simple pour attaques à distance", TypeArme, 25, 3}
	case "ArcLong":
		return Arme{"Arc long", "Arc puissant et précis", TypeArme, 35, 4}
	case "ArcElfe":
		return Arme{"Arc elfique", "Arc léger et rapide, très précis", TypeArme, 40, 3}

	default:
		fmt.Println("Arme inconnue, création d'une épée rouillée par défaut")
		return Arme{"Épée rouillée", "Vieille épée peu puissante", TypeArme, 15, 5}
	}
}

// Création d'une arme exclusive aux monstres selon son nom
func CreerArmeMonstre(nom string) ArmeMonstre {
	switch nom {
	// Armes bas niveau (1-2)
	case "GriffesSouillees":
		return ArmeMonstre{"Griffes souillées", "Griffes acérées couvertes d'impuretés", TypeArme, 18, 0, 3, 4}
	case "MassueBrute":
		return ArmeMonstre{"Massue brute", "Gros morceau de bois noueux", TypeArme, 22, 9, 2, 4}

	// Milieu de progression (3-6)
	case "LanceBrisee":
		return ArmeMonstre{"Lance brisée", "Lance ébréchée ramassée sur un champ de bataille", TypeArme, 26, 5, 4, 3}
	case "EpeeOsseuse":
		return ArmeMonstre{"Épée osseuse", "Lame en os taillé, instable mais mordante", TypeArme, 32, 6, 5, 5}
	case "HacheTronquee":
		return ArmeMonstre{"Hache tronquée", "Hache abîmée mais très lourde", TypeArme, 36, 11, 3, 6}

	// Haut niveaux (7-10), inférieures aux meilleures armes du joueur
	case "GlaiveSauvage":
		return ArmeMonstre{"Glaive sauvage", "Glaive brut reforgé à la hâte", TypeArme, 38, 8, 6, 6}
	case "MasseRituelle":
		return ArmeMonstre{"Masse rituelle", "Masse gravée de symboles inquiétants", TypeArme, 42, 12, 4, 7}
	case "FauxDeBrume":
		return ArmeMonstre{"Faux de brume", "Lame incurvée difficile à prévoir", TypeArme, 44, 7, 7, 6}

	default:
		fmt.Println("Arme de monstre inconnue, création de Griffes souillées par défaut")
		return ArmeMonstre{"Griffes souillées", "Griffes acérées couvertes d'impuretés", TypeArme, 18, 0, 3, 4}
	}
}

// Affiche les infos d'une arme de monstre
func AfficherArmeMonstre(o ArmeMonstre) {
	fmt.Printf("Nom : %s\nDescription : %s\nAttaque : %d\nPoids : %d\nInstabilité : %d\nSauvagerie : %d\n\n",
		o.Nom, o.Description, o.EffetAttaque, o.Poids, o.Instabilite, o.Sauvagerie)
}

// Fonction pour afficher les infos d'une arme
func AfficherArme(o Arme) {
	fmt.Printf("Nom : %s\nDescription : %s\nAttaque : %d\nPoids : %d\n\n",
		o.Nom, o.Description, o.EffetAttaque, o.Poids)
}
