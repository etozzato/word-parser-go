package wordtree

import (
	"C"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strings"
)

// Response is what you expect
type Response struct {
	ResponseID string
	Sentences  [][]string
}

//PostfixSets is what you expect
func PostfixSets(text, filter string) string {
	responses := matchingSentences(text, filter, true)
	jsonOutput, err := json.Marshal(responses)
	if err != nil {
		panic(err)
	}
	return string(jsonOutput)
}

//ResponseIds is what you expect
func ResponseIds(text, filter string) string {
	responses := matchingSentences(text, filter, false)
	output := make([]string, 0, len(responses))
	for _, value := range responses {
		output = append(output, value.ResponseID)
	}
	jsonOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return string(jsonOutput)
}

func matchingSentences(jsonString string, filter string, postfixRegExp bool) []Response {
	output := []Response{}
	dec := json.NewDecoder(strings.NewReader(jsonString))
	var r Response
	for {
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("DecodeJSON error:", err)
		}
		if len(filter) > 0 {
			var f Response
			for _, sentence := range r.Sentences {
				filtered := matchSentences(sentence, filter, postfixRegExp)
				if len(filtered) > 0 {
					f.Sentences = append(f.Sentences, filtered)
				}
			}
			if len(f.Sentences) > 0 {
				f.ResponseID = r.ResponseID
				output = append(output, f)
			}
		} else {
			output = append(output, r)
		}
	}
	return output
}

func matchSentences(slice []string, filter string, postfixRegExp bool) []string {
	var filtered []string
	re := regexp.MustCompile("(?i)\\b" + filter + "\\b")
	for _, str := range slice {
		subStrings := scanStr(str, "([^\\.\\!\\?]+[\\.\\!\\?]?)\\s*", false)
		for _, value := range subStrings {
			if re.MatchString(value) {
				if postfixRegExp {
					splits := re.Split(value, -1)
					postPhrase := strings.Join(splits[1:len(splits)], filter)
					sentences := scanStr(postPhrase, "([\\s\\.,;:-])", true)
					for _, sentence := range sentences {
						filtered = append(filtered, " "+sentence)
					}
				} else {
					filtered = append(filtered, value)
				}
			}
		}
	}
	return filtered
}

func scanStr(str string, pat string, isolate bool) []string {
	re := regexp.MustCompile(pat)
	if isolate {
		whitespaceRe := regexp.MustCompile("\u2980\\s+\u2980")
		tmpStr := re.ReplaceAllString(str, "\u2980$1\u2980")
		tmpStr = whitespaceRe.ReplaceAllString(tmpStr, "\u2980")
		words := strings.Split(tmpStr, "\u2980")
		output := words[:0]
		for _, value := range words {
			if value != "" && value != " " {
				output = append(output, value)
			}
		}
		return output
	}
	return strings.Split(re.ReplaceAllString(str, "$1\u2980"), "\u2980")
}
