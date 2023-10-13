# --- Go Development Stage ---
FROM golang:1.17 AS go-dev

# Set up the working directory
WORKDIR /app

# Install debugging tools
RUN go get github.com/go-delve/delve/cmd/dlv

# Install realize for live-reloading
RUN go get github.com/oxequa/realize

# Copy Go mod and sum files (cache dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Expose the default port for the application and the port for dlv
EXPOSE 8080 40000

# Use `realize` to start your development server
CMD ["realize", "start"]
