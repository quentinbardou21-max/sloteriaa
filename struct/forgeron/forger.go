package forgeron

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/mattn/go-tty"

	"github.com/quent/sloteriaa/struct/objet"
)

// Enable ANSI sequences on Windows consoles (VT processing)
func enableWindowsVT() {
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	h := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	_, _, _ = getConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode)))
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	_, _, _ = setConsoleMode.Call(uintptr(h), uintptr(mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING))
}

// Définition des matériaux utilisables à la forge
type Materiau string

const (
	Fer            Materiau = "Lingot de fer"
	Bois           Materiau = "Planche de bois"
	Cuir           Materiau = "Cuir tanné"
	EssenceMagique Materiau = "Essence magique"
	Or             Materiau = "Or"
)

// Coût d'une recette par matériau
type Cout map[Materiau]int

// Recette d'artisanat d'une arme humaine (non-monstre)
type Recette struct {
	CleArme    string
	NomAffiche string
	Cout       Cout
}

// Recette pour une armure
type RecetteArmure struct {
	CleArmure  string
	NomAffiche string
	Cout       Cout
}

// Inventaire des matériaux du joueur
type InventaireMateriaux map[Materiau]int

func (inv InventaireMateriaux) AAssez(c Cout) bool {
	for m, q := range c {
		if inv[m] < q {
			return false
		}
	}
	return true
}

func (inv InventaireMateriaux) Debiter(c Cout) {
	for m, q := range c {
		inv[m] -= q
		if inv[m] < 0 {
			inv[m] = 0
		}
	}
}

// Ordre d'affichage déterministe des matériaux
var materialOrder = []Materiau{Or, Fer, Bois, Cuir, EssenceMagique}

func orderedKeysFromCout(c Cout) []Materiau {
	ordered := make([]Materiau, 0, len(c))
	for _, m := range materialOrder {
		if _, ok := c[m]; ok {
			ordered = append(ordered, m)
		}
	}
	return ordered
}

// Catalogue des recettes pour armes humaines (clés compatibles avec objet.CreerArme)
func RecettesArmesHumaines() []Recette {
	return []Recette{
		{CleArme: "EpeeRouillee", NomAffiche: "Épée rouillée", Cout: Cout{Fer: 1, Or: 100}},
		{CleArme: "EpeeCourte", NomAffiche: "Épée courte", Cout: Cout{Fer: 2, Cuir: 1, Or: 180}},
		{CleArme: "EpeeFer", NomAffiche: "Épée en fer", Cout: Cout{Fer: 4, Cuir: 1, Or: 260}},
		{CleArme: "Hache", NomAffiche: "Hache lourde", Cout: Cout{Fer: 5, Bois: 2, Or: 300}},
		{CleArme: "HacheDeCombat", NomAffiche: "Hache de combat", Cout: Cout{Fer: 4, Bois: 2, Cuir: 1, Or: 380}},
		{CleArme: "HacheDeBataille", NomAffiche: "Hache de bataille", Cout: Cout{Fer: 7, Bois: 3, Or: 700}},
		{CleArme: "ArcBois", NomAffiche: "Arc en bois", Cout: Cout{Bois: 4, Cuir: 1, Or: 120}},
		{CleArme: "ArcLong", NomAffiche: "Arc long", Cout: Cout{Bois: 6, Cuir: 2, Or: 220}},
		{CleArme: "ArcElfe", NomAffiche: "Arc elfique", Cout: Cout{Bois: 5, Cuir: 2, EssenceMagique: 1, Or: 450}},
		{CleArme: "EpeeMagique", NomAffiche: "Épée magique", Cout: Cout{Fer: 5, EssenceMagique: 2, Or: 950}},
	}
}

// Catalogue des recettes d'armures
func RecettesArmures() []RecetteArmure {
	return []RecetteArmure{
		// Casques
		{CleArmure: "CasqueCuir", NomAffiche: "Casque en cuir", Cout: Cout{Cuir: 2, Or: 120}},
		{CleArmure: "CasqueCuirRenforce", NomAffiche: "Casque cuir renforcé", Cout: Cout{Cuir: 3, Fer: 1, Or: 180}},
		{CleArmure: "CasqueFer", NomAffiche: "Casque en fer", Cout: Cout{Fer: 3, Or: 260}},
		{CleArmure: "CasqueFerRenforce", NomAffiche: "Casque fer renforcé", Cout: Cout{Fer: 4, Cuir: 1, Or: 340}},

		// Plastrons
		{CleArmure: "PlastronCuir", NomAffiche: "Plastron cuir", Cout: Cout{Cuir: 3, Or: 180}},
		{CleArmure: "PlastronCuirRenforce", NomAffiche: "Plastron cuir renforcé", Cout: Cout{Cuir: 4, Fer: 1, Or: 260}},
		{CleArmure: "PlastronFer", NomAffiche: "Plastron fer", Cout: Cout{Fer: 5, Or: 420}},
		{CleArmure: "PlastronFerRenforce", NomAffiche: "Plastron fer renforcé", Cout: Cout{Fer: 6, Cuir: 1, Or: 800}},

		// Pantalons
		{CleArmure: "PantalonCuir", NomAffiche: "Pantalon cuir", Cout: Cout{Cuir: 2, Or: 150}},
		{CleArmure: "PantalonCuirRenforce", NomAffiche: "Pantalon cuir renforcé", Cout: Cout{Cuir: 3, Fer: 1, Or: 220}},
		{CleArmure: "PantalonFer", NomAffiche: "Pantalon fer", Cout: Cout{Fer: 4, Or: 320}},
		{CleArmure: "PantalonFerRenforce", NomAffiche: "Pantalon fer renforcé", Cout: Cout{Fer: 5, Cuir: 1, Or: 450}},

		// Chaussures
		{CleArmure: "BottesCuir", NomAffiche: "Bottes cuir", Cout: Cout{Cuir: 2, Or: 120}},
		{CleArmure: "BottesCuirRenforce", NomAffiche: "Bottes cuir renforcé", Cout: Cout{Cuir: 3, Fer: 1, Or: 180}},
		{CleArmure: "BottesFer", NomAffiche: "Bottes fer", Cout: Cout{Fer: 3, Or: 240}},
		{CleArmure: "BottesFerRenforce", NomAffiche: "Bottes fer renforcé", Cout: Cout{Fer: 4, Cuir: 1, Or: 320}},
	}
}

// Interface console: point d'entrée de la forge
func RunForge(inv InventaireMateriaux) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== Forge ===")
		afficherInventaire(inv)
		fmt.Print("Voulez-vous parler au forgeron ? (o/n): ")
		choix, _ := reader.ReadString('\n')
		choix = strings.TrimSpace(strings.ToLower(choix))
		fmt.Println()

		if choix == "n" || choix == "non" {
			fmt.Println("Vous quittez la forge.")
			return
		}
		if choix != "o" && choix != "oui" {
			fmt.Println("Réponse invalide. Réessayez.")
			continue
		}

		afficherCatalogue()
		fmt.Print("Entrez le numéro de l'arme à forger (ou vide pour annuler): ")
		entree, _ := reader.ReadString('\n')
		entree = strings.TrimSpace(entree)
		if entree == "" {
			fmt.Println("Retour au forgeron...")
			continue
		}
		idx, err := strconv.Atoi(entree)
		if err != nil {
			fmt.Println("Entrée invalide.")
			continue
		}

		recettes := RecettesArmesHumaines()
		if idx < 1 || idx > len(recettes) {
			fmt.Println("Numéro hors liste.")
			continue
		}

		recette := recettes[idx-1]
		if !inv.AAssez(recette.Cout) {
			fmt.Println("Vous n'avez pas assez de matériaux pour cette arme.")
			fmt.Println("Coût requis:")
			afficherCout(recette.Cout)
			fmt.Println()
			continue
		}

		inv.Debiter(recette.Cout)
		arme := objet.CreerArme(recette.CleArme)
		fmt.Println("Fabrication réussie ! Voici votre arme :")
		objet.AfficherArme(arme)
	}
}

func afficherInventaire(inv InventaireMateriaux) {
	if len(inv) == 0 {
		fmt.Println("Inventaire matériaux: (vide)")
		return
	}
	fmt.Println("Inventaire matériaux:")
	for _, m := range materialOrder {
		if q, ok := inv[m]; ok && q > 0 {
			fmt.Printf("  - %s x%d\n", m, q)
		}
	}
	fmt.Println()
}

func afficherCatalogue() {
	recettes := RecettesArmesHumaines()
	fmt.Println("--- Catalogue d'armes ---")
	for i, r := range recettes {
		fmt.Printf("%d) %s\n", i+1, r.NomAffiche)
		arme := objet.CreerArme(r.CleArme)
		fmt.Printf("   -> Attaque: %d | Poids: %d\n", arme.EffetAttaque, arme.Poids)
		fmt.Print("   Coût: ")
		afficherCoutInline(r.Cout)
		fmt.Println()
	}
	fmt.Println()
}

func afficherCout(c Cout) {
	for _, m := range orderedKeysFromCout(c) {
		q := c[m]
		fmt.Printf("  - %s x%d\n", m, q)
	}
}

func afficherCoutInline(c Cout) {
	items := make([]string, 0, len(c))
	for _, m := range orderedKeysFromCout(c) {
		q := c[m]
		items = append(items, fmt.Sprintf("%s x%d", m, q))
	}
	fmt.Print(strings.Join(items, ", "))
}

// Efface l'écran et replace le curseur (in-place, sans cls)
func clearScreenTUI() {
	// Home, then clear to end of screen
	fmt.Print("\033[H\033[J")
}

func enterAltScreen() { fmt.Print("\033[?1049h\033[H") }
func exitAltScreen()  { fmt.Print("\033[?1049l") }

func hideCursor() { fmt.Print("\033[?25l") }
func showCursor() { fmt.Print("\033[?25h") }

// Renders one full frame: UI + controls + optional modal message
func renderForgeTUIFrame(recettes []Recette, inv InventaireMateriaux, selection int, showArmors bool, modal string) {
	// Reposition cursor to top-left and redraw frame without full reset to reduce flicker
	fmt.Print("\033[H")
	renderForgeTUI(recettes, inv, selection, showArmors)
	// Controls line
	fmt.Printf("%sContrôles:%s ↑/↓ naviguer | C forger | T armes/armures | Q quitter\n", ansiCyan, ansiReset)
	// Modal message
	if modal != "" {
		fmt.Println("--------------------------------------------------")
		fmt.Println(modal)
		fmt.Println("(Appuyez sur une touche pour revenir)")
	}
	// Clear any remaining content below the current cursor position
	fmt.Print("\033[J")
}

func splitCostMaterialsGold(c Cout) (materials []string, gold int) {
	materials = []string{}
	for _, m := range orderedKeysFromCout(c) {
		if m == Or {
			gold = c[m]
			continue
		}
		materials = append(materials, fmt.Sprintf("%s x%d", m, c[m]))
	}
	return
}

// Build a full frame into a single string to minimize I/O and flicker
func renderFrameString(recettes []Recette, inv InventaireMateriaux, selection int, showArmors bool, modal string) string {
	var b strings.Builder
	w, _ := termSize()
	if runtime.GOOS == "windows" {
		w = 120
	}
	if w < 60 {
		w = 60
	}
	fmt.Fprintf(&b, "%s=== Forge ===%s\n", ansiCyan, ansiReset)

	gap := 2
	available := w - gap
	leftTotal := (available * 40) / 100
	if leftTotal < 26 {
		leftTotal = 26
	}
	rightTotal := available - leftTotal
	if rightTotal < 30 {
		rightTotal = 30
		leftTotal = available - rightTotal
		if leftTotal < 26 {
			leftTotal = 26
		}
	}
	leftInner := leftTotal - 2
	rightInner := rightTotal - 2
	if leftInner < 10 {
		leftInner = 10
	}
	if rightInner < 10 {
		rightInner = 10
	}

	left := []string{"Inventaire matériaux:"}
	printed := false
	for _, m := range materialOrder {
		if q, ok := inv[m]; ok && q > 0 {
			left = append(left, fmt.Sprintf("- %s x%d", m, q))
			printed = true
		}
	}
	if !printed {
		left = append(left, "(vide)")
	}
	left = append(left, "")

	if showArmors {
		left = append(left, "Liste: Armures")
		ars := RecettesArmures()
		if len(ars) == 0 {
			left = append(left, "(aucune recette)")
		} else {
			if selection >= len(ars) {
				selection = 0
			}
			for i, r := range ars {
				chev := ' '
				if i == selection {
					chev = '>'
				}
				line := fmt.Sprintf("%c %2d) %s", chev, i+1, r.NomAffiche)
				if i == selection {
					line = ansiBgWhite + ansiBlack + line + ansiReset
				}
				left = append(left, line)
			}
		}
	} else {
		left = append(left, "Liste: Armes")
		if len(recettes) == 0 {
			left = append(left, "(aucune recette)")
		} else {
			if selection >= len(recettes) {
				selection = 0
			}
			for i, r := range recettes {
				chev := ' '
				if i == selection {
					chev = '>'
				}
				line := fmt.Sprintf("%c %2d) %s", chev, i+1, r.NomAffiche)
				if i == selection {
					line = ansiBgWhite + ansiBlack + line + ansiReset
				}
				left = append(left, line)
			}
		}
	}
	leftBox := boxify("Forge", left, leftInner)

	right := []string{"Détails:"}
	if showArmors {
		ars := RecettesArmures()
		if len(ars) > 0 {
			if selection >= len(ars) {
				selection = 0
			}
			ch := ars[selection]
			arm := objet.CreerArmure(ch.CleArmure)
			right = append(right,
				fmt.Sprintf("Nom: %s", arm.Nom),
				fmt.Sprintf("Type: %s", arm.Type),
				fmt.Sprintf("Défense: %d", arm.EffetDefense),
				fmt.Sprintf("Poids: %d", arm.Poids),
				"Coût (matériaux):")
			mats, gold := splitCostMaterialsGold(ch.Cout)
			if len(mats) > 0 {
				right = append(right, strings.Join(mats, ", "))
			} else {
				right = append(right, "-")
			}
			if inv.AAssez(ch.Cout) {
				right = append(right, "[Disponible]")
			} else {
				right = append(right, "[Matériaux insuffisants]")
			}
			// Bottom-right gold cost inside right panel
			goldLine := fmt.Sprintf("Coût: %s", colorYellow(fmt.Sprintf("%d", gold)))
			right = append(right, padLeftANSI(goldLine, rightInner))
		} else {
			right = append(right, "(sélectionnez une armure)")
		}
	} else {
		if len(recettes) > 0 {
			if selection >= len(recettes) {
				selection = 0
			}
			ch := recettes[selection]
			arme := objet.CreerArme(ch.CleArme)
			right = append(right,
				fmt.Sprintf("Nom: %s", arme.Nom),
				fmt.Sprintf("Attaque: %d", arme.EffetAttaque),
				fmt.Sprintf("Poids: %d", arme.Poids),
				"Description:")
			right = append(right, arme.Description)
			right = append(right, "Coût (matériaux):")
			mats, gold := splitCostMaterialsGold(ch.Cout)
			if len(mats) > 0 {
				right = append(right, strings.Join(mats, ", "))
			} else {
				right = append(right, "-")
			}
			if inv.AAssez(ch.Cout) {
				right = append(right, "[Disponible]")
			} else {
				right = append(right, "[Matériaux insuffisants]")
			}
			// Bottom-right gold cost inside right panel
			goldLine := fmt.Sprintf("Coût: %s", colorYellow(fmt.Sprintf("%d", gold)))
			right = append(right, padLeftANSI(goldLine, rightInner))
		} else {
			right = append(right, "(sélectionnez une arme)")
		}
	}
	rightBox := boxify("Panneau", right, rightInner)

	leftBox, rightBox = equalizeBoxHeights(leftBox, rightBox, leftInner, rightInner)
	maxLines := len(leftBox)
	if len(rightBox) > maxLines {
		maxLines = len(rightBox)
	}
	for i := 0; i < maxLines; i++ {
		l := padRightANSI("", leftTotal)
		r := padRightANSI("", rightTotal)
		if i < len(leftBox) {
			l = padRightANSI(leftBox[i], leftTotal)
		}
		if i < len(rightBox) {
			r = padRightANSI(rightBox[i], rightTotal)
		}
		b.WriteString(l)
		b.WriteString(repeatRune(' ', gap))
		b.WriteString(r)
		b.WriteByte('\n')
	}

	// Controls
	fmt.Fprintf(&b, "%sContrôles:%s ↑/↓ naviguer | C forger | T armes/armures | Q quitter\n", ansiCyan, ansiReset)
	// Modal
	if modal != "" {
		b.WriteString("--------------------------------------------------\n")
		b.WriteString(modal)
		b.WriteByte('\n')
		b.WriteString("(Appuyez sur une touche pour revenir)\n")
	}
	return b.String()
}

// Rune-aware helpers
func runeCount(s string) int {
	return utf8.RuneCountInString(s)
}

func splitRunes(s string, n int) (string, string) {
	if n <= 0 {
		return "", s
	}
	count := 0
	for i := range s {
		if count == n {
			return s[:i], s[i:]
		}
		count++
	}
	// Loop finished. count == total rune count.
	if n >= count {
		// Requesting all runes or more: return whole string as head
		return s, ""
	}
	// Should not reach here normally, fallback safe
	return s, ""
}

func padRightRunes(s string, width int) string {
	l := runeCount(s)
	if l >= width {
		head, _ := splitRunes(s, width)
		return head
	}
	return s + repeatRune(' ', width-l)
}

// ANSI colors (used outside boxes to avoid alignment issues)
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiCyan    = "\033[36m"
	ansiYellow  = "\033[33m"
	ansiGreen   = "\033[32m"
	ansiRed     = "\033[31m"
	ansiBlack   = "\033[30m"
	ansiBgWhite = "\033[47m"
)

func colorYellow(s string) string { return ansiYellow + s + ansiReset }

func padLeftANSI(s string, width int) string {
	printable := runeCount(stripANSI(s))
	if printable >= width {
		return s
	}
	return repeatRune(' ', width-printable) + s
}

// Strip ANSI sequences for length calculations
func stripANSI(s string) string {
	res := make([]rune, 0, len(s))
	i := 0
	runes := []rune(s)
	for i < len(runes) {
		r := runes[i]
		if r == 0x1b && i+1 < len(runes) && runes[i+1] == '[' {
			// Skip until 'm' or end
			i += 2
			for i < len(runes) {
				if runes[i] == 'm' {
					i++
					break
				}
				i++
			}
			continue
		}
		res = append(res, r)
		i++
	}
	return string(res)
}

func padRightANSI(s string, width int) string {
	printable := runeCount(stripANSI(s))
	if printable >= width {
		// naive: no truncation with ANSI to avoid cutting sequences
		return s
	}
	return s + repeatRune(' ', width-printable)
}

// Utilities for layout
func termSize() (int, int) {
	// Avoid opening a TTY during rendering; rely on env or defaults
	if c := os.Getenv("COLUMNS"); c != "" {
		if w, err := strconv.Atoi(c); err == nil && w > 0 {
			if r := os.Getenv("LINES"); r != "" {
				if h, err := strconv.Atoi(r); err == nil && h > 0 {
					return w, h
				}
			}
			return w, 30
		}
	}
	return 100, 30
}

func repeatRune(r rune, n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = r
	}
	return string(b)
}

func wrapText(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}
	words := strings.Fields(s)
	lines := []string{}
	line := ""
	for _, w := range words {
		if runeCount(line) == 0 {
			if runeCount(w) <= width {
				line = w
			} else {
				head, tail := splitRunes(w, width)
				lines = append(lines, head)
				line = tail
			}
			continue
		}
		if runeCount(line)+1+runeCount(w) <= width {
			line = line + " " + w
		} else {
			lines = append(lines, line)
			if runeCount(w) <= width {
				line = w
			} else {
				head, tail := splitRunes(w, width)
				lines = append(lines, head)
				line = tail
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

func padRight(s string, width int) string {
	// Deprecated in favor of ANSI-aware padding when coloring
	if len(s) >= width {
		return s[:width]
	}
	return s + repeatRune(' ', width-len(s))
}

func boxify(title string, content []string, innerWidth int) []string {
	// innerWidth is inside borders, final width = innerWidth+2
	if innerWidth < 10 {
		innerWidth = 10
	}
	top := "+" + repeatRune('-', innerWidth) + "+"
	titleLine := "|" + padRightANSI(title, innerWidth) + "|"
	lines := []string{top, titleLine, "+" + repeatRune('-', innerWidth) + "+"}
	for _, c := range content {
		for _, ln := range wrapText(c, innerWidth) {
			lines = append(lines, "|"+padRightANSI(ln, innerWidth)+"|")
		}
	}
	lines = append(lines, top)
	return lines
}

func addEmptyRowsBeforeBottom(box []string, innerWidth int, count int) []string {
	if count <= 0 {
		return box
	}
	if len(box) < 2 {
		return box
	}
	bottomIdx := len(box) - 1
	empty := "|" + padRightRunes("", innerWidth) + "|"
	pads := make([]string, count)
	for i := range pads {
		pads[i] = empty
	}
	// Insert pads before bottom border
	box = append(box[:bottomIdx], append(pads, box[bottomIdx:]...)...)
	return box
}

func equalizeBoxHeights(leftBox []string, rightBox []string, leftInner int, rightInner int) ([]string, []string) {
	// Ensure both boxes end with their bottom border on the same line
	ll := len(leftBox)
	rl := len(rightBox)
	if ll == rl {
		return leftBox, rightBox
	}
	if ll < rl {
		leftBox = addEmptyRowsBeforeBottom(leftBox, leftInner, rl-ll)
	} else if rl < ll {
		rightBox = addEmptyRowsBeforeBottom(rightBox, rightInner, ll-rl)
	}
	return leftBox, rightBox
}

// Sorting
type SortMode int

const (
	SortNone      SortMode = iota
	SortByAttack           // armes
	SortByWeight           // both
	SortByCost             // both
	SortByDefense          // armures
)

func totalCost(c Cout) int {
	sum := 0
	for _, q := range c {
		sum += q
	}
	return sum
}

// Render without ANSI inside borders; highlight selection with colored line
func renderForgeTUI(recettes []Recette, inv InventaireMateriaux, selection int, showArmors bool) {
	// Do not full-clear here; caller positions cursor at home before render
	w, _ := termSize()
	if runtime.GOOS == "windows" {
		w = 120
	}
	if w < 60 {
		w = 60
	}
	fmt.Printf("%s=== Forge (TUI) ===%s\n", ansiCyan, ansiReset)

	gap := 2
	available := w - gap
	leftTotal := (available * 40) / 100
	if leftTotal < 26 {
		leftTotal = 26
	}
	rightTotal := available - leftTotal
	if rightTotal < 30 {
		rightTotal = 30
		leftTotal = available - rightTotal
		if leftTotal < 26 {
			leftTotal = 26
		}
	}
	leftInner := leftTotal - 2
	rightInner := rightTotal - 2
	if leftInner < 10 {
		leftInner = 10
	}
	if rightInner < 10 {
		rightInner = 10
	}

	left := []string{"Inventaire matériaux:"}
	printed := false
	for _, m := range materialOrder {
		if q, ok := inv[m]; ok && q > 0 {
			left = append(left, fmt.Sprintf("- %s x%d", m, q))
			printed = true
		}
	}
	if !printed {
		left = append(left, "(vide)")
	}
	left = append(left, "")

	if showArmors {
		left = append(left, "Liste: Armures")
		ars := RecettesArmures()
		if len(ars) == 0 {
			left = append(left, "(aucune recette)")
		} else {
			if selection >= len(ars) {
				selection = 0
			}
			for i, r := range ars {
				chev := ' '
				if i == selection {
					chev = '>'
				}
				line := fmt.Sprintf("%c %2d) %s", chev, i+1, r.NomAffiche)
				if i == selection {
					line = ansiBgWhite + ansiBlack + line + ansiReset
				}
				left = append(left, line)
			}
		}
	} else {
		left = append(left, "Liste: Armes")
		if len(recettes) == 0 {
			left = append(left, "(aucune recette)")
		} else {
			if selection >= len(recettes) {
				selection = 0
			}
			for i, r := range recettes {
				chev := ' '
				if i == selection {
					chev = '>'
				}
				line := fmt.Sprintf("%c %2d) %s", chev, i+1, r.NomAffiche)
				if i == selection {
					line = ansiBgWhite + ansiBlack + line + ansiReset
				}
				left = append(left, line)
			}
		}
	}
	leftBox := boxify("Forge", left, leftInner)

	right := []string{"Détails:"}
	if showArmors {
		ars := RecettesArmures()
		if len(ars) > 0 {
			if selection >= len(ars) {
				selection = 0
			}
			ch := ars[selection]
			arm := objet.CreerArmure(ch.CleArmure)
			right = append(right,
				fmt.Sprintf("Nom: %s", arm.Nom),
				fmt.Sprintf("Type: %s", arm.Type),
				fmt.Sprintf("Défense: %d", arm.EffetDefense),
				fmt.Sprintf("Poids: %d", arm.Poids),
				"Coût:")
			costLine := strings.Join(func() []string {
				a := []string{}
				for _, m := range orderedKeysFromCout(ch.Cout) {
					a = append(a, fmt.Sprintf("%s x%d", m, ch.Cout[m]))
				}
				return a
			}(), ", ")
			right = append(right, costLine)
			if inv.AAssez(ch.Cout) {
				right = append(right, "[Disponible]")
			} else {
				right = append(right, "[Matériaux insuffisants]")
			}
			// Bottom-right gold cost inside right panel
			_, gold := splitCostMaterialsGold(ch.Cout)
			goldLine := fmt.Sprintf("Coût: %s", colorYellow(fmt.Sprintf("%d", gold)))
			right = append(right, padLeftANSI(goldLine, rightInner))
		} else {
			right = append(right, "(sélectionnez une armure)")
		}
	} else {
		if len(recettes) > 0 {
			if selection >= len(recettes) {
				selection = 0
			}
			ch := recettes[selection]
			arme := objet.CreerArme(ch.CleArme)
			right = append(right,
				fmt.Sprintf("Nom: %s", arme.Nom),
				fmt.Sprintf("Attaque: %d", arme.EffetAttaque),
				fmt.Sprintf("Poids: %d", arme.Poids),
				"Description:")
			right = append(right, arme.Description)
			right = append(right, "Coût:")
			costLine := strings.Join(func() []string {
				a := []string{}
				for _, m := range orderedKeysFromCout(ch.Cout) {
					a = append(a, fmt.Sprintf("%s x%d", m, ch.Cout[m]))
				}
				return a
			}(), ", ")
			right = append(right, costLine)
			if inv.AAssez(ch.Cout) {
				right = append(right, "[Disponible]")
			} else {
				right = append(right, "[Matériaux insuffisants]")
			}
			// Bottom-right gold cost inside right panel
			_, gold := splitCostMaterialsGold(ch.Cout)
			goldLine := fmt.Sprintf("Coût: %s", colorYellow(fmt.Sprintf("%d", gold)))
			right = append(right, padLeftANSI(goldLine, rightInner))
		} else {
			right = append(right, "(sélectionnez une arme)")
		}
	}
	rightBox := boxify("Panneau", right, rightInner)

	leftBox, rightBox = equalizeBoxHeights(leftBox, rightBox, leftInner, rightInner)
	maxLines := len(leftBox)
	if len(rightBox) > maxLines {
		maxLines = len(rightBox)
	}
	for i := 0; i < maxLines; i++ {
		l := padRightANSI("", leftTotal)
		r := padRightANSI("", rightTotal)
		if i < len(leftBox) {
			l = padRightANSI(leftBox[i], leftTotal)
		}
		if i < len(rightBox) {
			r = padRightANSI(rightBox[i], rightTotal)
		}
		fmt.Print(l)
		fmt.Print(repeatRune(' ', gap))
		fmt.Println(r)
	}
}

// -------- TUI améliorée (navigation clavier) --------

func RunForgeTUI(inv InventaireMateriaux) {
	enableWindowsVT()
	recettes := RecettesArmesHumaines()
	if len(recettes) == 0 {
		fmt.Println("Aucune recette disponible.")
		return
	}

	enterAltScreen()
	defer exitAltScreen()
	reader := bufio.NewReader(os.Stdin)
	selection := 0
	showArmors := false
	for {
		frame := renderFrameString(recettes, inv, selection, showArmors, "")
		// Home only, then write frame
		fmt.Print("\033[H")
		fmt.Print(frame)
		fmt.Print("\033[J")
		// No extra prompt print to avoid duplicate lines
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		switch line {
		case "z":
			if selection > 0 {
				selection--
			}
		case "s":
			max := len(recettes)
			if showArmors {
				max = len(RecettesArmures())
			}
			if selection < max-1 {
				selection++
			}
		case "t":
			showArmors = !showArmors
			selection = 0
		case "c":
			if showArmors {
				ars := RecettesArmures()
				if len(ars) == 0 {
					continue
				}
				if selection >= len(ars) {
					selection = 0
				}
				choisie := ars[selection]
				if !inv.AAssez(choisie.Cout) {
					modal := "Matériaux insuffisants pour fabriquer cette armure."
					fmt.Print("\033[H")
					fmt.Print(renderFrameString(recettes, inv, selection, showArmors, modal))
					fmt.Print("\033[J")
					reader.ReadString('\n')
					continue
				}
				inv.Debiter(choisie.Cout)
				arm := objet.CreerArmure(choisie.CleArmure)
				modal := fmt.Sprintf("Fabrication réussie !\nNom: %s\nType: %s\nDéfense: %d | Poids: %d", arm.Nom, arm.Type, arm.EffetDefense, arm.Poids)
				fmt.Print("\033[H")
				fmt.Print(renderFrameString(recettes, inv, selection, showArmors, modal))
				fmt.Print("\033[J")
				reader.ReadString('\n')
			} else {
				if len(recettes) == 0 {
					continue
				}
				if selection >= len(recettes) {
					selection = 0
				}
				choisie := recettes[selection]
				if !inv.AAssez(choisie.Cout) {
					modal := "Matériaux insuffisants pour fabriquer cette arme."
					fmt.Print("\033[H")
					fmt.Print(renderFrameString(recettes, inv, selection, showArmors, modal))
					fmt.Print("\033[J")
					reader.ReadString('\n')
					continue
				}
				inv.Debiter(choisie.Cout)
				arme := objet.CreerArme(choisie.CleArme)
				modal := fmt.Sprintf("Fabrication réussie !\nNom: %s\nAttaque: %d | Poids: %d", arme.Nom, arme.EffetAttaque, arme.Poids)
				fmt.Print("\033[H")
				fmt.Print(renderFrameString(recettes, inv, selection, showArmors, modal))
				fmt.Print("\033[J")
				reader.ReadString('\n')
			}
		case "q":
			resetTerminalVisual()
			return
		}
	}
}

// -------- TUI temps réel avec flèches (sans Entrée) --------

// RunForgeInteractive lit les touches en temps réel (flèches, c, q) via go-tty
func RunForgeInteractive(inv InventaireMateriaux) {
	enableWindowsVT()
	recettes := RecettesArmesHumaines()
	if len(recettes) == 0 {
		fmt.Println("Aucune recette disponible.")
		return
	}
	t, err := tty.Open()
	if err != nil {
		fmt.Println("Impossible d'initialiser le TTY:", err)
		return
	}
	defer t.Close()
	enterAltScreen()
	defer exitAltScreen()
	hideCursor()
	defer showCursor()
	selection := 0
	showArmors := false
	modal := ""
	for {
		frame := renderFrameString(recettes, inv, selection, showArmors, modal)
		fmt.Print("\033[H")
		fmt.Print(frame)
		fmt.Print("\033[J")
		if modal != "" {
			_ = readKey(t)
			modal = ""
			continue
		}
		key := readKey(t)
		switch key {
		case "up":
			if selection > 0 {
				selection--
			}
		case "down":
			max := len(recettes)
			if showArmors {
				max = len(RecettesArmures())
			}
			if selection < max-1 {
				selection++
			}
		case "t":
			showArmors = !showArmors
			selection = 0
		case "c":
			if showArmors {
				ars := RecettesArmures()
				if len(ars) == 0 {
					continue
				}
				if selection >= len(ars) {
					selection = 0
				}
				choisie := ars[selection]
				if !inv.AAssez(choisie.Cout) {
					modal = "Matériaux insuffisants pour fabriquer cette armure."
					continue
				}
				inv.Debiter(choisie.Cout)
				arm := objet.CreerArmure(choisie.CleArmure)
				modal = fmt.Sprintf("Fabrication réussie !\nNom: %s\nType: %s\nDéfense: %d | Poids: %d", arm.Nom, arm.Type, arm.EffetDefense, arm.Poids)
			} else {
				if len(recettes) == 0 {
					continue
				}
				if selection >= len(recettes) {
					selection = 0
				}
				choisie := recettes[selection]
				if !inv.AAssez(choisie.Cout) {
					modal = "Matériaux insuffisants pour fabriquer cette arme."
					continue
				}
				inv.Debiter(choisie.Cout)
				arme := objet.CreerArme(choisie.CleArme)
				modal = fmt.Sprintf("Fabrication réussie !\nNom: %s\nAttaque: %d | Poids: %d", arme.Nom, arme.EffetAttaque, arme.Poids)
			}
		case "q":
			resetTerminalVisual()
			return
		}
	}
}

// readKey convertit les séquences de touches en identifiants simples
func readKey(t *tty.TTY) string {
	r, err := t.ReadRune()
	if err != nil {
		return ""
	}
	if r == 0x1b { // ESC
		// attente séquence CSI: ESC [ A/B/C/D
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
	// lettres de commande
	s := strings.ToLower(string(r))
	if s == "c" || s == "q" || s == "t" {
		return s
	}
	return ""
}

// Pause helpers
func messagePause(t *tty.TTY, msg string) {
	fmt.Println(msg)
	_, _ = t.ReadRune()
}

func messagePauseStd(msg string) {
	fmt.Println(msg)
	br := bufio.NewReader(os.Stdin)
	_, _ = br.ReadString('\n')
}

func resetTerminalVisual() {
	// Ensure cursor visible
	showCursor()
	// Leave alternate screen if active
	exitAltScreen()
	// Clear main screen with ANSI (no cls)
	fmt.Print("\033[H\033[J")
}
