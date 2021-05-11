# ctf-eth-env

### Prerequisites
* [Docker](https://www.docker.com/)
* [Golang](https://golang.org/)

### Initialize Ethereum POA Network
* create a sealer account
```bash
docker run -it --rm  -v `pwd`/geth-clique:/root/.ethereum ethereum/client-go account new
echo "your keystore password" > `pwd`/geth-clique/password.txt
```
* generate genesis configuration
```bash
go install github.com/ethereum/go-ethereum/cmd/puppeth@latest
puppeth

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
docker run -it --rm  -v `pwd`/geth-clique:/root/.ethereum ethereum/client-go init "/root/.ethereum/genesis.json"
```

### Run Ethereum with nginx proxy
```
docker compose up
```
