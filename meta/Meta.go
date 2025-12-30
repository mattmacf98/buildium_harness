package meta

import (
	"encoding/json"
	"fmt"
	"os"
)

type Meta struct {
	Stage      int    `json:"stage"`
	Entrypoint string `json:"entrypoint"`
	Path       string `json:"path"`
}

func NewMeta(path string) *Meta {
	metaBytes, err := os.ReadFile(path + "/meta.json")
	if err != nil {
		fmt.Println("Error reading meta file:", err)
		os.Exit(1)
	}
	var meta Meta
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		fmt.Println("Error unmarshalling meta file:", err)
		os.Exit(1)
	}
	meta.Path = path
	return &meta
}
