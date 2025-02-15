# Use Go 1.23 bookworm as base image
FROM golang:1.23-bookworm

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o re_partners_software_coding_challenge_go ./cmd/api

# Start the application
CMD ["/build/re_partners_software_coding_challenge_go"]
