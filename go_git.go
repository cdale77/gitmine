package main

import (
	//"compress/gzip"
	"bytes"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	//"reflect"
	//"net/url"
	"strconv"
	"strings"
	"time"
	//"strconv"
)

func main() {

	fullDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	getData(fullDate)
}

func getData(fullDate string) {

	urls := makeUrlArray(fullDate)
	fmt.Println(urls)

	for _, value := range urls {
		fmt.Println("getting url", value)
		resp, archiveErr := http.Get(value)

		if resp != nil {
			defer resp.Body.Close()
		}

		if archiveErr != nil {
			handleError("Error getting github archive:", archiveErr)
		}

	}

	/*
		contents, readErr := ioutil.ReadAll(resp.Body)

		if readErr != nil {
			handleError("Error writing response to file", readErr)
		}

		fname := makeFileName(fullDate)

		ioutil.WriteFile(fname, contents, 0644)
	*/
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
	buffer.WriteString("http//data.githubarchive.org/")
	buffer.WriteString(split[0]) //year
	buffer.WriteString("-")
	buffer.WriteString(split[1]) //month
	buffer.WriteString("-")
	buffer.WriteString(split[2]) //day

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
