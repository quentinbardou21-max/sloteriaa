package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/mattn/go-tty"
	"github.com/quent/sloteriaa/internal/personnage"
	"github.com/quent/sloteriaa/struct/monstre"
	"github.com/quent/sloteriaa/struct/objet"
)

// --- Couleurs ANSI ---
const (
	ansiReset   = "\u001b[0m"
	ansiRed     = "\u001b[31m"
	ansiGreen   = "\u001b[32m"
	ansiYellow  = "\u001b[33m"
	ansiCyan    = "\u001b[36m"
	ansiMagenta = "\u001b[35m"
	ansiBold    = "\u001b[1m"
)

func CombatTour(p *personnage.Personnage, m *monstre.Monstre) {
	fmt.Printf("%s attaque %s avec %s !\n", p.Nom, m.Nom, p.Attaque)

	// dégâts de base (exemple simplifié : selon l’arme ou classe du perso)
	degats := 20

	// on applique les dégâts
	m.HPActuels -= degats
	if m.HPActuels < 0 {
		m.HPActuels = 0
	}

	fmt.Printf("%s subit %d dégâts ! PV restants : %d/%d\n", m.Nom, degats, m.HPActuels, m.HPMax)

	// Vérif si le monstre est mort
	if m.HPActuels == 0 {
		fmt.Printf("%s est vaincu ! 🎉\n", m.Nom)
	}
}

// Combat oppose un personnage à un monstre en tenant compte
// de la puissance de leurs armes et de leurs armures.
// - armeJ: arme équipée par le joueur (objet.Arme)
// - armuresJ: armures équipées par le joueur ([]objet.Armure)
// Retourne true si le joueur gagne, false sinon.
func Combat(p *personnage.Personnage, armeJ objet.Arme, armuresJ []objet.Armure, m *monstre.Monstre) bool {
	seedOnce()
	// Calcul des statistiques effectives
	joueurAttaque := armeJ.EffetAttaque
	joueurDefense := sommeDefense(armuresJ)

	monstreAttaque := m.Attaque
	if m.PeutAvoirArme {
		monstreAttaque += m.Arme.EffetAttaque
	}
	monstreDefense := m.Defense + sommeDefense(m.Armures)

	fmt.Printf("\n⚔️  Combat: %s vs %s\n", p.Nom, m.Nom)
	fmt.Printf("Joueur → Atk:%d Def:%d | Monstre → Atk:%d Def:%d\n\n", joueurAttaque, joueurDefense, monstreAttaque, monstreDefense)

	// Boucle de combat tour par tour (joueur commence)
	for {
		// Tour du joueur
		if jetTouche(joueurAttaque, monstreDefense) {
			degatsJ := max(1, joueurAttaque-monstreDefense)
			m.HPActuels -= degatsJ
			if m.HPActuels < 0 {
				m.HPActuels = 0
			}
			fmt.Printf("%s frappe (%d) → %s: %d/%d\n", p.Nom, degatsJ, m.Nom, m.HPActuels, m.HPMax)
		} else {
			fmt.Printf("%s manque son attaque !\n", p.Nom)
		}
		if m.HPActuels == 0 {
			fmt.Printf("\n🎉 %s est vaincu !\n", m.Nom)
			return true
		}

		// Tour du monstre
		if jetTouche(monstreAttaque, joueurDefense) {
			degatsM := max(1, monstreAttaque-joueurDefense)
			p.PVActuels -= degatsM
			if p.PVActuels < 0 {
				p.PVActuels = 0
			}
			fmt.Printf("%s riposte (%d) → %s: %d/%d\n", m.Nom, degatsM, p.Nom, p.PVActuels, p.PVMax)
		} else {
			fmt.Printf("%s manque son attaque !\n", m.Nom)
		}
		if p.PVActuels == 0 {
			fmt.Printf("\n💀 %s est vaincu...\n", p.Nom)
			return false
		}
	}
}

func sommeDefense(armures []objet.Armure) int {
	total := 0
	for _, a := range armures {
		total += a.EffetDefense
	}
	return total
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --- Aléatoire / Jet pour toucher ---
var rngSeeded = false

func seedOnce() {
	if !rngSeeded {
		rand.Seed(time.Now().UnixNano())
		rngSeeded = true
	}
}

// jetTouche calcule une probabilité de toucher basée sur att vs def,
// clampée entre 10% et 95%, puis fait un jet aléatoire.
func jetTouche(attaque, defense int) bool {
	base := 0.75 // 75% de base
	delta := float64(attaque-defense) * 0.01
	chance := math.Max(0.10, math.Min(0.95, base+delta))
	return rand.Float64() < chance
}

// --- Combat interactif (TUI clavier) ---

type actionCombat struct {
	label        string
	modAtk       float64 // multiplicateur de dégâts sur l'attaque
	modHit       float64 // delta de probabilité de toucher (-0.15 ... +0.15)
	typeDefensif bool    // si true: pas de dégâts, applique un bonus défensif ce tour
	defBuffMul   float64 // multiplicateur pour le buff défense (par ex. 0.5 = +50%)
}

// RunCombatInteractive lance un combat avec sélection d'actions au clavier.
// Retourne true si le joueur gagne.
func RunCombatInteractive(p *personnage.Personnage, armeJ objet.Arme, armuresJ []objet.Armure, m *monstre.Monstre) bool {
	seedOnce()

	// Stats de base (liées au personnage)
	joueurAttaqueBase := calculerAttaqueJoueur(p, armeJ)
	joueurDefenseBase := calculerDefenseJoueur(p, armuresJ)

	// Options d'action selon l'arme
	actions := actionsForArme(armeJ)

	// TTY
	t, err := tty.Open()
	if err != nil {
		fmt.Println("Impossible d'initialiser le TTY:", err)
		// fallback: combat auto
		return Combat(p, armeJ, armuresJ, m)
	}
	defer t.Close()

	selection := 0
	defenseBuff := 0

	// Effets spéciaux persistants côté monstre
	bleedTurns := 0
	bleedPerTurn := 0
	shredTurns := 0
	shredAmount := 0 // réduction temporaire de défense
	monsterStunned := 0

	// Effets/états liés à l'épée
	riposteReady := false // activé si Parade avec épée; contre-attaque si l'ennemi rate

	// Mode simple par défaut pour stabilité d'affichage
	simpleMode := false
	headerPrinted := false

	for {
		// Appliquer effets persistants en début de tour
		if bleedTurns > 0 && m.HPActuels > 0 {
			m.HPActuels -= bleedPerTurn
			if m.HPActuels < 0 {
				m.HPActuels = 0
			}
			bleedTurns--
		}
		if shredTurns > 0 {
			shredTurns--
			if shredTurns == 0 {
				shredAmount = 0
			}
		}

		// Appliquer la transformation (régénération de vie pour loup-garou)
		appliquerTransformation(p)

		// Rendu écran
		defMonstreCourante := max(0, m.Defense+sommeDefense(m.Armures)-shredAmount)
		armeAffichage := armeJ.Nom
		if armeJ.EffetAttaque == 0 {
			// Afficher l'arme par défaut du personnage
			armeAffichage = getArmeParDefautNom(p)
		}

		if !simpleMode {
			clearScreenTUI()
			fmt.Println("=== Combat ===")
			fmt.Printf("Joueur: %s Nv %d  PV %d/%d  Arme %s", p.Nom, p.Niveau, p.PVActuels, p.PVMax, armeAffichage)
			if estEnFormeTransformee(p) {
				fmt.Printf(" \033[1;31m[🐺 LOUP-GAROU TRANSFORMÉ! 🐺]\033[0m")
			}
			fmt.Println()
			fmt.Printf("Monstre: %s Nv %d  PV %d/%d  Def %d\n", m.Nom, m.Niveau, m.HPActuels, m.HPMax, defMonstreCourante)
		} else if !headerPrinted {
			fmt.Println("=== Combat (mode simple) ===")
			fmt.Printf("Joueur: %s Nv %d  PV %d/%d  Arme %s", p.Nom, p.Niveau, p.PVActuels, p.PVMax, armeAffichage)
			if estEnFormeTransformee(p) {
				fmt.Printf(" \033[1;31m[🐺 LOUP-GAROU TRANSFORMÉ! 🐺]\033[0m")
			}
			fmt.Println()
			fmt.Printf("Monstre: %s Nv %d  PV %d/%d  Def %d\n", m.Nom, m.Niveau, m.HPActuels, m.HPMax, defMonstreCourante)
			headerPrinted = true
		}
		// Tags d'état
		if !simpleMode {
			if defenseBuff > 0 {
				fmt.Printf("[Défense +%d active pour ce tour]\n", defenseBuff)
			}
			if bleedTurns > 0 {
				fmt.Printf("[Saignement: -%d PV pendant %d tour(s)]\n", bleedPerTurn, bleedTurns)
			}
			if shredAmount > 0 && shredTurns > 0 {
				fmt.Printf("[Armure percée: -%d DEF (%d tour(s) restant)]\n", shredAmount, shredTurns)
			}
			if monsterStunned > 0 {
				fmt.Printf("[Étourdisssement: %s perd son prochain tour]\n", m.Nom)
			}
			if riposteReady {
				fmt.Printf("[Parade prête: contre-attaque si l'ennemi rate]\n")
			}
		}

		fmt.Println("\nChoisissez une action:")
		for i, a := range actions {
			prefix := "  "
			if i == selection {
				prefix = "> "
			}
			fmt.Printf("%s%s\n", prefix, a.label)
		}
		fmt.Println("\n[Flèches ↑/↓ pour naviguer, Entrée pour valider]")

		key := readKeyTTY(t)
		if key == "up" && selection > 0 {
			selection--
			continue
		}
		if key == "down" && selection < len(actions)-1 {
			selection++
			continue
		}
		if key != "enter" {
			continue
		}

		// Appliquer l'action choisie
		choix := actions[selection]
		joueurAttaque := joueurAttaqueBase
		joueurDefense := joueurDefenseBase + defenseBuff
		defenseBuff = 0 // reset

		// Vérifier la transformation en loup-garou
		etaitTransforme := estEnFormeTransformee(p)
		// Recalculer les stats après la transformation
		joueurAttaque = calculerAttaqueJoueur(p, armeJ)
		joueurDefense = calculerDefenseJoueur(p, armuresJ)
		estTransforme := estEnFormeTransformee(p)

		// Message de transformation
		if !etaitTransforme && estTransforme {
			fmt.Println()
			// Couleurs ANSI : Rouge sur fond jaune pour le cadre, texte en gras
			fmt.Println("\033[1;43;31m╔══════════════════════════════════════════════════════════════╗\033[0m")
			fmt.Printf("\033[1;43;31m║\033[0m  \033[1;33m🐺 %s se transforme en LOUP-GAROU! 🐺\033[0m  \033[1;43;31m║\033[0m\n", p.Nom)
			fmt.Println("\033[1;43;31m║\033[0m  \033[1;33mSes griffes deviennent plus puissantes et il régénère!\033[0m  \033[1;43;31m║\033[0m")
			fmt.Println("\033[1;43;31m╚══════════════════════════════════════════════════════════════╝\033[0m")
			fmt.Println()
		}

		// Bonus de classe pour certaines actions
		bonusClasse := calculerBonusClasse(p, choix)
		joueurAttaque += bonusClasse

		// Tour du joueur
		if choix.typeDefensif {
			// Buff défense jusqu'à la fin du tour ennemi
			mul := choix.defBuffMul
			if mul <= 0 {
				mul = 0.5 // défaut: +50%
			}
			defenseBuff = int(float64(joueurDefenseBase) * mul)
			fmt.Printf("\n%s se met en garde ! Défense augmentée pour ce tour (+%d).\n", p.Nom, defenseBuff)
			// Si épée, activer la riposte potentielle
			if strings.Contains(strings.ToLower(armeJ.Nom), "épée") || strings.Contains(strings.ToLower(armeJ.Nom), "epee") {
				riposteReady = true
			}
		} else {
			// Modificateurs d'attaque et de précision
			joueurAttaque = int(float64(joueurAttaque) * choix.modAtk)
			defCible := max(0, m.Defense+sommeDefense(m.Armures)-shredAmount)
			if jetToucheMod(joueurAttaque, defCible, choix.modHit) {
				degatsJ := max(1, joueurAttaque-defCible)
				m.HPActuels -= degatsJ
				if m.HPActuels < 0 {
					m.HPActuels = 0
				}
				fmt.Printf("\n%s frappe (%d) → %s: %d/%d\n", p.Nom, degatsJ, m.Nom, m.HPActuels, m.HPMax)

				// Effets spéciaux selon l'arme du joueur
				nomArme := strings.ToLower(armeJ.Nom)
				if strings.Contains(nomArme, "hache") {
					// Saignement: dégâts sur la durée, avec chance d'application
					if choix.modAtk > 1.0 { // lourde
						if rand.Float64() < 0.60 {
							bleedPerTurn = max(3, joueurAttaque/8)
							bleedTurns = 3
							fmt.Printf("→ Effet: %s inflige un saignement (%d PV pendant %d tour(s)).\n", p.Nom, bleedPerTurn, bleedTurns)
						}
					} else { // rapide/précise
						if rand.Float64() < 0.35 {
							bleedPerTurn = max(2, joueurAttaque/12)
							bleedTurns = 2
							fmt.Printf("→ Effet: %s inflige un saignement (%d PV pendant %d tour(s)).\n", p.Nom, bleedPerTurn, bleedTurns)
						}
					}
				} else if strings.Contains(nomArme, "arc") {
					// Perçage: réduit la défense de la cible temporairement, avec chance
					if choix.modAtk > 1.0 { // barrage
						if rand.Float64() < 0.55 {
							shredAmount = max(2, (m.Defense+sommeDefense(m.Armures))/5) // ~20%
							shredTurns = 2
							fmt.Printf("→ Effet: Armure de %s percée (-%d DEF, %d tour(s)).\n", m.Nom, shredAmount, shredTurns)
						}
					} else { // précis
						if rand.Float64() < 0.35 {
							shredAmount = max(1, (m.Defense+sommeDefense(m.Armures))/8) // ~12%
							shredTurns = 3
							fmt.Printf("→ Effet: Armure de %s percée (-%d DEF, %d tour(s)).\n", m.Nom, shredAmount, shredTurns)
						}
					}
				} else if strings.Contains(nomArme, "épée") || strings.Contains(nomArme, "epee") {
					// Brise-garde: attaque lourde peut étourdir (chance)
					if choix.modAtk > 1.0 { // lourde
						if rand.Float64() < 0.40 {
							monsterStunned = 1
							fmt.Printf("→ Effet: Brise-garde ! %s est étourdi et perd son prochain tour.\n", m.Nom)
						}
					}
				}
			} else {
				fmt.Printf("\n%s manque son attaque !\n", p.Nom)
			}
			if m.HPActuels == 0 {
				fmt.Printf("\n🎉 %s est vaincu !\n", m.Nom)
				return true
			}
		}

		// Tour du monstre
		if monsterStunned > 0 {
			fmt.Printf("%s est étourdi et ne peut pas agir !\n", m.Nom)
			monsterStunned = 0
		} else {
			monstreAttaque := m.Attaque
			if m.PeutAvoirArme {
				monstreAttaque += m.Arme.EffetAttaque
			}
			if jetTouche(monstreAttaque, joueurDefense) {
				degatsM := max(1, monstreAttaque-joueurDefense)
				p.PVActuels -= degatsM
				if p.PVActuels < 0 {
					p.PVActuels = 0
				}
				fmt.Printf("%s riposte (%d) → %s: %d/%d\n", m.Nom, degatsM, p.Nom, p.PVActuels, p.PVMax)
			} else {
				fmt.Printf("%s manque son attaque !\n", m.Nom)
				// Riposte si prête (parade à l'épée)
				if riposteReady && (strings.Contains(strings.ToLower(armeJ.Nom), "épée") || strings.Contains(strings.ToLower(armeJ.Nom), "epee")) {
					riposteReady = false
					contre := max(1, joueurAttaqueBase/2)
					m.HPActuels -= contre
					if m.HPActuels < 0 {
						m.HPActuels = 0
					}
					fmt.Printf("→ Riposte ! %s contre-attaque (%d) → %s: %d/%d\n", p.Nom, contre, m.Nom, m.HPActuels, m.HPMax)
					if m.HPActuels == 0 {
						fmt.Printf("\n🎉 %s est vaincu !\n", m.Nom)
						return true
					}
				}
			}
		}
		if p.PVActuels == 0 {
			fmt.Printf("\n💀 %s est vaincu...\n", p.Nom)
			return false
		}

		messagePauseTTY(t, "(Entrée pour continuer le combat)")
	}
}

func valOrZero(v int, ok bool) int {
	if ok {
		return v
	}
	return 0
}

func jetToucheMod(attaque, defense int, delta float64) bool {
	base := 0.75
	mod := float64(attaque-defense) * 0.01
	chance := math.Max(0.10, math.Min(0.95, base+mod+delta))
	return rand.Float64() < chance
}

// barreVieColor rend une barre de vie colorée 20 colonnes, avec compte PV
func barreVieColor(actuels, max int) string {
	if max <= 0 {
		return "[??????????]"
	}
	if actuels < 0 {
		actuels = 0
	}
	if actuels > max {
		actuels = max
	}
	barLength := 20
	filled := (actuels * barLength) / max
	empty := barLength - filled
	col := ansiGreen
	ratio := float64(actuels) / float64(max)
	if ratio < 0.33 {
		col = ansiRed
	} else if ratio < 0.66 {
		col = ansiYellow
	}
	return "[" + col + strings.Repeat("█", filled) + ansiReset + strings.Repeat("·", empty) + "]" +
		fmt.Sprintf(" %d/%d", actuels, max)
}

func readKeyTTY(t *tty.TTY) string {
	r, err := t.ReadRune()
	if err != nil {
		return ""
	}
	if r == 0x1b { // ESC
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
	if r == '\r' || r == '\n' {
		return "enter"
	}
	return strings.ToLower(string(r))
}

func messagePauseTTY(t *tty.TTY, msg string) {
	fmt.Println(msg)
	// attendre Entrée
	for {
		r, _ := t.ReadRune()
		if r == '\r' || r == '\n' {
			return
		}
	}
}

func clearScreenTUI() {
	fmt.Print("\033[H\033[2J")
}

// renderTopPanels dessine deux encadrés côte-à-côte (perso à gauche, monstre à droite)
func renderTopPanels(p *personnage.Personnage, armeJ objet.Arme, armuresJ []objet.Armure, atkJ, defJ int, m *monstre.Monstre, defM int) {
	left := buildPanelCyan(
		fmt.Sprintf("%s (%s)\nNv %d\nPV: %s\nArme: %s (Atk %d)\nArmures: %s",
			p.Nom, p.Classe, p.Niveau, barreVieColor(p.PVActuels, p.PVMax), armeJ.Nom, armeJ.EffetAttaque, formatArmures(armuresJ)),
	)
	right := buildPanelRed(
		fmt.Sprintf("%s\nNv %d\nPV: %s\nArme: %s\nDef: %d\nArmures: %s",
			m.Nom, m.Niveau, barreVieColor(m.HPActuels, m.HPMax),
			func() string {
				if m.PeutAvoirArme {
					return fmt.Sprintf("%s (Atk %d)", m.Arme.Nom, m.Arme.EffetAttaque)
				} else {
					return "Aucune"
				}
			}(),
			defM, formatArmures(m.Armures)),
	)

	// Aligner ligne par ligne
	linesL := strings.Split(left, "\n")
	linesR := strings.Split(right, "\n")
	maxLines := max(len(linesL), len(linesR))
	gap := 6
	for i := 0; i < maxLines; i++ {
		l := ""
		if i < len(linesL) {
			l = linesL[i]
		}
		r := ""
		if i < len(linesR) {
			r = linesR[i]
		}
		pad := gap
		if len(l) < 40 {
			pad += 40 - len(l)
		}
		fmt.Println(l + strings.Repeat(" ", pad) + r)
	}
}

func buildPanelCyan(content string) string {
	return buildPanel(content, ansiCyan)
}

func buildPanelRed(content string) string {
	return buildPanel(content, ansiRed)
}

func buildPanel(content, color string) string {
	lines := strings.Split(content, "\n")
	width := 0
	for _, ln := range lines {
		if len(ln) > width {
			width = len(ln)
		}
	}
	width += 2
	top := color + "+" + strings.Repeat("-", width) + "+" + ansiReset
	bot := color + "+" + strings.Repeat("-", width) + "+" + ansiReset
	out := strings.Builder{}
	out.WriteString(top)
	out.WriteString("\n")
	for _, ln := range lines {
		padding := width - len(ln)
		out.WriteString(color + "| " + ansiReset + ln + strings.Repeat(" ", padding-1) + color + "|" + ansiReset)
		out.WriteString("\n")
	}
	out.WriteString(bot)
	return out.String()
}

// formatArmures compacte la liste des armures pour l'affichage
func formatArmures(armures []objet.Armure) string {
	if len(armures) == 0 {
		return "Aucune"
	}
	b := strings.Builder{}
	for i, a := range armures {
		if i > 0 {
			b.WriteString(", ")
		}
		// Nom court + valeur DEF
		b.WriteString(a.Nom)
		b.WriteString(" (")
		b.WriteString(fmt.Sprintf("DEF %d", a.EffetDefense))
		b.WriteString(")")
	}
	return b.String()
}

// calculerAttaqueJoueur calcule l'attaque du joueur basée sur son niveau, classe et arme
func calculerAttaqueJoueur(p *personnage.Personnage, arme objet.Arme) int {
	// Attaque de base selon la classe
	attaqueClasse := 0
	switch p.Classe {
	case "Humain":
		attaqueClasse = 15
	case "Loups-Garou":
		// Transformation en loup-garou si PV < 30%
		if p.PVActuels < (p.PVMax * 30 / 100) {
			attaqueClasse = 25 // Buff d'attaque en forme transformée
		} else {
			attaqueClasse = 18
		}
	case "Bûcheron":
		attaqueClasse = 20
	default:
		attaqueClasse = 15
	}

	// Bonus de niveau (+2 par niveau)
	bonusNiveau := (p.Niveau - 1) * 2

	// Attaque de l'arme (ou arme par défaut si pas d'arme)
	attaqueArme := arme.EffetAttaque
	if attaqueArme == 0 {
		// Utiliser l'arme par défaut du personnage
		attaqueArme = getAttaqueParDefaut(p)
	}

	// Total
	return attaqueClasse + bonusNiveau + attaqueArme
}

// getAttaqueParDefaut retourne l'attaque de l'arme par défaut du personnage
func getAttaqueParDefaut(p *personnage.Personnage) int {
	switch p.Classe {
	case "Humain":
		return 10 // Épée
	case "Loups-Garou":
		// Griffes en forme transformée si PV bas, sinon épée
		if p.PVActuels < (p.PVMax * 30 / 100) {
			return 20 // Griffes (forme transformée) - plus puissant
		}
		return 8 // Épée (forme humaine)
	case "Bûcheron":
		return 12 // Hache
	default:
		return 5 // Coup de poing
	}
}

// calculerDefenseJoueur calcule la défense du joueur basée sur son niveau, classe et armures
func calculerDefenseJoueur(p *personnage.Personnage, armures []objet.Armure) int {
	// Défense de base selon la classe
	defenseClasse := 0
	switch p.Classe {
	case "Humain":
		defenseClasse = 8
	case "Loups-Garou":
		// Buff de défense en forme transformée si PV < 30%
		if p.PVActuels < (p.PVMax * 30 / 100) {
			defenseClasse = 12 // Buff de défense en forme transformée
		} else {
			defenseClasse = 6
		}
	case "Bûcheron":
		defenseClasse = 10
	default:
		defenseClasse = 8
	}

	// Bonus de niveau (+1 par niveau)
	bonusNiveau := (p.Niveau - 1) * 1

	// Défense des armures
	defenseArmures := sommeDefense(armures)

	// Total
	return defenseClasse + bonusNiveau + defenseArmures
}

// calculerBonusClasse applique des bonus selon la classe du personnage et l'action
func calculerBonusClasse(p *personnage.Personnage, choix actionCombat) int {
	bonus := 0

	switch p.Classe {
	case "Humain":
		// Humain: bonus sur les attaques rapides
		if strings.Contains(choix.label, "rapide") {
			bonus = 2
		}
	case "Loups-Garou":
		// Loups-Garou: bonus sur les attaques lourdes et enrage si PV bas
		if strings.Contains(choix.label, "lourde") {
			bonus = 3
		}
		// Enrage si PV < 30% (forme transformée)
		if p.PVActuels < (p.PVMax * 30 / 100) {
			bonus += 5
		}
	case "Bûcheron":
		// Bûcheron: bonus sur les attaques lourdes et défense
		if strings.Contains(choix.label, "lourde") {
			bonus = 4
		}
		if choix.typeDefensif {
			bonus = 2 // bonus défense
		}
	}

	return bonus
}

// estEnFormeTransformee vérifie si le loup-garou est en forme transformée
func estEnFormeTransformee(p *personnage.Personnage) bool {
	return p.Classe == "Loups-Garou" && p.PVActuels < (p.PVMax*30/100)
}

// appliquerTransformation applique les effets de transformation (régénération de vie)
func appliquerTransformation(p *personnage.Personnage) {
	if estEnFormeTransformee(p) {
		// Régénération de vie en forme transformée (5% des PV max par tour)
		regeneration := p.PVMax * 5 / 100
		if regeneration < 1 {
			regeneration = 1
		}
		p.PVActuels += regeneration
		if p.PVActuels > p.PVMax {
			p.PVActuels = p.PVMax
		}
	}
}

// getArmeParDefautNom retourne le nom de l'arme par défaut du personnage
func getArmeParDefautNom(p *personnage.Personnage) string {
	switch p.Classe {
	case "Humain":
		return "Épée"
	case "Loups-Garou":
		if p.PVActuels < (p.PVMax * 30 / 100) {
			return "Griffes (forme transformée)"
		}
		return "Épée (forme humaine)"
	case "Bûcheron":
		return "Hache"
	default:
		return "Coup de Poing"
	}
}

// actionsForArme retourne une liste d'actions adaptées à l'arme équipée
func actionsForArme(arme objet.Arme) []actionCombat {
	nom := strings.ToLower(arme.Nom)

	// Actions pour armes par défaut (pas d'arme équipée)
	actionsPoings := []actionCombat{
		{label: "Coup de poing rapide (+précision, -dégâts)", modAtk: 0.8, modHit: +0.05},
		{label: "Coup de poing lourd (-précision, +dégâts)", modAtk: 1.3, modHit: -0.10},
		{label: "Esquive (défense ce tour)", typeDefensif: true, defBuffMul: 0.4},
	}

	// Par défaut: épée / mêlée polyvalente
	actionsEpee := []actionCombat{
		{label: "Attaque rapide (+précision, -dégâts)", modAtk: 0.75, modHit: +0.10},
		{label: "Attaque lourde (-précision, +dégâts)", modAtk: 1.50, modHit: -0.15},
		{label: "Parade (défense ce tour)", typeDefensif: true, defBuffMul: 0.6},
	}

	// Hache: l'utilisateur souhaite la traiter comme distance (comme l'arc)
	actionsDistance := []actionCombat{
		{label: "Tir précis (+précision, -dégâts)", modAtk: 0.65, modHit: +0.15},
		{label: "Tir de barrage (-précision, +dégâts)", modAtk: 1.60, modHit: -0.20},
		{label: "Esquive (défense ce tour)", typeDefensif: true, defBuffMul: 0.5},
	}

	// Si pas d'arme (EffetAttaque = 0), utiliser les poings
	if arme.EffetAttaque == 0 {
		return actionsPoings
	}

	if strings.Contains(nom, "arc") {
		return actionsDistance
	}
	if strings.Contains(nom, "hache") {
		return actionsDistance
	}
	// fallback épée, couteaux, etc.
	return actionsEpee
}
