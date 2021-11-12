// retry implements request retry logic for the binance api.
package binancewrapper

import (
	"context"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	log "github.com/sirupsen/logrus"
)

// recvWindow retry schedule in ms which will retry a request 3 times with
// varying recvWindow intervals.
var recvWindowSchedule = []int64{
	5000,
	7000,
	10000,
}

// DoFunc is used to call the binance sdk service's Do method.
type DoFunc func(context.Context, ...futures.RequestOption) (interface{}, error)

// Retry will retry the request according to recvWindowSchedule.
func (b *binanceClient) Retry(ctx context.Context, err error, do DoFunc) (interface{}, error) {
	if retryable(err) {
		return b.retryWithRecvWindow(ctx, do)
	}
	return nil, err
}

// retryWithRecvWindow resyncs the system time with the server time and retries
// the request according to recvWindowSchedule.
func (b *binanceClient) retryWithRecvWindow(ctx context.Context, do DoFunc) (interface{}, error) {
	// Covers the first case of the request timestamp being 1000ms or more ahead
	// of the binance server's time.
	err := b.serverTimeSync(ctx)
	if err != nil {
		return nil, err
	}

	// Covers the second case of the request being outside of the recvWindow.
	for _, w := range recvWindowSchedule {
		res, err := do(ctx, futures.WithRecvWindow(w))
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

// serverTimeSync sets the time offset for each request to the binance server time.
func (b *binanceClient) serverTimeSync(ctx context.Context) error {
	serverTime, err := b.c.NewServerTimeService().Do(ctx)
	if err == nil {
		log.WithField("ServerTime", serverTime).Info("Updated time offset")
	}
	return err
}

// retryable returns whether or not the request should be retried based on the
// api error code. The request will only be retried if the api call failed with
// error codes:
// -1021 INVALID_TIMESTAMP
// 	- Timestamp for this request is outside of the recvWindow.
// 	- Timestamp for this request was 1000ms ahead of the server's time.
// 	See https://github.com/adshao/go-binance/issues/127
//
// -1007 TIMEOUT
// 	- Timeout waiting for response from backend server. Send status unknown;
// 	execution status unknown.
//
// Other cases with different retry strategies and schedules will need to
// eventually be added, such as for when we're rate limited:
// -1015 TOO_MANY_ORDERS
// -1003 TOO_MANY_REQUESTS
//
// See https://github.com/binance/binance-spot-api-docs/blob/master/errors.md
func retryable(err error) bool {
	if common.IsAPIError(err) {
		// TODO: go-binance sdk doesn't have api error types
		switch apiErr := errors.NewAPIError(err); apiErr.Code {
		case -1021: // INVALID_TIMESTAMP
			log.WithField("Code", apiErr.Code).
				Error("Binance API error: server time out of sync or recvWindow too small")
			return true
		case -1007: // TIMEOUT
			log.WithField("Code", apiErr.Code).
				Error("Binance API error: timeout waiting for response from backend server")
			return true
		default:
			log.WithField("Code", apiErr.Code).Error("Binance API error")
		}
	}
	return false
}
