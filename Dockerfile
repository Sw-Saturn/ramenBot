FROM golang:1.13.3
ENV GO111MODULE=on
CMD /bin/bash
WORKDIR /go/src/github.com/Sw-Saturn/ramenBot
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080

COPY . .

CMD ["go", "run", "main.go"]
