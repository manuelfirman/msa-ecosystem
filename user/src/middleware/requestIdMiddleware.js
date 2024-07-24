const requestIdMiddleware = (req, res, next) => {
    req.requestId = req.headers['x-request-id'];
    if (!req.requestId) {
        res.status(400).json({ error: 'Missing X-Request-ID header' });
        return;
    }
    req.traceInfo = req.headers['x-trace-info'];
    if (!req.traceInfo) {
        res.status(400).json({ error: 'Missing X-Trace-Info header' });
        return;
    }
    const traceInfo = req.traceInfo + ", " + `http://${process.env.SERVICE_NAME}:${process.env.PORT}${req.originalUrl}`;
    res.setHeader('X-Request-ID', req.requestId);
    res.setHeader('X-Trace-Info', traceInfo);
    next();
};

module.exports = requestIdMiddleware;