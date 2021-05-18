'use strict';

const path = require('path');
const Record = require('../models/records');
const uuid4 = require('uuid').v4;


const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('../utils/CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('../utils/AppUtil.js');

const channelName = 'mychannel';
const chaincodeName = 'record';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, '..', 'wallet');
const org1UserId = 'appUser';

const settings = setup();

function prettyJSONString(inputString) {
        return JSON.stringify(JSON.parse(inputString), null, 2);
}

async function setup() {
        try {
                // build an in memory object with the network configuration (also known as a connection profile)
                const ccp = buildCCPOrg1();

                // build an instance of the fabric ca services client based on
                // the information in the network configuration
                const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');

                // setup the wallet to hold the credentials of the application user
                const wallet = await buildWallet(Wallets, walletPath);

                // in a real application this would be done on an administrative flow, and only once
                await enrollAdmin(caClient, wallet, mspOrg1);

                // in a real application this would be done only when a new user was required to be added
                // and would be part of an administrative flow
                await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');

                // Create a new gateway instance for interacting with the fabric network.
                // In a real application this would be done as the backend server session is setup for
                // a user that has been verified.
                const gateway = new Gateway()

                return {gateway : gateway, ccp : ccp};
        } catch (error) {
                console.error(`******** FAILED to run the application: ${error}`);
        }

}
exports.setup = setup;

async function submit(record, settings) {
        try {



                try {
                        // setup the gateway instance
                        // The user will now be able to create connections to the fabric network and be able to
                        // submit transactions and query. All transactions submitted by this gateway will be
                        // signed by this user using the credentials stored in the wallet.
                        await settings.gateway.connect(settings.ccp, {
                                wallet,
                                identity: org1UserId,
                                discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
                        });

                        // Build a network instance based on the channel where the smart contract is deployed
                        const network = await settings.gateway.getNetwork(channelName);

                        // Get the contract from the network.
                        const contract = network.getContract(chaincodeName);

                        // Initialize a set of asset data on the channel using the chaincode 'InitLedger' function.
                        // This type of transaction would only be run once by an application the first time it was started after it
                        // deployed the first time. Any updates to the chaincode deployed later would likely not need to run
                        // an "init" type function.
                        console.log('\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger');
                        await contract.submitTransaction('InitLedger');
                        console.log('*** Result: committed');

                        // Let's try a query type operation (function).
                        // This will be sent to just one peer and the results will be shown.
                        console.log('\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger');
                        let result = await contract.evaluateTransaction('GetAllAssets');
                        console.log(`*** Result: ${prettyJSONString(result.toString())}`);

                        // Now let's try to submit a transaction.
                        // This will be sent to both peers and if both peers endorse the transaction, the endorsed proposal will be sent
                        // to the orderer to be committed by each of the peer's to the channel ledger.
                        console.log('\n--> Submit Transaction: CreateAsset, creates new asset with UUID, timestamp, owner and expiration arguments');
                        result = await contract.submitTransaction('CreateAsset', record.uuid, record.timestamp, record.owner, record.exp);
                        console.log('*** Result: committed');
                        if (`${result}` !== '') {
                                console.log(`*** Result: ${prettyJSONString(result.toString())}`);
                        }


                } finally {
                        // Disconnect from the gateway when the application is closing
                        // This will close all connections to the network
                        // gateway.disconnect();
                }
        } catch (error) {
                console.error(`******** FAILED to run the application: ${error}`);
        }
}
exports.submit = submit;


async function listAll(record) {
        try {
                // build an in memory object with the network configuration (also known as a connection profile)
                const ccp = buildCCPOrg1();

                // build an instance of the fabric ca services client based on
                // the information in the network configuration
                const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');

                // setup the wallet to hold the credentials of the application user
                const wallet = await buildWallet(Wallets, walletPath);

                // in a real application this would be done on an administrative flow, and only once
                await enrollAdmin(caClient, wallet, mspOrg1);

                // in a real application this would be done only when a new user was required to be added
                // and would be part of an administrative flow
                await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');

                // Create a new gateway instance for interacting with the fabric network.
                // In a real application this would be done as the backend server session is setup for
                // a user that has been verified.
                const gateway = new Gateway();

                try {
                        // setup the gateway instance
                        // The user will now be able to create connections to the fabric network and be able to
                        // submit transactions and query. All transactions submitted by this gateway will be
                        // signed by this user using the credentials stored in the wallet.
                        await gateway.connect(ccp, {
                                wallet,
                                identity: org1UserId,
                                discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
                        });

                        // Build a network instance based on the channel where the smart contract is deployed
                        const network = await gateway.getNetwork(channelName);

                        // Get the contract from the network.
                        const contract = network.getContract(chaincodeName);

                        // Initialize a set of asset data on the channel using the chaincode 'InitLedger' function.
                        // This type of transaction would only be run once by an application the first time it was started after it
                        // deployed the first time. Any updates to the chaincode deployed later would likely not need to run
                        // an "init" type function.
                        console.log('\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger');
                        await contract.submitTransaction('InitLedger');
                        console.log('*** Result: committed');

                        // Let's try a query type operation (function).
                        // This will be sent to just one peer and the results will be shown.
                        console.log('\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger');
                        let result = await contract.evaluateTransaction('GetAllAssets');
                        console.log(`*** Result: ${prettyJSONString(result.toString())}`);

                        // Now let's try to submit a transaction.
                        // This will be sent to both peers and if both peers endorse the transaction, the endorsed proposal will be sent
                        // to the orderer to be committed by each of the peer's to the channel ledger.
                        console.log('\n--> Submit Transaction: CreateAsset, creates new asset with UUID, timestamp, owner and expiration arguments');
                        result = await contract.submitTransaction('CreateAsset', record.uuid, record.timestamp, record.owner, record.exp);
                        console.log('*** Result: committed');
                        if (`${result}` !== '') {
                                console.log(`*** Result: ${prettyJSONString(result.toString())}`);
                        }


                } finally {
                        // Disconnect from the gateway when the application is closing
                        // This will close all connections to the network
                        gateway.disconnect();
                }
        } catch (error) {
                console.error(`******** FAILED to run the application: ${error}`);
        }
}
exports.listAll = listAll;

//const createHash = require('crypto');

//uuidv4(); // ⇨ '9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d'

//const hash = createHash('sha256');

//hash.update('some data to hash');
//console.log(hash.digest('hex'));

exports.index = function (req, res) {
        res.sendFile(path.resolve('views/records.html'));
};

exports.check = function (req, res) {
        Record.find({}).exec(function (err, records) {
                if (err) {
                        return res.send(500, err);
                }
                res.render('checkrecord', {
                        records: records.uuid
                });
        });
};

exports.create = function (req, res) {
        var newRecord = new Record(req.body);
        console.log(req.body);
        newRecord.save(function (err) {
                if (err) {
                        res.status(400).send('Unable to save record to database');
                } else {
                        submit(newRecord, settings);
                        res.redirect('/records/getrecord');
                }
        });
};

exports.list = function (req, res) {
        Record.find({}).exec(function (err, records) {
                if (err) {
                        return res.send(500, err);
                }
                res.render('getrecord', {
                        records: records
                });
        });
};


exports.createuuid = function (req, res) {
        var uuid = uuid4();
        res.send(uuid); // , "<a href="/records" role="button">Add Records</a>"
};