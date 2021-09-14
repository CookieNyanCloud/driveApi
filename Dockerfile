FROM golang:1.17.0-bullseye

RUN go version

ENV GOPATH=/

COPY ./ ./

RUN go mod download
EXPOSE 8090

RUN go build -o drive-api ./main.go

CMD ["./drive-api"]