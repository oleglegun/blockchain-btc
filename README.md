# Blockchain-BTC

A Go-based implementation of a simple blockchain designed for educational purposes. This project showcases the fundamental concepts of blockchain technology, including cryptographic hashing of transactions and blocks, peer-to-peer networking, and transaction handling. Note that this implementation is basic and not feature-complete, thus advanced features such as consensus mechanism for validator election and public ledger are beyond the scope of this project.

## Features

- **P2P Peer Exchange**: Each node broadcasts its peers to all nodes in the network, facilitating decentralized peer discovery.
- **P2P Transaction Propagation**: Transactions are broadcasted across the network in a decentralized and efficient manner (using Mempool).
- **Transaction Handling and Validation**: Nodes can create and broadcast transactions to the network, ensuring all transactions are validated before inclusion in a block, preventing **double-spending**.
- **Mempool**: A mempool is used for managing transactions before they are included in a block, reducing the overhead of re-broadcasting transactions.
- **Block/Tx/UTXO Storages**: All blockchain data entities are stored is separate memory stores, which can be easily extended by implementing a custom `Store` interface.
- **Protobuf Definitions**: Protocol buffers are used for defining the structure of messages exchanged between nodes.
- **gRPC-based Communication**: Nodes use gRPC for broadcasting transactions and blocks, enabling efficient and scalable communication.
- **Public Key Infrastructure (PKI)**: Transactions use a public key-based addressing system, enhancing security and traceability. Ed25519 signature algorithm is used for transaction/block signing.
- **Multi-node Network Bootstrapping**: The system supports a multi-node setup for testing and development, allowing easy network simulations.
- **Merkle Tree Calculation**: Each block contains a Merkle tree root hash of all transactions, ensuring blockchain data integrity and efficient verification.

## Installation

Clone the repository:

```sh
git clone https://github.com/oleglegun/blockchain-btc.git

cd blockchain-btc
```

Install dependencies:

```sh
go mod download
```

## Usage

To run the blockchain network with 3 nodes (default):

```sh
make run
```

To run the the network with an arbitrary number of nodes:

```sh
make build

./bin/blockchain -nodeCount=4
```

This will start blockchain network with a single (pre-elected) validator node. Transactions are broadcasted to the network each second.

## Project Structure

- `cmd/node/main.go`: Entry point for the blockchain node.
- `internal/cryptography`: Contains cryptographic utilities.
  - `keys.go`: Functions for key generation, signing, and verification.
  - `merkletree.go`: Implementation of Merkle tree for transaction verification.
- `internal/genproto`: Generated protobuf files.
  - `blockchain.pb.go`: Protobuf definitions for blockchain data structures.
  - `blockchain_grpc.pb.go`: gRPC service definitions for blockchain communication.
- `internal/node`: Core blockchain logic, including chain management and transaction handling.
  - `chain.go`: Blockchain chain management.
  - `mempool.go`: Memory pool for pending transactions.
  - `node.go`: Node operations and network communication.
  - `store.go`: Storage for blockchain data.
  - `utxo.go`: Unspent transaction output (UTXO) management.
- `internal/random`: Utilities for generating random data.
  - `random.go`: Functions for generating random hashes and blocks.
- `internal/types`: Extra behavior for the PB generated data structures (blocks, transactions).
  - `block.go`: Block data structure and related functions.
  - `transaction.go`: Transaction data structure and related functions.
- `proto/blockchain.proto`: Protobuf definitions for blockchain data structures and services.
- `Makefile`: Build, run, and test commands for the project.
- `go.mod`: Go module dependencies.
- `go.sum`: Checksums for module dependencies.

## Example output

```sh
make run

Running blockchain with 3 nodes
level=DEBUG node=:3001 msg=running... 
level=DEBUG node=:3001 msg="running validation loop" pubKey=a57f266d17767b307a7dda3e27dbccdfeb4545fe13813f13d0a5a201f4602399
level=DEBUG node=:3002 msg=running...
level=DEBUG node=:3002 msg="discovered new peers" peers=[localhost:3001]
level=DEBUG node=:3003 msg=running...
level=DEBUG node=:3003 msg="discovered new peers" peers=[localhost:3002]
level=DEBUG node=:3001 msg="connected nodes" count=1
level=DEBUG node=:3001 msg="new peer connected" peer=:3002
level=DEBUG node=:3002 msg="connected nodes" count=1
level=DEBUG node=:3002 msg="new peer connected" peer=:3001
level=DEBUG node=:3002 msg="connected nodes" count=2
level=DEBUG node=:3002 msg="new peer connected" peer=:3003
level=DEBUG node=:3003 msg="connected nodes" count=1
level=DEBUG node=:3003 msg="new peer connected" peer=:3002
level=DEBUG node=:3003 msg="discovered new peers" peers=[:3001]
level=DEBUG node=:3001 msg="connected nodes" count=2
level=DEBUG node=:3001 msg="new peer connected" peer=:3003
level=DEBUG node=:3003 msg="connected nodes" count=2
level=DEBUG node=:3003 msg="new peer connected" peer=:3001
level=DEBUG node=:3002 msg="received tx" from=[::1]:52277 tx=86a382becd80e93334869574b75079b470c2545321bf7c28ff06223c92d9e7c7
level=DEBUG node=:3001 msg="received tx" from=[::1]:52248 tx=86a382becd80e93334869574b75079b470c2545321bf7c28ff06223c92d9e7c7
level=DEBUG node=:3002 msg="received tx" from=[::1]:52277 tx=132311d3cf73751f49b7435461fface3b046ab5139288d7c126b0a5ceebfa2fd
level=DEBUG node=:3001 msg="received tx" from=[::1]:52248 tx=132311d3cf73751f49b7435461fface3b046ab5139288d7c126b0a5ceebfa2fd
level=DEBUG node=:3003 msg="received tx" from=[::1]:52304 tx=132311d3cf73751f49b7435461fface3b046ab5139288d7c126b0a5ceebfa2fd
level=DEBUG node=:3003 msg="received tx" from=[::1]:52304 tx=86a382becd80e93334869574b75079b470c2545321bf7c28ff06223c92d9e7c7
level=DEBUG node=:3001 msg="creating new block" txs=2
level=DEBUG node=:3002 msg="received tx" from=[::1]:52277 tx=a7d541420f5768b341b70259be2a5db5260641fd7e5383e1c0f953f8d69b49be
level=DEBUG node=:3003 msg="received tx" from=[::1]:52304 tx=a7d541420f5768b341b70259be2a5db5260641fd7e5383e1c0f953f8d69b49be
level=DEBUG node=:3001 msg="received tx" from=[::1]:52248 tx=a7d541420f5768b341b70259be2a5db5260641fd7e5383e1c0f953f8d69b49be
...
```