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
	RecipientName   string `json:"recipientName"`
	EmailSubject    string `json:"emailSubject"`
}
type Response struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var request Request
	err := new(utils.ErrorHandler)
	err = nil
	AWS_S3_REGION := os.Getenv("AWS_S3_REGION")
	AWS_S3_BUCKET := os.Getenv("AWS_S3_BUCKET")
	AWS_ROOT_CERT_S3_KEY := os.Getenv("AWS_ROOT_CERT_S3_KEY")
	AWS_CA_CHAIN_CERT_S3_KEY := os.Getenv("AWS_CA_CHAIN_CERT_S3_KEY")
	AWS_ROOT_CERT_KEY_S3_KEY := os.Getenv("AWS_ROOT_CERT_KEY_S3_KEY")
	var response Response
	for _, message := range sqsEvent.Records {
		fmt.Printf("Message Body queue: %s", message.Body)
		errJson := json.Unmarshal([]byte(message.Body), &request)
		if errJson != nil {
			devMessage := errJson.Error()
			err = &utils.ErrorHandler{DevMessage: devMessage}
		}
	}

	log.Println("starting, requestId--> ", request.RequestID)
	fmt.Println("request--> ", request)

	body := strings.NewReader(`{
		"status":"CONNECTING/DOWNLOADING FILES FROM S3"
	}`)
	resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
	if errParseTracker != nil {
		log.Println("requestId--> ", request.RequestID)
		log.Println("responseError--> ", resp)
		log.Println("errorParse--> ", errParseTracker)
		return nil
		// return errParseTracker
	}
	log.Println(request.RequestID, resp)
	cfg, errLoadDefaultConfig := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_S3_REGION))
	if errLoadDefaultConfig != nil {
		var devMessage string
		var clientMessage string
		var oe *smithy.OperationError
		if errors.As(errLoadDefaultConfig, &oe) {
			log.Printf("failed to loadconfig: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage = "Something went wrong while loading config"
		} else {
			devMessage = errLoadDefaultConfig.Error()
		}
		err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
	}
	if err != nil {
		body := strings.NewReader(`{
			"status":"FAILED",
			"lambdaError":"` + err.DevMessage + `"
		}`)
		resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
		if errParseTracker != nil {
			log.Println("requestId--> ", request.RequestID)
			log.Println("responseError--> ", resp)
			log.Println("errorParse--> ", errParseTracker)
			return nil

			// return errParseTracker
		}
		log.Println("requestId--> ", request.RequestID)
		log.Println("errorLambda--> ", err)
		return nil

	}
	// , func(o *s3.Options) {
	// 	o.Region = "ap-south-1"
	// }
	awsS3Client := s3.NewFromConfig(cfg)
	//ca certificate

	rootCertFileName := "rootCert"
	rootCertFile, errCreatefile := create(rootCertFileName)
	if errCreatefile != nil {
		var devMessage string
		var clientMessage string
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to create keyfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage = "Something went wrong while creating the Keyfile"
		} else {
			devMessage = errCreatefile.Error()

		}
		err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
	}
	if err != nil {
		body := strings.NewReader(`{
			"status":"FAILED",
			"lambdaError":"` + err.DevMessage + `"
		}`)
		resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
		if errParseTracker != nil {
			log.Println("requestId--> ", request.RequestID)
			log.Println("responseError--> ", resp)
			log.Println("errorParse--> ", errParseTracker)
			return nil

			// return errParseTracker
		}
		log.Println("requestId--> ", request.RequestID)
		log.Println("errorLambda--> ", err)
		return nil

	}
	defer rootCertFile.Close()

	// ca chain .pem
	caChainCertFileName := "caChainCert"
	caChainCertFile, errcaChainCreatefile := create(caChainCertFileName)
	if errcaChainCreatefile != nil {
		var devMessage string
		var clientMessage string
		var oe *smithy.OperationError
		if errors.As(errcaChainCreatefile, &oe) {
			log.Printf("failed to create keyfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage = "Something went wrong while creating the Keyfile"
		} else {
			devMessage = errcaChainCreatefile.Error()

		}
		err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
	}
	if err != nil {
		body := strings.NewReader(`{
			"status":"FAILED",
			"lambdaError":"` + err.DevMessage + `"
		}`)
		resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
		if errParseTracker != nil {
			log.Println("requestId--> ", request.RequestID)
			log.Println("responseError--> ", resp)
			log.Println("errorParse--> ", errParseTracker)
			return nil

			// return errParseTracker
		}
		log.Println("requestId--> ", request.RequestID)
		log.Println("errorLambda--> ", err)
		return nil

	}
	defer caChainCertFile.Close()
	// key .pem
	certKeyFileName := "certKey"
	certKeyFile, errCertKeyCreatefile := create(certKeyFileName)
	if errCertKeyCreatefile != nil {
		var devMessage string
		var clientMessage string
		var oe *smithy.OperationError
		if errors.As(errcaChainCreatefile, &oe) {
			log.Printf("failed to create keyfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage = "Something went wrong while creating the Keyfile"
		} else {
			devMessage = errCertKeyCreatefile.Error()

		}
		err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
	}
	if err != nil {
		body := strings.NewReader(`{
			"status":"FAILED",
			"lambdaError":"` + err.DevMessage + `"
		}`)
		resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
		if errParseTracker != nil {
			log.Println("requestId--> ", request.RequestID)
			log.Println("responseError--> ", resp)
			log.Println("errorParse--> ", errParseTracker)
			return nil

			// return errParseTracker
		}
		log.Println("requestId--> ", request.RequestID)
		log.Println("errorLambda--> ", err)
		return nil

	}
	defer certKeyFile.Close()
	// ----
	downloadRootCert := manager.NewDownloader(awsS3Client)
	numByteRootCert, errDownload := downloadRootCert.Download(context.TODO(), rootCertFile,
		&s3.GetObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(AWS_ROOT_CERT_S3_KEY),
		})
	if errDownload != nil {
		var devMessage string
		var clientMessage string
		var oe *smithy.OperationError
		if errors.As(errCreatefile, &oe) {
			log.Printf("failed to download rootfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			clientMessage = "Something went wrong while downloading the rootfile"
		} else {
			devMessage = errDownload.Error()

		}
		err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
		if err != nil {
			body := strings.NewReader(`{
				"status":"FAILED",
				"lambdaError":"` + err.DevMessage + `"
			}`)
			resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if errParseTracker != nil {
				log.Println("requestId--> ", request.RequestID)
				log.Println("responseError--> ", resp)
				log.Println("errorParse--> ", errParseTracker)
				return nil

				// return errParseTracker
			}
			log.Println("requestId--> ", request.RequestID)
			log.Println("errorLambda--> ", err)
			return nil

		}
	}

	log.Println(numByteRootCert)

	if err != nil {
		body := strings.NewReader(`{
			"status":"FAILED",
			"lambdaError":"` + err.DevMessage + `"
		}`)
		resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
		if errParseTracker != nil {
			log.Println("requestId--> ", request.RequestID)
			log.Println("responseError--> ", resp)
			log.Println("errorParse--> ", errParseTracker)
			return nil

			// return errParseTracker
		}
		log.Println("requestId--> ", request.RequestID)
		log.Println("errorLambda--> ", err)
		return nil

	}
	numByteCAChainCert, errDownload := downloadRootCert.Download(context.TODO(), caChainCertFile,
	&s3.GetObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(AWS_CA_CHAIN_CERT_S3_KEY),
	})
if errDownload != nil {
	var devMessage string
	var clientMessage string
	var oe *smithy.OperationError
	if errors.As(errCreatefile, &oe) {
		log.Printf("failed to download rootfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		clientMessage = "Something went wrong while downloading the rootfile"
	} else {
		devMessage = errDownload.Error()

	}
	err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
	if err != nil {
		body := strings.NewReader(`{
			"status":"FAILED",
			"lambdaError":"` + err.DevMessage + `"
		}`)
		resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
		if errParseTracker != nil {
			log.Println("requestId--> ", request.RequestID)
			log.Println("responseError--> ", resp)
			log.Println("errorParse--> ", errParseTracker)
			return nil

			// return errParseTracker
		}
		log.Println("requestId--> ", request.RequestID)
		log.Println("errorLambda--> ", err)
		return nil

	}
}

log.Println(numByteCAChainCert)
// --- key .pem
numByteCertKey, errDownload := downloadRootCert.Download(context.TODO(), certKeyFile,
&s3.GetObjectInput{
	Bucket: aws.String(AWS_S3_BUCKET),
	Key:    aws.String(AWS_ROOT_CERT_KEY_S3_KEY),
})
if errDownload != nil {
var devMessage string
var clientMessage string
var oe *smithy.OperationError
if errors.As(errCreatefile, &oe) {
	log.Printf("failed to download rootfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
	devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
	clientMessage = "Something went wrong while downloading the rootfile"
} else {
	devMessage = errDownload.Error()

}
err = &utils.ErrorHandler{DevMessage: devMessage, Message: clientMessage}
if err != nil {
	body := strings.NewReader(`{
		"status":"FAILED",
		"lambdaError":"` + err.DevMessage + `"
	}`)
	resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
	if errParseTracker != nil {
		log.Println("requestId--> ", request.RequestID)
		log.Println("responseError--> ", resp)
		log.Println("errorParse--> ", errParseTracker)
		return nil

		// return errParseTracker
	}
	log.Println("requestId--> ", request.RequestID)
	log.Println("errorLambda--> ", err)
	return nil

}
}

log.Println(numByteCertKey)
// ---- 
	processingBody := strings.NewReader(`{
		"status":"PROCESSING"
	}`)
	processingResp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, processingBody, &response)
	if errParseTracker != nil {
		log.Println("requestId--> ", request.RequestID)
		log.Println("responseError--> ", processingResp)
		log.Println("errorParse--> ", errParseTracker)
		return nil

		// return errParseTracker
	}
	log.Println("requestId--> ", request.RequestID, "response--> ", processingResp)
	check, errProfanityCheck := httpClient.ProfanityCheck(request.MessageBody, request.TrackerObjectId, request.ProjectCode)
	if !check || errProfanityCheck != nil{
		return nil
	}
	if request.CommsID == "1" {
		log.Println("starting fab sms/mail api")
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

		fabAPIResp, err := httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/email", body, &response, rootCertFile, caChainCertFile, certKeyFile)
		if err != nil {
			log.Println("fab sms/mail api response", fabAPIResp)
			log.Println("fab sms/mail api err", err)
			body := strings.NewReader(`{
				"status":"FAILED FAB API(connection error)",
				"mailError":"` + fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)) + `"
			}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if err != nil {
				log.Println("requestId--> ", request.RequestID)
				log.Println("response--> ", resp)
				log.Println("errorParse--> ", err)
				return nil

				// return err
			}
			log.Println("requestId--> ", request.RequestID)
			log.Println("response--> ", fabAPIResp)
			log.Println("errorLambda--> ", err)
			return nil

			// return err
		}
		log.Println(err)
		log.Println(request.RequestID, resp)
		if fabAPIResp.StatusCode != 200 {
			body := strings.NewReader(`{
					"status":"FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)",
					"mailError":"` + err.Error() + `"
				}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if err != nil {
				log.Println("requestId--> ", request.RequestID)
				log.Println("response--> ", resp)
				log.Println("errorParse--> ", err)
				return nil

				// return err
			}
			log.Println("requestId--> ", request.RequestID)
			log.Println("response--> ", resp)
			log.Println("errorLambda--> ", err)
			return nil

		} else {
			body := strings.NewReader(`{
					"status":"SUCCESS"
				}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if err != nil {
				log.Println("requestId--> ", request.RequestID)
				log.Println("response--> ", resp)
				log.Println("errorParse--> ", err)
				return nil

				// return err
			}

			log.Println("success, requestId--> ", request.RequestID, " response--> ", resp)
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
		resp, err := httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/sms", body, &response, rootCertFile, caChainCertFile, certKeyFile)
		if err != nil {
			body := strings.NewReader(`{
				"status":"FAILED",
				"mailError":"` + err.Error() + `"
			}`)
			resp, err := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
			if err != nil {
				return err
			}
			log.Println(request.RequestID, resp)
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
	return os.Create("/tmp/" + p)
}
