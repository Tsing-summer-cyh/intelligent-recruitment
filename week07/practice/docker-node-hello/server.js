const http = require('http');

const port = Number(process.env.PORT) || 3000;

http.createServer((req, res) => {
  res.writeHead(200, { 'Content-Type': 'text/plain; charset=utf-8' });
  res.end('Hello from Node.js in Docker\n');
}).listen(port, () => {
  console.log(`Listening on http://0.0.0.0:${port}`);
});