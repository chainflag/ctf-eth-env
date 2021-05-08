FROM ethereum/client-go:latest

RUN apk update && apk add bash curl coreutils jq

COPY genesis.json /geth/genesis.json

RUN geth --nodiscover --datadir "/geth/chain" init "/geth/genesis.json"

WORKDIR /
COPY start.sh /start.sh
RUN chmod 555 /start.sh

ENTRYPOINT [ "/start.sh" ]
