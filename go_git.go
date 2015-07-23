package main

import (
	//"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	//"strconv"
)

func main() {
	getData()
}

func getData() {
	url := "http://data.githubarchive.org/2015-01-01-12.json.gz"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("%s", err)
		os.Exit(1)
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("%s", err)
			os.Exit(1)
		}

		ioutil.WriteFile("dat1.gz", contents, 0644)

		fmt.Println("completed")
		fmt.Println(reflect.TypeOf(contents))
	}
}
