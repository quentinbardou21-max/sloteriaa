package objet

import "fmt"

type TypeArmure string

const (
	TypeCasque    TypeArmure = "Casque"
	TypePlastron  TypeArmure = "Plastron"
	TypePantalon  TypeArmure = "Pantalon"
	TypeChaussure TypeArmure = "Chaussures"
)

type Armure struct {
	Nom          string
	Description  string
	Type         TypeArmure
	EffetDefense int
	Poids        int
}

// Création d'une armure selon son nom
func CreerArmure(nom string) Armure {
	switch nom {
	// Casques
	case "CasqueCuir":
		return Armure{"Casque en cuir", "Casque léger en cuir", TypeCasque, 5, 2}
	case "CasqueCuirRenforce":
		return Armure{"Casque en cuir renforcé", "Casque plus résistant", TypeCasque, 8, 3}
	case "CasqueFer":
		return Armure{"Casque en fer", "Casque solide en fer", TypeCasque, 15, 5}
	case "CasqueFerRenforce":
		return Armure{"Casque en fer renforcé", "Casque très solide", TypeCasque, 20, 6}

	// Plastrons
	case "PlastronCuir":
		return Armure{"Plastron en cuir", "Protection souple pour le torse", TypePlastron, 10, 5}
	case "PlastronCuirRenforce":
		return Armure{"Plastron en cuir renforcé", "Plastron plus résistant", TypePlastron, 15, 6}
	case "PlastronFer":
		return Armure{"Plastron en fer", "Plastron métallique solide", TypePlastron, 30, 12}
	case "PlastronFerRenforce":
		return Armure{"Plastron en fer renforcé", "Plastron très solide", TypePlastron, 40, 14}

	// Pantalons
	case "PantalonCuir":
		return Armure{"Pantalon en cuir", "Pantalon léger offrant une protection modérée", TypePantalon, 8, 3}
	case "PantalonCuirRenforce":
		return Armure{"Pantalon en cuir renforcé", "Pantalon plus résistant", TypePantalon, 12, 4}
	case "PantalonFer":
		return Armure{"Pantalon en fer", "Pantalon blindé", TypePantalon, 20, 8}
	case "PantalonFerRenforce":
		return Armure{"Pantalon en fer renforcé", "Pantalon très solide", TypePantalon, 25, 10}

	// Chaussures
	case "BottesCuir":
		return Armure{"Bottes en cuir", "Bottes légères offrant un minimum de protection", TypeChaussure, 5, 2}
	case "BottesCuirRenforce":
		return Armure{"Bottes en cuir renforcé", "Bottes plus résistantes", TypeChaussure, 8, 3}
	case "BottesFer":
		return Armure{"Bottes en fer", "Bottes solides en fer", TypeChaussure, 15, 5}
	case "BottesFerRenforce":
		return Armure{"Bottes en fer renforcé", "Bottes très solides", TypeChaussure, 20, 6}

	default:
		fmt.Println("Armure inconnue, création d'un casque en cuir par défaut")
		return Armure{"Casque en cuir", "Casque léger en cuir", TypeCasque, 5, 2}
	}
}

// Affiche les infos d'une armure
func AfficherArmure(o Armure) {
	fmt.Printf("Nom : %s\nDescription : %s\nType : %s\nDéfense : %d\nPoids : %d\n\n",
		o.Nom, o.Description, o.Type, o.EffetDefense, o.Poids)
}

// Fonction pour calculer la défense totale d'un set d'armure
// Bonus de 10% si le joueur a équipé un set complet du même matériau
func CalculerDefenseTotale(casque, plastron, pantalon, chaussures Armure) int {
	total := casque.EffetDefense + plastron.EffetDefense + pantalon.EffetDefense + chaussures.EffetDefense

	// Vérifier si c'est un set complet (même matériau)
	if estSetComplet(casque, plastron, pantalon, chaussures) {
		bonus := total / 10 // 10% de bonus
		total += bonus
		fmt.Printf("Bonus de set complet appliqué : +%d défense\n", bonus)
	}

	return total
}

// Vérifie si les 4 pièces sont du même matériau
func estSetComplet(casque, plastron, pantalon, chaussures Armure) bool {
	// On considère que le matériau est le 2e mot du nom (Cuir, CuirRenforce, Fer, FerRenforce)
	getMat := func(o Armure) string {
		if len(o.Nom) > 0 {
			return o.Nom[len(o.Nom)-4:] // approche simple pour différencier Fer/FerRenforce
		}
		return ""
	}

	m1 := getMat(casque)
	m2 := getMat(plastron)
	m3 := getMat(pantalon)
	m4 := getMat(chaussures)

	return m1 == m2 && m2 == m3 && m3 == m4
}
