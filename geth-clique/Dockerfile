FROM ethereum/client-go:latest

RUN apk update && apk add bash curl coreutils jq

COPY genesis.json /geth/genesis.json
ADD dev-key.json /geth/keys/dev-key.json
RUN echo "" > /geth/keys/dev-key-password.txt

RUN geth --nodiscover --datadir "/geth/chain" --keystore "/geth/keys" init "/geth/genesis.json"

WORKDIR /
COPY start.sh /start.sh
RUN chmod 555 /start.sh

ENTRYPOINT [ "/start.sh" ]
