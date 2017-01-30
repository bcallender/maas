package main

import (
	"io/ioutil"
	"strings"
	"hash/fnv"
)

func check (e error) {
	if e != nil {
		panic(e)
	}
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}



var wordBoundary string = " "
const SEPARATOR string = "|"

func transformSlice(elements []string) (key string) {

	key = strings.Join(elements, SEPARATOR)
	return

}

func parseSentence(depth int, sentence string, transitions map[string][]string) {
	if(depth < 2 || depth > len(sentence)) {
		return
	}
	words := strings.Split(sentence, wordBoundary)
	finalIndex := len(words) - depth - 1
	for i:=0; i < finalIndex; i++ {
		keys := words[i:(i+depth )]
		value := words[i+depth]
		hashKey := transformSlice(keys)
		existingTransitions := transitions[hashKey]
		transitions[hashKey] = append(existingTransitions, value)
	}
}

func lines(filename string) []string  {
	f, err := ioutil.ReadFile(filename)
	check(err)
	var fileString = string(f)
	var lines []string = strings.Split(fileString, "\n")
	return lines
}



func main() {
	filename := "/home/brandon/dev/golang/src/github.com/bcallender/maas/story.txt"
	depth := 2
	lines := lines(filename)
	var transitions = make(map[string][]string)
	for _, sentence := range lines {
		parseSentence(depth, sentence, transitions)
	}

}
