package mocks

//go:generate moq -pkg mocks -out node_client.go ../ NodeClient
//go:generate moq -pkg mocks -out blockchain_client.go ../ BlockChainClient
//go:generate moq -pkg mocks -out control_client.go ../ ControlClient
//go:generate moq -pkg mocks -out mining_client.go ../ MiningClient
//go:generate moq -pkg mocks -out network_client.go ../ NetworkClient
//go:generate moq -pkg mocks -out transaction_client.go ../ TransactionClient
//go:generate moq -pkg mocks -out util_client.go ../ UtilClient
//go:generate moq -pkg mocks -out wallet_client.go ../ WalletClient

// Third party

//go:generate moq -pkg mocks -out zmq4_mock.go ../vendor/github.com/go-zeromq/zmq4 Socket
