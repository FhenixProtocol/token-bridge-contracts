import {
  L1Network,
  L1ToL2MessageGasEstimator,
  L1ToL2MessageStatus,
  L1TransactionReceipt,
  L2Network,
  L2TransactionReceipt,
} from '@arbitrum/sdk'
import { getBaseFee } from '@arbitrum/sdk/dist/lib/utils/lib'
import { JsonRpcProvider } from '@ethersproject/providers'
import { expect } from 'chai'
import { setupTokenBridgeInLocalEnv } from '../scripts/local-deployment/localDeploymentLib'
import {
  AeWETH__factory,
  BridgedUsdcCustomToken__factory,
  ERC20,
  ERC20__factory,
  IERC20Bridge__factory,
  IInbox__factory,
  IOwnable__factory,
  L1GatewayRouter__factory,
  L1OrbitCustomGateway__factory,
  L1OrbitERC20Gateway__factory,
  L1OrbitGatewayRouter__factory,
  L1USDCCustomGateway__factory,
  L2CustomGateway__factory,
  L2GatewayRouter__factory,
  L2USDCCustomGateway__factory,
  MockL1Usdc__factory,
  MockL2Usdc__factory,
  ProxyAdmin__factory,
  TestArbCustomToken__factory,
  TestCustomTokenL1__factory,
  TestERC20,
  TestERC20__factory,
  TestOrbitCustomTokenL1__factory,
  TransparentUpgradeableProxy__factory,
  UpgradeExecutor__factory,
} from '../build/types'
import { defaultAbiCoder } from 'ethers/lib/utils'
import { BigNumber, Wallet, ethers } from 'ethers'
import { exit } from 'process'

const config = {
  parentUrl: 'http://localhost:8547',
  childUrl: 'http://localhost:3347',
}

const LOCALHOST_L3_OWNER_KEY =
  '0xecdf21cb41c65afb51f91df408b7656e2c8739a5877f2814add0afd780cc210e'

let parentProvider: JsonRpcProvider
let childProvider: JsonRpcProvider

let deployerL1Wallet: Wallet
let deployerL2Wallet: Wallet

let userL1Wallet: Wallet
let userL2Wallet: Wallet

let _l1Network: L1Network
let _l2Network: L2Network

let token: TestERC20
let l2Token: ERC20
let nativeToken: ERC20 | undefined

describe('orbitTokenBridge', () => {
  // configure orbit token bridge
  before(async function () {
    parentProvider = new ethers.providers.JsonRpcProvider(config.parentUrl)
    childProvider = new ethers.providers.JsonRpcProvider(config.childUrl)

    const deployerKey = ethers.utils.sha256(
      ethers.utils.toUtf8Bytes('user_token_bridge_deployer')
    )
    deployerL1Wallet = new ethers.Wallet(deployerKey, parentProvider)
    deployerL2Wallet = new ethers.Wallet(deployerKey, childProvider)

    const { l1Network, l2Network } = await setupTokenBridgeInLocalEnv()

    _l1Network = l1Network
    _l2Network = l2Network

    // create user wallets and fund it
    const userKey = ethers.utils.sha256(ethers.utils.toUtf8Bytes('user_wallet'))
    userL1Wallet = new ethers.Wallet(userKey, parentProvider)
    userL2Wallet = new ethers.Wallet(userKey, childProvider)
    await (
      await deployerL1Wallet.sendTransaction({
        to: userL1Wallet.address,
        value: ethers.utils.parseEther('10.0'),
      })
    ).wait()

    const nativeTokenAddress = await getFeeToken(
      l2Network.ethBridge.inbox,
      parentProvider
    )
    nativeToken =
      nativeTokenAddress === ethers.constants.AddressZero
        ? undefined
        : ERC20__factory.connect(nativeTokenAddress, userL1Wallet)
  })

  it('should have deployed token bridge contracts', async function () {
    // get router as entry point
    const l1Router = L1OrbitGatewayRouter__factory.connect(
      _l2Network.tokenBridge.l1GatewayRouter,
      parentProvider
    )

    expect((await l1Router.defaultGateway()).toLowerCase()).to.be.eq(
      _l2Network.tokenBridge.l1ERC20Gateway.toLowerCase()
    )
  })

  it('can deposit token via default gateway', async function () {
    // fund user to be able to pay retryable fees
    if (nativeToken) {
      await (
        await nativeToken
          .connect(deployerL1Wallet)
          .transfer(userL1Wallet.address, ethers.utils.parseEther('1000'))
      ).wait()
      nativeToken.connect(userL1Wallet)
    }

    // create token to be bridged
    const tokenFactory = await new TestERC20__factory(userL1Wallet).deploy()
    token = await tokenFactory.deployed()
    await (await token.mint()).wait()

    // snapshot state before
    const userTokenBalanceBefore = await token.balanceOf(userL1Wallet.address)

    const gatewayTokenBalanceBefore = await token.balanceOf(
      _l2Network.tokenBridge.l1ERC20Gateway
    )
    const userNativeTokenBalanceBefore = nativeToken
      ? await nativeToken.balanceOf(userL1Wallet.address)
      : await parentProvider.getBalance(userL1Wallet.address)
    const bridgeNativeTokenBalanceBefore = nativeToken
      ? await nativeToken.balanceOf(_l2Network.ethBridge.bridge)
      : await parentProvider.getBalance(_l2Network.ethBridge.bridge)

    // approve token
    const depositAmount = 350
    await (
      await token.approve(_l2Network.tokenBridge.l1ERC20Gateway, depositAmount)
    ).wait()

    // calculate retryable params
    const maxSubmissionCost = nativeToken
      ? BigNumber.from(0)
      : BigNumber.from(584000000000)
    const callhook = '0x'

    const gateway = L1OrbitERC20Gateway__factory.connect(
      _l2Network.tokenBridge.l1ERC20Gateway,
      userL1Wallet
    )
    const outboundCalldata = await gateway.getOutboundCalldata(
      token.address,
      userL1Wallet.address,
      userL2Wallet.address,
      depositAmount,
      callhook
    )

    const l1ToL2MessageGasEstimate = new L1ToL2MessageGasEstimator(
      childProvider
    )
    const retryableParams = await l1ToL2MessageGasEstimate.estimateAll(
      {
        from: userL1Wallet.address,
        to: userL2Wallet.address,
        l2CallValue: BigNumber.from(0),
        excessFeeRefundAddress: userL1Wallet.address,
        callValueRefundAddress: userL1Wallet.address,
        data: outboundCalldata,
      },
      await getBaseFee(parentProvider),
      parentProvider
    )

    const gasLimit = retryableParams.gasLimit.mul(60)
    const maxFeePerGas = retryableParams.maxFeePerGas
    const tokenTotalFeeAmount = gasLimit.mul(maxFeePerGas).mul(2)

    // approve fee amount
    if (nativeToken) {
      await (
        await nativeToken.approve(
          _l2Network.tokenBridge.l1ERC20Gateway,
          tokenTotalFeeAmount
        )
      ).wait()
    }

    // bridge it
    const userEncodedData = nativeToken
      ? defaultAbiCoder.encode(
          ['uint256', 'bytes', 'uint256'],
          [maxSubmissionCost, callhook, tokenTotalFeeAmount]
        )
      : defaultAbiCoder.encode(
          ['uint256', 'bytes'],
          [maxSubmissionCost, callhook]
        )

    const router = nativeToken
      ? L1OrbitGatewayRouter__factory.connect(
          _l2Network.tokenBridge.l1GatewayRouter,
          userL1Wallet
        )
      : L1GatewayRouter__factory.connect(
          _l2Network.tokenBridge.l1GatewayRouter,
          userL1Wallet
        )

    const depositTx = await router.outboundTransferCustomRefund(
      token.address,
      userL1Wallet.address,
      userL2Wallet.address,
      depositAmount,
      gasLimit,
      maxFeePerGas,
      userEncodedData,
      { value: nativeToken ? BigNumber.from(0) : tokenTotalFeeAmount }
    )

    // wait for L2 msg to be executed
    await waitOnL2Msg(depositTx)

    ///// checks

    const l2TokenAddress = await router.calculateL2TokenAddress(token.address)
    l2Token = ERC20__factory.connect(l2TokenAddress, childProvider)
    expect(await l2Token.balanceOf(userL2Wallet.address)).to.be.eq(
      depositAmount
    )

    const userTokenBalanceAfter = await token.balanceOf(userL1Wallet.address)
    expect(userTokenBalanceBefore.sub(userTokenBalanceAfter)).to.be.eq(
      depositAmount
    )

    const gatewayTokenBalanceAfter = await token.balanceOf(
      _l2Network.tokenBridge.l1ERC20Gateway
    )
    expect(gatewayTokenBalanceAfter.sub(gatewayTokenBalanceBefore)).to.be.eq(
      depositAmount
    )

    const userNativeTokenBalanceAfter = nativeToken
      ? await nativeToken.balanceOf(userL1Wallet.address)
      : await parentProvider.getBalance(userL1Wallet.address)
    if (nativeToken) {
      expect(
        userNativeTokenBalanceBefore.sub(userNativeTokenBalanceAfter)
      ).to.be.eq(tokenTotalFeeAmount)
    } else {
      expect(
        userNativeTokenBalanceBefore.sub(userNativeTokenBalanceAfter)
      ).to.be.gte(tokenTotalFeeAmount.toNumber())
    }

    const bridgeNativeTokenBalanceAfter = nativeToken
      ? await nativeToken.balanceOf(_l2Network.ethBridge.bridge)
      : await parentProvider.getBalance(_l2Network.ethBridge.bridge)
    expect(
      bridgeNativeTokenBalanceAfter.sub(bridgeNativeTokenBalanceBefore)
    ).to.be.eq(tokenTotalFeeAmount)
  })

  xit('can withdraw token via default gateway', async function () {
    // fund userL2Wallet so it can pay for L2 withdraw TX
    await depositNativeToL2()

    // snapshot state before
    const userL1TokenBalanceBefore = await token.balanceOf(userL1Wallet.address)
    const userL2TokenBalanceBefore = await l2Token.balanceOf(
      userL2Wallet.address
    )
    const l1GatewayTokenBalanceBefore = await token.balanceOf(
      _l2Network.tokenBridge.l1ERC20Gateway
    )
    const l2TokenSupplyBefore = await l2Token.totalSupply()

    // start withdrawal
    const withdrawalAmount = 250
    const l2Router = L2GatewayRouter__factory.connect(
      _l2Network.tokenBridge.l2GatewayRouter,
      userL2Wallet
    )
    const withdrawTx = await l2Router[
      'outboundTransfer(address,address,uint256,bytes)'
    ](token.address, userL1Wallet.address, withdrawalAmount, '0x')
    const withdrawReceipt = await withdrawTx.wait()
    const l2Receipt = new L2TransactionReceipt(withdrawReceipt)

    // wait until dispute period passes and withdrawal is ready for execution
    await sleep(5 * 1000)

    const messages = await l2Receipt.getL2ToL1Messages(userL1Wallet)
    const l2ToL1Msg = messages[0]
    const timeToWaitMs = 1000
    await l2ToL1Msg.waitUntilReadyToExecute(childProvider, timeToWaitMs)

    // execute on L1
    await (await l2ToL1Msg.execute(childProvider)).wait()

    //// checks
    const userL1TokenBalanceAfter = await token.balanceOf(userL1Wallet.address)
    expect(userL1TokenBalanceAfter.sub(userL1TokenBalanceBefore)).to.be.eq(
      withdrawalAmount
    )

    const userL2TokenBalanceAfter = await l2Token.balanceOf(
      userL2Wallet.address
    )
    expect(userL2TokenBalanceBefore.sub(userL2TokenBalanceAfter)).to.be.eq(
      withdrawalAmount
    )

    const l1GatewayTokenBalanceAfter = await token.balanceOf(
      _l2Network.tokenBridge.l1ERC20Gateway
    )
    expect(
      l1GatewayTokenBalanceBefore.sub(l1GatewayTokenBalanceAfter)
    ).to.be.eq(withdrawalAmount)

    const l2TokenSupplyAfter = await l2Token.totalSupply()
    expect(l2TokenSupplyBefore.sub(l2TokenSupplyAfter)).to.be.eq(
      withdrawalAmount
    )
  })

  it('can deposit token via custom gateway', async function () {
    // fund user to be able to pay retryable fees
    if (nativeToken) {
      await (
        await nativeToken
          .connect(deployerL1Wallet)
          .transfer(userL1Wallet.address, ethers.utils.parseEther('1000'))
      ).wait()
    }

    // create L1 custom token
    const customL1TokenFactory = nativeToken
      ? await new TestOrbitCustomTokenL1__factory(deployerL1Wallet).deploy(
          _l2Network.tokenBridge.l1CustomGateway,
          _l2Network.tokenBridge.l1GatewayRouter
        )
      : await new TestCustomTokenL1__factory(deployerL1Wallet).deploy(
          _l2Network.tokenBridge.l1CustomGateway,
          _l2Network.tokenBridge.l1GatewayRouter
        )
    const customL1Token = await customL1TokenFactory.deployed()
    await (await customL1Token.connect(userL1Wallet).mint()).wait()

    // create L2 custom token
    if (nativeToken) {
      await depositNativeToL2()
    }
    const customL2TokenFactory = await new TestArbCustomToken__factory(
      deployerL2Wallet
    ).deploy(_l2Network.tokenBridge.l2CustomGateway, customL1Token.address)
    const customL2Token = await customL2TokenFactory.deployed()

    // prepare custom gateway registration params
    const router = nativeToken
      ? L1OrbitGatewayRouter__factory.connect(
          _l2Network.tokenBridge.l1GatewayRouter,
          userL1Wallet
        )
      : L1GatewayRouter__factory.connect(
          _l2Network.tokenBridge.l1GatewayRouter,
          userL1Wallet
        )
    const l1ToL2MessageGasEstimate = new L1ToL2MessageGasEstimator(
      childProvider
    )

    const routerData =
      L2GatewayRouter__factory.createInterface().encodeFunctionData(
        'setGateway',
        [[customL1Token.address], [_l2Network.tokenBridge.l2CustomGateway]]
      )
    const routerRetryableParams = await l1ToL2MessageGasEstimate.estimateAll(
      {
        from: _l2Network.tokenBridge.l1GatewayRouter,
        to: _l2Network.tokenBridge.l2GatewayRouter,
        l2CallValue: BigNumber.from(0),
        excessFeeRefundAddress: userL1Wallet.address,
        callValueRefundAddress: userL1Wallet.address,
        data: routerData,
      },
      await getBaseFee(parentProvider),
      parentProvider
    )

    const gatewayData =
      L2CustomGateway__factory.createInterface().encodeFunctionData(
        'registerTokenFromL1',
        [[customL1Token.address], [customL2Token.address]]
      )
    const gwRetryableParams = await l1ToL2MessageGasEstimate.estimateAll(
      {
        from: _l2Network.tokenBridge.l1CustomGateway,
        to: _l2Network.tokenBridge.l2CustomGateway,
        l2CallValue: BigNumber.from(0),
        excessFeeRefundAddress: userL1Wallet.address,
        callValueRefundAddress: userL1Wallet.address,
        data: gatewayData,
      },
      await getBaseFee(parentProvider),
      parentProvider
    )

    // approve fee amount
    const valueForGateway = gwRetryableParams.deposit.mul(BigNumber.from(2))
    const valueForRouter = routerRetryableParams.deposit.mul(BigNumber.from(2))
    if (nativeToken) {
      await (
        await nativeToken.approve(
          customL1Token.address,
          valueForGateway.add(valueForRouter)
        )
      ).wait()
    }

    // do the custom gateway registration
    const receipt = await (
      await customL1Token
        .connect(userL1Wallet)
        .registerTokenOnL2(
          customL2Token.address,
          gwRetryableParams.maxSubmissionCost,
          routerRetryableParams.maxSubmissionCost,
          gwRetryableParams.gasLimit.mul(2),
          routerRetryableParams.gasLimit.mul(2),
          BigNumber.from(100000000),
          valueForGateway,
          valueForRouter,
          userL1Wallet.address,
          {
            value: nativeToken
              ? BigNumber.from(0)
              : valueForGateway.add(valueForRouter),
          }
        )
    ).wait()

    /// wait for execution of both tickets
    const l1TxReceipt = new L1TransactionReceipt(receipt)
    const messages = await l1TxReceipt.getL1ToL2Messages(childProvider)
    const messageResults = await Promise.all(
      messages.map(message => message.waitForStatus())
    )
    if (
      messageResults[0].status !== L1ToL2MessageStatus.REDEEMED ||
      messageResults[1].status !== L1ToL2MessageStatus.REDEEMED
    ) {
      console.log(
        `Retryable ticket (ID ${messages[0].retryableCreationId}) status: ${
          L1ToL2MessageStatus[messageResults[0].status]
        }`
      )
      console.log(
        `Retryable ticket (ID ${messages[1].retryableCreationId}) status: ${
          L1ToL2MessageStatus[messageResults[1].status]
        }`
      )
      exit()
    }

    // snapshot state before
    const userTokenBalanceBefore = await customL1Token.balanceOf(
      userL1Wallet.address
    )
    const gatewayTokenBalanceBefore = await customL1Token.balanceOf(
      _l2Network.tokenBridge.l1CustomGateway
    )
    const userNativeTokenBalanceBefore = nativeToken
      ? await nativeToken.balanceOf(userL1Wallet.address)
      : await parentProvider.getBalance(userL1Wallet.address)
    const bridgeNativeTokenBalanceBefore = nativeToken
      ? await nativeToken.balanceOf(_l2Network.ethBridge.bridge)
      : await parentProvider.getBalance(_l2Network.ethBridge.bridge)

    // approve token
    const depositAmount = 110
    await (
      await customL1Token
        .connect(userL1Wallet)
        .approve(_l2Network.tokenBridge.l1CustomGateway, depositAmount)
    ).wait()

    // calculate retryable params
    const maxSubmissionCost = 0
    const callhook = '0x'

    const gateway = L1OrbitCustomGateway__factory.connect(
      _l2Network.tokenBridge.l1CustomGateway,
      userL1Wallet
    )
    const outboundCalldata = await gateway.getOutboundCalldata(
      customL1Token.address,
      userL1Wallet.address,
      userL2Wallet.address,
      depositAmount,
      callhook
    )

    const retryableParams = await l1ToL2MessageGasEstimate.estimateAll(
      {
        from: userL1Wallet.address,
        to: userL2Wallet.address,
        l2CallValue: BigNumber.from(0),
        excessFeeRefundAddress: userL1Wallet.address,
        callValueRefundAddress: userL1Wallet.address,
        data: outboundCalldata,
      },
      await getBaseFee(parentProvider),
      parentProvider
    )

    const gasLimit = retryableParams.gasLimit.mul(40)
    const maxFeePerGas = retryableParams.maxFeePerGas
    const tokenTotalFeeAmount = gasLimit.mul(maxFeePerGas).mul(2)

    // approve fee amount
    if (nativeToken) {
      await (
        await nativeToken.approve(
          _l2Network.tokenBridge.l1CustomGateway,
          tokenTotalFeeAmount
        )
      ).wait()
    }

    // bridge it
    const userEncodedData = nativeToken
      ? defaultAbiCoder.encode(
          ['uint256', 'bytes', 'uint256'],
          [maxSubmissionCost, callhook, tokenTotalFeeAmount]
        )
      : defaultAbiCoder.encode(
          ['uint256', 'bytes'],
          [BigNumber.from(334400000000), callhook]
        )

    const depositTx = await router.outboundTransferCustomRefund(
      customL1Token.address,
      userL1Wallet.address,
      userL2Wallet.address,
      depositAmount,
      gasLimit,
      maxFeePerGas,
      userEncodedData,
      { value: nativeToken ? BigNumber.from(0) : tokenTotalFeeAmount }
    )

    // wait for L2 msg to be executed
    await waitOnL2Msg(depositTx)

    ///// checks
    expect(await router.getGateway(customL1Token.address)).to.be.eq(
      _l2Network.tokenBridge.l1CustomGateway
    )

    const l2TokenAddress = await router.calculateL2TokenAddress(
      customL1Token.address
    )

    l2Token = ERC20__factory.connect(l2TokenAddress, childProvider)
    expect(await l2Token.balanceOf(userL2Wallet.address)).to.be.eq(
      depositAmount
    )

    const userTokenBalanceAfter = await customL1Token.balanceOf(
      userL1Wallet.address
    )
    expect(userTokenBalanceBefore.sub(userTokenBalanceAfter)).to.be.eq(
      depositAmount
    )

    const gatewayTokenBalanceAfter = await customL1Token.balanceOf(
      _l2Network.tokenBridge.l1CustomGateway
    )
    expect(gatewayTokenBalanceAfter.sub(gatewayTokenBalanceBefore)).to.be.eq(
      depositAmount
    )

    const userNativeTokenBalanceAfter = nativeToken
      ? await nativeToken.balanceOf(userL1Wallet.address)
      : await parentProvider.getBalance(userL1Wallet.address)
    if (nativeToken) {
      expect(
        userNativeTokenBalanceBefore.sub(userNativeTokenBalanceAfter)
      ).to.be.eq(tokenTotalFeeAmount)
    } else {
      expect(
        userNativeTokenBalanceBefore.sub(userNativeTokenBalanceAfter)
      ).to.be.gte(tokenTotalFeeAmount.toNumber())
    }
    const bridgeNativeTokenBalanceAfter = nativeToken
      ? await nativeToken.balanceOf(_l2Network.ethBridge.bridge)
      : await parentProvider.getBalance(_l2Network.ethBridge.bridge)
    expect(
      bridgeNativeTokenBalanceAfter.sub(bridgeNativeTokenBalanceBefore)
    ).to.be.eq(tokenTotalFeeAmount)
  })

  it('can upgrade from bridged USDC to native USDC', async function () {
    /// test applicable only for eth based chains
    if (nativeToken) {
      return
    }

    /// create new L1 usdc gateway behind proxy
    const proxyAdminFac = await new ProxyAdmin__factory(
      deployerL1Wallet
    ).deploy()
    const proxyAdmin = await proxyAdminFac.deployed()
    const l1USDCCustomGatewayFactory = await new L1USDCCustomGateway__factory(
      deployerL1Wallet
    ).deploy()
    const l1USDCCustomGatewayLogic = await l1USDCCustomGatewayFactory.deployed()
    const tupFactory = await new TransparentUpgradeableProxy__factory(
      deployerL1Wallet
    ).deploy(l1USDCCustomGatewayLogic.address, proxyAdmin.address, '0x')
    const tup = await tupFactory.deployed()
    const l1USDCCustomGateway = L1USDCCustomGateway__factory.connect(
      tup.address,
      deployerL1Wallet
    )
    console.log('L1USDCCustomGateway address: ', l1USDCCustomGateway.address)

    /// create new L2 usdc gateway behind proxy
    const proxyAdminL2Fac = await new ProxyAdmin__factory(
      deployerL2Wallet
    ).deploy()
    const proxyAdminL2 = await proxyAdminL2Fac.deployed()
    const l2USDCCustomGatewayFactory = await new L2USDCCustomGateway__factory(
      deployerL2Wallet
    ).deploy()
    const l2USDCCustomGatewayLogic = await l2USDCCustomGatewayFactory.deployed()
    const tupL2Factory = await new TransparentUpgradeableProxy__factory(
      deployerL2Wallet
    ).deploy(l2USDCCustomGatewayLogic.address, proxyAdminL2.address, '0x')
    const tupL2 = await tupL2Factory.deployed()
    const l2USDCCustomGateway = L2USDCCustomGateway__factory.connect(
      tupL2.address,
      deployerL2Wallet
    )
    console.log('L2USDCCustomGateway address: ', l2USDCCustomGateway.address)

    /// create l1 usdc behind proxy
    const l1UsdcFactory = await new MockL1Usdc__factory(
      deployerL1Wallet
    ).deploy()
    const l1UsdcLogic = await l1UsdcFactory.deployed()
    const tupL1UsdcFactory = await new TransparentUpgradeableProxy__factory(
      deployerL1Wallet
    ).deploy(l1UsdcLogic.address, proxyAdmin.address, '0x')
    const tupL1Usdc = await tupL1UsdcFactory.deployed()
    const l1Usdc = MockL1Usdc__factory.connect(
      tupL1Usdc.address,
      deployerL1Wallet
    )
    await (await l1Usdc.initialize()).wait()
    console.log('L1 USDC address: ', l1Usdc.address)

    /// create l2 usdc behind proxy
    const l2UsdcFactory = await new BridgedUsdcCustomToken__factory(
      deployerL2Wallet
    ).deploy()
    const l2UsdcLogic = await l2UsdcFactory.deployed()
    const tupL2UsdcFactory = await new TransparentUpgradeableProxy__factory(
      deployerL2Wallet
    ).deploy(l2UsdcLogic.address, proxyAdminL2.address, '0x')
    const tupL2Usdc = await tupL2UsdcFactory.deployed()
    const l2Usdc = BridgedUsdcCustomToken__factory.connect(
      tupL2Usdc.address,
      deployerL2Wallet
    )
    await (
      await l2Usdc.initialize(
        'Bridged USDC Orbit',
        l2USDCCustomGateway.address,
        l1Usdc.address
      )
    ).wait()
    console.log('L2 USDC address: ', l2Usdc.address)

    /// initialize gateways
    await (
      await l1USDCCustomGateway.initialize(
        l2USDCCustomGateway.address,
        _l2Network.tokenBridge.l1GatewayRouter,
        _l2Network.ethBridge.inbox,
        l1Usdc.address,
        l2Usdc.address,
        deployerL1Wallet.address
      )
    ).wait()
    console.log('L1 USDC custom gateway initialized')

    await (
      await l2USDCCustomGateway.initialize(
        l1USDCCustomGateway.address,
        _l2Network.tokenBridge.l2GatewayRouter,
        l1Usdc.address,
        l2Usdc.address,
        deployerL2Wallet.address
      )
    ).wait()
    console.log('L2 USDC custom gateway initialized')

    /// register USDC custom gateway
    const router = L1GatewayRouter__factory.connect(
      _l2Network.tokenBridge.l1GatewayRouter,
      deployerL1Wallet
    )
    const l2Router = L2GatewayRouter__factory.connect(
      _l2Network.tokenBridge.l2GatewayRouter,
      deployerL2Wallet
    )
    const maxGas = BigNumber.from(500000)
    const gasPriceBid = BigNumber.from(200000000)
    let maxSubmissionCost = BigNumber.from(257600000000)
    const registrationCalldata = router.interface.encodeFunctionData(
      'setGateways',
      [
        [l1Usdc.address],
        [l1USDCCustomGateway.address],
        maxGas,
        gasPriceBid,
        maxSubmissionCost,
      ]
    )
    const rollupOwner = new Wallet(LOCALHOST_L3_OWNER_KEY, parentProvider)
    const upExec = UpgradeExecutor__factory.connect(
      await IOwnable__factory.connect(
        _l2Network.ethBridge.rollup,
        deployerL1Wallet
      ).owner(),
      rollupOwner
    )
    const gwRegistrationTx = await upExec.executeCall(
      router.address,
      registrationCalldata,
      {
        value: maxGas.mul(gasPriceBid).add(maxSubmissionCost),
      }
    )
    await waitOnL2Msg(gwRegistrationTx)
    console.log('USDC custom gateway registered')

    /// check gateway registration
    expect(await router.getGateway(l1Usdc.address)).to.be.eq(
      l1USDCCustomGateway.address
    )
    expect(await l1USDCCustomGateway.depositsPaused()).to.be.eq(false)
    expect(await l2Router.getGateway(l1Usdc.address)).to.be.eq(
      l2USDCCustomGateway.address
    )
    expect(await l2USDCCustomGateway.withdrawalsPaused()).to.be.eq(false)

    /// do a deposit
    const depositAmount = ethers.utils.parseEther('2')
    await (await l1Usdc.transfer(userL1Wallet.address, depositAmount)).wait()
    await (
      await l1Usdc
        .connect(userL1Wallet)
        .approve(l1USDCCustomGateway.address, depositAmount)
    ).wait()
    maxSubmissionCost = BigNumber.from(334400000000)
    const depositTx = await router
      .connect(userL1Wallet)
      .outboundTransferCustomRefund(
        l1Usdc.address,
        userL2Wallet.address,
        userL2Wallet.address,
        depositAmount,
        maxGas,
        gasPriceBid,
        defaultAbiCoder.encode(['uint256', 'bytes'], [maxSubmissionCost, '0x']),
        { value: maxGas.mul(gasPriceBid).add(maxSubmissionCost) }
      )
    await waitOnL2Msg(depositTx)
    expect(await l2Usdc.balanceOf(userL2Wallet.address)).to.be.eq(depositAmount)
    expect(await l1Usdc.balanceOf(l1USDCCustomGateway.address)).to.be.eq(
      depositAmount
    )
    console.log('Deposited USDC')

    /// pause deposits
    await (await l1USDCCustomGateway.pauseDeposits()).wait()
    expect(await l1USDCCustomGateway.depositsPaused()).to.be.eq(true)

    /// pause withdrawals
    await (await l2USDCCustomGateway.pauseWithdrawals()).wait()
    expect(await l2USDCCustomGateway.withdrawalsPaused()).to.be.eq(true)

    /// transfer ownership to circle
    const circleWallet = ethers.Wallet.createRandom().connect(parentProvider)
    await (
      await deployerL1Wallet.sendTransaction({
        to: circleWallet.address,
        value: ethers.utils.parseEther('1'),
      })
    ).wait()

    await (await l1Usdc.setOwner(circleWallet.address)).wait()
    await (await l1USDCCustomGateway.setOwner(circleWallet.address)).wait()
    console.log('L1 USDC and L1 USDC gateway ownership transferred to circle')

    /// circle checks that deposits are paused, all in-flight deposits and withdrawals are processed

    /// add minter rights to usdc gateway so it can burn USDC
    await (
      await l1Usdc.connect(circleWallet).addMinter(l1USDCCustomGateway.address)
    ).wait()
    console.log('Minter rights added to USDC gateway')

    /// burn USDC
    await (
      await l1USDCCustomGateway.connect(circleWallet).burnLockedUSDC()
    ).wait()
    expect(await l1Usdc.balanceOf(l1USDCCustomGateway.address)).to.be.eq(0)
    expect(await l2Usdc.balanceOf(userL2Wallet.address)).to.be.eq(depositAmount)
    console.log('USDC burned')
  })
})

/**
 * helper function to fund user wallet on L2
 */
async function depositNativeToL2() {
  /// deposit tokens
  const amountToDeposit = ethers.utils.parseEther('2.0')
  await (
    await nativeToken!
      .connect(userL1Wallet)
      .approve(_l2Network.ethBridge.inbox, amountToDeposit)
  ).wait()

  const depositFuncSig = {
    name: 'depositERC20',
    type: 'function',
    stateMutability: 'nonpayable',
    inputs: [
      {
        name: 'amount',
        type: 'uint256',
      },
    ],
  }
  const inbox = new ethers.Contract(
    _l2Network.ethBridge.inbox,
    [depositFuncSig],
    userL1Wallet
  )

  const depositTx = await inbox.depositERC20(amountToDeposit)

  // wait for deposit to be processed
  const depositRec = await L1TransactionReceipt.monkeyPatchEthDepositWait(
    depositTx
  ).wait()
  await depositRec.waitForL2(childProvider)
}

async function waitOnL2Msg(tx: ethers.ContractTransaction) {
  const retryableReceipt = await tx.wait()
  const l1TxReceipt = new L1TransactionReceipt(retryableReceipt)
  const messages = await l1TxReceipt.getL1ToL2Messages(childProvider)

  // 1 msg expected
  const messageResult = await messages[0].waitForStatus()
  const status = messageResult.status
  expect(status).to.be.eq(L1ToL2MessageStatus.REDEEMED)
}

const getFeeToken = async (inbox: string, parentProvider: any) => {
  const bridge = await IInbox__factory.connect(inbox, parentProvider).bridge()

  let feeToken = ethers.constants.AddressZero

  try {
    feeToken = await IERC20Bridge__factory.connect(
      bridge,
      parentProvider
    ).nativeToken()
  } catch {}

  return feeToken
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
