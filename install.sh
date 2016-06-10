#!/bin/bash

# execute this script with ./install.sh

if [ "$OSTYPE" == "cygwin" ]; then
  echo "WINDOWS x64 installation not supported"
elif [ "$OSTYPE" == "win32" ]; then
  echo "WINDOWS x86 installation not supported"
elif [ "$OSTYPE" == "darwin" ]; then
  echo "Darwin installation not supported"
elif [ "$OSTYPE" == "linux-gnu" ]; then
  echo "> Linux installation :"

  mkdir /tmp/strettotemp
  cd /tmp/strettotemp
  mkdir usr

  wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
  tar -C ./usr -xzf go1.6.2.linux-amd64.tar.gz

  export GOROOT=/tmp/strettotemp/usr
  export PATH=$PATH:$GOROOT/go/bin
  export GOPATH=/tmp/strettotemp/strettoinstall

  mkdir strettoinstall

  echo "--> Application will now be installed"
  cd $GOPATH
  go get github.com/stretto-editor/stretto
  cd $GOPATH/src/github.com/stretto-editor/stretto
  go install
  mv stretto.json $HOME/.stretto.json
  mv Commands.md $HOME
  mv $GOPATH/bin/stretto $HOME
  cd $HOME
  rm -rf /tmp/strettotemp
  echo "----> Installation of Stretto completed"
  echo "-----> Stretto application is now in $HOME"
  echo "-----> You can now move stretto and Commands.md file to your /usr/bin"


fi
