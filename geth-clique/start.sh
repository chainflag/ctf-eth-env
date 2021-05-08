#!/bin/bash
exec geth \
	--nousb \
	--networkid "$(cat /geth/genesis.json | jq '.config.chainId')" \
	--datadir "/geth/chain" \
	--verbosity ${GETH_VERBOSITY:-2} --mine \
	--rpc --rpcapi admin,db,debug,eth,miner,net,personal,shh,txpool,web3 --rpcaddr 0.0.0.0 --rpcport 8545 --rpccorsdomain '*' --rpcvhosts '*' \
	--nodiscover \
	--targetgaslimit 6500000
