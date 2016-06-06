#!/bin/bash

# execute this script with ./install.sh

if [ "$OSTYPE" == "cygwin" ]; then
  echo "WINDOWS x64 installation not supported"
elif [ "$OSTYPE" == "win32" ]; then
  echo "WINDOWS x86 installation not supported"
elif [ "$OSTYPE" == "darwin" ]; then
  echo "Darwin installation not supported"
elif [ "$OSTYPE" == "linux-gnu" ]; then
  echo "Linux installation :"

  if [ -d  "$GOPATH" ]; then
    export PATH=$PATH:$GOPATH/bin > $HOME/.bashrc
    . $HOME/.bashrc

    cd $GOPATH
    go get github.com/stretto-editor/stretto
    cd $GOPATH/src/github.com/stretto-editor/stretto
    go install
    echo "-> Installation completed"
    echo "-> You can now run Stretto with "stretto" command"

  else

    echo "Could not complete installation :"
    echo "-> GOPATH not set"

  fi
fi
