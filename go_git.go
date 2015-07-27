package main

import (
	//"compress/gzip"
	"bytes"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	//"reflect"
	"net/url"
	"strings"
	"time"
	//"strconv"
)

func main() {

	fullDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	getData(fullDate)
}

func getData(fullDate string) {

	urlString := makeUrl(fullDate)

	fmt.Println(urlString)
	/*
		resp, archiveErr := http.Get(urlString)
		defer resp.Body.Close()
	*/

	req := &http.Request{
		Method: "GET",
		Host:   "data.githubarchive.org",
		URL: &url.URL{
			Host:   "ignored",
			Scheme: "http",
			Opaque: urlString,
		},
		Header: http.Header{
			"User-Agent": {"golang"},
		},
	}
	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("whoops")
		fmt.Println(err)
	}

	if resp != nil {
		fmt.Println(resp)
	}

	/*
			if archiveErr != nil {
				handleError("Error getting github archive:", archiveErr)
			}


		contents, readErr := ioutil.ReadAll(resp.Body)

		if readErr != nil {
			handleError("Error writing response to file", readErr)
		}

		fname := makeFileName(fullDate)

		ioutil.WriteFile(fname, contents, 0644)
	*/
}

func makeUrl(fullDate string) string {
	split := strings.Split(fullDate, "-")

	var buffer bytes.Buffer

	buffer.WriteString("//data.githubarchive.org/")
	buffer.WriteString(split[0]) //year
	buffer.WriteString("-")
	buffer.WriteString(split[1]) //month
	buffer.WriteString("-")
	buffer.WriteString(split[2])   //day
	buffer.WriteString("-{0..23}") //hours
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
