const express = require('express');
const fs = require('fs');
const path = require('path');

const app = express();

app.get('/', (req, res) => {
  const directoryPath = path.join(__dirname, 'cdn');

  fs.readdir(directoryPath, (err, files) => {
    if (err) {
      console.log('Unable to scan directory: ' + err);
      return res.status(500).send('Internal Server Error');
    }

    const fileList = files.map((file) => `<li><a href="/cdn/${file}">${file}</a></li>`).join('');

    const html = `
      <!DOCTYPE html>
      <html>
      <head>
        <title>File Server</title>
      </head>
      <body>
        <h1>File Server</h1>
        <ul>
          ${fileList}
        </ul>
      </body>
      </html>
    `;

    res.send(html);
  });
});

app.use('/cdn', express.static(path.join(__dirname, 'cdn')));

app.listen(8080, () => {
  console.log('Server is running on port 8080');
});
