// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

abstract contract MyErc20Context {
    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }
}

interface IMyErc20Errors {
    // 代币发送者错误。用于转账
    error MyErc20InvalidSender(address sender);
    // 代币接收者错误。用于转账
    error MyErc20InvalidReceiver(address receiver);

    // 代币发送者账户余额不足。用于转账
    error MyErc20InsufficientBalance(
        address sender,
        uint256 balance,
        uint256 value
    );

    // 代币授权操作中，授权者无效。用于授权代币
    error MyErc20InvalidApprover(address approver);
    // 代币授权操作中，接收者无效。用于授权代币
    error MyErc20InvalidSpender(address receiver);

    // 代币被授权人的津贴不够支付。用于转账
    error MyErc20InsufficientAllowance(
        address spender,
        uint256 allowance,
        uint256 value
    );
}

interface IMyErc20 {
    event Transfer(address indexed from, address indexed to, uint256 value);

    event Approve(
        address indexed owner,
        address indexed spender,
        uint256 value
    );
}

contract MyErc20 is MyErc20Context, IMyErc20Errors, IMyErc20 {
    // 存储账户的余额
    mapping(address => uint256) private _balance;
    // 存储代币的总供应量
    uint256 private _totalSupply;

    // 查询账户代币余额
    function balanceOf(address account) public view virtual returns (uint256) {
        return _balance[account];
    }

    // 代币转账
    function transfer(address to, uint256 value) public virtual returns (bool) {
        address from = _msgSender();
        _transfer(from, to, value);
        return true;
    }

    function _transfer(address from, address to, uint256 value) internal {
        if (from == address(0)) {
            revert MyErc20InvalidSender(from);
        }
        if (to == address(0)) {
            revert MyErc20InvalidReceiver(to);
        }
        _update(from, to, value);
    }

    // 包含代币的铸造(mint)和销毁(burn)功能
    function _update(address from, address to, uint256 value) internal {
        if (from == address(0)) {
            // 代币是从零地址"凭空"产生的(即铸造)，增加总供应量
            _totalSupply += value;
        } else {
            uint256 fromBalance = _balance[from];
            if (fromBalance < value) {
                // 代币发送者余额不足
                revert MyErc20InsufficientBalance(from, fromBalance, value);
            }
            unchecked {
                _balance[from] = fromBalance - value;
            }
        }

        if (to == address(0)) {
            unchecked {
                // 代币被发送到零地址(即销毁)，减少总供应量
                _totalSupply -= value;
            }
        } else {
            unchecked {
                _balance[to] += value;
            }
        }

        emit Transfer(from, to, value);
    }

    // 代币授权
    function approve(
        address spender,
        uint256 value
    ) public virtual returns (bool) {
        address owner = _msgSender();
        _approve(owner, spender, value);
        return true;
    }

    // 存储已授权的代币津贴
    mapping(address account => mapping(address spender => uint256 value))
        private _allowances;

    function _approve(address owner, address spender, uint256 value) internal {
        _approve(owner, spender, value, true);
    }

    function _approve(
        address owner,
        address spender,
        uint256 value,
        bool emitEvent
    ) internal {
        if (owner == address(0)) {
            revert MyErc20InvalidApprover(owner);
        }
        if (spender == address(0)) {
            revert MyErc20InvalidSpender(spender);
        }
        _allowances[owner][spender] = value;
        if (emitEvent) {
            emit Approve(owner, spender, value);
        }
    }

    // 代币代扣转账
    function transferFrom(
        address from,
        address to,
        uint256 value
    ) public virtual returns (bool) {
        address spender = _msgSender();
        _spendAllowance(from, spender, value);
        _transferFrom(from, to, value);
        return true;
    }

    function _spendAllowance(
        address owner,
        address spender,
        uint256 value
    ) internal {
        uint256 currentAllowance = _allowances[owner][spender];

        // 如果 currentAllowance == type(uint256).max 为无限授权 不需要浪费Gas走这里的逻辑
        if (currentAllowance < type(uint256).max) {
            if (currentAllowance < value) {
                revert MyErc20InsufficientAllowance(
                    spender,
                    currentAllowance,
                    value
                );
            }
            unchecked {
                _approve(owner, spender, currentAllowance - value, false);
            }
        }
    }

    function _transferFrom(address from, address to, uint256 value) internal {
        _transfer(from, to, value);
    }

    address private _owner;

    modifier onlyOwner() {
        require(_owner == _msgSender(), "Not: Owner");
        _;
    }

    constructor() {
        _owner = msg.sender;
    }

    // 合约所有者增发代币
    function mint(address to, uint256 value) public onlyOwner {
        _mint(to, value);
    }

    function _mint(address to, uint256 value) internal {
        if (to == address(0)) {
            revert MyErc20InvalidReceiver(to);
        }
        _update(address(0), to, value);
    }
}
// 已部署到 Sepolia 合约地址 0x3A9bF08a063270C07fd0dce59ecE6D1f462C939B