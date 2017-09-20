package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Project struct {
	Path             string
	IsTestProject    bool
	IsStartupProject bool
}

func (p Project) Dir() string {
	return filepath.Dir(p.Path)
}

func (p Project) Name() string {
	return strings.TrimSuffix(filepath.Base(p.Path), ".csproj")
}

type Solution struct {
	Path     string
	Projects []*Project
}

func (s Solution) TestProjects() (testProjects []*Project) {

	for _, p := range s.Projects {
		if p.IsTestProject {
			testProjects = append(testProjects, p)
		}
	}

	return
}

func (s Solution) StartupProjects() (startupProjects []*Project) {

	for _, p := range s.Projects {
		if p.IsStartupProject {
			startupProjects = append(startupProjects, p)
		}
	}

	return
}

type projectElement struct {
	Sdk               string                    `xml:"Sdk,attr"`
	PackageReferences []packageReferenceElement `xml:"ItemGroup>PackageReference"`
}

type packageReferenceElement struct {
	Include string `xml:"Include,attr"`
}

func LoadSolution() (*Solution, error) {

	solutionPath, err := findSolutionFile()

	if err != nil {
		return nil, fmt.Errorf("failed to find solution file: %v", err)
	}

	projectPaths, err := findProjectFiles(solutionPath)

	if err != nil {
		return nil, fmt.Errorf("failed to find project files: %v", err)
	}

	projects := make([]*Project, 0, 1)

	for _, projectPath := range projectPaths {

		project, err := loadProject(projectPath)

		if err != nil {
			return nil, fmt.Errorf("failed to load project file: %v", err)
		}

		projects = append(projects, project)
	}

	return &Solution{solutionPath, projects}, nil
}

func loadProject(projectPath string) (*Project, error) {

	file, err := os.Open(projectPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read project file: %v", err)
	}

	defer file.Close()

	project := Project{Path: projectPath}

	decoder := xml.NewDecoder(file)

	for {
		t, _ := decoder.Token()

		if t == nil {
			break
		}

		switch element := t.(type) {

		case xml.StartElement:

			if element.Name.Local == "Project" {

				var projectElement projectElement

				decoder.DecodeElement(&projectElement, &element)

				if projectElement.Sdk == "Microsoft.NET.Sdk.Web" {
					project.IsStartupProject = true
				}

				for _, referenceElement := range projectElement.PackageReferences {

					if referenceElement.Include == "MSTest.TestFramework" {
						project.IsTestProject = true
					}
				}
			}
		}
	}

	return &project, nil
}

func findSolutionFile() (string, error) {

	files, err := ioutil.ReadDir("./")

	if err != nil {
		return "", fmt.Errorf("failed to read current directory: %v", err)
	}

	solutions := make([]string, 0, 1)

	for _, file := range files {

		if filepath.Ext(file.Name()) == ".sln" {
			solutions = append(solutions, file.Name())
		}
	}

	if len(solutions) == 0 {
		return "", fmt.Errorf("no solution file found")
	}

	if len(solutions) > 1 {
		return "", fmt.Errorf("more than one solution file found")
	}

	return solutions[0], nil
}

func findProjectFiles(solutionPath string) ([]string, error) {

	file, err := os.Open(solutionPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read solution file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	projectPaths := make([]string, 0, 1)

	var projectRegEx = regexp.MustCompile("\"([^\"]*.csproj)\"")

	for scanner.Scan() {
		match := projectRegEx.FindStringSubmatch(scanner.Text())

		if match != nil {
			projectPaths = append(projectPaths, strings.Replace(match[1], "\\", string(os.PathSeparator), -1))
		}
	}

	return projectPaths, nil
}
