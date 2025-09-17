package main

import (
	"sloteriaa/struct/forgeron"
)

func main() {
	inv := forgeron.InventaireMateriaux{
		forgeron.Fer:            8,
		forgeron.Bois:           6,
		forgeron.Cuir:           4,
		forgeron.EssenceMagique: 2,
	}
	// Mode interactif avec flèches (sans Entrée)
	forgeron.RunForgeInteractive(inv)
}
