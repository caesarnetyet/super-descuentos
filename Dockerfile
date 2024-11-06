# syntax=docker/dockerfile:1

FROM golang:1.23.2

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /super-descuentos

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD [ "/super-descuentos" ]


#   genera archivo go.mod ----> go mod tidy
## COMANDOS DOCKER
#   build --------> docker build --tag super-descuentos .
#   list images --> docker image ls
#   run ----------> docker run super-descuentos || docker run --publish 8080:8080 super-descuentos