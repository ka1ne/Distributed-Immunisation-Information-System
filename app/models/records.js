const mongoose = require('mongoose');
const Schema = mongoose.Schema;

const Record = new Schema ({
        uuid: { type: String, required: true },
        timestamp: { type: Date, required: false },
        owner: { type: String, required: false },
        exp: { type: Date, required: false },
});

module.exports = mongoose.model('Record', Record)
