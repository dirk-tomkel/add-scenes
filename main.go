package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"

	"log"
	"os"

	"github.com/antchfx/xmlquery"
)

func main() {
	kdenliveFile := flag.String("k", "", "Dateiname der Kdenlive Projektdatei")
	csvFilename := flag.String("c", "", "Dateiname der pyscenedetect CSV-Datei")
	flag.Parse()

	// -- lese kdenliveFile ein --
	kf, err := os.Open(*kdenliveFile)
	if err != nil {
		log.Fatal("Fehler: Kann datei nicht einlesen:", err)
	}

	// -- XML parsen --
	doc, err := xmlquery.Parse(kf)
	if err != nil {
		log.Fatal("Fehler: Kann XML nicht parsen:", err)
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
	playlists := xmlquery.Find(doc, "//playlist[entry and @id != 'main_bin']")
	for _, playlist := range playlists {
		// lösche bestehende Einträge
		playlistId := playlist.SelectAttr("id")
		log.Printf("Lösche Einträge in %s.\n", playlistId)
		producerChain := xmlquery.FindOne(playlist, "//entry[@producer]").SelectAttr("producer")
		playlist.FirstChild = nil
		playlist.LastChild = nil

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

	// -- Erstelle die Kdenlive AUsgabe-Datei --
	finalOutputFilename := "scenes-" + *kdenliveFile
	outFile, err := os.Create(finalOutputFilename)
	if err != nil {
		log.Fatal("Fehler: kann Szenen-Datei nicht erstellen.", err)
	}
	defer outFile.Close()

	err = doc.WriteWithOptions(outFile)
	if err != nil {
		log.Fatal("Fehler: kann XML nicht schreiben.", err)
	}
}

func addEntry(playlist *xmlquery.Node, in, out, producer, id string) {
	entry := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "entry",
		Attr: []xmlquery.Attr{
			{Name: xml.Name{Local: "in"}, Value: in},
			{Name: xml.Name{Local: "out"}, Value: out},
			{Name: xml.Name{Local: "producer"}, Value: producer},
		},
	}

	property := &xmlquery.Node{
		Type: xmlquery.ElementNode,
		Data: "property",
		Attr: []xmlquery.Attr{
			{Name: xml.Name{Local: "name"}, Value: "kdenlive:id"},
		},
		FirstChild: &xmlquery.Node{
			Type: xmlquery.TextNode,
			Data: id,
		},
	}

	xmlquery.AddChild(entry, property)
	xmlquery.AddChild(playlist, entry)
}
