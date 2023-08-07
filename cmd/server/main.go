package main

import (
	"log"

	metadata "github.com/tashif-hoda/mock-ec2-metadata"
)

func main() {
	metadataService := metadata.NewMetaDataService()
	log.Fatal(metadataService.Serve())
}
