package main

import (
	"fmt"
	"os"
	"text/template"
)

const dockerTemplate = `# Build Image
FROM microsoft/aspnetcore-build:2.0 AS build
WORKDIR /build

# Cache package dependencies
COPY {{.Path}} ./
{{range .Projects -}}
COPY {{.Path}} {{.Dir}}/
{{end -}}
RUN dotnet restore

# Copy everything else, run the tests and publish the build
COPY ./ ./
{{range .TestProjects -}}
RUN dotnet test {{.Path}}
{{end -}}
RUN dotnet publish {{(index .StartupProjects 0).Path}} -c Release -o published

# Runtime Image
FROM microsoft/aspnetcore:2.0
WORKDIR /api

COPY --from=build /build/{{(index .StartupProjects 0).Dir}}/published ./
ENTRYPOINT ["dotnet", "{{(index .StartupProjects 0).Name}}.dll"]
`

func WriteDockerfile(solution *Solution) error {

	if len(solution.StartupProjects()) > 1 {
		return fmt.Errorf("multiple startup projects not supported")
	}

	template := template.New("Dockerfile")
	template, err := template.Parse(dockerTemplate)

	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile template %v", err)
	}

	file, err := os.Create("Dockerfile")

	if err != nil {
		return fmt.Errorf("unable to create or open Dockerfile: %v", err)
	}

	defer file.Close()

	template.Execute(file, solution)

	return nil
}
