function access(r) {
    var whitelist = [
        "net_version",
        "rpc_modules",
        "eth_chainId",
        "eth_getBalance",
        "eth_getCode",
        "eth_getStorageAt",
        "eth_call",
        "eth_getTransactionByHash",
        "eth_getTransactionReceipt",
        "eth_getTransactionCount",
        "eth_sendRawTransaction",
        "web3_clientVersion",
        "eth_estimateGas",
        "eth_gasPrice",
        "eth_blockNumber"
    ].map(method => method.toLowerCase());

    try {
        var payload = JSON.parse(r.requestBody.toLowerCase());
        if (payload.jsonrpc !== "2.0") {
            r.return(401, "jsonrpc version not supported\n");
            return;
        }
        if (!whitelist.includes(payload.method)) {
            r.return(401, "jsonrpc method is not allow\n");
            return;
        }
    } catch (error) {
        r.return(415, "Cannot parse payload into JSON\n");
        return;
    }

    r.internalRedirect('@jsonrpc');
}

export default { access }
