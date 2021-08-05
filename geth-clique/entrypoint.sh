#!/bin/sh

if ! [ -d "/root/.ethereum/geth" ]; then
  geth init "/root/.ethereum/genesis.json"
fi

geth --allow-insecure-unlock --networkid=`cat /root/.ethereum/genesis.json | jq '.config.chainId'` --unlock="0" --password="/root/.ethereum/password.txt" \
--nodiscover --mine --http --http.api=debug,eth,net,web3 --http.addr=0.0.0.0 --http.port=8545 --http.corsdomain='*' --http.vhosts='*'

/bin/sh "$@"
