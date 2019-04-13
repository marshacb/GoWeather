FROM golang

COPY . /go/src/go_weather

WORKDIR /go/src/go_weather

CMD ["go", "run", "main.go"]