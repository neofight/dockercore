package main

import (
	"fmt"
	"log"
)

func main() {

	solution, err := LoadSolution()

	if err != nil {
		log.Fatalf("failed to load solution: %v", err)
	}

	err = WriteDockerfile(solution)

	if err != nil {
		log.Fatalf("failed to write DockerFile: %v", err)
	}

	fmt.Println("Dockerfile successfully generated")
}
