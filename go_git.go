package main

import "fmt"

//import "reflect"
import "net/http"
import "io/ioutil"

import "os"

func main() {
	resp, err := http.Get("https://api.github.com/users/cdale77/repos")

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

		fmt.Println("%s\n", string(contents))

	}
}
