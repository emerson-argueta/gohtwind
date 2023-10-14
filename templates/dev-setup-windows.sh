# # install latest golang if not already installed on local windows machine
choco install golang
# # install latest node if not already installed on local windows machine
choco install nodejs

# # Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# # Install air for live-reloading
go install github.com/cosmtrek/air@latest
# # Install frontend dependencies
yarn install

# # start air
air
