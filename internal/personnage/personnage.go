package personnage

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"sloteriaa/struct/objet"
)

const (
	reset  = "\u001b[0m"
	red    = "\u001b[31m"
	green  = "\u001b[32m"
	yellow = "\u001b[33m"
	cyan   = "\u001b[36m"
	bold   = "\u001b[1m"
)

// Structure du personnage
type Personnage struct {
	Nom             string
	Classe          string
	Niveau          int
	PVMax           int
	PVActuels       int
	Inventaire      []string
	Argent          int
	Attaque         string
	Force           int
	Agilite         int
	Endurance       int
	ArmuresEquipees map[string]bool
	Materiaux       map[string]int // Matériaux de craft
	// Buffs temporaires
	BuffForce     int // Bonus temporaire de Force
	BuffAgilite   int // Bonus temporaire d'Agilité
	BuffEndurance int // Bonus temporaire d'Endurance
	BuffCombats   int // Nombre de combats restants pour les buffs
}

// ----------------- Initialisation -----------------
func initPersonnage(nom, classe string, niveau, pvmax, pvactuels, argent int, inventaire []string) Personnage {
	attaque := determineAttaque(classe, pvactuels, pvmax)

	// Convertir la clé d'arme en nom d'affichage
	armeNom := attaque
	switch attaque {
	case "Épée (forme humaine)":
		armeNom = "Épée rouillée"
	case "Griffes (forme transformée)":
		armeNom = "Griffes de loup-garou"
	case "Hache":
		armeNom = "Hache lourde"
	}

	// Créer l'inventaire avec l'arme de départ
	inventaireAvecArme := append(inventaire, armeNom)

	return Personnage{
		Nom:             nom,
		Classe:          classe,
		Niveau:          niveau,
		PVMax:           pvmax,
		PVActuels:       pvactuels,
		Inventaire:      inventaireAvecArme,
		Argent:          argent,
		Attaque:         attaque,
		Force:           5, // sera écrasé par les stats de classe
		Agilite:         5, // sera écrasé par les stats de classe
		Endurance:       5, // sera écrasé par les stats de classe
		ArmuresEquipees: make(map[string]bool),
		Materiaux:       make(map[string]int),
	}
}

func determineAttaque(classe string, pvActuels, pvMax int) string {
	switch classe {
	case "Humain":
		return "Épée"
	case "Loups-Garou":
		if pvMax > 0 && float64(pvActuels)/float64(pvMax) <= 0.3 {
			return "Griffes (forme transformée)"
		}
		return "Épée (forme humaine)"
	case "Bûcheron":
		return "Hache"
	default:
		return "Coup de Poing"
	}
}

// ----------------- Création interactive -----------------
var stdinReader = bufio.NewReader(os.Stdin)

func readLine(prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := stdinReader.ReadString('\n')
	if err != nil && !strings.HasSuffix(text, "\n") {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func nomValide(nom string) bool {
	if nom == "" {
		return false
	}
	for _, r := range nom {
		if !(unicode.IsLetter(r) || r == '-' || r == '\'' || r == ' ') {
			return false
		}
	}
	return true
}

func mettreMajuscule(nom string) string {
	parts := strings.FieldsFunc(strings.ToLower(nom), func(r rune) bool { return r == ' ' || r == '-' || r == '\'' })
	if len(parts) == 0 {
		return ""
	}
	for i, p := range parts {
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	rebuilt := strings.Builder{}
	last := rune(0)
	for _, r := range nom {
		if r == ' ' || r == '-' || r == '\'' {
			last = r
			rebuilt.WriteRune(r)
			continue
		}
		break
	}
	_ = last // keep to avoid unused hint; simple capitalization is sufficient
	return strings.Title(strings.ToLower(nom))
}

func choisirClasse() string {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        CLASSES DISPONIBLES                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Explication des statistiques
	fmt.Println("📊 EXPLICATION DES STATISTIQUES :")
	fmt.Println()
	fmt.Println("💪 FORCE :")
	fmt.Println("   • Augmente les dégâts d'attaque (+1 dégât tous les 2 points)")
	fmt.Println("   • Bonus sur toutes les armes")
	fmt.Println("   • Montée de niveau : +1 Force tous les 2 niveaux")
	fmt.Println()
	fmt.Println("🏃 AGILITÉ :")
	fmt.Println("   • Chance d'esquive : 2% par point d'agilité")
	fmt.Println("   • Bonus de dégâts sur armes rapides (épées/arcs) : +1 dégât tous les 3 points")
	fmt.Println("   • Chance de critique : 10 + Agilité (max 50%)")
	fmt.Println("   • Montée de niveau : +1 Agilité tous les 3 niveaux")
	fmt.Println()
	fmt.Println("❤️ ENDURANCE :")
	fmt.Println("   • Augmente les PV maximum (+10 PV par point)")
	fmt.Println("   • Montée de niveau : +1 Endurance tous les 2 niveaux")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("🏹 HUMAIN")
	fmt.Println("   • PV: 130 | Force: 6 | Agilité: 6 | Endurance: 6")
	fmt.Println("   • Équilibre parfait entre toutes les stats")
	fmt.Println("   • Pas de capacité spéciale")
	fmt.Println()
	fmt.Println("🐺 LOUPS-GAROU")
	fmt.Println("   • PV: 110 | Force: 5 | Agilité: 9 | Endurance: 5")
	fmt.Println("   • +4 Agilité (esquive, critiques, bonus armes rapides)")
	fmt.Println("   • TRANSFORMATION: À 30% PV ou moins, se transforme en loup")
	fmt.Println("     → +50% Attaque, +25% Vitesse, Griffes puissantes")
	fmt.Println()
	fmt.Println("🪓 BÛCHERON")
	fmt.Println("   • PV: 150 | Force: 8 | Agilité: 4 | Endurance: 7")
	fmt.Println("   • +3 Force, +2 Endurance, -2 Agilité")
	fmt.Println("   • BONUS HACHES: +5 dégâts avec toutes les haches")
	fmt.Println("   • Résistant et puissant, spécialisé en haches")
	fmt.Println()

	for {
		cls, _ := readLine("Choisissez une classe (Humain, Loups-Garou, Bûcheron) : ")
		switch strings.ToLower(strings.ReplaceAll(cls, " ", "")) {
		case "humain":
			return "Humain"
		case "loups-garou", "loupsgarou", "loupgarou":
			return "Loups-Garou"
		case "bûcheron", "bucheron":
			return "Bûcheron"
		default:
			fmt.Println("Classe invalide — entrez Humain, Loups-Garou ou Bûcheron.")
		}
	}
}

func CreationPersonnage() Personnage {
	var nom string
	for {
		n, _ := readLine("Nom (lettres, espaces, -, ') : ")
		if nomValide(n) {
			nom = n
			break
		}
		fmt.Println("Nom invalide. Utilisez lettres, espaces, tirets ou apostrophes.")
	}
	nom = mettreMajuscule(nom)
	classe := choisirClasse()

	var pvMax int
	var force, agilite, endurance int

	switch classe {
	case "Humain":
		pvMax = 130
		force, agilite, endurance = 6, 6, 6
	case "Loups-Garou":
		pvMax = 110
		force, agilite, endurance = 5, 9, 5
	case "Bûcheron":
		pvMax = 150
		force, agilite, endurance = 8, 4, 7
	}
	pvActuels := pvMax / 2
	niveau := 1
	inventaire := []string{}
	argentDepart := 100

	// Créer le personnage avec les stats de base
	p := initPersonnage(nom, classe, niveau, pvMax, pvActuels, argentDepart, inventaire)

	// Appliquer les stats spécifiques à la classe
	p.Force = force
	p.Agilite = agilite
	p.Endurance = endurance

	return p
}

// Fonction pour mettre à jour l'attaque du personnage selon sa transformation
func UpdatePlayerAttack(p *Personnage) {
	p.Attaque = determineAttaque(p.Classe, p.PVActuels, p.PVMax)
}

// ----------------- Affichage stylé -----------------
func repeat(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}

func barreDeVie(actuels, max int) string {
	if max <= 0 {
		return "[??????????]"
	}
	if actuels < 0 {
		actuels = 0
	}
	if actuels > max {
		actuels = max
	}
	barLength := 10
	filled := (actuels * barLength) / max
	empty := barLength - filled
	return "[" + green + repeat("█", filled) + red + repeat("░", empty) + reset + "]"
}

func AfficherInfos(p Personnage) {
	innerWidth := 46
	top := cyan + "╔" + repeat("═", innerWidth+2) + "╗" + reset
	mid := cyan + "╠" + repeat("═", innerWidth+2) + "╣" + reset
	bot := cyan + "╚" + repeat("═", innerWidth+2) + "╝" + reset

	fmt.Println(top)
	fmt.Println(cyan + "║ " + reset + padVisible(bold+yellow+"STATS DU PERSONNAGE"+reset, innerWidth) + cyan + " ║" + reset)
	fmt.Println(mid)

	line := func(label, value string) {
		content := fmt.Sprintf("%s : %s", label, value)
		fmt.Println(cyan + "║ " + reset + padVisible(content, innerWidth) + cyan + " ║" + reset)
	}

	line("Nom", p.Nom)

	// Afficher la classe avec ses caractéristiques spéciales
	classeInfo := p.Classe
	switch p.Classe {
	case "Humain":
		classeInfo += " (Équilibré)"
	case "Loups-Garou":
		classeInfo += " (Agile, Transformation)"
		if float64(p.PVActuels)/float64(p.PVMax) <= 0.3 {
			classeInfo += " [TRANSFORMÉ]"
		}
	case "Bûcheron":
		classeInfo += " (Fort, Résistant)"
	}
	line("Classe", classeInfo)
	line("Niveau", fmt.Sprintf("%d", p.Niveau))
	pv := fmt.Sprintf("%d/%d %s", p.PVActuels, p.PVMax, barreDeVie(p.PVActuels, p.PVMax))
	line("PV", pv)
	// Weapon attack and total damage
	wepAtk := 0
	if p.Attaque != "" {
		if a, ok := findWeaponByNameOrKey(p.Attaque); ok {
			wepAtk = a.EffetAttaque
		}
	}
	line("Attaque", p.Attaque)
	line("Force", fmt.Sprintf("%d (+%d) = %d", p.Force, wepAtk, p.Force+wepAtk))
	// Defense from equipped armors
	defTotal := computeDefense(p)
	line("Défense", fmt.Sprintf("%d", defTotal))
	line("Argent", fmt.Sprintf("%d pièces", p.Argent))

	fmt.Println(bot)
}

var ansiRegexp = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

func visibleLen(s string) int {
	return len([]rune(stripANSI(s)))
}

func padVisible(s string, width int) string {
	vis := visibleLen(s)
	if vis >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vis)
}

// --- Helpers to resolve equipped items to stats ---
func findWeaponByNameOrKey(n string) (objet.Arme, bool) {
	keys := []string{
		"EpeeRouillee", "EpeeFer", "EpeeMagique", "EpeeCourte",
		"Hache", "HacheDeCombat", "HacheDeBataille",
		"ArcBois", "ArcLong", "ArcElfe",
	}
	needle := strings.ToLower(strings.TrimSpace(n))
	for _, k := range keys {
		a := objet.CreerArme(k)
		if strings.EqualFold(needle, k) || strings.EqualFold(needle, a.Nom) {
			return a, true
		}
	}
	return objet.Arme{}, false
}

func findArmorByDisplayOrKey(n string) (objet.Armure, bool) {
	keys := []string{
		"CasqueCuir", "CasqueCuirRenforce", "CasqueFer", "CasqueFerRenforce",
		"PlastronCuir", "PlastronCuirRenforce", "PlastronFer", "PlastronFerRenforce",
		"PantalonCuir", "PantalonCuirRenforce", "PantalonFer", "PantalonFerRenforce",
		"BottesCuir", "BottesCuirRenforce", "BottesFer", "BottesFerRenforce",
	}
	needle := strings.ToLower(strings.TrimSpace(n))
	for _, k := range keys {
		ar := objet.CreerArmure(k)
		if strings.EqualFold(needle, k) || strings.EqualFold(needle, ar.Nom) {
			return ar, true
		}
	}
	return objet.Armure{}, false
}

func computeDefense(p Personnage) int {
	if len(p.ArmuresEquipees) == 0 {
		return 0
	}
	total := 0
	for name, equipped := range p.ArmuresEquipees {
		if !equipped {
			continue
		}
		if ar, ok := findArmorByDisplayOrKey(name); ok {
			total += ar.EffetDefense
		}
	}
	return total
}
