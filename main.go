package main

import (
	"github.com/speckleworks/gospeckle/cmd"
)

// PrintJson is a utility function to pretty print a struct as a
// JSON
// func PrintJson(body interface{}) {
// 	var prettyJSON bytes.Buffer
// 	// var jsonData []byte
// 	jsonData, _ := json.Marshal(body)

// 	error := json.Indent(&prettyJSON, jsonData, "", "\t")
// 	if error != nil {
// 		log.Println("JSON parse error: ", error)
// 		return
// 	}

// 	log.Println(string(prettyJSON.Bytes()))
// }

func main() {
	cmd.Execute()
}
