FROM golang:1.24.1

# Set up work directory
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the app
RUN go build -o server .

# Expose the port
EXPOSE 8070

# Run the server
CMD ["./server"]
