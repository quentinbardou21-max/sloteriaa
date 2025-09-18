package main

import (
	"fmt"
	"strings"

	"sloteriaa/struct/forgeron"
	"sloteriaa/struct/objet"
)

func EnterForgeSimple(gs *GameState) {
	for {
		header := fmt.Sprintf("Forge — Or %d", gs.Joueur.Argent)
		idx, cancelled := selectWithArrows(header, []string{"Forger une arme", "Forger une armure", "Sortir de la forge"})
		if cancelled || idx == 2 {
			return
		}
		switch idx {
		case 0:
			forgeSelectWeapon(gs)
		case 1:
			forgeSelectArmor(gs)
		}
	}
}

func forgeSelectWeapon(gs *GameState) {
	recs := forgeron.RecettesArmesHumaines()
	if len(recs) == 0 {
		fmt.Println("Aucune recette.")
		attendreEntree()
		return
	}
	opts := make([]string, len(recs))
	for i, r := range recs {
		gold := r.Cout[forgeron.Or]
		matsStr := formatMaterials(r.Cout)
		if matsStr != "-" {
			// Afficher les stats de l'arme
			arme := objet.CreerArme(r.CleArme)
			opts[i] = fmt.Sprintf("%s (ATK %d) — Coût: %d or | Mat: %s", r.NomAffiche, arme.EffetAttaque, gold, matsStr)
		} else {
			opts[i] = fmt.Sprintf("%s — Coût: %d or", r.NomAffiche, gold)
		}
	}
	for {
		sel, cancelled := selectWithArrows("Choisissez une arme:", opts)
		if cancelled {
			return
		}
		r := recs[sel]
		craftWithCost(gs, r.Cout, func() {
			a := objet.CreerArme(r.CleArme)
			gs.Joueur.Inventaire = append(gs.Joueur.Inventaire, a.Nom)
			fmt.Printf("Forgé: %s (ATK %d, Poids %d) — Ajouté à l'inventaire\n", a.Nom, a.EffetAttaque, a.Poids)
			attendreEntree()
		})
		// boucle: rester dans la liste des armes après craft
	}
}

func forgeSelectArmor(gs *GameState) {
	recs := forgeron.RecettesArmures()
	if len(recs) == 0 {
		fmt.Println("Aucune recette.")
		attendreEntree()
		return
	}
	opts := make([]string, len(recs))
	for i, r := range recs {
		gold := r.Cout[forgeron.Or]
		matsStr := formatMaterials(r.Cout)
		if matsStr != "-" {
			// Afficher les stats de l'armure
			armure := objet.CreerArmure(r.CleArmure)
			opts[i] = fmt.Sprintf("%s (DEF %d) — Coût: %d or | Mat: %s", r.NomAffiche, armure.EffetDefense, gold, matsStr)
		} else {
			// Afficher les stats de l'armure même sans matériaux
			armure := objet.CreerArmure(r.CleArmure)
			opts[i] = fmt.Sprintf("%s (DEF %d) — Coût: %d or", r.NomAffiche, armure.EffetDefense, gold)
		}
	}
	for {
		sel, cancelled := selectWithArrows("Choisissez une armure:", opts)
		if cancelled {
			return
		}
		r := recs[sel]
		craftWithCost(gs, r.Cout, func() {
			ar := objet.CreerArmure(r.CleArmure)
			gs.Joueur.Inventaire = append(gs.Joueur.Inventaire, ar.Nom)
			fmt.Printf("Forgé: %s (DEF %d, Poids %d) — Ajouté à l'inventaire\n", ar.Nom, ar.EffetDefense, ar.Poids)
			attendreEntree()
		})
		// boucle: rester dans la liste des armures après craft
	}
}

// Helper: check mats + gold, debit, then run success action
func craftWithCost(gs *GameState, cout forgeron.Cout, onSuccess func()) {
	// Convert mats to forgeron inventory for checking and debiting
	inv := matsToForgeron(gs.Mats)
	// Build materials list and gold cost
	mats := []string{}
	gold := 0
	for _, m := range []forgeron.Materiau{forgeron.Fer, forgeron.Bois, forgeron.Cuir, forgeron.EssenceMagique, forgeron.Or} {
		if q, ok := cout[m]; ok {
			if m == forgeron.Or {
				gold = q
				continue
			}
			mats = append(mats, fmt.Sprintf("%s x%d", m, q))
		}
	}
	// Check materials availability
	hasMats := inv.AAssez(cout)
	hasGold := gs.Joueur.Argent >= gold
	if !hasMats || !hasGold {
		fmt.Println("Ressources insuffisantes.")
		fmt.Println("Coût: ")
		if len(mats) > 0 {
			fmt.Println("  Matériaux:", strings.Join(mats, ", "))
		} else {
			fmt.Println("  Matériaux: -")
		}
		fmt.Printf("  Or: %d\n", gold)
		attendreEntree()
		return
	}
	// Debit
	inv.Debiter(cout)
	gs.Mats = forgeronToMats(inv)
	gs.Joueur.Argent -= gold
	// Success
	onSuccess()
}

// formatMaterials returns a human string for non-gold material costs
func formatMaterials(cout forgeron.Cout) string {
	parts := []string{}
	order := []forgeron.Materiau{forgeron.Fer, forgeron.Bois, forgeron.Cuir, forgeron.EssenceMagique}
	for _, m := range order {
		if q, ok := cout[m]; ok && q > 0 {
			parts = append(parts, fmt.Sprintf("%s x%d", m, q))
		}
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, ", ")
}
