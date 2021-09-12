FROM golang:1.16.0-bullseye

RUN go version

ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go build -o driveApi ./cmd/main.go

CMD ["./driveApi"]