package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type CommitCommit struct {
	Message string
}

type CommitPayload struct {
	Size    int
	Commits []CommitCommit
}

type Commit struct {
	Id      string
	Type    string
	Payload CommitPayload
}

func main() {

	fullDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	getData(fullDate)
}

func parseFile(fName string) {
	// https://groups.google.com/forum/#!topic/golang-nuts/GjIkryuCyAY
	fileOS, err := os.Open(fName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rollupParse: Can't open %s: error: %s\n", fName, err)
		os.Exit(1)
	}

	fileGzip, err := gzip.NewReader(fileOS)
	if err != nil {
		fmt.Printf("The file %v is not in gzip format.\n", fName)
		os.Exit(1)
	}

	fileRead := bufio.NewReader(fileGzip)
	i := 0
	for {
		line, err := fileRead.ReadString('\n')
		if err != nil {
			fmt.Println("Did not get a new line.")
			break
		}

		var commit Commit

		jsonErr := json.Unmarshal([]byte(line), &commit)
		if jsonErr != nil {
			fmt.Println("Could not parse json.")
			fmt.Println(jsonErr)
		} else {
			if commit.Type == "PushEvent" {
				fmt.Println("Found a PushEvent commit")
				if commit.Payload.Size > 0 {
					fmt.Println(commit.Payload.Commits[0])

				}
			}
		}

		i++
	}
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

		parseFile(fname)

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
