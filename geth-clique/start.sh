#!/bin/bash
exec geth \
	--networkid "$(cat /geth/genesis.json | jq '.config.chainId')" \
	--datadir "/geth/chain" \
	--keystore "/geth/keys" \
	--password "/geth/keys/dev-key-password.txt" \
	--verbosity ${GETH_VERBOSITY:-2} --mine \
	--http --http.api admin,debug,eth,miner,net,personal,txpool,web3 --http.addr 0.0.0.0 --http.port 8545 --http.corsdomain '*' --http.vhosts '*' \
	--nodiscover
	--targetgaslimit 6500000
