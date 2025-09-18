package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const saveFile = "save.json"

func SaveGame(gs *GameState) error {
	if gs == nil {
		return errors.New("aucune partie en cours à sauvegarder")
	}
	file, err := os.Create(saveFile)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(gs)
}

func LoadGame() (*GameState, error) {
	file, err := os.Open(saveFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var gs GameState
	dec := json.NewDecoder(file)
	if err := dec.Decode(&gs); err != nil {
		return nil, err
	}
	return &gs, nil
}

func DeleteSave() error {
	if err := os.Remove(saveFile); err != nil {
		return err
	}
	fmt.Println("Sauvegarde supprimée.")
	return nil
}
