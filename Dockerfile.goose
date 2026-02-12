# Dockerfile.goose
FROM golang:1.25

# Set the working directory
WORKDIR /app

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the application code into the container
COPY . .

# Set the entrypoint for the container
ENTRYPOINT ["goose"]
