package wordcloud

import (
	"C"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

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
func ParseWords(responseText, stopWordsStr string) string {

	stopWords := strings.Split(stopWordsStr, "*")
	stopWordsMap := make(map[string]bool)
	for _, value := range stopWords {
		stopWordsMap[value] = true
	}

	re := regexp.MustCompile("([\\.\\!\\?,\\(\\)\\s\\*]+)")
	words := strings.Split(re.ReplaceAllString(strings.ToLower(responseText), "*"), "*")

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
	jsonOutput, err := json.Marshal(sortableWeightedCloud[0:50])
	if err != nil {
		panic(err)
	}

	return string(jsonOutput)
}
