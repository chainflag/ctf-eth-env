function access(r) {
    try {
        var payload = JSON.parse(r.requestBody);
        if (!isAllowedMethod(payload.method)) {
            r.return(401, "jsonrpc method is not allow\n");
            return;
        }
    } catch (error) {
        r.return(415, "Cannot parse payload into JSON\n");
        return;
    }

    r.internalRedirect('@jsonrpc');
}

function isAllowedMethod (method) {
    return true;
}

export default {access}
