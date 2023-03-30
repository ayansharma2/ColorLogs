package main

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"os"
	"strings"
	"time"
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
		if strings.Trim(string(c), "") == "" {
			continue
		}
		var jsonString map[string]interface{}
		err = json.Unmarshal([]byte(first), &jsonString)
		if err != nil && bracketCount == 0 && string(c) == "{" && !isInString && strings.Trim(first, "") != "" {
			fmt.Println(color.With(color.Blue, "\n****************\n"+first+"\n****************\n"))
			first = ""
		}
		if string(c) == "{" && !isInString {
			bracketCount++
		} else if string(c) == "}" && !isInString {
			bracketCount--
		} else if string(c) == "\"" && len(first)-1 >= 0 && first[len(first)-1] != '\\' {
			isInString = !isInString
		}
		first += string(c)
		first = strings.TrimPrefix(first, "\n")

		if bracketCount == 0 && first != "" {
			currentTime := time.Now()
			var jsonString map[string]interface{}
			err := json.Unmarshal([]byte(first), &jsonString)
			if jsonString != nil {
				jsonString["ts.formatted"] = currentTime.Format(`2/01/2006 15:04:05`)
			}
			b, err1 := json.MarshalIndent(jsonString, "", "   ")
			finalVal := strings.Replace(string(b), `\n`, "\n\t", -1)
			finalVal = strings.Replace(finalVal, `\t`, "\t", -1)

			if string(b) != "null" {
				if strings.HasPrefix(first, "{\"level\":\"error\"") {
					fmt.Println(color.With(color.Red, finalVal))
				} else if strings.HasPrefix(first, "{\"level\":\"warn\"") {
					fmt.Println(color.With(color.Yellow, finalVal))
				} else {
					fmt.Println(string(b))
				}
			}
			if err == nil && err1 == nil {
				first = ""
			}
		}
	}
}
