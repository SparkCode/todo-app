import http from 'http';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PORT = process.env.PORT || 3000;

const server = http.createServer((req, res) => {
  let filePath = '';
  
  // Route handling
  if (req.url === '/' || req.url === '/index.html') {
    filePath = path.join(__dirname, '..', 'dist', 'index.html');
  } else {
    // Serve static assets from the dist folder
    filePath = path.join(__dirname, '..', 'dist', req.url || '');
  }
  
  // Determine content type based on file extension
  const extname = path.extname(filePath);
  let contentType = 'text/html; charset=utf-8';
  
  switch (extname) {
    case '.js':
      contentType = 'application/javascript';
      break;
    case '.css':
      contentType = 'text/css';
      break;
    case '.json':
      contentType = 'application/json';
      break;
    case '.png':
      contentType = 'image/png';
      break;
    case '.jpg':
      contentType = 'image/jpg';
      break;
  }
  
  // Read and serve the file
  fs.readFile(filePath, (err: NodeJS.ErrnoException | null, content: Buffer) => {
    if (err) {
      if (err.code === 'ENOENT') {
        res.statusCode = 404;
        res.setHeader('Content-Type', 'text/html; charset=utf-8');
        res.end('<h1>404 - File Not Found</h1>');
      } else {
        res.statusCode = 500;
        res.setHeader('Content-Type', 'text/html; charset=utf-8');
        res.end(`<h1>500 - Server Error</h1><p>${err.code}</p>`);
      }
    } else {
      res.statusCode = 200;
      res.setHeader('Content-Type', contentType);
      res.end(content);
    }
  });
});

server.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}/`);
});

