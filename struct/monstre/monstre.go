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
		// FALLBACK DEBUG: Ne devrait jamais apparaître en jeu normal
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

// Affiche les infos d'un monstre
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

// Structure pour les monstres de donjon (version simplifiée)
type MonsterDungeon struct {
	Nom     string
	PV      int
	Attaque int
	Defense int
	Type    string
	// Statuts du monstre
	Stunned  bool
	Poisoned bool
	Burned   bool
	Bleeding bool
	Shielded bool
	// Durée des statuts
	StunTurns   int
	PoisonTurns int
	BurnTurns   int
	BleedTurns  int
	ShieldTurns int
}

// Génère un monstre spécialisé pour un tier de donjon donné
func CreerMonstreDungeon(tier int) MonsterDungeon {
	// Multiplicateurs de difficulté basés sur le tier
	var baseHPMult, baseAtkMult, baseDefMult float64
	switch tier {
	case 1:
		baseHPMult, baseAtkMult, baseDefMult = 1.0, 1.0, 1.0
	case 2:
		baseHPMult, baseAtkMult, baseDefMult = 1.3, 1.2, 1.1
	case 3:
		baseHPMult, baseAtkMult, baseDefMult = 1.6, 1.4, 1.2
	case 4:
		baseHPMult, baseAtkMult, baseDefMult = 2.0, 1.6, 1.3
	default:
		baseHPMult, baseAtkMult, baseDefMult = 2.5, 2.0, 1.5
	}

	// Type basé sur le tier
	mtype := "Bête"
	switch {
	case tier >= 4:
		mtype = "Vétéran"
	case tier >= 3:
		mtype = "Guerrier"
	case tier >= 2:
		mtype = "Adepte"
	}

	// Choisir un monstre spécialisé selon le tier
	var nom string
	var baseHP, baseAtk, baseDef int

	switch tier {
	case 1:
		// Tier 1: 3 monstres variés (équilibrés pour joueur non équipé)
		switch rand.Intn(3) {
		case 0: // Gobelin agile - faible défense, forte attaque
			nom = "Gobelin agile"
			baseHP = 50
			baseAtk = 10
			baseDef = 1
		case 1: // Rat géant - équilibré
			nom = "Rat géant"
			baseHP = 70
			baseAtk = 8
			baseDef = 2
		case 2: // Squelette - tank
			nom = "Squelette"
			baseHP = 90
			baseAtk = 6
			baseDef = 4
		}
	case 2:
		// Tier 2: 3 monstres variés (pour joueur avec équipement basique)
		switch rand.Intn(3) {
		case 0: // Assassin - très forte attaque, très faible défense
			nom = "Assassin"
			baseHP = 80
			baseAtk = 18
			baseDef = 2
		case 1: // Bandit - équilibré
			nom = "Bandit"
			baseHP = 110
			baseAtk = 13
			baseDef = 6
		case 2: // Garde - tank
			nom = "Garde"
			baseHP = 140
			baseAtk = 10
			baseDef = 10
		}
	case 3:
		// Tier 3: 3 monstres variés (pour joueur avec équipement intermédiaire)
		switch rand.Intn(3) {
		case 0: // Berserker - attaque extrême, défense faible
			nom = "Berserker"
			baseHP = 150
			baseAtk = 50
			baseDef = 5
		case 1: // Orc - équilibré
			nom = "Orc"
			baseHP = 180
			baseAtk = 35
			baseDef = 20
		case 2: // Troll - tank massif
			nom = "Troll"
			baseHP = 250
			baseAtk = 25
			baseDef = 30
		}
	case 4:
		// Tier 4: 3 monstres variés (pour joueur avec équipement avancé)
		switch rand.Intn(3) {
		case 0: // Assassin maître - attaque mortelle
			nom = "Assassin maître"
			baseHP = 180
			baseAtk = 70
			baseDef = 8
		case 1: // Chevalier - équilibré puissant
			nom = "Chevalier"
			baseHP = 220
			baseAtk = 45
			baseDef = 35
		case 2: // Golem - tank ultime
			nom = "Golem"
			baseHP = 350
			baseAtk = 30
			baseDef = 50
		}
	default:
		// Tier 5+: Boss et créatures légendaires (pour joueur avec équipement légendaire)
		switch rand.Intn(3) {
		case 0: // Dragon ancien - attaque légendaire
			nom = "Dragon ancien"
			baseHP = 250
			baseAtk = 90
			baseDef = 20
		case 1: // Liche - équilibré magique
			nom = "Liche"
			baseHP = 350
			baseAtk = 65
			baseDef = 45
		case 2: // Titan - tank légendaire
			nom = "Titan"
			baseHP = 500
			baseAtk = 45
			baseDef = 70
		}
	}

	// Appliquer les multiplicateurs de tier
	adjustedHP := int(float64(baseHP) * baseHPMult)
	adjustedAtk := int(float64(baseAtk) * baseAtkMult)
	adjustedDef := int(float64(baseDef) * baseDefMult)

	return MonsterDungeon{
		Nom: nom, PV: adjustedHP, Attaque: adjustedAtk, Type: mtype, Defense: adjustedDef,
		// Initialiser les statuts à false
		Stunned: false, Poisoned: false, Burned: false, Bleeding: false, Shielded: false,
		StunTurns: 0, PoisonTurns: 0, BurnTurns: 0, BleedTurns: 0, ShieldTurns: 0,
	}
}
