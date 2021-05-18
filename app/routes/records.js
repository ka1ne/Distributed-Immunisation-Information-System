const express = require('express');
const router = express.Router();
const record = require('../controllers/records');

router.get('/', function(req, res){
    record.index(req,res);
});

router.post('/addrecord', function(req, res) {
    record.create(req,res);
});

router.get('/genuuid', function(req, res) {
    record.createuuid(req,res);
});

router.post('/checkrecord', function(req, res) {
    record.check(req,res);
});

module.exports = router;
