package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/cdale77/go_git/Godeps/_workspace/src/github.com/melvinmt/firebase"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type StoredCommit struct {
	Date    string
	Login   string
	Avatar  string
	Message string
	Url     string
}

type CommitCommit struct {
	Message string
	Url     string
}

type CommitPayload struct {
	Size    int
	Commits []CommitCommit
}

type CommitActor struct {
	Login      string
	Avatar_url string
}

type Commit struct {
	Type       string
	Created_at string
	Actor      CommitActor
	Payload    CommitPayload
}

func main() {
	fullDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	getData(fullDate)
}

func storeCommit(commit Commit) bool {
	fmt.Println("storing commit:")
	fmt.Println(commit)
	authToken := os.Getenv("FIREBASE_SECRET")

	url := os.Getenv("FIREBASE_URL")

	fireBase := firebase.NewReference(url).Auth(authToken)

	var storedCommit StoredCommit
	storedCommit.Date = commit.Created_at
	storedCommit.Login = commit.Actor.Login
	storedCommit.Avatar = commit.Actor.Avatar_url
	// This only stores the first commit message. Fix it so it stores one of these
	// for every message
	//storedCommit.Message = commit.Payload.Commits[0]
	storedCommit.Url = commit.Payload.Commits[0].Url

	err := fireBase.Push(storedCommit)
	if err != nil {
		fmt.Println("Firebase error")
		fmt.Println(err)
		return false
	} else {
		fmt.Println("Firebase success")
		return true
	}
}

// There must be a better way to do this. Probably sort cussWords alpha
// and use a lookup table.
func isDirty(message string) bool {

	result := false

	cussWords := []string{
		"fuck",
		"bitch",
		"stupid",
		"tits",
		"asshole",
		"cocksucker",
		"cunt",
		"hell",
		"douche",
		"testicle",
		"twat",
		"bastard",
		"sperm",
		"shit",
		"dildo",
		"wanker",
		"prick",
		"penis",
		"vagina",
		"whore"}

	messageWords := strings.Split(message, " ")

	for _, cussWord := range cussWords {
		for _, word := range messageWords {
			if word == cussWord {
				result = true
			}
		}

	}
	return result
}

func parseCommit(line string) {
	var commit Commit

	jsonErr := json.Unmarshal([]byte(line), &commit)
	if jsonErr != nil {
		fmt.Println("Could not parse json.")
		fmt.Println(jsonErr)
	}

	if commit.Type == "PushEvent" && commit.Payload.Size > 0 {
		// TO-DO this is only checking the first commit. Probably build
		// the storage struct here and pass it to the storeCommit()
		// function
		if isDirty(commit.Payload.Commits[0].Message) {
			storeCommit(commit)
		}
	}

}

func parseFile(fName string) {
	// https://groups.google.com/forum/#!topic/golang-nuts/GjIkryuCyAY
	// TODO: standardize use of file api
	// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file
	fileOS, err := os.Open(fName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open %s: error: %s\n", fName, err)
		os.Exit(1)
	}

	//https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file
	// close fi on exit and check for its returned error
	defer func() {
		if err := fileOS.Close(); err != nil {
			panic(err)
		}
	}()

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
			fmt.Println("Error reading file.")
			fmt.Println(err)
			break
		}

		parseCommit(line)

		i++
	}

	os.Remove(fName)
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
