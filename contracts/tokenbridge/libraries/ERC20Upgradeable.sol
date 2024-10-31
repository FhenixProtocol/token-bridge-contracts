// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.8.0) (token/ERC20/ERC20.sol)

pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ContextUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import { PermissionedV2, PermissionV2 } from "./fhenix/PermissionedV2.sol";
import { FHE, euint128, inEuint128 } from "@fhenixprotocol/contracts/FHE.sol";

/**
 * @dev Implementation of the {IERC20} interface.
 *
 * This implementation is based on OpenZeppelin implementation (https://github.com/OpenZeppelin/openzeppelin-contracts-upgradeable/blob/v4.8.3/contracts/token/ERC20/ERC20Upgradeable.sol)
 * with a small modification. OZ implementation removed `decimals` storage variable when they did upgrade to Solidity 0.8. Since we're already using older OZ implementation, here we're
 * adding back the `decimals` storage variable along with the appropriate getter and setter. That way we avoid changes in storage layout that would've happen due to OZ removing token
 * decimals storage variable.
 *
 */
contract ERC20Upgradeable is
    Initializable,
    ContextUpgradeable,
    IERC20Upgradeable,
    IERC20MetadataUpgradeable,
    PermissionedV2
{
    mapping(address => uint256) private _balances;
    mapping(address => euint128) internal _encBalances;

    mapping(address => mapping(address => uint256)) private _allowances;
    mapping(address => mapping(address => euint128)) internal _encAllowances;

    uint256 private _totalSupply;
    euint128 internal _encTotalSupply = FHE.asEuint128(0);

    string private _name;
    string private _symbol;
    uint8 private _decimals;

    // TODO: Enable this once a fixed address is established for the PermitV2 contract
    // TODO: Move this from the constructor to an initializer function in PermissionedV2 contract
    constructor() PermissionedV2(address(120), "ERC20") {}

    /**
     * @dev Sets the values for {name} and {symbol}.
     *
     * The default value of {decimals} is 18. To select a different value for
     * {decimals} you should overload it.
     *
     * All two of these values are immutable: they can only be set once during
     * construction.
     */
    function __ERC20_init(string memory name_, string memory symbol_) internal onlyInitializing {
        __ERC20_init_unchained(name_, symbol_);
    }

    function __ERC20_init_unchained(
        string memory name_,
        string memory symbol_
    ) internal onlyInitializing {
        _name = name_;
        _symbol = symbol_;
        _decimals = 18;
    }

    /**
     * @dev Returns the name of the token.
     */
    function name() public view virtual override returns (string memory) {
        return _name;
    }

    /**
     * @dev Returns the symbol of the token, usually a shorter version of the
     * name.
     */
    function symbol() public view virtual override returns (string memory) {
        return _symbol;
    }

    /**
     * @dev Returns the number of decimals used to get its user representation.
     * For example, if `decimals` equals `2`, a balance of `505` tokens should
     * be displayed to a user as `5.05` (`505 / 10 ** 2`).
     *
     * Tokens usually opt for a value of 18, imitating the relationship between
     * Ether and Wei. This is the value {ERC20} uses, unless {_setupDecimals} is
     * called;
     *
     * NOTE: This information is only used for _display_ purposes: it in
     * no way affects any of the arithmetic of the contract, including
     * {IERC20-balanceOf} and {IERC20-transfer}.
     */
    function decimals() public view virtual override returns (uint8) {
        return _decimals;
    }

    /**
     * @dev See {IERC20-totalSupply}.
     */
    function totalSupply() public view virtual override returns (uint256) {
        return _totalSupply;
    }

    /**
     * @dev Returns the encrypted value of tokens in existence.
     *
     * @dev Designed to be used as part of a write tx
     * @dev Can only be called by an authorized router (see PermitV2 documentation)
     */
    function encTotalSupply(
        uint256 _permitId
    ) public virtual override withPermitRouter(_permitId) returns (euint128) {
        return _encTotalSupply;
    }

    /**
     * @dev Returns the encrypted value of tokens in existence, sealed for the caller.
     */
    function sealedTotalSupply(
        PermissionV2 calldata permission
    ) public view virtual override withPermission(permission) returns (string memory) {
        return _encTotalSupply.seal(permission.sealingKey);
    }

    /**
     * @dev See {IERC20-balanceOf}.
     */
    function balanceOf(address account) public view virtual override returns (uint256) {
        return _balances[account];
    }

    /**
     * @dev Returns the value of the encrypted tokens owned by `account`
     * @dev Designed to be used as part of a write tx
     * @dev Can only be called by an authorized router (see PermitV2 documentation)
     */
    function encBalanceOf(
        uint256 permitId
    ) public virtual override withPermitRouter(permitId) returns (euint128) {
        return _encBalances[permitIssuer];
    }

    /**
     * @dev Returns the value of the encrypted tokens owned by the issuer of the PermitNft, sealed for the caller
     */
    function sealedBalanceOf(
        PermissionV2 calldata permission
    ) public view virtual override withPermission(permission) returns (string memory) {
        return _encBalances[permission.issuer].seal(permission.sealingKey);
    }

    /**
     * @dev See {IERC20-transfer}.
     *
     * Requirements:
     *
     * - `to` cannot be the zero address.
     * - the caller must have a balance of at least `amount`.
     */
    function transfer(address to, uint256 amount) public virtual override returns (bool) {
        address owner = _msgSender();
        _transfer(owner, to, amount);
        return true;
    }

    /**
     * @dev Moves a `value` amount of tokens from the caller's account to `to`.
     * Accepts the value as inEuint128, more convenient for calls from EOAs.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     */
    function encTransfer(
        address to,
        inEuint128 calldata ieAmount
    ) public virtual override returns (bool) {
        _encTransfer(msg.sender, to, FHE.asEuint128(ieAmount));
        return true;
    }

    /**
     * @dev See {IERC20-allowance}.
     */
    function allowance(
        address owner,
        address spender
    ) public view virtual override returns (uint256) {
        return _allowances[owner][spender];
    }

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     *
     * @dev Designed to be used as part of a write tx
     * @dev Can only be called by an authorized router (see PermitV2 documentation)
     */
    function encAllowance(
        uint256 permitId,
        address spender
    ) public virtual override withPermitRouter(permitId) returns (euint128) {
        return _encAllowances[permitIssuer][spender];
    }

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default. Sealed for the caller.
     *
     * Permission issuer must be either the owner or spender.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function sealedAllowance(
        PermissionV2 calldata permission,
        address owner,
        address spender
    ) public view virtual override withPermission(permission) returns (string memory) {
        if (permission.issuer != owner && permission.issuer != spender) {
            revert ERC20NotOwnerOrSpender();
        }
        return _encAllowances[owner][spender].seal(permission.sealingKey);
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * NOTE: If `amount` is the maximum `uint256`, the allowance is not updated on
     * `transferFrom`. This is semantically equivalent to an infinite approval.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address spender, uint256 amount) public virtual override returns (bool) {
        address owner = _msgSender();
        _approve(owner, spender, amount);
        return true;
    }

    /**
     * @dev Sets `ieAmount` tokens as the allowance of `spender` over the
     * caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits an {EncApproved} event.
     */
    function encApprove(
        address spender,
        inEuint128 calldata ieAmount
    ) public virtual override returns (bool) {
        _encApprove(msg.sender, spender, FHE.asEuint128(ieAmount));
        return true;
    }

    /**
     * @dev See {IERC20-transferFrom}.
     *
     * Emits an {Approval} event indicating the updated allowance. This is not
     * required by the EIP. See the note at the beginning of {ERC20}.
     *
     * NOTE: Does not update the allowance if the current allowance
     * is the maximum `uint256`.
     *
     * Requirements:
     *
     * - `from` and `to` cannot be the zero address.
     * - `from` must have a balance of at least `amount`.
     * - the caller must have allowance for ``from``'s tokens of at least
     * `amount`.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) public virtual override returns (bool) {
        address spender = _msgSender();
        _spendAllowance(from, spender, amount);
        _transfer(from, to, amount);
        return true;
    }

    /**
     * @dev Moves `ieAmount` tokens from `from` to `to` using the
     * allowance mechanism. `value` is then deducted from the caller's
     * allowance. Accepts the value as inEuint128, more convenient for calls from EOAs.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {EncTransfer} event.
     */
    function encTransferFrom(
        address from,
        address to,
        inEuint128 calldata ieAmount
    ) public virtual override returns (bool) {
        euint128 encSpent = _encSpendAllowance(from, msg.sender, FHE.asEuint128(ieAmount));
        _encTransfer(from, to, encSpent);
        return true;
    }

    /**
     * @dev Atomically increases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function increaseAllowance(address spender, uint256 addedValue) public virtual returns (bool) {
        address owner = _msgSender();
        _approve(owner, spender, allowance(owner, spender) + addedValue);
        return true;
    }

    /**
     * @dev Atomically decreases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     * - `spender` must have allowance for the caller of at least
     * `subtractedValue`.
     */
    function decreaseAllowance(
        address spender,
        uint256 subtractedValue
    ) public virtual returns (bool) {
        address owner = _msgSender();
        uint256 currentAllowance = allowance(owner, spender);
        require(currentAllowance >= subtractedValue, "ERC20: decreased allowance below zero");
        unchecked {
            _approve(owner, spender, currentAllowance - subtractedValue);
        }

        return true;
    }

    /**
     * @dev Moves `amount` of tokens from `from` to `to`.
     *
     * This internal function is equivalent to {transfer}, and can be used to
     * e.g. implement automatic token fees, slashing mechanisms, etc.
     *
     * Emits a {Transfer} event.
     *
     * Requirements:
     *
     * - `from` cannot be the zero address.
     * - `to` cannot be the zero address.
     * - `from` must have a balance of at least `amount`.
     */
    function _transfer(address from, address to, uint256 amount) internal virtual {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");

        _beforeTokenTransfer(from, to, amount);

        uint256 fromBalance = _balances[from];
        require(fromBalance >= amount, "ERC20: transfer amount exceeds balance");
        unchecked {
            _balances[from] = fromBalance - amount;
            // Overflow not possible: the sum of all balances is capped by totalSupply, and the sum is preserved by
            // decrementing then incrementing.
            _balances[to] += amount;
        }

        emit Transfer(from, to, amount);

        _afterTokenTransfer(from, to, amount);
    }

    /**
     * @dev Moves `eAmount` tokens from the caller's account to `to`.
     * Accepts the value as euint128, more convenient for calls from other contracts
     *
     * Returns an `euint128` of the true amount transferred.
     *
     * Emits an {EncTransfer} event.
     */
    function _encTransfer(address from, address to, euint128 eAmount) internal returns (euint128) {
        if (from == address(0)) {
            revert ERC20InvalidSender(address(0));
        }
        if (to == address(0)) {
            revert ERC20InvalidReceiver(address(0));
        }

        // Make sure the sender has enough tokens.
        eAmount = FHE.select(eAmount.lte(_encBalances[from]), eAmount, FHE.asEuint128(0));

        _encBeforeTokenTransfer(from, to, eAmount);

        // Add to the balance of `to` and subtract from the balance of `from`.
        _encBalances[to] = _encBalances[to] + eAmount;
        _encBalances[from] = _encBalances[from] - eAmount;

        emit EncTransfer(from, to);

        _encAfterTokenTransfer(from, to, eAmount);

        return eAmount;
    }

    /** @dev Creates `amount` tokens and assigns them to `account`, increasing
     * the total supply.
     *
     * Emits a {Transfer} event with `from` set to the zero address.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     */
    function _mint(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: mint to the zero address");

        _beforeTokenTransfer(address(0), account, amount);

        _totalSupply += amount;
        unchecked {
            // Overflow not possible: balance + amount is at most totalSupply + amount, which is checked above.
            _balances[account] += amount;
        }
        emit Transfer(address(0), account, amount);

        _afterTokenTransfer(address(0), account, amount);
    }

    /**
     * @dev Creates `eAmount` encrypted tokens and assigns them to `to`.
     * Increases `encTotalSupply` by `eAmount`
     * Accepts the value as euint128, more convenient for calls from other contracts
     *
     * Emits an {EncTransfer} event with `from` set to the zero address.
     */
    function _encMint(address to, euint128 eAmount) internal {
        if (to == address(0)) {
            revert ERC20InvalidReceiver(address(0));
        }

        _encBeforeTokenTransfer(address(0), to, eAmount);

        _encBalances[to] = _encBalances[to] + eAmount;
        _encTotalSupply = _encTotalSupply + eAmount;

        emit EncTransfer(address(0), to);

        _encAfterTokenTransfer(address(0), to, eAmount);
    }

    /**
     * @dev Destroys `amount` tokens from `account`, reducing the
     * total supply.
     *
     * Emits a {Transfer} event with `to` set to the zero address.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     * - `account` must have at least `amount` tokens.
     */
    function _burn(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: burn from the zero address");

        _beforeTokenTransfer(account, address(0), amount);

        uint256 accountBalance = _balances[account];
        require(accountBalance >= amount, "ERC20: burn amount exceeds balance");
        unchecked {
            _balances[account] = accountBalance - amount;
            // Overflow not possible: amount <= accountBalance <= totalSupply.
            _totalSupply -= amount;
        }

        emit Transfer(account, address(0), amount);

        _afterTokenTransfer(account, address(0), amount);
    }

    /**
     * @dev Destroys `eAmount` encrypted tokens from `to`.
     * Decreases `encTotalSupply` by `eAmount`
     * Accepts the value as euint128, more convenient for calls from other contracts
     *
     * Emits an {EncTransfer} event with `to` set to the zero address.
     */
    function _encBurn(address from, euint128 eAmount) internal returns (euint128) {
        if (from == address(0)) {
            revert ERC20InvalidSender(address(0));
        }

        eAmount = FHE.select(_encBalances[msg.sender].gte(eAmount), eAmount, FHE.asEuint128(0));

        _encBeforeTokenTransfer(from, address(0), eAmount);

        _encBalances[from] = _encBalances[from] - eAmount;
        _encTotalSupply = _encTotalSupply - eAmount;

        emit EncTransfer(from, address(0));

        _encAfterTokenTransfer(from, address(0), eAmount);

        return eAmount;
    }

    /**
     * @dev Sets `amount` as the allowance of `spender` over the `owner` s tokens.
     *
     * This internal function is equivalent to `approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an {Approval} event.
     *
     * Requirements:
     *
     * - `owner` cannot be the zero address.
     * - `spender` cannot be the zero address.
     */
    function _approve(address owner, address spender, uint256 amount) internal virtual {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    function _encApprove(address owner, address spender, euint128 eAmount) internal {
        if (owner == address(0)) {
            revert ERC20InvalidApprover(address(0));
        }
        if (spender == address(0)) {
            revert ERC20InvalidSpender(address(0));
        }
        _encAllowances[owner][spender] = eAmount;
    }

    /**
     * @dev Sets {decimals} to a value other than the default one of 18.
     *
     * WARNING: This function should only be called from the constructor. Most
     * applications that interact with token contracts will not expect
     * {decimals} to ever change, and may work incorrectly if it does.
     */
    function _setupDecimals(uint8 decimals_) internal virtual {
        _decimals = decimals_;
    }

    /**
     * @dev Updates `owner` s allowance for `spender` based on spent `amount`.
     *
     * Does not update the allowance amount in case of infinite allowance.
     * Revert if not enough allowance is available.
     *
     * Might emit an {Approval} event.
     */
    function _spendAllowance(address owner, address spender, uint256 amount) internal virtual {
        uint256 currentAllowance = allowance(owner, spender);
        if (currentAllowance != type(uint256).max) {
            require(currentAllowance >= amount, "ERC20: insufficient allowance");
            unchecked {
                _approve(owner, spender, currentAllowance - amount);
            }
        }
    }

    function _encSpendAllowance(
        address owner,
        address spender,
        euint128 eAmount
    ) internal virtual returns (euint128) {
        euint128 eCurrentAllowance = _encAllowances[owner][spender];
        euint128 eSpent = FHE.min(eCurrentAllowance, eAmount);
        _encApprove(owner, spender, (eCurrentAllowance - eSpent));

        return eSpent;
    }

    /**
     * @dev Hook that is called before any transfer of tokens. This includes
     * minting and burning.
     *
     * Calling conditions:
     *
     * - when `from` and `to` are both non-zero, `amount` of ``from``'s tokens
     * will be transferred to `to`.
     * - when `from` is zero, `amount` tokens will be minted for `to`.
     * - when `to` is zero, `amount` of ``from``'s tokens will be burned.
     * - `from` and `to` are never both zero.
     *
     * To learn more about hooks, head to xref:ROOT:extending-contracts.adoc#using-hooks[Using Hooks].
     */
    function _beforeTokenTransfer(address from, address to, uint256 amount) internal virtual {}

    /**
     * @dev Hook that is called before any transfer of encrypted tokens. This includes
     * minting and burning.
     */
    function _encBeforeTokenTransfer(address from, address to, euint128 eAmount) internal virtual {}

    /**
     * @dev Hook that is called after any transfer of tokens. This includes
     * minting and burning.
     *
     * Calling conditions:
     *
     * - when `from` and `to` are both non-zero, `amount` of ``from``'s tokens
     * has been transferred to `to`.
     * - when `from` is zero, `amount` tokens have been minted for `to`.
     * - when `to` is zero, `amount` of ``from``'s tokens have been burned.
     * - `from` and `to` are never both zero.
     *
     * To learn more about hooks, head to xref:ROOT:extending-contracts.adoc#using-hooks[Using Hooks].
     */
    function _afterTokenTransfer(address from, address to, uint256 amount) internal virtual {}

    /**
     * @dev Hook that is called after any transfer of encrypted tokens. This includes
     * minting and burning.
     */
    function _encAfterTokenTransfer(address from, address to, euint128 eAmount) internal virtual {}

    /**
     * @dev This empty reserved space is put in place to allow future versions to add new
     * variables without shifting down storage in the inheritance chain.
     * See https://docs.openzeppelin.com/contracts/4.x/upgradeable#storage_gaps
     */
    uint256[44] private __gap;
}
