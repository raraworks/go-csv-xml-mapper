package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	csvFilePath := flag.String("csv", "products.csv", "Shopify products CSV file path")
	xmlFilePath := flag.String("xml", "products.xml", "Salidzini products XML file path")
	flag.Parse()
	readCsvAndWriteXmlFile(*csvFilePath, *xmlFilePath)
}

func readCsvAndWriteXmlFile(filePath string, outputFilePath string) {
	//open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	//create file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal("Unable to create output file "+outputFilePath, err)
	}

	csvReader := csv.NewReader(file)
	outputFile.WriteString(xml.Header)
	outputFile.WriteString("<root>\n")
	//read csv line by line
	isFirstLineRead := false
	for {
		//if first line - get headers

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Unable to read input file "+filePath, err)
		}
		//if all good do smth
		if isFirstLineRead == false {
			isFirstLineRead = true
			continue
		}
		if record[52] != "active" {
			continue
		}
		price, _ := strconv.ParseFloat(record[20], 64)
		inStock, _ := strconv.Atoi(record[17])
		line := &SalidziniProduct{
			Name:         record[1],
			Link:         "https://www.rocketbaby.lv/products/" + record[0],
			Price:        price,
			Image:        record[25],
			CategoryFull: record[5],
			//CategoryLink: record[4],
			Brand: record[3],
			//Model:        record[6],
			//Color:        record[7],
			Mpn: record[14],
			//Gtin:         record[9],
			InStock: inStock,
		}
		out, _ := xml.MarshalIndent(line, "", "  ")
		_, err = outputFile.WriteString(string(out))
		if err != nil {
			log.Fatal("Unable to write to output file "+outputFilePath, err)
		}
	}
	outputFile.WriteString("</root>\n")

}

type SalidziniProduct struct {
	XMLName      xml.Name `xml:"item"`
	Name         string   `xml:"name"`
	Link         string   `xml:"link"`
	Price        float64  `xml:"price"`
	Image        string   `xml:"image"`
	CategoryFull string   `xml:"category_full"`
	CategoryLink string   `xml:"category_link"`
	Brand        string   `xml:"brand"`
	Model        string   `xml:"model"`
	Color        string   `xml:"color"`
	Mpn          string   `xml:"mpn"`
	Gtin         string   `xml:"gtin"`
	InStock      int      `xml:"in_stock"`
}
