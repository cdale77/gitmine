package main

import (
	//"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"reflect"
	"bytes"
	"time"
	//"strconv"
)

func main() {
	getData()
}

func getData() {
	year := time.Now().AddDate(0, 0, -1).Format("2006")
	date := time.Now().AddDate(0, 0, -1).Format("01-02")
	fmt.Println(year)
	fmt.Println(date)

	var buffer bytes.Buffer
	buffer.WriteString("http://data.githubarchive.org/")
	buffer.WriteString(year)
	buffer.WriteString("-01-")
	buffer.WriteString(date)
	buffer.WriteString(".json.gz")

	url := buffer.String()
	fmt.Println(url)

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Error getting github archive:", err)
		os.Exit(1)
	} else {

		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error writing response to file", err)
			os.Exit(1)
		}

		var buffer bytes.Buffer
		buffer.WriteString("data-")
		buffer.WriteString(year)
		buffer.WriteString("-")
		buffer.WriteString(date)
		buffer.WriteString(".gz")

		fname := buffer.String()

		ioutil.WriteFile(fname, contents, 0644)
	}
}
