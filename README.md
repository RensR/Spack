# Spack

Spack parses Solidity structs and packs the fields efficiently to reduce the
number of storage slots they use. It also adds struct packing comments to clearly indicate
how the fields are packed.

It can deal with comments and whitespace in the struct definition, and will 
preserve them in the output. It handles unknown types by assuming they cannot be 
packed, treating they as `bytes32`.

## Disclaimer

This code is a work in progress and can contain bugs. Use it at your own risk.
Feature request and bug reports are welcome.

### Example

input

```solidity
    struct RequestMeta {
        uint64 completedRequests;
        Custom.Datatype data;
        address requestingContract;
        uint72 adminFee; // in wei
        address subscriptionOwner;
        bytes32 flags; // 32 bytes of flags
        uint96 availableBalance; // in wei. 0 if not specified.
        uint64 subscriptionId;
        uint64 initiatedRequests;// number of requests initiated by this contract
        uint32 callbackGasLimit;
        uint16 dataVersion;
    }
```

output

```solidity
    struct RequestMeta {
        Custom.Datatype data; //                     
        bytes32 flags; //                  32 bytes of flags
        address requestingContract; // ──┐
        uint96 availableBalance; // ─────┘ in wei. 0 if not specified.
        address subscriptionOwner; // ───┐
        uint64 completedRequests; //     │
        uint32 callbackGasLimit; // ─────┘
        uint72 adminFee; // ─────────────┐ in wei
        uint64 subscriptionId; //        │
        uint64 initiatedRequests; //     │ number of requests initiated by this contract
        uint16 dataVersion; // ──────────┘
    }
```

## Quickstart

```bash
go build && ./spack examples/RequestMeta.txt
```

## TODO

- [ ] Add more flexible command line options
- [ ] Add tests
- [ ] Improve error handling
