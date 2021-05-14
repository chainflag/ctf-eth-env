# ctf-eth-env

As the data in the permissionless blockchain is public, dishonest CTF players of blockchain(smart contract) challenge can plagiarize someone else's solutions by querying the block information, which causes unfairness.  

So the idea of this project is to disable several RPC methods (e.g. `eth_getBlockByHash`, `eth_getBlockByNumber`) of an Ethereum POA node and make it as challenging server environment so that players cannot know the transaction IDs of others. This solution is to use [Nginx](https://www.nginx.com/) as a reverse proxy and set up a whitelist of Ethereum RPC methods using [njs](https://nginx.org/en/docs/njs/) to control access to the upstream Ethereum POA nodes.

### Prerequisites
* [Docker](https://www.docker.com/)
* [Golang](https://golang.org/)

### Initialize Ethereum POA Network
* create a sealer account
```bash
$ docker run -it --rm  -v `pwd`/geth-clique:/root/.ethereum ethereum/client-go account new
$ echo "your keystore password" > `pwd`/geth-clique/password.txt
```
* generate genesis configuration
```bash
$ go install github.com/ethereum/go-ethereum/cmd/puppeth@latest
$ puppeth

Please specify a network name to administer (no spaces, hyphens or capital letters please)
> genesis

What would you like to do? (default = stats)
 1. Show network stats
 2. Configure new genesis
 3. Track new remote server
 4. Deploy network components
> 2

What would you like to do? (default = create)
 1. Create new genesis from scratch
 2. Import already existing genesis
> 1

Which consensus engine to use? (default = clique)
 1. Ethash - proof-of-work
 2. Clique - proof-of-authority
> 2

How many seconds should blocks take? (default = 15)
> 5

Which accounts are allowed to seal? (mandatory at least one)
> 0x # Enter the account address created in the previous step

Which accounts should be pre-funded? (advisable at least one)
> 0x # Enter the account address created in the previous step

Should the precompile-addresses (0x1 .. 0xff) be pre-funded with 1 wei? (advisable yes)
> no

Specify your chain/network ID if you want an explicit one (default = random)
> 

What would you like to do? (default = stats)
 1. Show network stats
 2. Manage existing genesis
 3. Track new remote server
 4. Deploy network components
> 2

 1. Modify existing configurations
 2. Export genesis configurations
 3. Remove genesis configuration
> 2

Which folder to save the genesis specs into? (default = current)
  Will create genesis.json, genesis-aleth.json, genesis-harmony.json, genesis-parity.json
> geth-clique
```

* Initialize a new chain
```bash
$ docker run -it --rm  -v `pwd`/geth-clique:/root/.ethereum ethereum/client-go init "/root/.ethereum/genesis.json"
```

### Run Ethereum with nginx proxy
```bash
$ export NETWORKID=`cat geth-clique/genesis.json | jq '.config.chainId'`
$ docker compose up
```
