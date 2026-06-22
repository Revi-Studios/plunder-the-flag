package lib

import (
	"encoding/gob"
	"log"
	"os"
)

type GameSave struct {
	Debug     bool
	DebugMenu bool

	WorldScale float64
	TextScale  float64
}

func (save *GameSave) Save() {
	log.Println("game: trying to save...")
	file, err := os.Create("save.dat")
	if err != nil {
		panic("Error: game failed to save:" + err.Error())
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	encoder.Encode(save)
	log.Println("game: saved data.")
}

func (save *GameSave) Load() {
	log.Println("game: loading save file...")
	file, err := os.OpenFile("save.dat", os.O_CREATE, 0666)
	if err != nil {
		panic("Error: failed to load save file: " + err.Error())
	}
	info, err := file.Stat()
	if err != nil {
		panic("Error: failed to load save file info: " + err.Error())
	}

	if info.Size() == 0 {
		log.Println("game: save file emty, none is loaded")
		return
	}

	defer file.Close()

	var data GameSave

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		panic("Error: failed to decode save file: " + err.Error())
	}
	*save = data
	log.Println("game: loaded save file:", data)
}
