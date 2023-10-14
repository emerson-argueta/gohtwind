# install latest golang if not already installed on local macos machine
brew install go
# install latest node if not already installed on local macos machine
brew install node

# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# Install realize for live-reloading
go install github.com/cosmtrek/air@latest
# Install frontend dependencies
yarn install

# start air
air

