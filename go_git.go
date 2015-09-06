package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	fullDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	getData(fullDate)
	//unpackFile("data-2015-09-04-10.gz")
}

func unpackFile(filename string) {
	fi, err := os.Open(filename)
	if err != nil {
		handleError("Error opening file", err)
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		handleError("gunzipping the file", err)
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		handleError("error reading gunzip stream", err)
	}

	//fmt.Println(s)
}

func getData(fullDate string) {

	urls := makeUrlArray(fullDate)

	for i, value := range urls {

		fmt.Println("fetching url", value)

		resp, archiveErr := http.Get(value)

		if resp != nil {
			defer resp.Body.Close()
		}

		if archiveErr != nil {
			handleError("Error getting github archive", archiveErr)
		}

		contents, readErr := ioutil.ReadAll(resp.Body)

		if readErr != nil {
			handleError("Error converting response", readErr)
		}

		fname := makeFileName(fullDate, i)

		fileErr := ioutil.WriteFile(fname, contents, 0644)

		if fileErr != nil {
			handleError("Error writing response to file", fileErr)
		}
	}
}

func makeUrlArray(fullDate string) [24]string {

	baseUrl := makeUrlBase(fullDate)
	urlEnd := ".json.gz"

	var urls [24]string

	for i := 0; i < 24; i++ {

		var buffer bytes.Buffer
		buffer.WriteString(baseUrl)
		buffer.WriteString("-")
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(urlEnd)
		url := buffer.String()

		urls[i] = url
	}

	return urls
}

func makeUrlBase(fullDate string) string {
	split := strings.Split(fullDate, "-")

	var buffer bytes.Buffer
	buffer.WriteString("http://data.githubarchive.org/")
	buffer.WriteString(split[0]) //year
	buffer.WriteString("-")
	buffer.WriteString(split[1]) //month
	buffer.WriteString("-")
	buffer.WriteString(split[2]) //day

	return buffer.String()
}

func makeFileName(fullDate string, i int) string {
	var buffer bytes.Buffer
	buffer.WriteString("data-")
	buffer.WriteString(fullDate)
	buffer.WriteString("-")
	buffer.WriteString(strconv.Itoa(i))
	buffer.WriteString(".gz")

	return buffer.String()

}

func handleError(message string, err error) {
	fmt.Println(message, err)
	os.Exit(1)
}
