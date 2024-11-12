# syntax=docker/dockerfile:1

FROM golang:1.23.2

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Ejecutar las pruebas antes de construir la aplicación
RUN go test -v ./...  # Ejecutar las pruebas de Go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /super-descuentos

EXPOSE 8080

# Run
CMD ["/super-descuentos"]






#   genera archivo go.mod ----> go mod tidy
## COMANDOS DOCKER
#   build --------> docker build --tag super-descuentos .
#   list images --> docker image ls
#   run ----------> docker run super-descuentos || docker run --publish 8080:8080 super-descuentos