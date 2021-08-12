FROM ethereum/client-go:v1.10.6

RUN apk add --no-cache jq

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8545 8546

ENTRYPOINT ["/entrypoint.sh"]
