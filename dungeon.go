package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"sloteriaa/internal/personnage"
	"sloteriaa/struct/forgeron"
	"sloteriaa/struct/monstre"
	"sloteriaa/struct/objet"

	"github.com/eiannone/keyboard"
)

type Monster struct {
	Nom     string
	PV      int
	PVMax   int
	Attaque int
	Type    string
	Defense int
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

// Structure pour les attaques spéciales
type SpecialAttack struct {
	Nom         string
	Description string
	Damage      int
	Effects     []StatusEffect
	Cooldown    int
	CurrentCD   int
}

// Structure pour les effets de statut
type StatusEffect struct {
	Type        string // "stun", "poison", "burn", "bleed", "shield"
	Duration    int
	Damage      int // dégâts par tour (pour poison, burn, bleed)
	Description string
}

// Structure pour les drops d'objets
type ItemDrop struct {
	ItemName string
	Chance   int // pourcentage de chance (0-100)
}

// Structure pour les drops de matériaux
type MaterialDrop struct {
	Material forgeron.Materiau
	Quantity int
	Chance   int // pourcentage de chance (0-100)
}

// Attaques spéciales du joueur
var playerSpecialAttacks = []SpecialAttack{
	{
		Nom:         "Coup de poing",
		Description: "Attaque basique sans effet spécial",
		Damage:      0, // utilise l'attaque normale
		Effects:     []StatusEffect{},
		Cooldown:    0,
		CurrentCD:   0,
	},
	{
		Nom:         "Coup étourdissant",
		Description: "Assomme l'ennemi (étourdit 1 tour)",
		Damage:      -5, // -5 dégâts mais étourdit
		Effects:     []StatusEffect{{Type: "stun", Duration: 1, Damage: 0, Description: "Étourdi"}},
		Cooldown:    3,
		CurrentCD:   0,
	},
	{
		Nom:         "Coup empoisonné",
		Description: "Empoisonne l'ennemi (5% PV max/tour pendant 3 tours)",
		Damage:      -3,
		Effects:     []StatusEffect{{Type: "poison", Duration: 3, Damage: 5, Description: "Empoisonné"}},
		Cooldown:    4,
		CurrentCD:   0,
	},
	{
		Nom:         "Coup de feu",
		Description: "Brûle l'ennemi (3% PV max/tour pendant 4 tours)",
		Damage:      -2,
		Effects:     []StatusEffect{{Type: "burn", Duration: 4, Damage: 3, Description: "Brûlé"}},
		Cooldown:    5,
		CurrentCD:   0,
	},
	{
		Nom:         "Coup saignant",
		Description: "Fait saigner l'ennemi (4% PV max/tour pendant 2 tours)",
		Damage:      -1,
		Effects:     []StatusEffect{{Type: "bleed", Duration: 2, Damage: 4, Description: "Saigne"}},
		Cooldown:    3,
		CurrentCD:   0,
	},
}

// Tables de drop par niveau de donjon
var dungeonDrops = map[int]struct {
	Items     []ItemDrop
	Materials []MaterialDrop
}{
	1: { // Donjon Cuir
		Items: []ItemDrop{
			{"EpeeRouillee", 15},
			{"ArcBois", 10},
			{"CasqueCuir", 8},
			{"PlastronCuir", 8},
			{"PantalonCuir", 8},
			{"BottesCuir", 8},
		},
		Materials: []MaterialDrop{
			{forgeron.Cuir, 1, 40},
			{forgeron.Bois, 1, 30},
			{forgeron.Fer, 1, 20},
			{forgeron.Or, 50, 60},
		},
	},
	2: { // Donjon Cuir Renforcé
		Items: []ItemDrop{
			{"EpeeCourte", 12},
			{"ArcLong", 8},
			{"Hache", 10},
			{"CasqueCuirRenforce", 6},
			{"PlastronCuirRenforce", 6},
			{"PantalonCuirRenforce", 6},
			{"BottesCuirRenforce", 6},
		},
		Materials: []MaterialDrop{
			{forgeron.Cuir, 2, 35},
			{forgeron.Fer, 1, 45},
			{forgeron.Bois, 2, 25},
			{forgeron.Or, 100, 70},
		},
	},
	3: { // Donjon Fer
		Items: []ItemDrop{
			{"EpeeFer", 10},
			{"HacheDeCombat", 8},
			{"ArcElfe", 6},
			{"CasqueFer", 5},
			{"PlastronFer", 5},
			{"PantalonFer", 5},
			{"BottesFer", 5},
		},
		Materials: []MaterialDrop{
			{forgeron.Fer, 2, 50},
			{forgeron.Cuir, 1, 30},
			{forgeron.Bois, 1, 20},
			{forgeron.EssenceMagique, 1, 15},
			{forgeron.Or, 200, 80},
		},
	},
	4: { // Donjon Fer Renforcé
		Items: []ItemDrop{
			{"EpeeMagique", 8},
			{"HacheDeBataille", 6},
			{"CasqueFerRenforce", 4},
			{"PlastronFerRenforce", 4},
			{"PantalonFerRenforce", 4},
			{"BottesFerRenforce", 4},
		},
		Materials: []MaterialDrop{
			{forgeron.Fer, 3, 40},
			{forgeron.EssenceMagique, 2, 25},
			{forgeron.Cuir, 2, 35},
			{forgeron.Bois, 2, 30},
			{forgeron.Or, 300, 90},
		},
	},
}

// Attaques spéciales des monstres
var monsterSpecialAttacks = map[string][]SpecialAttack{
	"Gobelin agile": {
		{Nom: "Griffes rapides", Description: "Attaque rapide", Damage: 0, Effects: []StatusEffect{}, Cooldown: 0, CurrentCD: 0},
		{Nom: "Coup sournois", Description: "Étourdit l'ennemi", Damage: -3, Effects: []StatusEffect{{Type: "stun", Duration: 1, Damage: 0, Description: "Étourdi"}}, Cooldown: 4, CurrentCD: 0},
	},
	"Rat géant": {
		{Nom: "Morsure", Description: "Attaque basique", Damage: 0, Effects: []StatusEffect{}, Cooldown: 0, CurrentCD: 0},
		{Nom: "Morsure empoisonnée", Description: "Empoisonne l'ennemi", Damage: -2, Effects: []StatusEffect{{Type: "poison", Duration: 2, Damage: 2, Description: "Empoisonné"}}, Cooldown: 3, CurrentCD: 0},
	},
	"Squelette": {
		{Nom: "Coup d'os", Description: "Attaque basique", Damage: 0, Effects: []StatusEffect{}, Cooldown: 0, CurrentCD: 0},
		{Nom: "Malédiction", Description: "Affaiblit l'ennemi", Damage: 0, Effects: []StatusEffect{{Type: "bleed", Duration: 3, Damage: 1, Description: "Maudit"}}, Cooldown: 5, CurrentCD: 0},
	},
}

func EnterDungeon(gs *GameState) {
	for {
		idx, cancelled := selectWithArrows("Donjon — choisissez une salle:", []string{
			"Couloir bas-niveau",
			"Couloir novice (lvl 5)",
			"Couloir intermédiaire (lvl 10)",
			"Antre avancée (lvl 15)",
			"Salle du Boss (lvl 20)",
			"Sortir du donjon",
		})
		if cancelled {
			return
		}
		switch idx {
		case 0:
			if gs.Level < requiredLevelForTier(1) {
				fmt.Printf("Niveau insuffisant (niveau requis: %d).\n", requiredLevelForTier(1))
				attendreEntree()
				continue
			}
			fightRoom(gs, 1)
		case 1:
			if gs.Level < 5 {
				fmt.Println("Niveau insuffisant (niveau requis: 5).")
				attendreEntree()
				continue
			}
			fightRoom(gs, 2)
		case 2:
			if gs.Level < 10 {
				fmt.Println("Niveau insuffisant (niveau requis: 10).")
				attendreEntree()
				continue
			}
			fightRoom(gs, 3)
		case 3:
			if gs.Level < 15 {
				fmt.Println("Niveau insuffisant (niveau requis: 15).")
				attendreEntree()
				continue
			}
			fightRoom(gs, 4)
		case 4:
			if gs.Level < 20 {
				fmt.Println("Niveau insuffisant (niveau requis: 20).")
				attendreEntree()
				continue
			}
			bossFight(gs)
		case 5:
			return
		}
	}
}

func fightRoom(gs *GameState, tier int) {
	if gs.Joueur.PVActuels <= 0 {
		gs.Joueur.PVActuels = gs.Joueur.PVMax
	}
	mon := generateMonster(tier)
	fmt.Printf("Un %s apparaît ! (PV %d, ATK %d)\n", mon.Nom, mon.PV, mon.Attaque)
	playerHP := gs.Joueur.PVActuels
	monsterStunned := false
	playerGuard := false
	playerStunned := false
	fled := false
	attackedOnce := false
	for mon.PV > 0 && playerHP > 0 {
		// Appliquer les effets de statut du monstre
		applyStatusEffects(&mon, mon.PVMax)

		// Décrémenter les cooldowns des monstres
		if attacks, exists := monsterSpecialAttacks[mon.Nom]; exists {
			for i := range attacks {
				if attacks[i].CurrentCD > 0 {
					attacks[i].CurrentCD--
				}
			}
		}

		renderBattle(gs, mon, playerHP, monsterStunned, playerGuard, playerStunned)
		dmg, didStun, didGuard, didUse := 0, false, false, false
		if playerStunned {
			fmt.Println("Vous êtes étourdi et perdez votre tour !")
		} else {
			dmg, didStun, didGuard, didUse = playerAction(gs)
		}
		if dmg == -1 {
			if attackedOnce {
				fmt.Println("Vous avez déjà attaqué. Vous ne pouvez plus fuir !")
			} else {
				fmt.Println("Vous prenez la fuite !")
				fled = true
				break
			}
		}
		if dmg > 0 {
			mon.PV -= dmg
			if mon.PV < 0 {
				mon.PV = 0
			}
			fmt.Printf("Vous infligez %d dégâts. (PV monstre %d)\n", dmg, mon.PV)
			attackedOnce = true
		}
		if didStun {
			monsterStunned = true
			fmt.Println("Le monstre est étourdi pour 1 tour !")
			attackedOnce = true
		}
		playerGuard = didGuard
		if didUse {
			fmt.Printf("PV: %d/%d\n", gs.Joueur.PVActuels, gs.Joueur.PVMax)
			playerHP = gs.Joueur.PVActuels
		}
		if mon.PV <= 0 {
			break
		}
		if monsterStunned {
			fmt.Println("Le monstre est étourdi et ne peut pas attaquer.")
			monsterStunned = false
		} else {
			mincoming, specialName, stunPlayer := monsterAction(mon)
			reduction := gs.Joueur.Endurance / 3
			if reduction > mincoming-1 {
				reduction = mincoming - 1
			}
			mincoming -= reduction
			if playerGuard {
				mincoming = mincoming / 2
			}
			if mincoming < 1 {
				mincoming = 1
			}
			playerHP -= mincoming
			if specialName != "" {
				fmt.Printf("%s utilise %s et inflige %d (PV %d/%d)\n", mon.Nom, specialName, mincoming, max0(playerHP), gs.Joueur.PVMax)
				if stunPlayer {
					fmt.Println("Vous êtes étourdi pour 1 tour !")
				}
			} else {
				fmt.Printf("Le %s vous touche pour %d (PV %d/%d)\n", mon.Nom, mincoming, max0(playerHP), gs.Joueur.PVMax)
			}

			// Mettre à jour l'attaque du joueur (pour la transformation du loup-garou)
			gs.Joueur.PVActuels = playerHP
			personnage.UpdatePlayerAttack(&gs.Joueur)
			playerStunned = stunPlayer
		}
		playerGuard = false

		// Pause d'un tour pour que l'action soit visible
		fmt.Println("(Appuyez sur Entrée pour continuer)")
		attendreEntree()
	}
	if playerHP <= 0 {
		fmt.Println("Vous tombez inconscient... Vous êtes ramené à la ville.")
		gs.Joueur.PVActuels = gs.Joueur.PVMax
		fmt.Println("(Appuyez sur Entrée pour revenir)")
		attendreEntree()
		return
	}
	gs.Joueur.PVActuels = playerHP
	if fled {
		fmt.Println("Vous avez fui. Aucune récompense.")
		fmt.Println("(Appuyez sur Entrée pour revenir)")
		attendreEntree()
		return
	}
	fmt.Println("Victoire !")
	reward(gs, tier)
	gainXP(gs, xpForTier(tier))
	fmt.Println("(Appuyez sur Entrée pour revenir)")
	attendreEntree()
}

func bossFight(gs *GameState) {
	if gs.Joueur.PVActuels <= 0 {
		gs.Joueur.PVActuels = gs.Joueur.PVMax
	}
	fmt.Println("Vous entrez dans la salle interdite... Votre mère, métamorphosée, se dresse devant vous !")
	mon := Monster{Nom: "Mère métamorphe", PV: 400, Attaque: 35, Type: "Boss"}
	playerHP := gs.Joueur.PVActuels
	monsterStunned := false
	playerGuard := false
	playerStunned := false
	fled := false
	attackedOnce := false
	for mon.PV > 0 && playerHP > 0 {
		renderBattle(gs, mon, playerHP, monsterStunned, playerGuard, playerStunned)
		dmg, didStun, didGuard, didUse := 0, false, false, false
		if playerStunned {
			fmt.Println("Vous êtes étourdi et perdez votre tour !")
		} else {
			dmg, didStun, didGuard, didUse = playerAction(gs)
		}
		if dmg == -1 {
			if attackedOnce {
				fmt.Println("Vous avez déjà attaqué. Vous ne pouvez plus fuir !")
			} else {
				fmt.Println("Vous prenez la fuite !")
				fled = true
				break
			}
		}
		if dmg > 0 {
			mon.PV -= dmg
			if mon.PV < 0 {
				mon.PV = 0
			}
			fmt.Printf("Vous infligez %d dégâts. (PV monstre %d)\n", dmg, mon.PV)
			attackedOnce = true
		}
		if didStun {
			monsterStunned = true
			fmt.Println("Le monstre est étourdi pour 1 tour !")
			attackedOnce = true
		}
		playerGuard = didGuard
		if didUse {
			fmt.Printf("PV: %d/%d\n", gs.Joueur.PVActuels, gs.Joueur.PVMax)
			playerHP = gs.Joueur.PVActuels
		}
		if mon.PV <= 0 {
			break
		}
		if monsterStunned {
			fmt.Println("Le monstre est étourdi et ne peut pas attaquer.")
			monsterStunned = false
		} else {
			mincoming, specialName, stunPlayer := monsterAction(mon)
			reduction := gs.Joueur.Endurance / 3
			if reduction > mincoming-1 {
				reduction = mincoming - 1
			}
			mincoming -= reduction
			if playerGuard {
				mincoming = mincoming / 2
			}
			if mincoming < 1 {
				mincoming = 1
			}
			playerHP -= mincoming
			if specialName != "" {
				fmt.Printf("%s utilise %s et inflige %d (PV %d/%d)\n", mon.Nom, specialName, mincoming, max0(playerHP), gs.Joueur.PVMax)
				if stunPlayer {
					fmt.Println("Vous êtes étourdi pour 1 tour !")
				}
			} else {
				fmt.Printf("%s vous touche pour %d (PV %d/%d)\n", mon.Nom, mincoming, max0(playerHP), gs.Joueur.PVMax)
			}

			// Mettre à jour l'attaque du joueur (pour la transformation du loup-garou)
			gs.Joueur.PVActuels = playerHP
			personnage.UpdatePlayerAttack(&gs.Joueur)
			playerStunned = stunPlayer
		}
		playerGuard = false
	}
	if playerHP <= 0 {
		fmt.Println("Vous tombez... Le destin attend une autre tentative.")
		gs.Joueur.PVActuels = gs.Joueur.PVMax
		fmt.Println("(Appuyez sur Entrée pour revenir)")
		attendreEntree()
		return
	}
	gs.Joueur.PVActuels = playerHP
	if fled {
		fmt.Println("Vous avez fui. Aucune récompense.")
		fmt.Println("(Appuyez sur Entrée pour revenir)")
		attendreEntree()
		return
	}
	// Messages de fin avant l'animation
	fmt.Println("Votre mère reprend forme humaine. Ses yeux redeviennent doux. Elle vous serre dans ses bras.")
	fmt.Println("(Appuyez sur Entrée)")
	attendreEntree()
	fmt.Println("Merci, mon enfant... Tu m'as sauvée.")
	fmt.Println("(Appuyez sur Entrée)")
	attendreEntree()

	// Animation de fin
	showEndingAnimation()

	// Récompenses
	reward(gs, 5)
	gainXP(gs, xpForBoss())
	fmt.Println("(Appuyez sur Entrée pour revenir)")
	attendreEntree()
}

// Petite animation de fin: la métamorphe redevient humaine puis FIN
func showEndingAnimation() {
	frames := []string{
		"La métamorphe chancelle...",
		"Sa silhouette vacille...",
		"Ses traits se recomposent...",
		"La bête disparaît, une femme apparaît...",
	}
	for _, f := range frames {
		clearHome()
		clearScreenAll()
		fmt.Println(f)
		fmt.Println("(Appuyez sur Entrée)")
		attendreEntree()
	}
	clearHome()
	clearScreenAll()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("    ███████╗ ██╗ ███╗   ██╗")
	fmt.Println("    ██╔════╝ ██║ ████╗  ██║")
	fmt.Println("    █████╗   ██║ ██╔██╗ ██║")
	fmt.Println("    ██╔══╝   ██║ ██║╚██╗██║")
	fmt.Println("    ██║      ██║ ██║ ╚████║")
	fmt.Println("    ╚═╝      ╚═╝ ╚═╝  ╚═══╝")
	fmt.Println()
	fmt.Println()
	fmt.Println("                    ╔══════════════════════════════════════╗")
	fmt.Println("                    ║              VICTOIRE !              ║")
	fmt.Println("                    ║         L'AVENTURE SE TERMINE        ║")
	fmt.Println("                    ╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("                        (Appuyez sur Entrée)")
	fmt.Println()
	fmt.Println()
	attendreEntree()
}

// XP scaling: more XP for higher tier rooms and boss
func xpForTier(tier int) int {
	switch tier {
	case 1:
		return 30
	case 2:
		return 60
	case 3:
		return 120
	case 4:
		return 240
	default:
		if tier < 1 {
			return 0
		}
		return 30 * tier * tier
	}
}

func xpForBoss() int {
	return 600
}

// Level gating per tier
func requiredLevelForTier(tier int) int {
	switch tier {
	case 1:
		return 1
	case 2:
		return 5
	case 3:
		return 10
	case 4:
		return 15
	default:
		if tier < 1 {
			return 1
		}
		return 1 + (tier-1)*5
	}
}

func renderBattle(gs *GameState, mon Monster, playerHP int, monsterStunned bool, playerGuard bool, playerStunned bool) {
	clearHome()
	clearScreenAll()

	// Statuts textuels
	pStatus := ""
	if playerGuard {
		pStatus = " [Garde]"
	}
	if playerStunned {
		pStatus += " [Étourdi]"
	}
	eStatus := displayStatusEffects(&mon)

	// Valeurs calculées
	pAtk := computePlayerAttack(gs)
	pDef := computePlayerDefense(gs)
	pHP := fmt.Sprintf("%d/%d", max0(playerHP), gs.Joueur.PVMax)
	eHP := fmt.Sprintf("%d", max0(mon.PV))
	weap := gs.Joueur.Attaque
	if weap == "" {
		weap = "(mains nues)"
	}

	// Taille du terminal (fallback = 100)
	termW := 100
	if c := os.Getenv("COLUMNS"); c != "" {
		if v, err := strconv.Atoi(c); err == nil && v > 0 {
			termW = v
		}
	}

	// Colonnes : gauche ~45%, droite le reste
	leftWidth := (termW * 45) / 100
	if leftWidth < 30 {
		leftWidth = 30
	}

	// Ligne 1 : noms (Joueur  |  Ennemi)
	leftLine := fmt.Sprintf("Joueur: %s%s", gs.Joueur.Nom, pStatus)
	rightLine := fmt.Sprintf("Ennemi: %s (%s)%s", mon.Nom, mon.Type, eStatus)
	fmt.Printf("%-*s %s\n", leftWidth, leftLine, rightLine)

	// Ligne 2 : stats (PV / Att / Def)
	leftStats := fmt.Sprintf("PV: %s | Att: %d | Def: %d", pHP, pAtk, pDef)
	rightStats := fmt.Sprintf("PV: %s | Att: %d | Def: %d", eHP, mon.Attaque, mon.Defense)
	fmt.Printf("%-*s %s\n", leftWidth, leftStats, rightStats)

	// Ligne 3 : arme / force (on garde à gauche)
	leftWeapon := fmt.Sprintf("Force: %d | Arme: %s", gs.Joueur.Force, truncate(weap, 20))
	fmt.Printf("%-*s\n\n", leftWidth, leftWeapon)
}

func truncate(s string, n int) string {
	if n <= 0 || len([]rune(s)) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n])
}

// Somme la défense des armures équipées du joueur
func computePlayerDefense(gs *GameState) int {
	if len(gs.Joueur.ArmuresEquipees) == 0 {
		return 0
	}
	keys := []string{
		"CasqueCuir", "CasqueCuirRenforce", "CasqueFer", "CasqueFerRenforce",
		"PlastronCuir", "PlastronCuirRenforce", "PlastronFer", "PlastronFerRenforce",
		"PantalonCuir", "PantalonCuirRenforce", "PantalonFer", "PantalonFerRenforce",
		"BottesCuir", "BottesCuirRenforce", "BottesFer", "BottesFerRenforce",
	}
	total := 0
	for _, k := range keys {
		ar := objet.CreerArmure(k)
		if gs.Joueur.ArmuresEquipees[ar.Nom] {
			total += ar.EffetDefense
		}
	}
	return total
}

// Menu d'action du joueur: retourne (dégâts infligés, a étourdi, a gardé, a utilisé une option/objet)
func playerAction(gs *GameState) (int, bool, bool, bool) {
	// Créer les options d'attaque avec les attaques spéciales
	opts := []string{"Attaquer"}
	for _, attack := range playerSpecialAttacks {
		if attack.CurrentCD <= 0 {
			opts = append(opts, fmt.Sprintf("%s - %s", attack.Nom, attack.Description))
		} else {
			opts = append(opts, fmt.Sprintf("%s (CD: %d) - %s", attack.Nom, attack.CurrentCD, attack.Description))
		}
	}
	opts = append(opts, "Parade", "Potion", "Fuir")

	idx, cancelled := battleSelectWithArrows("Choisissez une action:", opts)
	if cancelled || idx == len(opts)-1 {
		// fuite
		return -1, false, false, false
	}

	// Décrémenter les cooldowns
	for i := range playerSpecialAttacks {
		if playerSpecialAttacks[i].CurrentCD > 0 {
			playerSpecialAttacks[i].CurrentCD--
		}
	}
	switch idx {
	case 0: // Attaquer standard avec petit critique
		base := computePlayerAttack(gs)
		critChance := 10 + gs.Joueur.Agilite
		if critChance > 50 {
			critChance = 50
		}
		if rand.Intn(100) < critChance {
			base = int(float64(base) * 1.5)
			fmt.Println("Coup critique !")
		}
		return base, false, false, false
	case 1, 2, 3, 4, 5: // attaques spéciales
		attackIdx := idx - 1
		if attackIdx < len(playerSpecialAttacks) && playerSpecialAttacks[attackIdx].CurrentCD <= 0 {
			attack := playerSpecialAttacks[attackIdx]
			base := computePlayerAttack(gs)
			damage := base + attack.Damage
			if damage < 0 {
				damage = 0
			}

			// Mettre le cooldown
			playerSpecialAttacks[attackIdx].CurrentCD = attack.Cooldown

			// Retourner les effets
			hasStun := false
			for _, effect := range attack.Effects {
				if effect.Type == "stun" {
					hasStun = true
					break
				}
			}

			return damage, hasStun, false, false
		}
		return 0, false, false, false
	case 6: // parade (index ajusté)
		fmt.Println("Vous vous mettez en garde. Les prochains dégâts seront réduits.")
		return 0, false, true, false
	case 7: // potion (index ajusté)
		return menuPotion(gs)
	default:
		return 0, false, false, false
	}
}

func generateMonster(tier int) Monster {
	// Utilise la fonction de génération de monstres de donjon
	monsterDungeon := monstre.CreerMonstreDungeon(tier)

	// Convertit MonsterDungeon en Monster
	return Monster{
		Nom:     monsterDungeon.Nom,
		PV:      monsterDungeon.PV,
		PVMax:   monsterDungeon.PV,
		Attaque: monsterDungeon.Attaque,
		Defense: monsterDungeon.Defense,
		Type:    monsterDungeon.Type,
		// Initialiser les statuts
		Stunned: monsterDungeon.Stunned, Poisoned: monsterDungeon.Poisoned, Burned: monsterDungeon.Burned, Bleeding: monsterDungeon.Bleeding, Shielded: monsterDungeon.Shielded,
		StunTurns: monsterDungeon.StunTurns, PoisonTurns: monsterDungeon.PoisonTurns, BurnTurns: monsterDungeon.BurnTurns, BleedTurns: monsterDungeon.BleedTurns, ShieldTurns: monsterDungeon.ShieldTurns,
	}
}

func computePlayerAttack(gs *GameState) int {
	name := gs.Joueur.Attaque
	// Inclure les buffs temporaires dans le calcul de la force
	totalForce := gs.Joueur.Force + gs.Joueur.BuffForce
	base := 12 + totalForce/2

	// Bonus de transformation du loup-garou
	transformationBonus := 1.0
	if gs.Joueur.Classe == "Loups-Garou" {
		if float64(gs.Joueur.PVActuels)/float64(gs.Joueur.PVMax) <= 0.3 {
			transformationBonus = 1.5 // +50% d'attaque
			fmt.Println("🐺 Le loup-garou se transforme ! Puissance décuplée !")
		}
	}

	if name == "" {
		return int(float64(base) * transformationBonus)
	}

	// try map to known weapons
	keys := []string{"EpeeRouillee", "EpeeFer", "EpeeMagique", "EpeeCourte", "Hache", "HacheDeCombat", "HacheDeBataille", "ArcBois", "ArcLong", "ArcElfe"}
	for _, k := range keys {
		w := objet.CreerArme(k)
		if strings.EqualFold(name, w.Nom) {
			// Bonus d'agilité pour les armes rapides (épées et arcs)
			agilityBonus := 0
			if strings.Contains(k, "Epee") || strings.Contains(k, "Arc") {
				totalAgilite := gs.Joueur.Agilite + gs.Joueur.BuffAgilite
				agilityBonus = totalAgilite / 3 // +1 dégât tous les 3 points d'agilité
			}

			// Bonus spécial du Bûcheron avec les haches
			bucheronBonus := 0
			if gs.Joueur.Classe == "Bûcheron" && strings.Contains(k, "Hache") {
				bucheronBonus = 5 // +5 dégâts avec les haches
			}

			return int(float64(w.EffetAttaque+gs.Joueur.Force/2+agilityBonus+bucheronBonus) * transformationBonus)
		}
	}

	// Griffes du loup-garou transformé
	if strings.Contains(name, "Griffes") {
		clawDamage := 25 + gs.Joueur.Force/2
		return int(float64(clawDamage) * transformationBonus)
	}

	return int(float64(base+3) * transformationBonus)
}

// Fonction pour gérer les drops d'objets et matériaux
func processDrops(gs *GameState, tier int) {
	if tier > 4 {
		tier = 4 // Boss utilise les drops du niveau 4
	}

	drops, exists := dungeonDrops[tier]
	if !exists {
		return
	}

	fmt.Println("\n🎁 Butin trouvé :")

	// Pas de drops d'équipements ici (uniquement matériaux)

	// Drops de matériaux
	for _, materialDrop := range drops.Materials {
		if rand.Intn(100) < materialDrop.Chance {
			gs.Mats[string(materialDrop.Material)] += materialDrop.Quantity
			fmt.Printf("  📦 %s x%d\n", materialDrop.Material, materialDrop.Quantity)
		}
	}

	// Récompense d'or de base
	baseGold := tier * 50
	gs.Joueur.Argent += baseGold
	fmt.Printf("  💰 %d or\n", baseGold)
}

// Fonctions utilitaires pour identifier les types d'objets
func isWeapon(itemName string) bool {
	weapons := []string{"EpeeRouillee", "EpeeCourte", "EpeeFer", "EpeeMagique", "Hache", "HacheDeCombat", "HacheDeBataille", "ArcBois", "ArcLong", "ArcElfe"}
	for _, w := range weapons {
		if w == itemName {
			return true
		}
	}
	return false
}

func isArmor(itemName string) bool {
	armors := []string{"CasqueCuir", "CasqueCuirRenforce", "CasqueFer", "CasqueFerRenforce", "PlastronCuir", "PlastronCuirRenforce", "PlastronFer", "PlastronFerRenforce", "PantalonCuir", "PantalonCuirRenforce", "PantalonFer", "PantalonFerRenforce", "BottesCuir", "BottesCuirRenforce", "BottesFer", "BottesFerRenforce"}
	for _, a := range armors {
		if a == itemName {
			return true
		}
	}
	return false
}

func reward(gs *GameState, tier int) {
	// Or de base selon le tier
	baseGold := tier * 20
	gs.Joueur.Argent += baseGold
	fmt.Printf("💰 Vous obtenez %d or !\n", baseGold)

	// Items de loot des monstres (taux de drop bas)
	if rand.Intn(100) < 25 { // 25% de chance d'obtenir un item vendable (loot de monstre)
		item := getRandomLootItemForTier(tier)
		gs.Joueur.Inventaire = append(gs.Joueur.Inventaire, item)
		fmt.Printf("🗡️ Vous obtenez %s !\n", item)
	}

	// Matériaux (taux de drop bas)
	for _, mat := range getMaterialsForTier(tier) {
		// taux bas: 30% par matériau listé
		if rand.Intn(100) < 30 {
			gs.Joueur.Materiaux[mat]++
			fmt.Printf("📦 Vous obtenez %s !\n", mat)
		}
	}
}

func gainXP(gs *GameState, amount int) {
	gs.XP += amount
	for gs.XP >= gs.Level*50 {
		gs.XP -= gs.Level * 50
		gs.Level++
		gs.Joueur.Niveau = gs.Level
		// Équilibrage simple: +3 PV max, +1 Force tous les niveaux, +1 Endurance tous les 2 niveaux, +1 Agilité tous les 3 niveaux
		gs.Joueur.PVMax += 3
		gs.Joueur.Force += 1
		if gs.Level%2 == 0 {
			gs.Joueur.Endurance += 1
		}
		if gs.Level%3 == 0 {
			gs.Joueur.Agilite += 1
		}
		gs.Joueur.PVActuels = gs.Joueur.PVMax
		msg := fmt.Sprintf("Niveau %d atteint ! PV max +3, Force +1", gs.Level)
		if gs.Level%2 == 0 {
			msg += ", Endurance +1"
		}
		if gs.Level%3 == 0 {
			msg += ", Agilité +1"
		}
		fmt.Println(msg)
	}
}

func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

// IA du monstre: retourne (dégâts, nom attaque spéciale, étourdir joueur)
func monsterAction(mon Monster) (int, string, bool) {
	// Vérifier si le monstre a des attaques spéciales
	if attacks, exists := monsterSpecialAttacks[mon.Nom]; exists && len(attacks) > 0 {
		// Choisir une attaque spéciale disponible
		availableAttacks := []SpecialAttack{}
		for _, attack := range attacks {
			if attack.CurrentCD <= 0 {
				availableAttacks = append(availableAttacks, attack)
			}
		}

		// Si des attaques spéciales sont disponibles, les utiliser 30% du temps
		if len(availableAttacks) > 0 && rand.Intn(100) < 30 {
			attack := availableAttacks[rand.Intn(len(availableAttacks))]

			// Mettre le cooldown
			for i := range attacks {
				if attacks[i].Nom == attack.Nom {
					attacks[i].CurrentCD = attack.Cooldown
					break
				}
			}

			// Calculer les dégâts
			damage := mon.Attaque + attack.Damage
			if damage < 0 {
				damage = 0
			}

			// Appliquer les effets de statut au joueur
			hasStun := false
			for _, effect := range attack.Effects {
				if effect.Type == "stun" {
					hasStun = true
					break
				}
			}

			return damage, attack.Nom, hasStun
		}
	}

	// Attaques normales
	roll := rand.Intn(100)
	switch {
	case roll < 50:
		return mon.Attaque, "", false
	case roll < 80:
		return mon.Attaque * 12 / 10, "Fracas lourd", false
	default:
		// peur: chance d'étourdir
		return mon.Attaque * 7 / 10, "Peur viscérale", rand.Intn(100) < 35
	}
}

// Fonction utilitaire pour obtenir le maximum
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Fonction pour appliquer les effets de statut
func applyStatusEffects(mon *Monster, maxHP int) {
	// Appliquer les dégâts de statut
	if mon.Poisoned && mon.PoisonTurns > 0 {
		poisonDamage := max(1, maxHP*5/100) // 5% des PV max, minimum 1
		mon.PV -= poisonDamage
		mon.PoisonTurns--
		if mon.PoisonTurns <= 0 {
			mon.Poisoned = false
		}
	}
	if mon.Burned && mon.BurnTurns > 0 {
		burnDamage := max(1, maxHP*3/100) // 3% des PV max, minimum 1
		mon.PV -= burnDamage
		mon.BurnTurns--
		if mon.BurnTurns <= 0 {
			mon.Burned = false
		}
	}
	if mon.Bleeding && mon.BleedTurns > 0 {
		bleedDamage := max(1, maxHP*4/100) // 4% des PV max, minimum 1
		mon.PV -= bleedDamage
		mon.BleedTurns--
		if mon.BleedTurns <= 0 {
			mon.Bleeding = false
		}
	}

	// Gérer l'étourdissement
	if mon.Stunned && mon.StunTurns > 0 {
		mon.StunTurns--
		if mon.StunTurns <= 0 {
			mon.Stunned = false
		}
	}

	// Gérer le bouclier
	if mon.Shielded && mon.ShieldTurns > 0 {
		mon.ShieldTurns--
		if mon.ShieldTurns <= 0 {
			mon.Shielded = false
		}
	}
}

// Fonction pour afficher les statuts actifs
func displayStatusEffects(mon *Monster) string {
	statuses := []string{}
	if mon.Stunned {
		statuses = append(statuses, "[Étourdi]")
	}
	if mon.Poisoned {
		statuses = append(statuses, "[Empoisonné]")
	}
	if mon.Burned {
		statuses = append(statuses, "[Brûlé]")
	}
	if mon.Bleeding {
		statuses = append(statuses, "[Saigne]")
	}
	if mon.Shielded {
		statuses = append(statuses, "[Bouclier]")
	}
	if len(statuses) == 0 {
		return ""
	}
	return " " + strings.Join(statuses, " ")
}

func battleSelectWithArrows(header string, options []string) (int, bool) {
	if err := keyboard.Open(); err != nil {
		return 0, false
	}
	defer keyboard.Close()
	index := 0

	// affichage initial
	fmt.Println()
	for i, opt := range options {
		prefix := "  "
		if i == index {
			prefix = "> "
		}
		fmt.Printf("%s%s\n", prefix, opt)
	}
	fmt.Println() // ligne vide après le menu (curseur "tampon")

	lines := len(options)

	for {
		ch, key, err := keyboard.GetKey()
		if err != nil {
			return index, false
		}

		// remonte et efface le menu
		for i := 0; i < lines+1; i++ { // +1 pour la ligne vide
			fmt.Print("\033[A\033[2K")
		}

		// navigation
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
			// redraw final avant de sortir
			fmt.Println()
			for i, opt := range options {
				prefix := "  "
				if i == index {
					prefix = "> "
				}
				fmt.Printf("%s%s\n", prefix, opt)
			}
			fmt.Println()
			return index, false
		case keyboard.KeyEsc:
			return 0, true
		default:
			if ch == 'q' || ch == 'Q' {
				return 0, true
			}
			if ch == '\r' || ch == '\n' {
				fmt.Println()
				for i, opt := range options {
					prefix := "  "
					if i == index {
						prefix = "> "
					}
					fmt.Printf("%s%s\n", prefix, opt)
				}
				fmt.Println()
				return index, false
			}
		}

		// redessine le menu
		for i, opt := range options {
			prefix := "  "
			if i == index {
				prefix = "> "
			}
			fmt.Printf("%s%s\n", prefix, opt)
		}
		fmt.Println() // ligne vide de sécurité
	}
}

func printEnemyStats(mobs []Monster) {
	// trouver la largeur max du nom
	maxLen := 0
	for _, mob := range mobs {
		if len(mob.Nom) > maxLen {
			maxLen = len(mob.Nom)
		}
	}

	// afficher chaque mob aligné
	for _, mob := range mobs {
		fmt.Printf("%-*s HP:%-5d ATK:%-5d DEF:%-5d\n",
			maxLen, mob.Nom, mob.PV, mob.Attaque, mob.Defense)
	}
}

// Menu de sélection des potions en combat
func menuPotion(gs *GameState) (int, bool, bool, bool) {
	// Créer la liste des potions disponibles dans l'inventaire
	potionsDisponibles := []string{}
	descriptions := []string{}

	// Vérifier quelles potions sont dans l'inventaire
	potionsPossibles := []string{"potion", "potion majeure", "potion force", "potion agilite", "potion endurance", "antidote", "elixir vie"}

	for _, potion := range potionsPossibles {
		// Vérifier si le joueur a cette potion
		for _, item := range gs.Joueur.Inventaire {
			if strings.EqualFold(item, potion) {
				potionsDisponibles = append(potionsDisponibles, potion)
				// Ajouter la description
				switch potion {
				case "potion":
					descriptions = append(descriptions, "Potion (+20 PV)")
				case "potion majeure":
					descriptions = append(descriptions, "Potion majeure (+50 PV)")
				case "potion force":
					descriptions = append(descriptions, "Potion de force (+2 Force, 3 combats)")
				case "potion agilite":
					descriptions = append(descriptions, "Potion d'agilité (+2 Agilité, 3 combats)")
				case "potion endurance":
					descriptions = append(descriptions, "Potion d'endurance (+2 Endurance, 3 combats)")
				case "antidote":
					descriptions = append(descriptions, "Antidote (Guérit statuts)")
				case "elixir vie":
					descriptions = append(descriptions, "Élixir de vie (+100 PV)")
				}
				break
			}
		}
	}

	// Si aucune potion disponible
	if len(potionsDisponibles) == 0 {
		fmt.Println("❌ Vous n'avez aucune potion !")
		attendreEntree()
		return 0, false, false, false
	}

	// Afficher le menu de sélection des potions
	idx, cancelled := battleSelectWithArrows("Choisissez une potion:", descriptions)
	if cancelled {
		return 0, false, false, false
	}

	// Utiliser la potion sélectionnée
	potionChoisie := potionsDisponibles[idx]
	before := gs.Joueur.PVActuels

	// Appeler la fonction d'utilisation appropriée
	switch potionChoisie {
	case "potion":
		utiliserPotion(&gs.Joueur)
	case "potion majeure":
		utiliserPotionMajeure(&gs.Joueur)
	case "potion force":
		utiliserPotionForce(&gs.Joueur)
	case "potion agilite":
		utiliserPotionAgilite(&gs.Joueur)
	case "potion endurance":
		utiliserPotionEndurance(&gs.Joueur)
	case "antidote":
		utiliserAntidote(&gs.Joueur)
	case "elixir vie":
		utiliserElixirVie(&gs.Joueur)
	}

	// Retourner si une potion a été utilisée (pour les potions de soin)
	used := gs.Joueur.PVActuels > before || potionChoisie == "potion force" || potionChoisie == "potion agilite" || potionChoisie == "potion endurance" || potionChoisie == "antidote"
	return 0, false, false, used
}

// Gère la diminution des buffs temporaires après un combat
func gererBuffsApresCombat(gs *GameState) {
	if gs.Joueur.BuffCombats > 0 {
		gs.Joueur.BuffCombats--
		if gs.Joueur.BuffCombats <= 0 {
			// Les buffs expirent
			if gs.Joueur.BuffForce > 0 {
				fmt.Printf("💪 L'effet de la potion de force s'estompe...\n")
				gs.Joueur.BuffForce = 0
			}
			if gs.Joueur.BuffAgilite > 0 {
				fmt.Printf("🏃 L'effet de la potion d'agilité s'estompe...\n")
				gs.Joueur.BuffAgilite = 0
			}
			if gs.Joueur.BuffEndurance > 0 {
				fmt.Printf("❤️ L'effet de la potion d'endurance s'estompe...\n")
				// Restaurer les PV max à la normale
				gs.Joueur.PVMax -= gs.Joueur.BuffEndurance * 10
				if gs.Joueur.PVActuels > gs.Joueur.PVMax {
					gs.Joueur.PVActuels = gs.Joueur.PVMax
				}
				gs.Joueur.BuffEndurance = 0
			}
		}
	}
}

// Matériaux par tier de donjon
func getMaterialsForTier(tier int) []string {
	switch tier {
	case 1:
		return []string{"cuir", "pierre", "bois"}
	case 2:
		return []string{"cuir renforcé", "fer", "bois dur"}
	case 3:
		return []string{"fer renforcé", "pierre précieuse", "os ancien"}
	case 4:
		return []string{"fer renforcé", "gemme", "écailles"}
	case 5: // Boss
		return []string{"gemme de pouvoir", "cristal de mana", "écailles de dragon"}
	default:
		return []string{"cuir", "pierre"}
	}
}

// Items de loot des monstres par tier
func getRandomLootItemForTier(tier int) string {
	lootItems := map[int][]string{
		1: {"Griffes souillées", "Massue brute"}, // Tier 1: items basiques
		2: {"Lance brisée", "Épée osseuse", "Hache tronquée"},
		3: {"Glaive sauvage", "Masse rituelle", "Faux de brume"},
		4: {"Glaive sauvage", "Masse rituelle", "Faux de brume"}, // Tier 4: items de qualité
		5: {"Glaive sauvage", "Masse rituelle", "Faux de brume"}, // Boss: items de qualité
	}

	if items, exists := lootItems[tier]; exists {
		return items[rand.Intn(len(items))]
	}
	return "Griffes souillées"
}
