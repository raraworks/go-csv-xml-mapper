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
	columnMap := make(map[string]int)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Unable to read input file "+filePath, err)
		}
		//if all good do smth
		//if first line - map header value to column index
		if isFirstLineRead == false {
			for i, v := range record {
				columnMap[v] = i
			}
			isFirstLineRead = true
			continue
		}
		//exclude inactive products
		if record[columnMap["Status"]] != "active" {
			continue
		}
		//exclude delivery product
		if record[columnMap["Type"]] == "delivery" {
			continue
		}
		price, _ := strconv.ParseFloat(record[columnMap["Variant Price"]], 64)
		//TODO: create ternary operator function
		inStock, _ := strconv.Atoi(record[columnMap["Variant Inventory Qty"]])
		if inStock < 1 {
			inStock = 1
		}
		line := &SalidziniProduct{
			Name:                record[columnMap["Title"]],
			Link:                "https://www.rocketbaby.lv/products/" + record[columnMap["Handle"]],
			Price:               price,
			Image:               record[columnMap["Image Src"]],
			CategoryFull:        record[columnMap["Type"]],
			Brand:               record[columnMap["Vendor"]],
			Mpn:                 record[columnMap["Variant SKU"]],
			InStock:             inStock,
			DeliveryDpdPakuBode: 2.99,
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
	XMLName             xml.Name `xml:"item"`
	Name                string   `xml:"name"`
	Link                string   `xml:"link"`
	Price               float64  `xml:"price"`
	Image               string   `xml:"image"`
	CategoryFull        string   `xml:"category_full"`
	CategoryLink        string   `xml:"category_link"`
	Brand               string   `xml:"brand"`
	Model               string   `xml:"model"`
	Color               string   `xml:"color"`
	Mpn                 string   `xml:"mpn"`
	DeliveryDpdPakuBode float64  `xml:"delivery_dpd_paku_bode"`
	InStock             int      `xml:"in_stock"`
}
