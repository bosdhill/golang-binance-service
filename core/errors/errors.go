package errors

import (
	err "errors"

	"github.com/adshao/go-binance/v2/common"
)

func NewAPIError(e error) *common.APIError {
	if common.IsAPIError(e) {
		apiArror, _ := e.(*common.APIError)
		return apiArror
	}
	return nil
}

func NewNoUSDTBalance() error {
	return err.New("no USDT balance")
}
