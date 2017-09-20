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

	fmt.Printf("Solution: %v\n", solution.Path)

	for _, project := range solution.Projects {

		if project.IsTestProject {
			fmt.Printf("%v - Test Project\n", project.Path)
		}

		if project.IsStartupProject {
			fmt.Printf("%v - Startup Project\n", project.Path)
		}
	}
}
