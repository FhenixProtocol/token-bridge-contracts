// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import {L2ERC20Gateway} from "contracts/tokenbridge/arbitrum/gateway/L2ERC20Gateway.sol";
import {StandardArbERC20} from "contracts/tokenbridge/arbitrum/StandardArbERC20.sol";
import {
    BeaconProxyFactory,
    ClonableBeaconProxy
} from "contracts/tokenbridge/libraries/ClonableBeaconProxy.sol";
import {UpgradeableBeacon} from "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
import {ITokenGateway} from "contracts/tokenbridge/libraries/gateway/ITokenGateway.sol";

contract L2ERC20GatewayTest is Test {
    L2ERC20Gateway public l2Gateway;

    address public l2BeaconProxyFactory;
    address public router = makeAddr("router");
    address public l1Counterpart = makeAddr("l1Counterpart");

    function setUp() public virtual {
        l2Gateway = new L2ERC20Gateway();

        // create beacon
        StandardArbERC20 standardArbERC20 = new StandardArbERC20();
        UpgradeableBeacon beacon = new UpgradeableBeacon(address(standardArbERC20));
        l2BeaconProxyFactory = address(new BeaconProxyFactory());
        BeaconProxyFactory(l2BeaconProxyFactory).initialize(address(beacon));

        L2ERC20Gateway(l2Gateway).initialize(l1Counterpart, router, l2BeaconProxyFactory);
    }

    /* solhint-disable func-name-mixedcase */
    function test_calculateL2TokenAddress() public {
        address l1Token = makeAddr("l1Token");
        assertEq(
            l2Gateway.getUserSalt(l1Token), keccak256(abi.encode(l1Token)), "Invalid user salt"
        );
    }

    function test_cloneableProxyHash() public {
        assertEq(
            l2Gateway.cloneableProxyHash(),
            keccak256(type(ClonableBeaconProxy).creationCode),
            "Invalid proxy hash"
        );
    }

    function test_getOutboundCalldata() public {
        address token = makeAddr("token");
        address from = makeAddr("from");
        address to = makeAddr("to");
        uint256 amount = 200;
        bytes memory data = new bytes(340);

        bytes memory expected = abi.encodeWithSelector(
            ITokenGateway.finalizeInboundTransfer.selector,
            token,
            from,
            to,
            amount,
            abi.encode(0, data)
        );
        bytes memory actual = l2Gateway.getOutboundCalldata(token, from, to, amount, data);

        assertEq(actual, expected, "Invalid outbound calldata");
    }

    function test_getUserSalt() public {
        address l1Token = makeAddr("l1Token");
        assertEq(
            l2Gateway.getUserSalt(l1Token), keccak256(abi.encode(l1Token)), "Invalid user salt"
        );
    }

    function test_initialize() public {
        L2ERC20Gateway gateway = new L2ERC20Gateway();
        L2ERC20Gateway(gateway).initialize(l1Counterpart, router, l2BeaconProxyFactory);

        assertEq(gateway.counterpartGateway(), l1Counterpart, "Invalid counterpartGateway");
        assertEq(gateway.router(), router, "Invalid router");
        assertEq(gateway.beaconProxyFactory(), l2BeaconProxyFactory, "Invalid beacon");
    }

    function test_initialize_revert_InvalidBeacon() public {
        L2ERC20Gateway gateway = new L2ERC20Gateway();
        vm.expectRevert("INVALID_BEACON");
        L2ERC20Gateway(gateway).initialize(l1Counterpart, router, address(0));
    }

    function test_initialize_revert_BadRouter() public {
        L2ERC20Gateway gateway = new L2ERC20Gateway();
        vm.expectRevert("BAD_ROUTER");
        L2ERC20Gateway(gateway).initialize(l1Counterpart, address(0), l2BeaconProxyFactory);
    }

    function test_initialize_revert_InvalidCounterpart() public {
        L2ERC20Gateway gateway = new L2ERC20Gateway();
        vm.expectRevert("INVALID_COUNTERPART");
        L2ERC20Gateway(gateway).initialize(address(0), router, l2BeaconProxyFactory);
    }

    function test_initialize_revert_AlreadyInit() public {
        L2ERC20Gateway gateway = new L2ERC20Gateway();
        L2ERC20Gateway(gateway).initialize(l1Counterpart, router, l2BeaconProxyFactory);
        vm.expectRevert("ALREADY_INIT");
        L2ERC20Gateway(gateway).initialize(l1Counterpart, router, l2BeaconProxyFactory);
    }
}
