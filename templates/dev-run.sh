# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# Install air for live-reloading
go install github.com/cosmtrek/air@latest
# Install frontend dependencies
yarn install
go mod tidy
# start air
air
