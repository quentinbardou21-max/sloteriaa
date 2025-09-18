package main

import (
	"fmt"
	"strings"

	"sloteriaa/internal/personnage"
	"sloteriaa/struct/forgeron"
	"sloteriaa/struct/objet"
)

type GameState struct {
	Joueur personnage.Personnage
	Mats   map[string]int
	XP     int
	Level  int
}

func StartGameNew() {
	// Utiliser la fonction de création de personnage qui gère correctement l'inventaire
	p := personnage.CreationPersonnage()

	// Admin mode if name is "costa"
	isAdmin := strings.EqualFold(p.Nom, "costa")
	if isAdmin {
		p.Force = 9999
		p.Endurance = 9999
		p.Agilite = 9999
		p.PVMax = 999999
		p.PVActuels = p.PVMax
		// Equip best weapon
		weap := objet.CreerArme("EpeeMagique")
		p.Attaque = weap.Nom
		p.Inventaire = append(p.Inventaire, weap.Nom)
		// Equip strong armors
		if p.ArmuresEquipees == nil {
			p.ArmuresEquipees = make(map[string]bool)
		}
		for _, k := range []string{"CasqueFerRenforce", "PlastronFerRenforce", "PantalonFerRenforce", "BottesFerRenforce"} {
			ar := objet.CreerArmure(k)
			if ar.Nom != "" {
				p.ArmuresEquipees[ar.Nom] = true
				p.Inventaire = append(p.Inventaire, ar.Nom)
			}
		}
		p.Argent = 9999999
	}

	gs := GameState{
		Joueur: p,
		Mats: map[string]int{
			"Or":                            10000,
			string(forgeron.Fer):            8,
			string(forgeron.Bois):           6,
			string(forgeron.Cuir):           4,
			string(forgeron.EssenceMagique): 2,
		},
		XP:    0,
		Level: 1,
	}
	if isAdmin {
		gs.Level = 20
		gs.XP = 0
	}
	// Boost materials/gold further for admin
	if isAdmin {
		gs.Mats["Or"] = 9999999
		gs.Mats[string(forgeron.Fer)] = 9999
		gs.Mats[string(forgeron.Bois)] = 9999
		gs.Mats[string(forgeron.Cuir)] = 9999
		gs.Mats[string(forgeron.EssenceMagique)] = 9999
	}

	showIntroLore(&gs)
	if !isAdmin {
		giveStartingEquipment(&gs)
	}
	// Auto-save right after creating the new game
	_ = SaveGame(&gs)
	enterAltScreen()
	defer exitAltScreen()
	worldLoop(&gs)
}

func StartGameFromSave(gs *GameState) {
	if gs == nil {
		return
	}
	// Ensure we re-save on resume to keep schema up to date
	_ = SaveGame(gs)
	enterAltScreen()
	defer exitAltScreen()
	worldLoop(gs)
}

// --- Conversion helpers for materials ---
func matsToForgeron(inv map[string]int) forgeron.InventaireMateriaux {
	res := make(forgeron.InventaireMateriaux)
	for k, v := range inv {
		res[forgeron.Materiau(k)] = v
	}
	return res
}

func forgeronToMats(inv forgeron.InventaireMateriaux) map[string]int {
	res := make(map[string]int)
	for k, v := range inv {
		res[string(k)] = v
	}
	return res
}

func showIntroLore(gs *GameState) {
	// Cinématique d'intro par scènes, validation manuelle
	scenes := [][]string{
		{
			"… Silence du matin …",
			"Vous vous réveillez. Aujourd'hui, vous avez 18 ans.",
			"Le monde semble plus lourd, et plus vaste à la fois.",
		},
		{
			"Votre père vous attend, ému.",
			"Il vous tend une petite boîte en bois poli…",
			"À l'intérieur: une arme, et 100 pièces d'or.",
		},
		{
			"Sous le couvercle, une lettre soigneusement pliée…",
			"Lettre de votre mère:",
			"\"Mon enfant, je suis une métamorphe. Pour ne mettre personne en danger,\"",
			"\"je me suis exilée dans le Donjon. C'est risqué.\"",
			"\"Si je ne suis pas revenue, viens me sauver dans le Donjon. — Maman\"",
		},
		{
			"Votre père baisse la tête.",
			"— Je n'ai jamais cessé d'espérer. Va… ramène-la.",
		},
	}
	for _, lines := range scenes {
		clearHome()
		clearScreenAll()
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()
		fmt.Println("(Appuyez sur Entrée)")
		attendreEntree()
	}
}

func giveStartingEquipment(gs *GameState) {
	// L'équipement de départ est déjà géré dans personnage.go
	// Pas besoin d'ajouter d'arme ici pour éviter la duplication
}

func worldLoop(gs *GameState) {
	for {
		header := fmt.Sprintf("Ville de Sloteria — Niveau %d (XP %d) — Or %d", gs.Level, gs.XP, gs.Joueur.Argent)
		opts := []string{"Aller à la Forge", "Aller au Marché", "Entrer dans le Donjon", "Inventaire", "Stats du personnage", "Sauvegarder", "Quitter le jeu"}
		choice, cancelled := selectWithArrows(header, opts)
		if cancelled {
			showCursor() // Réaffiche le curseur avant de quitter
			fmt.Println("À bientôt !")
			return
		}
		switch choice {
		case 0:
			EnterForge(gs)
		case 1:
			EnterShop(gs)
		case 2:
			EnterDungeon(gs)
		case 3:
			afficherInventaireInteractif(&gs.Joueur)
		case 4:
			clearScreen()
			personnage.AfficherInfos(gs.Joueur)
			attendreEntree()
			clearScreen()
		case 5:
			if err := SaveGame(gs); err != nil {
				fmt.Printf("Erreur de sauvegarde: %s\n", err)
			} else {
				fmt.Println("Sauvegarde effectuée.")
			}
			attendreEntree()
		case 6:
			showCursor() // Réaffiche le curseur avant de quitter
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}
