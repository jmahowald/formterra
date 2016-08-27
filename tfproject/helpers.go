package tfproject

import (
	"log"
	"path/filepath"
)

//go:generate go-bindata -pkg tfproject -o assets.go assets/

func loadAsset(path string) []byte {
	assetPath := filepath.Join("assets", path)
	templateBytes, err := Asset(filepath.Join("assets", path))
	if err != nil {
		log.Fatalln("Unable to retrieve asset file %s", assetPath, err)
	}
	return templateBytes
}
