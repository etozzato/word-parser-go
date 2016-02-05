package main

import (
	"C"

	wordcloud "./wordcloud"
	wordtree "./wordtree"
)

////////////////
// word cloud //
////////////////

//export parseWords
func parseWords(pResponseText, pStopWords *C.char) string {
	return wordcloud.ParseWords(C.GoString(pResponseText), C.GoString(pStopWords))
}

///////////////
// word tree //
///////////////

//export postfixSets
func postfixSets(pText, pFilter *C.char) string {
	return wordtree.PostfixSets(C.GoString(pText), C.GoString(pFilter))
}

//export responseIds
func responseIds(pText, pFilter *C.char) string {
	return wordtree.ResponseIds(C.GoString(pText), C.GoString(pFilter))
}

func main() {}
