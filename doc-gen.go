/**
This Executable creates a .Docx with a Image, Name, SPP, and coords from a csv file
using Unidoc
*/

// Copyright 2017 Baliance. All rights reserved.
package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"

	"github.com/unidoc/unioffice/schema/soo/wml"
)

func main() {
	doc := document.New()

	csvFile, _ := os.Open("Todas las capas3.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		table := doc.AddTable()
		// width of the page
		table.Properties().SetWidthPercent(100)
		// table.Properties().SetLayout(wml.ST_TblLayoutTypeAutofit)
		borders := table.Properties().Borders()
		// thin borders
		borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)

		row := table.AddRow()
		cell := row.AddCell()
		cell.Properties().SetColumnSpan(2)

		if line[8] != "" {
			img1, err := common.ImageFromFile("./Todas las capas/images/" + line[8])
			if err != nil {
				log.Fatalf("unable to create image: %s", err)
			}
			img1ref, err := doc.AddImage(img1)
			if err != nil {
				log.Fatalf("unable to add image to document: %s", err)
			}

			inl, err := cell.AddParagraph().AddRun().AddDrawingInline(img1ref)
			if err != nil {
				log.Fatalf("unable to add inline image: %s", err)
			}
			inl.SetSize(2*measurement.Inch, 2.5*measurement.Inch)

			/*
				anchored, err := cell.AddParagraph().AddRun().AddDrawingAnchored(img1ref)
				if err != nil {
					log.Fatalf("unable to add anchored image: %s", err)
				}
				anchored.SetName("Gopher")
				anchored.SetSize(2*measurement.Inch, 2.5*measurement.Inch)
				anchored.SetOrigin(wml.WdST_RelFromHPage, wml.WdST_RelFromVTopMargin)
				anchored.SetHAlignment(wml.WdST_AlignHCenter)
				anchored.SetTextWrapSquare(wml.WdST_WrapTextBothSides)
			*/

		} else {
			cell.AddParagraph()
		}

		row = table.AddRow()
		cell = row.AddCell()
		cell.Properties().SetColumnSpan(2)
		cell.AddParagraph().AddRun().AddText(line[0])

		row = table.AddRow()
		cell = row.AddCell()
		cell.Properties().SetColumnSpan(2)
		cell.AddParagraph().AddRun().AddText(line[9])

		row = table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText(line[1])
		row.AddCell().AddParagraph().AddRun().AddText(line[2])

		doc.AddParagraph()
	}

	if err := doc.Validate(); err != nil {
		log.Fatalf("error during validation: %s", err)
	}
	doc.SaveToFile("tables.docx")

}
