This project implements a smart contract for managing financial assets using Hyperledger Fabric. It allows for asset creation, querying, and transaction history retrieval.
To deploy, clone the repo, navigate to fabric-samples/test-network, and start the network using ./network.sh up.
Chaincode can be deployed with ./network.sh deployCC -ccn assetTransfer -ccp ../asset-transfer-go -ccl go.
Licensed under MIT.