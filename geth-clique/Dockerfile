# > docker build -t geth .
# > docker run -it -p 8545:8545 -v `pwd`/example:/root/.ethereum geth
FROM ethereum/client-go:v1.10.6

RUN apk add --no-cache jq

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8545

ENTRYPOINT ["/entrypoint.sh"]
