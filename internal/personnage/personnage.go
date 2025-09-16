package personnage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
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
	Nom        string
	Classe     string
	Niveau     int
	PVMax      int
	PVActuels  int
	Inventaire []string
	Argent     int
	Attaque    string
}

// ----------------- Initialisation -----------------
func initPersonnage(nom, classe string, niveau, pvmax, pvactuels, argent int, inventaire []string) Personnage {
	attaque := determineAttaque(classe, pvactuels, pvmax)
	return Personnage{
		Nom:        nom,
		Classe:     classe,
		Niveau:     niveau,
		PVMax:      pvmax,
		PVActuels:  pvactuels,
		Inventaire: inventaire,
		Argent:     argent,
		Attaque:    attaque,
	}
}

func determineAttaque(classe string, pvActuels, pvMax int) string {
	switch classe {
	case "Humain":
		return "Épée"
	case "Loups-Garou":
		if pvMax > 0 && float64(pvActuels)/float64(pvMax) < 0.3 {
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
	switch classe {
	case "Humain":
		pvMax = 120
	case "Loups-Garou":
		pvMax = 100
	case "Bûcheron":
		pvMax = 140
	}
	pvActuels := pvMax / 2
	niveau := 1
	inventaire := []string{"Potion", "Potion", "Potion"}
	argentDepart := 100

	return initPersonnage(nom, classe, niveau, pvMax, pvActuels, argentDepart, inventaire)
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
	fmt.Println(cyan + "╔═══════════════════════════════════╗" + reset)
	fmt.Printf(cyan+"║"+reset+" %-46s "+cyan+"║\n"+reset, bold+yellow+"FEUILLE DE PERSONNAGE"+reset)
	fmt.Println(cyan + "╠═══════════════════════════════════╣" + reset)

	fmt.Printf(cyan+"║"+reset+" Nom    : %-32s "+cyan+"║\n"+reset, bold+p.Nom+reset)
	fmt.Printf(cyan+"║"+reset+" Classe : %-24s "+cyan+"║\n"+reset, p.Classe)
	fmt.Printf(cyan+"║"+reset+" Niveau : %-24d "+cyan+"║\n"+reset, p.Niveau)
	fmt.Printf(cyan+"║"+reset+" PV     : %-3d/%-3d %-31s"+cyan+"║\n"+reset,
		p.PVActuels, p.PVMax, barreDeVie(p.PVActuels, p.PVMax))
	fmt.Printf(cyan+"║"+reset+" Attaque : %-23s "+cyan+"║\n"+reset, p.Attaque)
	fmt.Printf(cyan+"║"+reset+" Argent : %-4d pièces              "+cyan+"║\n"+reset, p.Argent)

	fmt.Println(cyan + "╠═══════════════════════════════════╣" + reset)
	fmt.Println(cyan + "║" + reset + " Inventaire :                      " + cyan + "║" + reset)
	if len(p.Inventaire) == 0 {
		fmt.Println(cyan + "║" + reset + "   (vide)                           " + cyan + "║" + reset)
	} else {
		for _, item := range p.Inventaire {
			if len(item) > 30 {
				item = item[:30]
			}
			fmt.Printf(cyan+"║"+reset+"   - %-30s"+cyan+"║\n"+reset, item)
		}
	}
	fmt.Println(cyan + "╚═══════════════════════════════════╝" + reset)
}
