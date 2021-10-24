FROM trufflesuite/ganache-cli:v6.12.2

EXPOSE 8545

ENTRYPOINT ["node", "/app/ganache-core.docker.cli.js", "--accounts", "1", "--defaultBalanceEther", "100000000", "--fork", "https://cloudflare-eth.com"]
