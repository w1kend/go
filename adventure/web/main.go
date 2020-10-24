package main

import (
	"flag"
	"fmt"
	"github.com/w1kend/go/adventure"
	"log"
	"net/http"
	"os"
)

func main() {
	filename := flag.String("file", "story.json", "json file with the story")
	flag.Parse()

	reader, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	story, err := adventure.JsonStory(reader)
	if err != nil {
		panic(err)
	}

	handler := adventure.NewHandler(story)
	fmt.Println("Starting the server")

	log.Fatal(http.ListenAndServe(":3000", handler))
}
