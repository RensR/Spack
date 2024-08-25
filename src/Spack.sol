// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.24;

contract Spack {
    struct RequestMeta {
        uint64 completedRequests;
        DataType data;
        address requestingContract;
        uint72 adminFee; // in wei
        address subscriptionOwner;
        bytes32 flags; // 32 bytes of flags
        uint96 availableBalance; // in wei. 0 if not specified.
        uint64 subscriptionId;
        uint64 initiatedRequests; // number of requests initiated by this contract
        uint32 callbackGasLimit;
        uint16 dataVersion;
    }

    struct DataType {
        uint96 timestamp;
        address sender;
    }
}
