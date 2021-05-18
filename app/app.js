const express = require('express');
const app = express();
const router = express.Router();
const db = require('./db');
const recordRouter = require('./routes/records');

const path = __dirname + '/views/';
const port = 3000;

app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');
app.use(express.urlencoded({ extended: true }));
app.use(express.static(path));
app.use('/records', recordRouter);

app.listen(port, function () {
  console.log('Example app listening on port 3000!')
})