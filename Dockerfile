FROM golang:latest

# Install protoc-gen-go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
ENV PATH="${PATH}:/go/bin"

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the code
COPY . .

# Install Air
RUN go install github.com/air-verse/air@latest

# Use Air as entrypoint (expects a .air.toml in project root)
ENTRYPOINT ["air"]
