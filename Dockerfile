FROM golang

WORKDIR /usr/app

COPY . .

CMD ["go", "run", "src/main.go"]