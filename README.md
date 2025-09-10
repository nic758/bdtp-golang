# Blockchain data transfer protocol
The BDTP is a lightweight protocol designed to facilitate secure and reliable communication and data storage on blockchains. It treats blockchain addresses as pointers, enabling efficient data transfer between clients and servers on supported chains.<br>

This Protocol helps developers store data seamlessly on chain and is ideal for projects that require immutable and verifiable data storage linked to blockchain transactions. It will do everything to use the blockchain as a user.

## Currently supported chain
- Polygon
- Waves

## Requirements
Go version >= 1.17
## Installation
1.  
    ```sh
    git clone git@github.com:nic758/bdtp-golang.git
    go mod tidy
    ```
2. Create a wallet and set up a network with the following details: 
    1. RPC URL: `https://rpc-amoy.polygon.technology/`
    2. ID chain: `80002`
    3. Symbol: `POL`
    4. Block explorer URL: `https://www.oklink.com/amoy`
3. Copy the wallet recovery phrase and set export 
    ```sh
    export BDTP_SEED="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    ```
4. Create a Etherscan account and generate an API key, then set export
    ```sh
    export POLYSCAN_API_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    ```
5. To fund your wallet with Polygon Amoy test tokens, go to https://faucet.polygon.technology/<br/>

To run a server: 
```sh 
go run ./main.go server 
```
To upload a file to the blockchain via BDTP **with a server running**: 
```sh 
go run ./main.go forge {PATH}
```