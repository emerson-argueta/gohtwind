# # install latest golang if not already installed on local windows machine
choco install golang

# # Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# # Install air for live-reloading
go install github.com/cosmtrek/air@latest
# # Install frontend dependencies

# # start air
air
