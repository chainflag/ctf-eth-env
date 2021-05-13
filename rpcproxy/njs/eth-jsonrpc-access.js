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
        "eth_sendRawTransaction",
        "web3_clientVersion"
    ]
    
    try {
        var payload = JSON.parse(r.requestBody);
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
