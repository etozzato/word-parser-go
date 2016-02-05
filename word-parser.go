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
