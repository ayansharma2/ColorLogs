package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"os"
	"strings"
)

func main() {
	var first string
	var bracketCount = 0
	var isInString = false

	for {
		var c = make([]byte, 1)
		_, err := os.Stdin.Read(c)
		if err != nil {
			return
		}
		var jsonString map[string]interface{}
		json.Unmarshal([]byte(first), &jsonString)
		_, err = json.MarshalIndent(jsonString, "", "   ")
		if err != nil && string(c) == "{" && !isInString {
			fmt.Println("Error Found :", err)
			fmt.Println(color.With(color.Blue, first))
			first = ""
		} else if string(c) == "{" && !isInString {
			bracketCount++
		} else if string(c) == "}" && !isInString {
			bracketCount--
		} else if string(c) == "\"" && first[len(first)-1] != '\\' {
			isInString = !isInString
		}
		first += string(c)
		first = strings.TrimPrefix(first, "\n")

		if bracketCount == 0 && first != "" {
			if strings.HasPrefix(first, "{\"level\":\"error\"") {
				var jsonString map[string]interface{}
				json.Unmarshal([]byte(first), &jsonString)
				b, _ := json.MarshalIndent(jsonString, "", "   ")
				finalVal := strings.Replace(string(b), `\n`, "\n\t", -1)
				finalVal = strings.Replace(finalVal, `\t`, "\t", -1)
				fmt.Println(color.With(color.Red, finalVal))
			} else if strings.HasPrefix(first, "{\"level\":\"warn\"") {
				var jsonString map[string]interface{}
				json.Unmarshal([]byte(first), &jsonString)
				b, _ := json.MarshalIndent(jsonString, "", "   ")
				finalVal := strings.Replace(string(b), `\n`, "\n\t", -1)
				finalVal = strings.Replace(finalVal, `\t`, "\t", -1)
				fmt.Println(color.With(color.Yellow, finalVal))
			} else {
				var jsonString map[string]interface{}
				json.Unmarshal([]byte(first), &jsonString)
				b, _ := json.MarshalIndent(jsonString, "", "   ")
				fmt.Println(string(b))
			}
			first = ""
		}
	}
}
