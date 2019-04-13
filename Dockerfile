FROM golang:latest

WORKDIR /src

COPY . .

RUN go get github.com/denisenkom/go-mssqldb
RUN go get github.com/go-chi/chi
RUN go get github.com/go-chi/chi/middleware
RUN go get github.com/rs/cors

RUN go build main.go

CMD ["./main"]