# golang-binance-service

# Testing

To run tests:
```
godotenv -f .env.test go test -v ./...
```

# Endpoints

`GET` `/v1/user/account` 

Returns the user's perpetual futures account info.

`GET` `/v1/user/balance`

Returns the user's perpetual futures `usdtBalance`.

`POST` `/v1/user/order` 

Creates either a `LIMIT`, `MARKET`, or `STOP_LOSS` order, depending on the order
type provided. 

Returns a response with the following format:
```
{
    "clientOrderId": "testOrder",
    "cumQty": "0",
    "cumQuote": "0",
    "executedQty": "0",
    "orderId": 22542179,
    "avgPrice": "0.00000",
    "origQty": "10",
    "price": "0",
    "reduceOnly": false,
    "side": "BUY",
    "positionSide": "SHORT",
    "status": "NEW",
    "stopPrice": "9300",        // please ignore when order type is TRAILING_STOP_MARKET
    "closePosition": false,   // if Close-All
    "symbol": "BTCUSDT",
    "timeInForce": "GTC",
    "type": "TRAILING_STOP_MARKET",
    "origType": "TRAILING_STOP_MARKET",
    "activatePrice": "9020",    // activation price, only return with TRAILING_STOP_MARKET order
    "priceRate": "0.3",         // callback rate, only return with TRAILING_STOP_MARKET order
    "updateTime": 1566818724722,
    "workingType": "CONTRACT_PRICE",
    "priceProtect": false            // if conditional order trigger is protected   
}
```

Creates a new futures order for the user in Hedge Mode, from the trading signal format:
```
SYMBOLUSDTPERP 
SIDE at ORDER_TYPE
STOP LOSS: XXXX 
TAKE PROFIT: [XXXX] 
SIZE: XXX 
```

* `SYMBOL` is the crypto symbol 

* `SIZE` is the percentage of the user's usdt account balance with maximum leverage. `SIZE` will be used to calculate the usdt `quantity`:
    ```
    quantity = SIZE * usdtBalance * leverageMultiplier 
    ```
    Assuming that `leverageMultiplier` is 10.

* `SIDE` is either `LONG`(buy) or `SHORT`(sell)

* `ORDER_TYPE` is either `MARKET` or `LIMIT`. This will open a `SIDE` position of size `quantity`. 

    * NOTE: there will be a 30 minute time limit for the user to accept a trade for a `MARKET` order. 

* `STOP LOSS` is the decimal representing the limit at which to completely close the position (size `quantity`)

* `TAKE PROFIT` is the list of decimals representing different take profit limit levels (`TP_limits`) sorted in ascending order which are used to partially close the position. The position size of each TP limit level is:

    ```
    quantity / len(TP_limits)
    ```
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
