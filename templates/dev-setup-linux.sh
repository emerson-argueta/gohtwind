# install go if not already installed on local linux machine
# check if go is installed
if ! command -v go &> /dev/null
then
  echo "go could not be found"
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
fi

 #install nodejs if not already installed on local linux machine
if ! command -v node &> /dev/null
then
  echo "node could not be found"
  curl -sL https://deb.nodesource.com/setup_14.x | sudo -E bash - && sudo apt-get install -y nodejs
fi

# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
# Install air for live-reloading
go install github.com/cosmtrek/air@latest
# Install frontend dependencies
yarn install

# start air
air
