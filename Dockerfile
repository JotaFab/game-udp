
FROM golang:latest



RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
ENV PATH="${PATH}:/go/bin"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
copy . .


RUN go install github.com/air-verse/air@latest


ENTRYPOINT ["air"]
