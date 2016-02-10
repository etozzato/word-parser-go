package wordtree

import (
	"C"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strings"
)

// Response is the main struct
type Response struct {
	ResponseID string
	Sentences  [][]string
}

//PostfixSets is the exported routine
// Entry point of the algorithm: it receives an array of Responses in JSON format {"ResponseID":"1","Sentences":[['sentence1', 'sentence2']]} and the heaviest word of the set, which is the root of the tree.
func PostfixSets(JSONResponses, root string) string {
	responses := parseJSONResponses(JSONResponses)
	postfixSets := postfixSetsFromResponses(responses, root)
	JSONPostFixSets, err := json.Marshal(postfixSets)
	if err != nil {
		panic(err)
	}
	return string(JSONPostFixSets)
}

// parses JSONResponses into a []Response
func parseJSONResponses(JSONResponses string) []Response {
	output := []Response{}
	decoder := json.NewDecoder(strings.NewReader(JSONResponses))
	for {
		var r Response
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("DecodeJSON error:", err)
		}
		output = append(output, r)
	}
	return output
}

// removes response.Sentences that don't match root, removes root from the matched sentences
func postfixSetsFromResponses(responses []Response, root string) []Response {
	output := []Response{}
	filterRe := regexp.MustCompile("(?i)\\b" + root + "\\b([^.!?]*[.!?])")
	tokensRe := regexp.MustCompile("([\\s.,;:-])")
	for _, response := range responses {
		var postPhrases [][]string
		for _, sentence := range response.Sentences[0] {
			filterReStringSubmatch := filterRe.FindAllStringSubmatch(sentence, -1)
			if len(filterReStringSubmatch) > 0 {
				// FindAllStringSubmatch returns [[:match, :capure]], we need to keep only :capure - so we delete the first element of the slice.
				filterReStringSubmatch[0] = append(filterReStringSubmatch[0][:0], filterReStringSubmatch[0][1:]...)
				filterReStringSubmatch[0] = parsePostPhrase(filterReStringSubmatch[0][0], tokensRe)
				postPhrases = append(postPhrases, filterReStringSubmatch...)
			}
		}
		if len(postPhrases) > 0 {
			output = append(
				output,
				Response{ResponseID: response.ResponseID, Sentences: postPhrases})
		}
	}
	return output
}

// splits the postPhrase into individual words
func parsePostPhrase(sentence string, re *regexp.Regexp) []string {
	tokens := strings.Split(re.ReplaceAllString(sentence, "\u2980$1\u2980"), "\u2980")
	output := []string{}
	for _, token := range tokens {
		// unless the token is empty
		if strings.TrimSpace(token) != "" {
			// unless the token is :punctuation, we prepend a space
			if !re.MatchString(token) {
				token = " " + token
			}
			output = append(output, token)
		}
	}
	return output
}
