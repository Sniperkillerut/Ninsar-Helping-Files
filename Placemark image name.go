/***
This executable will search inside an extracted Goole Maps KML file with folder image
and correlate all Placemarks with an Image Name and return a CSV file with them
*/

package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Kml is the struct which contains the complete array of all Points in the file
type Kml struct {
	XMLName  xml.Name `xml:"kml"`
	Document Document `xml:"Document"`
}

// Document is the struct that contains the Kml Document array
type Document struct {
	XMLName xml.Name `xml:"Document"`
	Folders []Folder `xml:"Folder"`
}

//Folder contains the Kml folders inside a document
type Folder struct {
	XMLName    xml.Name    `xml:"Folder"`
	PlaceMarks []Placemark `xml:"Placemark"`
}

// Placemark is a simple struct which contains the kml placemarks
type Placemark struct {
	XMLName      xml.Name     `xml:"Placemark"`
	Name         string       `xml:"name"`
	ExtendedData ExtendedData `xml:"ExtendedData"`
}

//ExtendedData is a struct that contains all the Extended data of our Placemark
type ExtendedData struct {
	XMLName    xml.Name   `xml:"ExtendedData"`
	SchemaData SchemaData `xml:"SchemaData`
}

//SchemaData is a struct that contains Simpledata array
type SchemaData struct {
	XMLName     xml.Name `xml:"SchemaData"`
	SimpleDatas []string `xml:"SimpleData"`
}

func main() {

	// Open our kmlFile
	kmlFile, err := os.Open("doc.kml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened dock.kml")
	// defer the closing of our kmlFile so that we can parse it later on
	defer kmlFile.Close()

	// read our opened kmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(kmlFile)

	// we initialize our Kml array
	var kml Kml
	// we unmarshal our byteArray which contains our
	// kmlFiles content into 'kml' which we defined above
	xml.Unmarshal(byteValue, &kml)

	//we create a CSV to be filled with the result of the operation
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// we iterate through every Placemark within our Documents array and
	// extract the image name and the placemark name
	//then we add them to a new CSV file so that it is simplier to link/correlate the imgName with the Placemark.Name
	for i := 0; i < len(kml.Document.Folders); i++ {
		for j := 0; j < len(kml.Document.Folders[i].PlaceMarks); j++ {
			if !strings.Contains(kml.Document.Folders[i].PlaceMarks[j].Name, "Track") && !strings.Contains(kml.Document.Folders[i].PlaceMarks[j].Name, "V574") {
				// fmt.Println("Placemark name: " + kml.Document.Folders[i].PlaceMarks[j].Name)
				temp := strings.TrimSuffix(kml.Document.Folders[i].PlaceMarks[j].ExtendedData.SchemaData.SimpleDatas[1], "\" /><br />")
				temp = strings.TrimPrefix(temp, "<img src=\"images/")
				// fmt.Println("Placemark img: " + temp)
				temp2 := []string{kml.Document.Folders[i].PlaceMarks[j].Name, temp}
				err := writer.Write(temp2)
				if err != nil {
					log.Fatal("Cannot write to file", err)
				}
			}

		}
	}
	fmt.Println("Done!")

}
