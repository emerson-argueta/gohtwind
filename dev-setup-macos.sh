# install latest golang if not already installed on local macos machine
brew install go

# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# Install realize for live-reloading
go install github.com/cosmtrek/air@latest

# Copy Go mod and sum files (cache dependencies)
go mod tidy
go build

