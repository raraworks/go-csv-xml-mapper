package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"os"
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
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Unable to read input file "+filePath, err)
		}
		//if all good do smth
		line := &SalidziniProduct{ProductName: record[0]}
		out, _ := xml.MarshalIndent(line, "", "  ")
		_, err = outputFile.WriteString(string(out))
		if err != nil {
			log.Fatal("Unable to write to output file "+outputFilePath, err)
		}
	}
	outputFile.WriteString("</root>\n")

}

type SalidziniProduct struct {
	XMLName     xml.Name `xml:"item"`
	ProductName string   `xml:"name"`
}
