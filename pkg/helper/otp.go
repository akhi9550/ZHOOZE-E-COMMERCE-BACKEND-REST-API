package helper

import (
	"errors"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

var tw *twilio.RestClient

func TwilioSetup(username string, password string) {
	tw = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})
}
func TwilioSendOTP(phone string, serviceID string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo("+91" + phone)
	params.SetChannel("sms")
	res, err := tw.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {
		return "", err
	}
	return *res.Sid, nil
}
func TwilioVerifyOTP(serviceID string, code string, phone string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo("+91" + phone)
	params.SetCode(code)
	res, err := tw.VerifyV2.CreateVerificationCheck(serviceID, params)
	fmt.Println("res status", *res.Status)
	if err != nil {
		return err
	}
	if *res.Status == "approved" {
		return nil
	}
	return errors.New("failed to validate otp")
}
