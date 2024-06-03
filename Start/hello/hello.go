package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"example.com/greetings"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"Gladys", "Samantha", "Darrin"}

	messages, err := greetings.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
