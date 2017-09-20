## dockercore
A tool that generates a Dockerfile for a given .NET Core solution.

### Installation Instructions

For the moment installation is from source. If you have a Go environment set up, you can install dockercore with the following command:

`go get github.com/neofight/dockercore`

### Usage

From the root folder of a .NET Core solution simply run dockercore without any arguments:

`dockercore`

A Dockerfile will be generated, then if you have Docker installed you should be able to build and run the solution quite easily:

`docker build -t {name} .`  
`docker run -d -p 80:80 {name}`

{name} should be substituted with the name of your project.

