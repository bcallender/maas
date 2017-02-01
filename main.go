package main

import "gopkg.in/mgo.v2/bson"

func main() {
	filename := "src/github.com/bcallender/maas/story.txt"
	c := make(chan bson.ObjectId)
	for i:= 0; i<25; i++ {
		go ParseAndGenerate(string(i), filename, c)
	}
	for i:= 0; i<25; i++ {
		v :=  <-c
		val := v.String()
		println(val)
	}





}
