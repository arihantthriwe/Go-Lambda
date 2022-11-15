package lambdas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golambda/httpClient"
	"golambda/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	// "github.com/aws/smithy-go/document/json"
)

type Request struct {
	CommsID         string `json:"commsId"`
	ProjectCode     string `json:"projectCode"`
	RequestID       string `json:"requestId"`
	RecipientEmail  string `json:"recipientEmail"`
	CountryCode     string `json:"countryCode"`
	MobileNumber    string `json:"mobileNumber"`
	MessageBody     string `json:"messageBody"`
	TrackerObjectId string `json:"trackerObjectId"`
	RecepientName   string `json:"recepientName"`
	EmailSubject    string `json:"emailSubject"`
}
type Response struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var request Request
	AWS_S3_REGION := os.Getenv("AWS_S3_REGION")
	AWS_S3_BUCKET := os.Getenv("AWS_S3_BUCKET")
	AWS_CA_KEY := os.Getenv("AWS_CA_KEY")
	AWS_CERT_KEY := os.Getenv("AWS_CERT_KEY")
    AWS_INTER_KEY := os.Getenv("AWS_INTER_KEY")

	cfg, errLoadDefaultConfig := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_S3_REGION))
	if errLoadDefaultConfig != nil {
		var oe *smithy.OperationError
		if errors.As(errLoadDefaultConfig, &oe) {
			log.Printf("failed to loadconfig: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while loading config"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errLoadDefaultConfig
	}

	awsS3Client := s3.NewFromConfig(cfg)
	var certificate string

	// certificate file
	certFile, errCreatefile := create(certificate)
	if errCreatefile != nil {
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to create Certfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while creating the file"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errCreatefile
	}
	defer certFile.Close()

	downloader := manager.NewDownloader(awsS3Client)
	numBytes, errDownload := downloader.Download(context.TODO(), certFile,
		&s3.GetObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(AWS_CERT_KEY),
		})
	if errDownload != nil {
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to download certFile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while downloading the certfile"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errDownload
	}

	log.Println(numBytes)

	// inter file
	var inter string
	interFile, errCreatefile := create(inter)
	if errCreatefile != nil {
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to create keyfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while creating the Keyfile"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errCreatefile
	}
	defer interFile.Close()

	download := manager.NewDownloader(awsS3Client)
	numByte, errDownload := download.Download(context.TODO(), interFile,
		&s3.GetObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(AWS_INTER_KEY),
		})
	if errDownload != nil {
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to download interfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while downloading the interfile"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errDownload
	}
	log.Println(numByte)

	//ca certificate

	var ca string
	caFile, errCreatefile := create(ca)
	if errCreatefile != nil {
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to create keyfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while creating the Keyfile"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errCreatefile
	}
	defer caFile.Close()

	downloadroot := manager.NewDownloader(awsS3Client)
	numByteroot, errDownload := downloadroot.Download(context.TODO(), caFile,
		&s3.GetObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(AWS_CA_KEY),
		})
	if errDownload != nil {
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to download rootfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage := fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage := "Something went wrong while downloading the rootfile"
			err := &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
			return err
		}
		return errDownload
	}
	log.Println(numByteroot)

	for _, message := range sqsEvent.Records {
		fmt.Printf("Message Body queue: %s", message.Body)
		err := json.Unmarshal([]byte(message.Body), &request)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println(request.RequestID, "start")
	fmt.Println(request, "request")
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
				"toAddress": "` + request.RecipientEmail + `",
				"fromAddress": "FAB <donotreply@bankfab.com>",
				"emailSubject": "FAB OTP",
				"emailBodyContent": "` + request.MessageBody + `",
				"emailBodyContentType": "text/html"
			}
		}`)
		// apiKey := os.Getenv("FAB_ONE_API_KEY")
		// client := mandrill.ClientWithKey(apiKey)

		// message := &mandrill.Message{}
		// message.AddRecipient(request.RecipientEmail, request.RecepientName, "to")
		// message.FromEmail = os.Getenv("FAB_ONE_FROM_EMAIL")
		// message.FromName = os.Getenv("FAB_ONE_FROM_NAME")
		// message.Subject = request.EmailSubject
		// message.HTML = request.MessageBody

		// _, err := client.MessagesSend(message)

		resp, err := httpClient.NormalClient("POST", "https://kproxy-sit.risk-middleware-dev.mesouth1.bankfab.com/communication/v1/send/email", body, &response, certFile,interFile, caFile)
		if err != nil {
			return err
		}
		log.Println(err)
		log.Println(request.RequestID, resp)
		if resp.StatusCode != 200 {
			body := strings.NewReader(`{
					"status":"FAILED",
					"mailError":"` + err.Error() + `"
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
	} else if request.CommsID == "2" {
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
		resp, err := httpClient.NormalClient("POST", "https://kproxy-sit.risk-middleware-dev.mesouth1.bankfab.com/communication/v1/send/sms", body, &response, certFile, interFile, caFile)
		if err != nil {
			return err
		}
		log.Println(request.RequestID, resp)
		if resp.StatusCode != 200 {
			body := strings.NewReader(`{
				"status":"FAILED"
				"mailError":"` + err.Error() + `"
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

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir("/tmp/"), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
