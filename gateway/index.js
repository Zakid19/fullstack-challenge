const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();

app.use('/products', createProxyMiddleware({ target: 'http://product-service:3000', changeOrigin: true }));
app.use('/orders', createProxyMiddleware({ target: 'http://order-service:4000', changeOrigin: true }));

app.listen(8080, () => console.log('gateway on 8080'));
