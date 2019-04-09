Getting Started
========================
Set up GOPATH 
Inside of go_weather directory run the go build command.
Once the binary is built use the ./go_weather command to execute.

Follow instructions to download Microsoft SQL Server Docker Image of choice from https://hub.docker.com/_/microsoft-mssql-server


Testing
=============
To run Unit tests, install ginkgo and gomega on your machine by running 'go get https://github.com/onsi/ginkgo' and 'go get https://github.com/onsi/gomega' respectively.
Once installed run the gingko -r command inside of the go_weather root directory.