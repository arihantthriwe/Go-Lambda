package lambdas

import (
	"context"
	"encoding/json"
	"fmt"
	"golambda/httpClient"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	// "github.com/aws/smithy-go/document/json"
)

type Request struct {
	CommsID         string    `json:"commsId"`
	ProjectCode     string `json:"projectCode"`
	RequestID       string `json:"requestId"`
	Recipient       string `json:"recipient"`
	CountryCode     string `json:"countryCode"`
	MobileNumber    string `json:"mobileNumber"`
	MessageBody     string `json:"messageBody"`
	TrackerObjectId string `json:"trackerObjectId"`
}
type Response struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var request Request

	for _, message := range sqsEvent.Records {
		fmt.Printf("Message Body queue: %s", message.Body)
		err := json.Unmarshal([]byte(message.Body), &request)
		if err != nil {
			log.Fatalln(err)
		} 
	}
	log.Println(request.RequestID, "start")
	fmt.Println(request, "request")
	// log.Println(ctx, "context")
	//log.Fatalln("exiting")
	var response Response
	body := strings.NewReader(`{
		"status":"PROCESSING"
	}`)
	resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
	if err != nil {
		return err
	}
	log.Println(request.RequestID, resp)

	if request.CommsID == "1" {
		body := strings.NewReader(`{
		"applicationArea": {
			"correlationId": "FT87745646i465",
			"interfaceID": "ESB",
			"countryOfOrigin": "AE",
			"senderId": "GCN",
			"senderUserId": "5823XIG",
			"senderAuthorizationID": "KgL9STHWNPtrhPfXnbX5DEUf5j6lIfiI",
			"senderReferenceID": "IPI1234567890-123",
			"transactionId": "FT87741234543",
			"transactionDateTime": "` + fmt.Sprint(time.Now()) + `",
			"transactionTimeZone": "(GMT+4:00) Asia/Dubai",
			"language": "EN",
			"creationDateTime":"` + fmt.Sprint(time.Now()) + `"
		},
		"dataArea": {
			"toAddress": "` + request.Recipient + `",
			"fromAddress": "FAB <donotreply@bankfab.com>",
			"emailSubject": "FAB OTP",
			"emailBodyContent": "` + request.MessageBody + `",
			"emailBodyContentType": "text/html"
		}
	}`)
			resp, err := httpClient.NormalClient("POST", "https://kproxy-sit.risk-middleware-dev.mesouth1.bankfab.com/communication/v1/send/email", body, &response)
			if err != nil {
				return err
			}
			log.Println(request.RequestID, resp)
			if resp.StatusCode != 200 {
				body := strings.NewReader(`{
					"status":"FAILED"
				}`)
				resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
				if err != nil {
					return err
				}
				log.Println(request.RequestID, resp)
			} else {
				body := strings.NewReader(`{
					"status":"SUCCESS"
				}`)
				resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
				if err != nil {
					return err
				}
				log.Println(request.RequestID, resp)
			}
	}else if request.CommsID == "2" {
		body := strings.NewReader(`{
			"applicationArea": {
				"correlationId": "FT87745646i465",
				"interfaceID": "ESB",
				"countryOfOrigin": "AE",
				"senderId": "GCN",
				"senderUserId": "5823XIG",
				"senderAuthorizationID": "KgL9STHWNPtrhPfXnbX5DEUf5j6lIfiI",
				"senderReferenceID": "IPI1234567890-123",
				"transactionId": "FT87741234543",
				"transactionDateTime": "2020-02-19T10:06:26.026Z",
				"transactionTimeZone": "(GMT+4:00) Asia/Dubai",
				"language": "EN",
				"creationDateTime": "2020-06-20T09:12:28Z"
			},
			"dataArea": {
				"mobileNumber": "` + request.MobileNumber + `",
				"messageText": "` + request.MessageBody + `",
				"messageType": "OTP_EVENT",
				"originatorName": "FAB"
			}
		}`)
		resp, err := httpClient.NormalClient("POST", "https://kproxy-sit.risk-middleware-dev.mesouth1.bankfab.com/communication/v1/send/sms", body, &response)
		if err != nil {
			return err
		}
		log.Println(request.RequestID, resp)
		if resp.StatusCode != 200 {
			body := strings.NewReader(`{
				"status":"FAILED"
			}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if err != nil {
				return err
			}
			log.Println(request.RequestID, resp)
		} else {
			body := strings.NewReader(`{
				"status":"SUCCESS"
			}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if err != nil {
				return err
			}
			log.Println(request.RequestID, resp)

		}
	}
	response.Code = "200"
	response.Error = "no error"
	log.Println(request.RequestID, response)

	return nil
}
