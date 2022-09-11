#!/bin/bash

if [[ -n ${GOPATH} ]]; then
  echo "GOPATH env var found!"
  go install ./
  echo "${GOPATH}/bin/net-conn-proxy"
else
  echo "GOPATH env var not found, alternative installation..."
  go build ./

  if [ "$(uname)" == "Darwin" ]; then
    # Do something under Mac OS X platform
    echo "Mac OX X detected..."
  elif [ "$(expr substr "$(uname -s)" 1 5)" == "Linux" ]; then
    # Do something under GNU/Linux platform
    echo "GNU/Linux detected..."
    chmod +x ./net-conn-proxy
    sudo cp ./net-conn-proxy /usr/local/bin/net-conn-proxy
    echo "/usr/local/bin/net-conn-proxy"
  elif [ "$(expr substr "$(uname -s)" 1 10)" == "MINGW32_NT" ] || [ "$(expr substr "$(uname -s)" 1 10)" == "MINGW64_NT" ]; then
    # Do something under 32/64 bits Windows NT platform
    echo "Windows NT detected..."
    cp ./net-conn-proxy.exe /C/Windows/System32/net-conn-proxy.exe
    echo "/C/Windows/System32/net-conn-proxy.exe"
  fi
fi
