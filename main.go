package main

import "io.github.starfreck/valkey-go-pub-sub/advance"

func main() {

	// Advance Example Usage: http://localhost:8080/publish?message="test"
	advance.Call("127.0.0.1:6379")

	// simple.Call("127.0.0.1:6379", "ch1", "Hey I was published :)")
}
