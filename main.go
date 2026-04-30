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
	var buffer string
	var bracketCount = 0
	var isInString = false

	for {
		c := make([]byte, 1)
		_, err := os.Stdin.Read(c)
		if err != nil {
			return
		}

		ch := string(c)

		// Skip empty chars
		if strings.TrimSpace(ch) == "" {
			continue
		}

		// Track JSON structure
		if ch == "{" && !isInString {
			bracketCount++
		} else if ch == "}" && !isInString {
			bracketCount--
		} else if ch == `"` && (len(buffer) == 0 || buffer[len(buffer)-1] != '\\') {
			isInString = !isInString
		}

		buffer += ch
		buffer = strings.TrimPrefix(buffer, "\n")

		// If full JSON object read
		if bracketCount == 0 && buffer != "" {

			var jsonMap map[string]interface{}
			err := json.Unmarshal([]byte(buffer), &jsonMap)

			if err != nil {
				// Print raw if not valid JSON
				fmt.Println(color.With(color.Blue, "\n****************\n"+buffer+"\n****************\n"))
				buffer = ""
				continue
			}

			// ✅ Convert ts → formatted time
			if tsVal, ok := jsonMap["ts"]; ok {
				if tsFloat, ok := tsVal.(float64); ok {
					tsMillis := int64(tsFloat)

					// Convert milliseconds → time
					t := time.UnixMilli(tsMillis)

					// Convert to IST
					loc, err := time.LoadLocation("Asia/Kolkata")
					if err == nil {
						t = t.In(loc)
					}

					jsonMap["ts.formatted"] = t.Format("02/01/2006 15:04:05")
				}
			}

			// Pretty print JSON
			b, err := json.MarshalIndent(jsonMap, "", "   ")
			if err != nil {
				fmt.Println(buffer)
				buffer = ""
				continue
			}

			output := string(b)
			output = strings.ReplaceAll(output, `\n`, "\n\t")

			// Color based on log level
			if level, ok := jsonMap["level"].(string); ok {
				switch level {
				case "error":
					fmt.Println(color.With(color.Red, output))
				case "warn":
					fmt.Println(color.With(color.Yellow, output))
				default:
					fmt.Println(output)
				}
			} else {
				fmt.Println(output)
			}

			buffer = ""
		}
	}
}
