# GOLANG Word Parser (CGO/Rails/RUBY/FFI)

```ruby
# see ./ruby-test/cloud.rb
# the Cloud ruby class delegates the method parseWords
# to wordcloud.go though the "main" file word-parser.go

class Cloud < Base
  attr_accessor :stop_words, :word_cloud, :heaviest_word, :responses, :comments

  def initialize(responses, comments, stop_words)
    self.responses  = responses
    self.comments   = comments
    self.stop_words = stop_words.join("*")

    parse_words
  end

  private

  def parse_words
    self.word_cloud = JSON.parse(golangParseWords(sentences.join("*"), stop_words)) rescue []
    self.heaviest_word = word_cloud.first['Text'] if word_cloud.any?
  end

  def sentences
    (responses + comments).map(&:text)
  end

end

```

```go
// see ./word-parser.go

package main

import (
	"C"

	wordcloud "./wordcloud"
	wordtree "./wordtree"
)

// word cloud //

//export golangParseWords
func golangParseWords(pResponseText, pStopWords *C.char) string {
	return wordcloud.ParseWords(C.GoString(pResponseText), C.GoString(pStopWords))
}

// word tree //

//export golangPostfixSets
func golangPostfixSets(pText, pFilter *C.char) string {
	return wordtree.PostfixSets(C.GoString(pText), C.GoString(pFilter))
}

func main() {}
```

```go
// see ./wordcloud/wordcloud.go
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
```

```go
// see ./wordcloud/wordtree.go
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

```

### Test the cloud

```
~/go/word-parser  ruby ./ruby-test/cloud.rb
heaviest_word: south
word_cloud: [{"Text"=>"south", "Weight"=>28}, {"Text"=>"africa", "Weight"=>22}, {"Text"=>"people", "Weight"=>20}, {"Text"=>"country", "Weight"=>20}, {"Text"=>"government", "Weight"=>12}, {"Text"=>"rights", "Weight"=>12}, {"Text"=>"change", "Weight"=>10}, {"Text"=>"majority", "Weight"=>10}, {"Text"=>"african", "Weight"=>10}, {"Text"=>"africans", "Weight"=>10}, {"Text"=>"order", "Weight"=>10}, {"Text"=>"political", "Weight"=>10}, {"Text"=>"democracy", "Weight"=>8}, {"Text"=>"vision", "Weight"=>8}, {"Text"=>"freedom", "Weight"=>8}, {"Text"=>"hope", "Weight"=>8}, {"Text"=>"citizens", "Weight"=>6}, {"Text"=>"national", "Weight"=>6}, {"Text"=>"island", "Weight"=>6}, {"Text"=>"spirit", "Weight"=>6}, {"Text"=>"elections", "Weight"=>6}, {"Text"=>"adopted", "Weight"=>6}, {"Text"=>"nation", "Weight"=>6}, {"Text"=>"constitution", "Weight"=>6}, {"Text"=>"democratic", "Weight"=>6}, {"Text"=>"fought", "Weight"=>6}, {"Text"=>"cape", "Weight"=>6}, {"Text"=>"today", "Weight"=>6}, {"Text"=>"create", "Weight"=>4}, {"Text"=>"reconstruction", "Weight"=>4}, {"Text"=>"better", "Weight"=>4}, {"Text"=>"bill", "Weight"=>4}, {"Text"=>"work", "Weight"=>4}, {"Text"=>"task", "Weight"=>4}, {"Text"=>"robben", "Weight"=>4}, {"Text"=>"congress", "Weight"=>4}, {"Text"=>"coloured", "Weight"=>4}, {"Text"=>"assist", "Weight"=>4}, {"Text"=>"peninsula", "Weight"=>4}, {"Text"=>"based", "Weight"=>4}, {"Text"=>"religious", "Weight"=>4}, {"Text"=>"requires", "Weight"=>4}, {"Text"=>"centuries", "Weight"=>4}, {"Text"=>"unity", "Weight"=>4}, {"Text"=>"bring", "Weight"=>4}, {"Text"=>"fighters", "Weight"=>4}, {"Text"=>"united", "Weight"=>4}, {"Text"=>"society", "Weight"=>4}, {"Text"=>"speak", "Weight"=>4}, {"Text"=>"resistance", "Weight"=>4}]
```
### Test the tree

```
~/go/word-parser  ruby ./ruby-test/tree.rb
response_ids: ["2021707","2021708","2021709","2021710","2021712","2021713"]
postfix_sets: [{"ResponseID":"2021707","Sentences":[[" when"," people"," make"," trees"]]},{"ResponseID":"2021708","Sentences":[[" eating"," pizza"]]},{"ResponseID":"2021709","Sentences":[[" me"," ,"," tender"," ."]]},{"ResponseID":"2021710","Sentences":[[" you"],[" to"," eat"," pizza"," for"," breakfast"],[" is"," all"," about"]]},{"ResponseID":"2021712","Sentences":[[" you"," to"," help"," me"," more"," !"]]},{"ResponseID":"2021713","Sentences":[[" to"," eat"," with"," friends"," and"," family"]]}]
```
