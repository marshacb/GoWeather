FROM golang

WORKDIR /usr/app

COPY . .

CMD ["go", "run", "main.go"]