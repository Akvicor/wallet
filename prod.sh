#!/bin/bash

if [ -f "/data/config.toml" ]; then
  echo "config.toml file found"
else
  ./wallet example -p "/data/" -c > /data/config.toml
fi

if [ -f "/data/wallet.db" ]; then
  echo "wallet.db found"
  ./wallet migrate -c /data/config.toml
fi

./wallet server -c /data/config.toml

