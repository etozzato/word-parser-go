package wordtree

import (
  "C"
  "encoding/json"
  "io"
  "log"
  "regexp"
  "strings"
)
import "sort"

// Response is the input struct
type Response struct {
  ResponseID int
  Text       string
}

// PostfixSet is output struct
type PostfixSet struct {
  ResponseID int
  Sentences  [][]string
}

//PostfixSets is the exported routine
// Entry point of the algorithm: it receives an array of Responses in JSON format "[{\"ResponseID\":645233,\"Text\":\"Sailing with my dad on a windless day on the willamette river. My happy place.\"}]" and the heaviest word of the set, which is the root of the tree.
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
  err := decoder.Decode(&output)
  if err != nil && err != io.EOF {
    log.Fatal("DecodeJSON error:", err)
  }
  return output
}

// removes response.Sentences that don't match root, removes root from the matched sentences []Response
func postfixSetsFromResponses(responses []Response, root string) []PostfixSet {
  responsesMap := make(map[int][][]string)
  output := []PostfixSet{}
  responseIDs := []int{}

  filterRe := regexp.MustCompile("(?i)\\b" + root + "\\b([^.!?]*[.!?])")
  tokensRe := regexp.MustCompile("([\\s.,;:\\-])")
  for _, response := range responses {
    filterReStringSubmatch := filterRe.FindAllStringSubmatch(response.Text, -1)
    if len(filterReStringSubmatch) > 0 {
      for _, stringSubmatch := range filterReStringSubmatch {
        postPhrases := parsePostPhrase(stringSubmatch[1], tokensRe)
        if len(postPhrases) > 0 {
          responsesMap[response.ResponseID] = append(responsesMap[response.ResponseID], postPhrases)
        }
      }
    }
  }

  for responseID, tokens := range responsesMap {
    if len(tokens) > 0 {
      responseIDs = append(responseIDs, responseID)
    }
  }

  sort.Ints(responseIDs)

  for _, responseID := range responseIDs {
    output = append(output, PostfixSet{ResponseID: responseID, Sentences: responsesMap[responseID]})
  }

  return output
}

// // splits the postPhrase into individual words
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
