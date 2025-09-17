package monstre

import (
	"fmt"
	"math/rand"

	"sloteriaa/struct/objet"
)

type Monstre struct {
	Nom           string
	HPMax         int
	Attaque       int
	Defense       int
	Arme          objet.ArmeMonstre
	Armures       []objet.Armure
	Niveau        int
	PeutAvoirArme bool
}

// Attribution d'une arme selon le niveau et le type de monstre
func armePourMonstre(niveau int, peutAvoirArme bool) objet.ArmeMonstre {
	if !peutAvoirArme {
		return objet.ArmeMonstre{} // Arme vide si le monstre ne peut pas en avoir
	}

	switch {
	case niveau <= 2:
		armesDispos := []string{"GriffesSouillees", "MassueBrute"}
		return objet.CreerArmeMonstre(armesDispos[rand.Intn(len(armesDispos))])
	case niveau <= 4:
		armesDispos := []string{"LanceBrisee", "EpeeOsseuse"}
		return objet.CreerArmeMonstre(armesDispos[rand.Intn(len(armesDispos))])
	case niveau <= 6:
		armesDispos := []string{"HacheTronquee", "EpeeOsseuse"}
		return objet.CreerArmeMonstre(armesDispos[rand.Intn(len(armesDispos))])
	case niveau <= 8:
		armesDispos := []string{"GlaiveSauvage", "MasseRituelle"}
		return objet.CreerArmeMonstre(armesDispos[rand.Intn(len(armesDispos))])
	default: // niveau 9-10
		armesDispos := []string{"MasseRituelle", "FauxDeBrume"}
		return objet.CreerArmeMonstre(armesDispos[rand.Intn(len(armesDispos))])
	}
}

// Attribution d'armures aléatoires
func armuresPourMonstre(niveau int) []objet.Armure {
	liste := []objet.Armure{}

	if rand.Intn(2) == 0 {
		liste = append(liste, objet.CreerArmure("CasqueCuir"))
	}
	if niveau >= 3 && rand.Intn(2) == 0 {
		liste = append(liste, objet.CreerArmure("PlastronCuirRenforce"))
	}
	if niveau >= 5 && rand.Intn(2) == 0 {
		liste = append(liste, objet.CreerArmure("PantalonFer"))
	}
	if niveau >= 7 && rand.Intn(2) == 0 {
		liste = append(liste, objet.CreerArmure("CasqueFerRenforce"))
	}
	if niveau >= 9 {
		liste = append(liste,
			objet.CreerArmure("PlastronFerRenforce"),
			objet.CreerArmure("BottesFerRenforce"),
		)
	}
	return liste
}

// Création d'un monstre selon son niveau
func CreerMonstre(niveau int) Monstre {
	var nom string
	var hpMax int
	var defense int
	var attaque int
	var peutAvoirArme bool

	switch niveau {
	case 1:
		nom = "Rat géant"
		hpMax = 100 + rand.Intn(10)
		defense = 3
		attaque = 8 + rand.Intn(3)
		peutAvoirArme = true
	case 2:
		nom = "Gobelin"
		hpMax = 110 + rand.Intn(20)
		defense = 5
		attaque = 10 + rand.Intn(5)
		peutAvoirArme = true
	case 3:
		nom = "Bandit"
		hpMax = 120 + rand.Intn(20)
		defense = 7
		attaque = 12 + rand.Intn(5)
		peutAvoirArme = true
	case 4:
		nom = "Orc"
		hpMax = 130 + rand.Intn(20)
		defense = 10
		attaque = 15 + rand.Intn(5)
		peutAvoirArme = true
	case 5:
		nom = "Gnoll"
		hpMax = 140 + rand.Intn(20)
		defense = 12
		attaque = 18 + rand.Intn(5)
		peutAvoirArme = true
	case 6:
		nom = "Troll"
		hpMax = 160 + rand.Intn(20)
		defense = 15
		attaque = 20 + rand.Intn(5)
		peutAvoirArme = true
	case 7:
		nom = "Ogre"
		hpMax = 180 + rand.Intn(20)
		defense = 18
		attaque = 25 + rand.Intn(8)
		peutAvoirArme = true
	case 8:
		nom = "Élémentaire de pierre"
		hpMax = 210 + rand.Intn(20)
		defense = 22
		attaque = 33 + rand.Intn(8)
		peutAvoirArme = false
	case 9:
		nom = "Chevalier maudit"
		hpMax = 210 + rand.Intn(20)
		defense = 25
		attaque = 30 + rand.Intn(8)
		peutAvoirArme = true
	case 10:
		nom = "Dragon"
		hpMax = 260 + rand.Intn(20)
		defense = 30
		attaque = 45 + rand.Intn(10)
		peutAvoirArme = false
	default:
		nom = "Créature inconnue"
		hpMax = 100 + rand.Intn(50)
		defense = 5 + rand.Intn(5)
		attaque = 10 + rand.Intn(10)
		peutAvoirArme = true
	}

	return Monstre{
		Nom:           nom,
		HPMax:         hpMax,
		Attaque:       attaque,
		Defense:       defense,
		Arme:          armePourMonstre(niveau, peutAvoirArme),
		Armures:       armuresPourMonstre(niveau),
		Niveau:        niveau,
		PeutAvoirArme: peutAvoirArme,
	}
}

// Affiche les infos d’un monstre
func AfficherMonstre(m Monstre) {
	fmt.Printf("Nom : %s\nHP : %d\nAttaque : %d\nDéfense : %d\n",
		m.Nom, m.HPMax, m.Attaque, m.Defense)

	if m.PeutAvoirArme {
		fmt.Printf("Arme (monstre) : %s (Atk %d, Instab %d, Sauv %d)\n", m.Arme.Nom, m.Arme.EffetAttaque, m.Arme.Instabilite, m.Arme.Sauvagerie)
	} else {
		fmt.Println("Arme : Aucune")
	}

	if len(m.Armures) > 0 {
		fmt.Println("Armures équipées :")
		for _, a := range m.Armures {
			fmt.Printf("  - %s (Déf %d)\n", a.Nom, a.EffetDefense)
		}
	} else {
		fmt.Println("Aucune armure")
	}

	fmt.Printf("Niveau : %d\n\n", m.Niveau)
}
