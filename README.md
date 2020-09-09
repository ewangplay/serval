# Serval

This is an implementation of the [W3C Decentralized Identifiers (DIDs) Specification](https://www.w3.org/TR/did-core/).

## Prerequisites

### Blockchain Network

In order to run the serval service, you must already have a Hyperledger Fabric Blockchain Network(version >= 2.0).

You can deploy a test network by using scripts that are provided in the [fabric-samples repository](https://github.com/hyperledger/fabric-samples). Please refer to [Using the Fabric test network](https://hyperledger-fabric.readthedocs.io/en/release-2.2/test_network.html).

### DID Chaincode

After bring up the blockchain network, you should deploy the DID smart contract to a channel (eg. mychannel).

If you are using the test network provided by fabric-sample, you can easily install the smart contract using the 'network.sh' script, like this:
```
./network.sh deployCC -ccn did -ccp /path/to/chaincode/did/go -ccv 1.0 -ccl go -ccep "AND('Org1MSP.member)','Org2MSP.member'"
```
This will request one signature from each of the two principals. Under this endorsement policy, the 'endorsingPeers' configure item must be set at least two principals, like:
```
endorsingPeers:
    - peer0.org1.example.com:7051
    - peer0.org2.example.com:9051
```

If using another endorsement policy, such as:
```
./network.sh deployCC -ccn did -ccp /path/to/chaincode/did/go -ccv 1.0 -ccl go -ccep "OR('Org1MSP.member)','Org2MSP.member'"
```
This will request one signature from either one of the two principals. Under this endorsement policy, the 'endorsingPeers' configure item can be set either one of the two principals, like:
```
endorsingPeers:
    - peer0.org1.example.com:7051
```

When the blockchain network and DID smart contract are all ready, you should update configuare. 

The 'blockchain' section in 'serval.yaml' represents the blockchain configure, you should update these items based on actual values.

In the 'blockchain' directory, save the blockchain connection files, you should replace these files according to the actual situation.

### Application Key

Application Key is used to sign / verify DID Document.

Before running serval service, Application Key should be generated first and set in configure file.

## How to build, install and run

### build service 

```
make srvc
```

### code checking and unit testing

```
make checks
```

### install service

```
make install
```

### run service

```
/opt/serval/bin/serval --config /opt/serval/etc/serval.yaml
```