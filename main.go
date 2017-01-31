package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"math/rand"
	"unicode/utf8"
	"unicode"
	"time"
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
)

func check (e error) {
	if e != nil {
		panic(e)
	}
}


type Model struct {
	Seeds []string
	Transitions map[string][]string
	Title string
	Depth int
}


var wordBoundary string = " "
const SEPARATOR string = "|"

func transformSlice(elements []string) (key string) {

	key = strings.Join(elements, SEPARATOR)
	return

}

func parseSentence(depth int, sentence []string, transitions map[string][]string) []string {
	if(depth < 2 || depth > len(sentence)) {
		return []string {}
	}
	var seeds []string
	words := sentence
	finalIndex := len(words) - depth - 1
	for i:=0; i < finalIndex; i++ {
		keys := words[i:(i+depth )]
		value := words[i+depth]
		hashKey := transformSlice(keys)
		if(len(hashKey) > 0) {
			existingTransitions := transitions[hashKey]
			transitions[hashKey] = append(existingTransitions, value)
			if x, _ := utf8.DecodeRuneInString(hashKey); unicode.IsUpper(x) {
				seeds = append(seeds, hashKey)
			}
		}
	}
	return seeds
}

func words(filename string) []string  {
	f, err := ioutil.ReadFile(filename)
	check(err)
	var fileString = string(f)
	var words []string = strings.Fields(fileString)
	return words
}

func replaceSeparator(replaceable string) (replaced string) {
	replaced = strings.Replace(replaceable, SEPARATOR ," ", -1)
	return
}

func generate(transitions map[string][]string, seeds []string, depth int) (generated string) {
	var generator []string
	rand.Seed(time.Now().UnixNano())
	seed := rand.Int() % len(seeds)
	seedWord := seeds[seed]
	possibleWords := transitions[seedWord]
	chosenIndex := rand.Int() % len(possibleWords)
	nextWord := transitions[seedWord][chosenIndex]
	printable := replaceSeparator(nextWord)
	generator = append(generator, replaceSeparator(seedWord))
	generator = append(generator, printable)
	full := strings.Join(generator, " ")
	last := generator[(len(generator) -1)]


	for !strings.HasSuffix(last, ".") && !strings.HasSuffix(last, "?") && !strings.HasSuffix(last, "!")  {


		prefix := strings.Fields(full)
		seedWord = strings.Join(prefix[len(prefix)-depth:], " ")
		hashKey := transformSlice(strings.Split(seedWord, wordBoundary))
		possibleWords = transitions[hashKey]
		chosenIndex = rand.Int() % len(possibleWords)
		nextWord = transitions[hashKey][chosenIndex]
		printable = replaceSeparator(nextWord)
		generator = append(generator, printable)
		full = strings.Join(generator, " ")
		last = generator[(len(generator) -1)]


	}

	generated = strings.Join(generator, wordBoundary)
	return


}

func lookupAndGenerate(title string) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("models")
	result := Model{}
	err = c.Find(bson.M{"title": title}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	generated := generate(result.Transitions, result.Seeds, result.Depth)
	fmt.Println(generated)
}

func parseAndGenerate(title string, filename string) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("models")
	depth := 2
	words := words(filename)
	var transitions = make(map[string][]string)
	seeds  := parseSentence(depth, words, transitions)
	entity := Model{seeds, transitions, title, depth}
	err = c.Insert(&entity)
	if err != nil {
		log.Fatal(err)
	}
	generated := generate(transitions, seeds, depth)
	fmt.Println(generated)
}



func main() {
	//filename := "/home/brandon/dev/golang/src/github.com/bcallender/maas/story.txt"
	//parseAndGenerate("A Tale of Two Cities", filename)
	lookupAndGenerate("A Tale of Two Cities")

}
