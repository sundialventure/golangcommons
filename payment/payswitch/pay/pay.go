package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty"
	"github.com/spf13/cast"
)

//Modes ...
const (
	ModeTest = "test"
	ModeLive = "live"
)

//Networs ...
const (
	MTN        = "MTN"
	Vodafone   = "VDF"
	AirtelTigo = "ATL"
)

//Switch ...
type Switch struct {
	Mode              string
	liveEndpoint      string
	testEndpoint      string
	MobileMoneyParams *MobileMoneyParams
	MerchantID        string
	APIKey            string
	APIUser           string
	Debug             bool
}

//MobileMoneyParams ....
type MobileMoneyParams struct {
	PhoneNumber     string
	Amount          float32
	RefNum          string
	Network         string
	VodaVoucherCode string
	Description     string
}

//MomoResponse ...
type MomoResponse struct {
	TransactionID string      `json:"transaction_id"`
	Status        string      `json:"status"`
	Code          interface{} `json:"code"`
	Reason        string      `json:"reason"`
}

//DoMobileMoney ...
func (p *Switch) DoMobileMoney() (MomoResponse, error) {
	var momoResponse MomoResponse

	if err := p._validateMomoRequest(); err != nil {
		return momoResponse, err
	}

	client := resty.New()
	client.SetDebug(p.Debug)

	client.SetBasicAuth(p.APIUser, p.APIKey)

	var err error
	var response *resty.Response

	if strings.ToLower(p.Mode) == ModeTest {
		response, err = client.R().SetHeader("Content-Type", "application/json").SetBody(map[string]interface{}{
			"amount":            p.formatAmount(p.MobileMoneyParams.Amount),
			"processing_code":   "000200",
			"transaction_id":    p.MobileMoneyParams.RefNum,
			"desc":              p.MobileMoneyParams.Description,
			"merchant_id":       p.MerchantID,
			"subscriber_number": p.MobileMoneyParams.PhoneNumber,
			"r-switch":          p.MobileMoneyParams.Network,
			"voucher_code":      p.MobileMoneyParams.VodaVoucherCode,
		}).Post("https://test.theteller.net/v1.1/transaction/process")
	} else {
		response, err = client.R().SetHeader("Content-Type", "application/json").SetBody(map[string]interface{}{
			"amount":            p.formatAmount(p.MobileMoneyParams.Amount),
			"processing_code":   "000200",
			"transaction_id":    p.MobileMoneyParams.RefNum,
			"desc":              p.MobileMoneyParams.Description,
			"merchant_id":       p.MerchantID,
			"subscriber_number": p.MobileMoneyParams.PhoneNumber,
			"r-switch":          p.MobileMoneyParams.Network,
			"voucher_code":      p.MobileMoneyParams.VodaVoucherCode,
		}).Post("https://prod.theteller.net/v1.1/transaction/process")
	}

	if err != nil {
		return momoResponse, err
	}

	if err := json.Unmarshal(response.Body(), &momoResponse); err != nil {
		return momoResponse, err
	}

	return momoResponse, nil
}

func (p *Switch) _validateMomoRequest() error {
	if p.Mode == "" {
		return errors.New("No Mode selected")
	}

	if p.APIKey == "" {
		return errors.New("No API Key set")
	}

	if p.APIUser == "" {
		return errors.New("No API User Set")
	}

	if p.MerchantID == "" {
		return errors.New("No Merchanht ID set")
	}

	if len(p.MobileMoneyParams.RefNum) != 12 {
		return errors.New("Reference number should always be 12")
	}

	networks := []string{MTN, Vodafone, AirtelTigo}

	if !arrays.StringExistInSlice(p.MobileMoneyParams.Network, networks...) {
		return errors.New("Invalid network selected")
	}

	//lets validate the phone number...
	//phone number, can either be 12 or 10..
	if len(p.MobileMoneyParams.PhoneNumber) != 10 && len(p.MobileMoneyParams.PhoneNumber) != 12 {
		return errors.New("Invalid phone number, Phone number should be length of 10 or 12")
	}

	if p.MobileMoneyParams.Amount < 0.01 {
		return errors.New("Amount should be greater than 0.1")
	}

	return nil
}

func (p *Switch) formatAmount(amt float32) string {
	stringAmt := cast.ToString(amt)
	amts := strings.Split(stringAmt, ".")
	var formattedAmount string
	if len(amts) == 1 {
		formattedAmount = fmt.Sprintf("%010v", amt) + "00"
	} else if len(amts) == 2 {
		formattedAmount = arrays.LeftPad2Len(amts[0], "0", 10) + arrays.RightPad2Len(amts[1], "0", 2)
	}

	if len(formattedAmount) != 12 {
		panic("amount length should be not more than 12")
	}
	return formattedAmount
}
