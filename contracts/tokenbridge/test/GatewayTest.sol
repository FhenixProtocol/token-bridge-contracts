// SPDX-License-Identifier: Apache-2.0

/*
 * Copyright 2020, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

pragma solidity ^0.6.11;

import "../ethereum/router/L1CustomGateway.sol";
import "../ethereum/router/L1ERC20Gateway.sol";
import "../arbitrum/router/L2CustomGateway.sol";
import "../arbitrum/router/L2ERC20Gateway.sol";

contract L1GatewayTester is L1ERC20Gateway {
    function isCounterpartGateway() internal view virtual override returns (bool) {
        return msg.sender == counterpartGateway;
    }

    function createOutboundTx(
        address _user,
        uint256 _maxSubmissionCost,
        uint256 _maxGas,
        uint256 _gasPriceBid,
        bytes memory _data
    ) internal virtual override returns (uint256) {
        (bool success, bytes memory retdata) = counterpartGateway.call(_data);
        require(success, "OUTBOUND_REVERT");
        return 1337;
    }
}

contract L2GatewayTester is L2ERC20Gateway {
    function isCounterpartGateway() internal view virtual override returns (bool) {
        return msg.sender == counterpartGateway;
    }

    function createOutboundTx(bytes memory _data) internal virtual override returns (uint256) {
        (bool success, bytes memory retdata) = counterpartGateway.call(_data);
        require(success, "OUTBOUND_REVERT");
        return 1337;
    }

    function arbgasReserveIfCallRevert() internal pure virtual override returns (uint256) {
        return 50000;
    }
}

contract L1CustomGatewayTester is L1CustomGateway {
    event EventErr(bytes errmsg);

    function isCounterpartGateway() internal view virtual override returns (bool) {
        return msg.sender == counterpartGateway;
    }

    function createOutboundTx(
        address _user,
        uint256 _maxSubmissionCost,
        uint256 _maxGas,
        uint256 _gasPriceBid,
        bytes memory _data
    ) internal virtual override returns (uint256) {
        (bool success, bytes memory retdata) = counterpartGateway.call(_data);
        emit EventErr(retdata);
        // require(success, retdata);
        return 1337;
    }
}

contract L2CustomGatewayTester is L2CustomGateway {
    function isCounterpartGateway() internal view virtual override returns (bool) {
        return msg.sender == counterpartGateway;
    }

    function createOutboundTx(bytes memory _data) internal virtual override returns (uint256) {
        (bool success, bytes memory retdata) = counterpartGateway.call(_data);
        require(success, "OUTBOUND_REVERT");
        return 1337;
    }

    function arbgasReserveIfCallRevert() internal pure virtual override returns (uint256) {
        return 50000;
    }
}
