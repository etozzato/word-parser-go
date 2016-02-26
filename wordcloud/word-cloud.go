package wordcloud

import (
  "C"
  "encoding/json"
  "io"
  "log"
  "regexp"
  "sort"
  "strings"
)

// Response is the input struct
type Response struct {
  ResponseID int
  Text       string
}

// WeightedWord is the base of the structure
type WeightedWord struct {
  Text   string
  Weight int
}

// WordCloud is the collection of WeightedWord
type WordCloud []WeightedWord

// Implementation of the sortable interface
func (wc WordCloud) Len() int           { return len(wc) }
func (wc WordCloud) Swap(i, j int)      { wc[i], wc[j] = wc[j], wc[i] }
func (wc WordCloud) Less(i, j int) bool { return wc[i].Weight > wc[j].Weight }

//ParseWords makes the cloud
func ParseWords(JSONResponses, stopWordsStr string) string {

  stopWords := strings.Split(stopWordsStr, " ")
  stopWordsMap := make(map[string]bool)
  for _, value := range stopWords {
    stopWordsMap[value] = true
  }

  responses := parseJSONResponses(JSONResponses)
  responseText := ""
  for _, response := range responses {
    responseText += response.Text + " "
  }

  re := regexp.MustCompile("([\\.\\!\\?,\\(\\)\\s\\*]+)")
  words := strings.Split(re.ReplaceAllString(strings.ToLower(responseText), "\u2980"), "\u2980")

  weightedCloud := make(map[string]int)
  for _, word := range words {
    if !stopWordsMap[word] {
      weightedCloud[word]++
    }
  }

  sortableWeightedCloud := WordCloud{}
  for key, value := range weightedCloud {
    sortableWeightedCloud = append(sortableWeightedCloud, WeightedWord{Text: key, Weight: value})
  }

  sort.Sort(sortableWeightedCloud)
  max := 50
  if len(sortableWeightedCloud) < 50 {
    max = len(sortableWeightedCloud)
  }

  jsonOutput, err := json.Marshal(sortableWeightedCloud[0:max])
  if err != nil {
    panic(err)
  }

  return string(jsonOutput)
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
