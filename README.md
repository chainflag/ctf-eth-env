# ctf-eth-env

An ethereum node environment for CTF challenges.

The solution is to use [Nginx](https://www.nginx.com/) as a reverse proxy and set up the whitelist of Ethereum RPC methods by using [njs](https://nginx.org/en/docs/njs/) for access control to the upstream Ethereum POA nodes, and thus implement an anti-plagiarism server-side environment.

## Background
It is unfair that some CTF blockchain challenge players can cheat by searching back the blockchain history, where all the transactions of those who have solved the challenges are recorded. These dishonest players can solve the challenges simply by replaying the transactions. The root cause of this problem is that all data in the permissionless blockchain is public and everyone can fetch it by querying the specified RPC methods.  

So the idea of this project is to disable several RPC methods (e.g. `eth_getBlockByHash`, `eth_getBlockByNumber`) of an Ethereum POA node and then use it as the challenge server-side environment. In this way, players on the client side have no longer any access to the transaction IDs of others. 

## Usage
1. Clone the repository
```
git clone https://github.com/chainflag/ctf-eth-env.git
cd ctf-eth-env
```

2. Create a sealer account
```bash
docker run -it --rm  -v `pwd`/config:/root/.ethereum ethereum/client-go account new
echo "your keystore password" > `pwd`/config/password.txt
```

3. Generate genesis config
```
go get github.com/chainflag/ctf-eth-env/genesis-builder
genesis-builder --address "address created in previous step"
```

4. Run docker container
```bash
docker-compose up -d
```

**Open Ports**

| Service                 | Port
| ----------------------- | -----
| json-rpc with whitelist | 8545      
| ether faucet            | 8080

## Related Project
* [eth-faucet](https://github.com/chainflag/eth-faucet)
* [eth-challenge-base](https://github.com/chainflag/eth-challenge-base)

## Contributing

PRs accepted.

## License

The MIT License (MIT)
