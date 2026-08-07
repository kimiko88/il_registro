package signatures

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"time"
)

type CadManifest struct {
	XMLName              xml.Name           `xml:"ManifestoConservazioneCAD"`
	Versione             string             `xml:"versione,attr"`
	CodiceScuola         string             `xml:"CodiceScuola"`
	AnnoScolastico       string             `xml:"AnnoScolastico"`
	DataArchiviazione    string             `xml:"DataArchiviazione"`
	ImprontaCrittografica string             `xml:"ImprontaCrittograficaSHA256"`
	DocumentiConservati  []CadDocumentoItem `xml:"Documenti>Documento"`
}

type CadDocumentoItem struct {
	ID               string `xml:"ID"`
	Tipologia        string `xml:"Tipologia"`
	Oggetto          string `xml:"Oggetto"`
	ImprontaHash     string `xml:"ImprontaHash"`
	DataSottoscrizione string `xml:"DataSottoscrizione"`
}

func GenerateCadPreservationPackage(ctx context.Context, schoolID string, year string) ([]byte, error) {
	if year == "" {
		year = "2025/2026"
	}

	docItems := []CadDocumentoItem{
		{ID: "REG-2025-001", Tipologia: "RegistroDiClasse", Oggetto: "Registro di Classe Anno " + year, ImprontaHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", DataSottoscrizione: time.Now().Format(time.RFC3339)},
		{ID: "SCR-2025-Q1", Tipologia: "VerbaleScrutinio", Oggetto: "Verbale Scrutinio Primo Quadrimestre", ImprontaHash: "f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2", DataSottoscrizione: time.Now().Format(time.RFC3339)},
		{ID: "CERT-2025-ALL", Tipologia: "CertificazioneCompetenze", Oggetto: "Certificazioni delle Competenze DM 742", ImprontaHash: "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3", DataSottoscrizione: time.Now().Format(time.RFC3339)},
	}

	manifest := CadManifest{
		Versione:             "1.0",
		CodiceScuola:         schoolID,
		AnnoScolastico:       year,
		DataArchiviazione:    time.Now().Format(time.RFC3339),
		ImprontaCrittografica: fmt.Sprintf("%x", sha256.Sum256([]byte(year+schoolID))),
		DocumentiConservati:  docItems,
	}

	xmlData, err := xml.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to generate CAD XML manifest: %w", err)
	}
	fullXML := append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), xmlData...)

	// Create ZIP package
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add ManifestoConservazione.xml
	f, err := zipWriter.Create("ManifestoConservazione.xml")
	if err != nil {
		return nil, err
	}
	_, _ = f.Write(fullXML)

	// Add IndiceDati.csv
	f2, err := zipWriter.Create("IndiceDati.csv")
	if err != nil {
		return nil, err
	}
	csvContent := "ID;Tipologia;ImprontaHash;DataSottoscrizione\n"
	for _, d := range docItems {
		csvContent += fmt.Sprintf("%s;%s;%s;%s\n", d.ID, d.Tipologia, d.ImprontaHash, d.DataSottoscrizione)
	}
	_, _ = f2.Write([]byte(csvContent))

	_ = zipWriter.Close()
	return buf.Bytes(), nil
}
