package lambdas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golambda/httpClient"
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
			updateTracker(&request, &response, "FAILED", "", devMessage)
			return nil
		}
	}

	log.Println("starting, requestId--> ", request.RequestID)
	log.Println("request--> ", request)

	updateTracker(&request, &response, "CONNECTING/DOWNLOADING FILES FROM S3", "", "")

	cfg, errLoadDefaultConfig := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWS_S3_REGION))
	if errLoadDefaultConfig != nil {
		var devMessage string
		var oe *smithy.OperationError
		if errors.As(errLoadDefaultConfig, &oe) {
			log.Printf("failed to load config: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			devMessage = errLoadDefaultConfig.Error()
		}
		updateTracker(&request, &response, "FAILED", "", devMessage)
		return nil
	}
	awsS3Client := s3.NewFromConfig(cfg)
	//ca certificate

	rootCertFileName := "rootCert"
	rootCertFile, errCreateRootCertFile := create(rootCertFileName)
	if errCreateRootCertFile != nil {
		updateTracker(&request, &response, "FAILED", "", errCreateRootCertFile.Error())
		return nil
	}

	defer rootCertFile.Close()

	// ca chain .pem
	caChainCertFileName := "caChainCert"
	caChainCertFile, errCreateCaChainCertFile := create(caChainCertFileName)
	if errCreateCaChainCertFile != nil {
		updateTracker(&request, &response, "FAILED", "", errCreateCaChainCertFile.Error())
		return nil
	}

	defer caChainCertFile.Close()

	// key .pem
	certKeyFileName := "certKey"
	certKeyFile, errCreateCertKeyFile := create(certKeyFileName)
	if errCreateCertKeyFile != nil {
		updateTracker(&request, &response, "FAILED", "", errCreateCertKeyFile.Error())
		return nil
	}

	defer certKeyFile.Close()

	// Downloading certificates and keys from s3
	downloadManagerAWS := manager.NewDownloader(awsS3Client)

	// root certificate
	errDownloadRootCert := downloadFileFromS3(downloadManagerAWS, rootCertFile, AWS_S3_BUCKET, AWS_ROOT_CERT_S3_KEY)
	if errDownloadRootCert != "" {
		updateTracker(&request, &response, "FAILED", "", errDownloadRootCert)
		return nil
	}


	// ca-chain certificate
	errDownloadCAChainCert := downloadFileFromS3(downloadManagerAWS, caChainCertFile, AWS_S3_BUCKET, AWS_CA_CHAIN_CERT_S3_KEY)
	if errDownloadCAChainCert != "" {
		updateTracker(&request, &response, "FAILED", "", errDownloadCAChainCert)
		return nil
	}


	// --- key .pem
	errDownloadCertKey := downloadFileFromS3(downloadManagerAWS, certKeyFile, AWS_S3_BUCKET, AWS_ROOT_CERT_KEY_S3_KEY)
	if errDownloadCertKey != "" {
		updateTracker(&request, &response, "FAILED", "", errDownloadCertKey)
		return nil
	}


	updateTracker(&request, &response, "PROCESSING", "", "")

	// --- Profanity Check ---
	// check, errProfanityCheck := httpClient.ProfanityCheck(request.MessageBody, request.TrackerObjectId, request.ProjectCode)
	// if !check || errProfanityCheck != nil{
	// 	profanityCheckFailedBody := strings.NewReader(`{
	// 		"status": "FAILED FAB API(profanity-check failed)"
	// 	}`)
	// 	trackerResp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, profanityCheckFailedBody, &response)
	// 	if errParseTracker != nil {
	// 		return nil
	// 	}
	// 	log.Println("requestId--> ", request.RequestID, "response--> ", trackerResp)
	// 	return nil
	// }

	if request.CommsID == "1" {
		log.Println("starting fab mail api")
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
			log.Println("fab mail api response", fabAPIResp)
			log.Println("fab mail api err", err)
			updateTracker(&request, &response, "FAILED FAB API(connection error)", fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)), "")
			return nil
		}
		if fabAPIResp.StatusCode != 200 {
			updateTracker(&request, &response, "FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)", err.Error(), "")
			return nil
		} else {
			updateTracker(&request, &response, "SUCCESS", "", "")
			log.Println("success, requestId--> ", request.RequestID, " response--> ", fabAPIResp)
			return nil
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
		fabAPIResp, err := httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/sms", body, &response, rootCertFile, caChainCertFile, certKeyFile)
		if err != nil {
			log.Println("fab sms api response", fabAPIResp)
			log.Println("fab sms api err", err)
			updateTracker(&request, &response, "FAILED FAB API(connection error)", fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)), "")
			return nil
		}
		if fabAPIResp.StatusCode != 200 {
			updateTracker(&request, &response, "FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)", err.Error(), "")
			return nil
		} else {
			updateTracker(&request, &response, "SUCCESS", "", "")
			log.Println("success, requestId--> ", request.RequestID, " response--> ", fabAPIResp)
			return nil
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
func updateTracker(request *Request, response *Response, status, mailError, lambdaError string) {
	body := strings.NewReader(`{
		"status":"` + status + `"
	}`)
	if mailError != "" && lambdaError == "" {
		body = strings.NewReader(`{
			"status":"` + status + `",
			"mailError": "` + mailError + `"
		}`)
	} else if lambdaError != "" && mailError == "" {
		body = strings.NewReader(`{
			"status":"` + status + `",
			"lambdaError": "` + lambdaError + `"
		}`)
	} else if lambdaError != "" && mailError != "" {
		body = strings.NewReader(`{
			"status":"` + status + `",
			"lambdaError": "` + lambdaError + `",
			"mailError":"` + mailError + `"
		}`)
	}

	resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectId, body, &response)
	if errParseTracker != nil {
		log.Println("Parse - PUT - requestId--> ", request.RequestID, "Parse - PUT - response--> ", resp, "Parse - PUT - error--> ", errParseTracker)
	}
}
func downloadFileFromS3(downloadManagerAWS *manager.Downloader, file *os.File, bucket, key string) (err string) {
	_, errDownload := downloadManagerAWS.Download(context.TODO(), file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if errDownload != nil {
			var devMessage string
			var oe *smithy.OperationError
			if errors.As(errDownload, &oe) {
				log.Printf("failed to download rootfile: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
				devMessage = fmt.Sprintf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			} else {
				devMessage = errDownload.Error()
			}
			return devMessage
		}
	return ""
}