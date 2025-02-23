This swagger site hosts the service interface definitions for the Ethereum Serenity JSON API.

## Service Overview

At the core of Ethereum Serenity lies the "Beacon Chain", a proof-of-stake based blockchain managing a set of validators and "shards" across a network, which are stateful, smaller blockchains used for transaction and smart contract execution. The Beacon Chain reconciles consensus results across shards to provide the backbone of Ethereum Serenity. The services below allow for retrieval of current and historical information relating to validators, consensus, blocks, and staking participation.

| Package | Service     | Version  | Description                                                                                                                                                                                                                                |
| ------- | ----------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| eth     | BeaconChain | v1alpha1 | This service is used to retrieve critical data relevant to the Ethereum Beacon Chain, including the most recent head block, current pending deposits, the chain state and more.                                                            |
| eth     | Node        | v1alpha1 | The Node service returns information about the Ethereum node itself, including versioning and general information as well as network sync status and a list of services currently implemented on the node.                                 |
| eth     | Validator   | v1alpha1 | This API provides the information a validator needs to retrieve throughout its lifecycle, including received assignments from the network, its current index in the state, as well the rewards and penalties that have been applied to it. |

### JSON Mapping

The majority of field primitive types for Ethereum are either `bytes` or `uint64`. The canonical JSON mapping for those fields are a Base64 encoded string for `bytes`, or a string representation of `uint64`. Since JavaScript loses precision for values over [MAX_SAFE_INTEGER](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER), uint64 must be a JSON string in order to accommodate those values. If the field value not changed and is still set to protobuf's default, this field will be omitted from the JSON encoding entirely.

For more details on JSON mapping for other types, view the relevant section in the [proto3 language guide](https://developers.google.com/protocol-buffers/docs/proto3#json).

### Helpful gRPC Tooling and Resources

- [Google's API Style Guide](https://cloud.google.com/apis/design/)
- [Language reference for protoc3](https://developers.google.com/protocol-buffers/docs/proto3)
- [Awesome gRPC](https://github.com/grpc-ecosystem/awesome-grpc)
- [Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
- [Transcoding gRPC to JSON/HTTP using Envoy](https://blog.jdriven.com/2018/11/transcoding-grpc-to-http-json-using-envoy/)