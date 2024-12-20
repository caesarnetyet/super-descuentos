# syntax=docker/dockerfile:1
FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v  ./relational/... ./rest/...

RUN CGO_ENABLED=0 GOOS=linux go build -o /super-descuentos

EXPOSE 8080

CMD ["/super-descuentos"]


#   genera archivo go.mod ----> go mod tidy
## COMANDOS DOCKER
#   tests --------> docker build -f Dockerfile.test -t super-descuentos-test .
#   build --------> docker build -t super-descuentos .
#   list images --> docker image ls
#   run ----------> docker run -p 8080:8080 super-descuentos


# docker build -t super-descuentos .
# docker images
# docker run -p 8080:8080 super-descuentos