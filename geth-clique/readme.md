# Initialize a private Ethereum POA Network
## Pre-requisities
* [Docker](https://www.docker.com/)
* [Go](https://golang.org/)

## Generate genesis configuration
1. create a sealer account
```bash
docker run -it --rm  -v `pwd`/data:/root/.ethereum ethereum/client-go account new
```
2. generate genesis using puppeth
```bash
go install github.com/ethereum/go-ethereum/cmd/puppeth@latest
puppeth
```

## Initialize datadir with genesis
```bash
docker run -it --rm  -v `pwd`/data:/root/.ethereum ethereum/client-go init "/root/.ethereum/genesis.json"
```
## Start POA network(test)
```bash
docker run -it -v `pwd`/data:/root/.ethereum -p 8545:8545 ethereum/client-go \
    --networkid "$(cat `pwd`/data/genesis.json | jq '.config.chainId')" \
    --unlock "0" --password "`pwd`/data/password.txt" \
    --mine \
    --http --http.api debug,eth,net,personal,web3 --http.addr 0.0.0.0 --http.port 8545 --http.corsdomain '*' --http.vhosts '*' \
    --nodiscover
```
