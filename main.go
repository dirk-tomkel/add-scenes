package main

import (
	"encoding/csv"
	"flag"

	"log"
	"os"

	"github.com/beevik/etree"
)

func main() {
	kdenliveFile := flag.String("k", "", "Dateiname der Kdenlive Projektdatei")
	csvFilename := flag.String("c", "", "Dateiname der pyscenedetect CSV-Datei")
	flag.Parse()

	// -- lese kdenliveFile ein --
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(*kdenliveFile); err != nil {
		log.Fatal("Fehler: Kann datei nicht einlesen:", err)
	}

	// -- lese CSV datei ein --
	cf, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatal("Fehler: Kann CSV Datei nicht einlesen", err)
	}
	defer cf.Close()

	reader := csv.NewReader(cf)

	// WICHTIG: Erlaube unterschiedliche Spaltenanzahl
	reader.FieldsPerRecord = -1
	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Fehler: Kann CSV parsen", err)
	}

	// Finde Playlisten mit Einträgen
	var kdenliveId string
	playlists := doc.FindElements("//playlist[entry]")
	for _, playlist := range playlists {
		// lösche bestehende Einträge
		playlistId := playlist.SelectAttrValue("id", "keine Playlist")
		if playlistId != "main_bin" {
			log.Printf("Lösche Einträge in %s.\n", playlistId)
			producerChain := playlist.FindElement("//entry[@producer]").SelectAttrValue("producer", "keine Ahnung")
			playlist.Child = nil

			// -- Schreibe Szenen-Einträge aus CSV Datei in Playlists
			// ignoreire die ersten beiden zeilen
			var sceneCount int
			for _, row := range csvData[2:] {
				start := row[2]
				ende := row[5]
				addEntry(playlist, start, ende, producerChain, kdenliveId)
				sceneCount++
			}

			log.Printf("Füge %d Szenen zu producer=%s hinzu.\n", sceneCount, producerChain)
		}
	}

	// -- Erstelle die Kdenlive AUsgabe-Datei --
	finalOutputFilename := "scenes-" + *kdenliveFile
	doc.Indent(2)
	err = doc.WriteToFile(finalOutputFilename)
	if err != nil {
		log.Fatal("Fehler: kann XML nicht schreiben.", err)
	}
}

func addEntry(playlist *etree.Element, in, out, producer, id string) {
	// 1. <entry in="..." out="..." producer="...">
	entry := playlist.CreateElement("entry")
	entry.CreateAttr("in", in)
	entry.CreateAttr("out", out)
	entry.CreateAttr("producer", producer)

	// 2. <property name="kdenlive:id"><id-value></property>
	property := entry.CreateElement("property")
	property.CreateAttr("name", "kdenlive:id")

	// Textinhalt der property
	property.CreateText(id)
}
