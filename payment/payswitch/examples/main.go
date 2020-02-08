package payment

import (
	"golang.org/x/exp/errors/fmt"
	"sundialventure.com/common/payment/payswitch/pay"
)

func main() {
	m1 := MobileMoneyParams{
		PhoneNumber: "0243922636",
		Amount:      1.22,
		RefNum:      "098765432127",
		Network:     pay.MTN,
		Description: "sample test",
	}

	pp := Switch{
		Mode:              ModeTest,
		APIKey:            "YzY3M2EwMWRjNTY5MDlkNmY2ZWQwMjJkNjYyYTIwNjg=",
		APIUser:           "yestech5d2da3cf962b6",
		MerchantID:        "TTM-00000679",
		MobileMoneyParams: &m1,
	}

	if resp, err := pp.DoMobileMoney(); err != nil {
		panic(err)
	} else {
		fmt.Println(resp)
	}
}
