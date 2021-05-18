const mongoose = require('mongoose');
const Schema = mongoose.Schema;

const Record = new Schema ({
        uuid: { type: String, required: true },
        timestamp: { type: Date, required: true },
        owner: { type: String, required: true },
        exp: { type: Number, required: true },
});

module.exports = mongoose.model('Record', Record)
