const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const app = express();
const PORT = 4000;

app.use(express.json());

// Connect to Fabric Network
async function connectToFabric() {
    const ccpPath = path.resolve(__dirname, 'connection-org1.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'appUser',
        discovery: { enabled: true, asLocalhost: true },
    });

    return gateway;
}

// Create Asset API
app.post('/createAsset', async (req, res) => {
    const { dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks } = req.body;

    try {
        const gateway = await connectToFabric();
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('assetTransfer');

        await contract.submitTransaction('CreateAsset', dealerID, msisdn, mpin, balance.toString(), status, transAmount.toString(), transType, remarks);
        res.status(200).send(`Asset with DealerID ${dealerID} created successfully`);

        await gateway.disconnect();
    } catch (error) {
        res.status(500).send(`Error creating asset: ${error}`);
    }
});

// Query Asset API
app.get('/queryAsset/:dealerID', async (req, res) => {
    const { dealerID } = req.params;

    try {
        const gateway = await connectToFabric();
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('assetTransfer');

        const result = await contract.evaluateTransaction('QueryAsset', dealerID);
        res.status(200).send(`Asset data: ${result.toString()}`);

        await gateway.disconnect();
    } catch (error) {
        res.status(500).send(`Error querying asset: ${error}`);
    }
});

app.listen(PORT, () => {
    console.log(`API server running on http://localhost:${PORT}`);
});
