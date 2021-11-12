# golang-binance-service

# Setup, Build, and Run

You might need to update `go.mod`:
```
go mod tidy
go mod vendor
```

Create an `.env` file with:
```
PORT=4200
USE_TESTNET=true
DEBUG=true
```

Start the server on port 4200 with:
```
make build run PORT=4200
```
# Testing

Go to https://testnet.binancefuture.com/ and create an account. After creating an account, you can make an API key and API secret
below the candlestick chart. You will use these in postman and unit tests.

## Postman
The postman collection is under `postman/collections` and the postman environment is under `postman/environments`. There are example requests for each endpoint.

I suggest having the https://testnet.binancefuture.com/ open in one window to close any positions you open or cancel any orders made while testing the endpoints.

The set up is basically the same as https://github.com/binance/binance-api-postman

## Unit tests

In order to run tests, create an `.env.test` file with:
```
FUTURES_API_KEY="XXXXXXX"
FUTURES_API_SECRET="XXXXXX"
```
in the `./libs/binancewrapper`, `./libs/test`, and `./controllers/v1/user` directories.

To run all tests:
```
godotenv -f .env.test go test -v ./...
```

To run a specific test not cached:
```
godotenv -f .env.test go test -count=1 -run TestCreateLimitOrder -v ./libs/binancewrapper
```

# Viewing Go Doc of code
```
go get -v  golang.org/x/tools/cmd/godoc
godoc -http:6060
```
Then go to http://localhost:6060/pkg/github.com/bosdhill/golang-binance-service/ (its kind of fucked up rn)

# Endpoints

## `GET` `/v1/user/balance`

Returns the user's perpetual futures `usdtBalance`.

Example request body:
```
{
    "api_key": "{{binance-api-key}}",
    "api_secret": "{{binance-api-secret}}"
}
```

Example response body:
```
{
    "accountAlias": "sRmYFzoCuXFz",
    "asset": "USDT",
    "balance": "98076.98393216",
    "crossWalletBalance": "98076.98393216",
    "crossUnPnl": "0.00000000",
    "availableBalance": "98076.98393216",
    "maxWithdrawAmount": "98076.98393216"
}
```

## `GET` `/v1/user/account`

Returns the user's perpetual futures account info.

Example request body:
```
{
    "api_key": "{{binance-api-key}}",
    "api_secret": "{{binance-api-secret}}"
}
```

Example response body:
```
{
    "assets": [
        {
            "asset": "BNB",
            "initialMargin": "0.00000000",
            "maintMargin": "0.00000000",
            "marginBalance": "0.00000000",
            "maxWithdrawAmount": "0.00000000",
            "openOrderInitialMargin": "0.00000000",
            "positionInitialMargin": "0.00000000",
            "unrealizedProfit": "0.00000000",
            "walletBalance": "0.00000000"
        },
        {
            "asset": "USDT",
            "initialMargin": "0.00000000",
            "maintMargin": "0.00000000",
            "marginBalance": "98076.98393216",
            "maxWithdrawAmount": "98076.98393216",
            "openOrderInitialMargin": "0.00000000",
            "positionInitialMargin": "0.00000000",
            "unrealizedProfit": "0.00000000",
            "walletBalance": "98076.98393216"
        },
        {
            "asset": "BUSD",
            "initialMargin": "0.00000000",
            "maintMargin": "0.00000000",
            "marginBalance": "0.00000000",
            "maxWithdrawAmount": "0.00000000",
            "openOrderInitialMargin": "0.00000000",
            "positionInitialMargin": "0.00000000",
            "unrealizedProfit": "0.00000000",
            "walletBalance": "0.00000000"
        }
    ],
    "canDeposit": true,
    "canTrade": true,
    "canWithdraw": true,
    "feeTier": 0,
    "maxWithdrawAmount": "98076.98393216",
    "positions": [
        {
            "isolated": false,
            "leverage": "20",
            "initialMargin": "0",
            "maintMargin": "0",
            "openOrderInitialMargin": "0",
            "positionInitialMargin": "0",
            "symbol": "RAYUSDT",
            "unrealizedProfit": "0.00000000",
            "entryPrice": "0.0",
            "maxNotional": "25000",
            "positionSide": "BOTH",
            "positionAmt": "0.0",
            "notional": "0",
            "isolatedWallet": "0",
            "updateTime": 0
        },
        ...
        {
            "isolated": false,
            "leverage": "20",
            "initialMargin": "0",
            "maintMargin": "0",
            "openOrderInitialMargin": "0",
            "positionInitialMargin": "0",
            "symbol": "CTSIUSDT",
            "unrealizedProfit": "0.00000000",
            "entryPrice": "0.0",
            "maxNotional": "25000",
            "positionSide": "BOTH",
            "positionAmt": "0",
            "notional": "0",
            "isolatedWallet": "0",
            "updateTime": 0
        }
    ],
    "totalInitialMargin": "0.00000000",
    "totalMaintMargin": "0.00000000",
    "totalMarginBalance": "98076.98393216",
    "totalOpenOrderInitialMargin": "0.00000000",
    "totalPositionInitialMargin": "0.00000000",
    "totalUnrealizedProfit": "0.00000000",
    "totalWalletBalance": "98076.98393216",
    "updateTime": 0
}
```

## `POST` `/v1/user/order`

Creates either a `LIMIT`, `MARKET`, or `STOP_LOSS` order, depending on the order type provided.

Example request body:
```
{
    "user": {
        "api_key": "{{binance-api-key}}",
        "api_secret": "{{binance-api-secret}}"
    },
    "order": {
        "type": "MARKET",
        "symbol": "BTCUSDT",
        "side": "BUY",
        "percentage": 0.01
    }
}
```

Example response body:

```
{
    "symbol": "BTCUSDT",
    "orderId": 2869718120,
    "clientOrderId": "G9Wqjy1RisSjYLDhR4rzYi",
    "price": "0",
    "origQty": "0.152",
    "executedQty": "0",
    "cumQuote": "0",
    "reduceOnly": false,
    "status": "NEW",
    "stopPrice": "0",
    "timeInForce": "GTC",
    "type": "MARKET",
    "side": "BUY",
    "updateTime": 1636705431064,
    "workingType": "CONTRACT_PRICE",
    "activatePrice": "",
    "priceRate": "",
    "avgPrice": "0.00000",
    "positionSide": "BOTH",
    "closePosition": false,
    "priceProtect": false
}
```

Examples for `LIMIT` and `STOP_MARKET` are in the postman collection.

## Issue with Buy limit and Take Profit
If order is not filled, take profit might be triggered immediately.
Fill or kill. 
If order not filled, don't execute

## Issue with Binance server time synchronization 

Sometimes the system time can fall out of sync with the binance server time, for example:
```
--- FAIL: TestChangeSymbolLeverage (36.89s)
    leverage_test.go:68: <APIError> code=-1021, msg=Timestamp for this request is outside of the recvWindow.
``` 

There is a fix implemented for this in the go-binance sdk: https://github.com/adshao/go-binance/issues/127

This would require querying the binance server time once when creating a client and using as a time offset it in each request. 

Since it will eventually get out of sync again, we will retry a failed request (one with error code -1021) with the updated 
time offset after quering the binance server time. 
