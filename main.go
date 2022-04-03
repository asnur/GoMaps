package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	tablewriter "github.com/olekukonko/tablewriter"
)

type Response struct {
	Result []Result `json:"results"`
	Status string   `json:"status"`
}

type Result struct {
	Name     string `json:"name"`
	Address  string `json:"formatted_address"`
	Geometry Geom   `json:"geometry"`
}

type Geom struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
}

func main() {
	InputData()
}

func InputData() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("---------------------")
	fmt.Println("Google Maps Scraper by Asnur")
	fmt.Println("---------------------")
	fmt.Printf("Input Latitude Coordinates: ")
	lat, _ := reader.ReadString('\n')
	fmt.Printf("Input Longitude Coordinates: ")
	lng, _ := reader.ReadString('\n')
	fmt.Printf("Input Query: ")
	query, _ := reader.ReadString('\n')
	fmt.Println("---------------------")
	fetchData(lat, lng, query)
}

func fetchData(lat, lng, query string) {
	lats := TrimSpaceNewlineInString(lat)
	lngs := TrimSpaceNewlineInString(lng)
	queries := TrimSpaceNewlineInString(query)
	url := "https://maps.googleapis.com/maps/api/place/textsearch/json?query=" + strings.ReplaceAll(queries, " ", "") + "&location=" + lats + "," + lngs + "&key=AIzaSyBIIfuR8-AJIrG2tScD4zW3Fmm4Ret3wX4"
	fmt.Println("Fetching Data...")
	var data Response
	response, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	json.Unmarshal(bytes, &data)
	var tableData [][]string
	for _, result := range data.Result {
		var row []string
		row = append(row, result.Name)
		row = append(row, result.Address)
		row = append(row, fmt.Sprintf("%f", result.Geometry.Location.Lat))
		row = append(row, fmt.Sprintf("%f", result.Geometry.Location.Lng))
		tableData = append(tableData, row)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Address", "Latitude", "Longitude"})
	for _, v := range tableData {
		table.Append(v)
	}
	table.Render()
}

func TrimSpaceNewlineInString(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(s, "")
}
