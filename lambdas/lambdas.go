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
	CommsID  int `json:"commsId"`
	UserData struct {
		ObjectID           string    `json:"objectId"`
		FirstName          string    `json:"firstName"`
		LastName           string    `json:"lastName"`
		Username           string    `json:"username"`
		MobileNumber       string    `json:"mobileNumber"`
		CountryCode        string    `json:"countryCode"`
		MembershipID       []string  `json:"membershipId"`
		IsEmailVerified    string    `json:"isEmailVerified"`
		IsMobileVerfied    string    `json:"isMobileVerfied"`
		TermsAndCondition1 bool      `json:"termsAndCondition1"`
		TermsAndCondition2 bool      `json:"termsAndCondition2"`
		Handicap           int       `json:"handicap"`
		CreatedAt          time.Time `json:"createdAt"`
		UpdatedAt          time.Time `json:"updatedAt"`
		Addresses          []struct {
			NickName string `json:"nickName"`
			Address  string `json:"address"`
			Emirates string `json:"emirates"`
			Region   string `json:"region"`
			ZipCode  string `json:"zipCode"`
		} `json:"addresses"`
		EmailAlerts  bool `json:"emailAlerts"`
		MobileAlerts bool `json:"mobileAlerts"`
		UserSavings  []struct {
			MembershipID      string `json:"membershipId"`
			BenefitGroupWorth int    `json:"benefitGroupWorth"`
			Amount            int    `json:"amount"`
		} `json:"userSavings"`
		DataID            string `json:"dataId"`
		FromWhere         int    `json:"fromWhere"`
		Otp               string `json:"otp"`
		Recipient         string `json:"recipient"`
		Type              int    `json:"type"`
		WalkthroughStatus bool   `json:"walkthroughStatus"`
		FavouriteConfigs  []struct {
			MembershipID string   `json:"membershipId"`
			Configs      []string `json:"configs"`
		} `json:"favouriteConfigs"`
	} `json:"userData"`
	ProjectCode  string `json:"projectCode"`
	TemplateCode string `json:"templateCode"`
	RequestID    string `json:"requestId"`
	BookingData  struct {
		ObjectID             string `json:"objectId"`
		CourierPickupAddress struct {
			FullAddress struct {
				NickName string `json:"nickName"`
				Address  string `json:"address"`
				Emirates string `json:"emirates"`
				Region   string `json:"region"`
				ZipCode  string `json:"zipCode"`
			} `json:"fullAddress"`
			Emirates struct {
				Location string `json:"location"`
			} `json:"emirates"`
			Region struct {
				Location string `json:"location"`
			} `json:"region"`
			FullName string `json:"fullName"`
			Mobile   string `json:"mobile"`
			Landmark string `json:"landmark"`
			Landline string `json:"landline"`
		} `json:"courierPickupAddress"`
		CourierDropAddress struct {
			FullAddress struct {
				NickName string `json:"nickName"`
				Address  string `json:"address"`
				Emirates string `json:"emirates"`
				Region   string `json:"region"`
				ZipCode  string `json:"zipCode"`
			} `json:"fullAddress"`
			Emirates struct {
				Location string `json:"location"`
			} `json:"emirates"`
			Region struct {
				Location string `json:"location"`
			} `json:"region"`
			FullName string `json:"fullName"`
			Mobile   string `json:"mobile"`
			Landmark string `json:"landmark"`
			Landline string `json:"landline"`
		} `json:"courierDropAddress"`
		Date                     string `json:"date"`
		Time                     string `json:"time"`
		BenefitGroupItemObjectID string `json:"benefitGroupItemObjectId"`
		BookingStatus            string `json:"bookingStatus"`
		Service                  struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Category    struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Logo string `json:"logo"`
			} `json:"category"`
			Partner struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Logo string `json:"logo"`
			} `json:"partner"`
			Facility struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Logo string `json:"logo"`
			} `json:"facility"`
			CostPrice        string `json:"costPrice"`
			SellingPrice     string `json:"sellingPrice"`
			TaxType          string `json:"taxType"`
			TaxPercentage    string `json:"taxPercentage"`
			Currency         string `json:"currency"`
			LongDescription  string `json:"longDescription"`
			ShortDescription string `json:"shortDescription"`
			Images           struct {
				CoverImageMobile   string `json:"coverImageMobile"`
				CoverImageWeb      string `json:"coverImageWeb"`
				DisplayImageMobile string `json:"displayImageMobile"`
				DisplayImageWeb    string `json:"displayImageWeb"`
			} `json:"images"`
			ThirdPartyWebsite        string `json:"thirdPartyWebsite"`
			Faq                      string `json:"faq"`
			DefaultRedemptionType    int    `json:"defaultRedemptionType"`
			DefaultRedemptionProcess string `json:"defaultRedemptionProcess"`
			DefaultBenefitOffered    string `json:"defaultBenefitOffered"`
			DefaultTermsAndCondition string `json:"defaultTermsAndCondition"`
			Tat                      struct {
				Minimum struct {
					Value int `json:"value"`
					Unit  int `json:"unit"`
				} `json:"minimum"`
				Maximum struct {
					Value int `json:"value"`
					Unit  int `json:"unit"`
				} `json:"maximum"`
				Cancel struct {
					Value int `json:"value"`
					Unit  int `json:"unit"`
				} `json:"cancel"`
			} `json:"tat"`
			Mql struct {
				IsActive            bool `json:"isActive"`
				MqlAlertLevelFirst  int  `json:"mqlAlertLevelFirst"`
				MqlAlertLevelSecond int  `json:"mqlAlertLevelSecond"`
				MqlAlertLevelThird  int  `json:"mqlAlertLevelThird"`
			} `json:"mql"`
			DurationUnit  int `json:"durationUnit"`
			DurationValue int `json:"durationValue"`
			IsActive      int `json:"isActive"`
		} `json:"service"`
		Status    int       `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"bookingData"`
	BookingTypeID int `json:"bookingTypeId"`
}
type Response struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var request Request
	
	for _, message := range sqsEvent.Records {
		// fmt.Printf("Message Body queue: %s", message.Body)
		json.Unmarshal([]byte(message.Body), &request)
	}
	log.Println(request.RequestID,"start")
	fmt.Println(request, "request")
	// log.Println(ctx, "context")
	//log.Fatalln("exitting")
	var response Response
	body := strings.NewReader(`{
		"status":"progessing"
	}`)
	resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.RequestID, body, &response)
	if err != nil {
		return err
	}
	log.Println(request.RequestID, resp)

	if request.CommsID == 1 {
		if request.BookingTypeID == 1 && request.TemplateCode == "" {

		} else if request.BookingTypeID == 2 && request.TemplateCode == "" {
			log.Println(request.RequestID, request.BookingTypeID)

		} else if request.BookingTypeID == 3 && request.TemplateCode == "" {
			log.Println(request.RequestID, request.BookingTypeID)

		} else if request.BookingTypeID == 4 && request.TemplateCode == "" {
			log.Println(request.RequestID, request.BookingTypeID)

		} else if request.BookingTypeID == 5 && request.TemplateCode == "" {

			log.Println(request.RequestID, request.BookingTypeID)
		} else if request.BookingTypeID == 6 && request.TemplateCode == "" {
			log.Println(request.RequestID, request.BookingTypeID)

		} else if request.BookingTypeID == 7 && request.TemplateCode == "" {
			log.Println(request.RequestID, request.BookingTypeID)

		} else if request.BookingTypeID == 8 && request.TemplateCode == "" {
			log.Println(request.RequestID, request.BookingTypeID)

		} else {
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
			"toAddress": "` + request.UserData.Username + `",
			"fromAddress": "FAB <donotreply@bankfab.com>",
			"emailSubject": "FAB OTP",
			"emailBodyContent": "` + request.TemplateCode + `",
			"emailBodyContentType": "text/html"
		}
	}`)
			resp, err := httpClient.NormalClient("POST", "https://kproxy-sit.risk-middleware-dev.mesouth1.bankfab.com/communication/v1/send/email", body, &response)
			if err != nil {
				return err
			}
			log.Println(request.RequestID, resp)
			if resp.StatusCode != 200 && resp.StatusCode != 201 {
				body := strings.NewReader(`{
					"status":"failure"
				}`)
				resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracking", body, &response)
				if err != nil {
					return err
				}
				log.Println(request.RequestID, resp)
			} else {
				body := strings.NewReader(`{
					"status":"success"
				}`)
				resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracking", body, &response)
				if err != nil {
					return err
				}
				log.Println(request.RequestID, resp)
			}
		}
	}
	if request.CommsID == 2 && request.TemplateCode == "sms" {
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
				"mobileNumber": "` + request.UserData.MobileNumber + `",
				"messageText": "` + request.TemplateCode + `",
				"messageType": "OTP_EVENT",
				"originatorName": "FAB"
			}
		}`)
		resp, err := httpClient.NormalClient("POST", "https://kproxy-sit.risk-middleware-dev.mesouth1.bankfab.com/communication/v1/send/sms", body, &response)
		if err != nil {
			return err
		}
		log.Println(request.RequestID, resp)
		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			body := strings.NewReader(`{
				"status":"failure"
			}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracking", body, &response)
			if err != nil {
				return err
			}
			log.Println(request.RequestID, resp)
		} else {
			body := strings.NewReader(`{
				"status":"success"
			}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracking", body, &response)
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
