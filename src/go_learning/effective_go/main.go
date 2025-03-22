package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SampleJson struct {
	Text string `json:"text"`
}

func DeferExample() {
	//File open
	f, err := os.Open("sample.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var s *SampleJson
	json.NewDecoder(f).Decode(&s)
	fmt.Println(s.Text)
}

func main() {
	DeferExample()
}
