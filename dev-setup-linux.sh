# install go if not already installed on local linux machine
# check if go is installed
if ! command -v go &> /dev/null
then
  echo "go could not be found"
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
fi

# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# Install air for live-reloading
go install github.com/cosmtrek/air@latest
# Install frontend dependencies

# start air
air
