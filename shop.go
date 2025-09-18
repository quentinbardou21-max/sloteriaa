package main

import (
	"fmt"

	"sloteriaa/struct/forgeron"
)

var shopPricesBuy = map[string]int{
	"potion":           40,  // +20 PV
	"potion majeure":   80,  // +50 PV
	"potion force":     60,  // +2 Force temporaire (3 combats)
	"potion agilite":   60,  // +2 Agilité temporaire (3 combats)
	"potion endurance": 60,  // +2 Endurance temporaire (3 combats)
	"antidote":         30,  // Guérit poison/brûlure/saignement
	"elixir vie":       120, // +100 PV
}

var materialPrices = map[forgeron.Materiau]int{
	forgeron.Fer:            20,
	forgeron.Bois:           10,
	forgeron.Cuir:           15,
	forgeron.EssenceMagique: 100,
}

var lootSellPrices = map[string]int{
	"Griffes souillées": 15,
	"Massue brute":      20,
	"Lance brisée":      18,
	"Épée osseuse":      24,
	"Hache tronquée":    22,
	"Glaive sauvage":    26,
	"Masse rituelle":    28,
	"Faux de brume":     30,
}

func EnterShop(gs *GameState) {
	for {
		header := fmt.Sprintf("Marché — Or %d", gs.Joueur.Argent)
		idx, cancelled := selectWithArrows(header, []string{"Acheter matériaux", "Acheter potions", "Vendre objets", "Sortir du marché"})
		if cancelled {
			return
		}
		switch idx {
		case 0:
			buyMaterials(gs)
		case 1:
			buyConsumables(gs)
		case 2:
			sellLoot(gs)
		case 3:
			return
		}
	}
}

func buyMaterials(gs *GameState) {
	keys := []forgeron.Materiau{forgeron.Fer, forgeron.Bois, forgeron.Cuir, forgeron.EssenceMagique}
	opts := make([]string, 0, len(keys))
	for _, m := range keys {
		opts = append(opts, fmt.Sprintf("%s (%d or)", m, materialPrices[m]))
	}
	idx, cancelled := selectWithArrows("Matériaux à acheter:", opts)
	if cancelled {
		return
	}
	m := keys[idx]
	price := materialPrices[m]
	q := promptQuantity("Quantité:")
	if q <= 0 {
		return
	}
	cost := price * q
	if gs.Joueur.Argent < cost {
		fmt.Println("Pas assez d'or.")
		return
	}
	gs.Joueur.Argent -= cost
	gs.Mats[string(m)] += q
	fmt.Printf("Acheté %d x %s.\n", q, m)
	// rester dans le sous-menu matériaux
	buyMaterials(gs)
}

func buyConsumables(gs *GameState) {
	items := []string{"potion", "potion majeure", "potion force", "potion agilite", "potion endurance", "antidote", "elixir vie"}
	opts := make([]string, 0, len(items))
	for _, it := range items {
		label := fmt.Sprintf("%s (%d or)", it, shopPricesBuy[it])
		// Ajouter la description de chaque potion
		switch it {
		case "potion":
			label += " (+20 PV)"
		case "potion majeure":
			label += " (+50 PV)"
		case "potion force":
			label += " (+2 Force, 3 combats)"
		case "potion agilite":
			label += " (+2 Agilité, 3 combats)"
		case "potion endurance":
			label += " (+2 Endurance, 3 combats)"
		case "antidote":
			label += " (Guérit statuts)"
		case "elixir vie":
			label += " (+100 PV)"
		}
		opts = append(opts, label)
	}
	idx, cancelled := selectWithArrows("Consommables:", opts)
	if cancelled {
		return
	}
	item := items[idx]
	price := shopPricesBuy[item]
	q := promptQuantity("Quantité:")
	if q <= 0 {
		return
	}
	cost := price * q
	if gs.Joueur.Argent < cost {
		fmt.Println("Pas assez d'or.")
		return
	}
	gs.Joueur.Argent -= cost
	for i := 0; i < q; i++ {
		gs.Joueur.Inventaire = append(gs.Joueur.Inventaire, item)
	}
	fmt.Printf("Acheté %d x %s.\n", q, item)
	// rester dans le sous-menu consommables
	buyConsumables(gs)
}

func sellLoot(gs *GameState) {
	sellableIdx := []int{}
	opts := []string{}
	for idx, name := range gs.Joueur.Inventaire {
		if _, ok := lootSellPrices[name]; ok {
			opts = append(opts, fmt.Sprintf("%s (vend %d or)", name, lootSellPrices[name]))
			sellableIdx = append(sellableIdx, idx)
		}
	}
	if len(sellableIdx) == 0 {
		fmt.Println("Rien à vendre.")
		fmt.Println("(Appuyez sur Entrée pour revenir)")
		attendreEntree()
		return
	}
	choice, cancelled := selectWithArrows("Vendre:", opts)
	if cancelled {
		return
	}
	invIdx := sellableIdx[choice]
	name := gs.Joueur.Inventaire[invIdx]
	price := lootSellPrices[name]
	gs.Joueur.Inventaire = append(gs.Joueur.Inventaire[:invIdx], gs.Joueur.Inventaire[invIdx+1:]...)
	gs.Joueur.Argent += price
	fmt.Printf("Vendu %s pour %d or.\n", name, price)
	// rester dans le sous-menu vente
	sellLoot(gs)
}

func promptQuantity(label string) int {
	opts := make([]string, 20)
	for i := 1; i <= 20; i++ {
		opts[i-1] = fmt.Sprintf("%s %d", label, i)
	}
	idx, _ := selectWithArrows("", opts)
	return idx + 1
}
