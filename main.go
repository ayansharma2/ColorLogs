package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TwiN/go-color"
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
			var jsonString map[string]interface{}
			err := json.Unmarshal([]byte(first), &jsonString)

			if err == nil && jsonString != nil {
				if tsVal, ok := jsonString["ts"]; ok {
					// JSON numbers come as float64
					tsFloat, ok := tsVal.(float64)
					if ok {
						tsMillis := int64(tsFloat)

						// Convert milliseconds → time.Time
						t := time.UnixMilli(tsMillis)

						// Convert to IST (optional)
						loc, _ := time.LoadLocation("Asia/Kolkata")
						t = t.In(loc)

						jsonString["ts.formatted"] = t.Format("02/01/2006 15:04:05")
					}
				}
			}
			b, err1 := json.MarshalIndent(jsonString, "", "   ")
			finalVal := strings.Replace(string(b), `\n`, "\n\t", -1)
			finalVal = strings.Replace(finalVal, `\t`, "\t", -1)

			if string(b) != "null" {
				level, _ := jsonString["level"].(string)
				switch strings.ToLower(level) {
				case "error", "err", "fatal", "critical":
					fmt.Println(color.With(color.Red, finalVal))
				case "warn", "warning":
					fmt.Println(color.With(color.Yellow, finalVal))
				default:
					fmt.Println(finalVal)
				}
			}
			if err == nil && err1 == nil {
				first = ""
			}
		}
	}
}
