package main

import (
	//"compress/gzip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"reflect"
	"strings"
	"time"
	//"strconv"
)

func main() {
	getData()
}

func getData() {
	fullDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	url := makeUrl(fullDate)
	resp, archiveErr := http.Get(url)
	defer resp.Body.Close()

	if archiveErr != nil {
		handleError("Error getting github archive:", archiveErr)
	}

	contents, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		handleError("Error writing response to file", readErr)
	}

	fname := makeFileName(fullDate)

	ioutil.WriteFile(fname, contents, 0644)
}

func makeUrl(fullDate string) string {
	split := strings.Split(fullDate, "-")

	var buffer bytes.Buffer

	buffer.WriteString("http://data.githubarchive.org/")
	buffer.WriteString(split[0]) //year
	buffer.WriteString("-01-")
	buffer.WriteString(split[1]) //month
	buffer.WriteString("-")
	buffer.WriteString(split[2]) //day
	buffer.WriteString(".json.gz")

	return buffer.String()
}

func makeFileName(fullDate string) string {
	var buffer bytes.Buffer
	buffer.WriteString("data-")
	buffer.WriteString(fullDate)
	buffer.WriteString(".gz")

	return buffer.String()
}

func handleError(message string, err error) {
	fmt.Println("Error getting github archive:", err)
	os.Exit(1)
}
