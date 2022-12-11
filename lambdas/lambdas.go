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

	b64 "encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	// "github.com/aws/smithy-go/document/json"
)

type Request struct {
	CommsID          string `json:"commsId"`
	ProjectCode      string `json:"projectCode"`
	CountryCode      string `json:"countryCode"`
	MobileNumber     string `json:"mobileNumber"`
	RequestID        string `json:"requestId"`
	SmsMessageBody   string `json:"smsMessageBody"`
	TrackerObjectID  string `json:"trackerObjectId"`
	RecipientEmail   string `json:"recipientEmail"`
	RecipientName    string `json:"recipientName"`
	EmailMessageBody string `json:"emailMessageBody"`
	EmailSubject     string `json:"emailSubject"`
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
	// log.Println("starting KMS Encryption&Decryption")
	// KMSEncrypt(request, response)
	// return nil
	// for _, message := range sqsEvent.Records {
	fmt.Printf("Message Body queue: %s", sqsEvent.Records[0].Body)
	errJson := json.Unmarshal([]byte(sqsEvent.Records[0].Body), &request)
	if errJson != nil {
		devMessage := errJson.Error()
		updateTracker(&request, &response, "FAILED", "", devMessage)
		return nil
	}
	// }
	request.EmailMessageBody = strings.Replace(request.EmailMessageBody, "'", `"`, -1)
	request.EmailMessageBody = strings.Replace(request.EmailMessageBody, "^", "\n", -1)
	request.SmsMessageBody = strings.Replace(request.SmsMessageBody, "'", `"`, -1)
	request.SmsMessageBody = strings.Replace(request.SmsMessageBody, "^", "\n", -1)
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
	var a ApplicationArea
	a.CorrelationID = "8134481cfe144e1e959cdaa3a569e120"
	a.InterfaceID = "SMERewards"
	a.CountryOfOrigin = "AE"
	a.SenderID = "RWD"
	a.SenderUserID = "5823XIG"
	a.TransactionID = "FBACC0003651563885814959"
	a.TransactionDateTime = time.Now().String()
	a.TransactionTimeZone = "(GMT+4:00) Asia/Dubai"
	a.Language = "EN"
	a.CreationDateTime = time.Now().Format("2006-01-02T15:04:05Z")
	a.SenderLocation = "UAE"
	if request.CommsID == "1" {
		log.Println("starting fab mail api")
		var req FabCommunicationsMail
		var d DataAreaMail
		req.ApplicationArea = a
		d.ToAddress = request.RecipientEmail
		d.FromAddress = "donotreply@bankfab.com"
		d.FromEntityName = "FAB"
		d.EmailSubject = request.EmailSubject
		d.EmailBodyContent = request.EmailMessageBody
		d.EmailBodyContentType = "text/html"
		req.DataArea = d
		byteBody, err := json.Marshal(req)
		if err != nil {
			updateTracker(&request, &response, "FAILED", "", err.Error())
			return nil
		}

		// fmt.Println("Body:", string(byteBody))
		body := strings.NewReader(string(byteBody))
		// body := strings.NewReader(`{
		// 	"applicationArea": {
		// 	  "correlationId": "8134481cfe144e1e959cdaa3a569e120",
		// 	  "interfaceID": "SMERewards",
		// 	  "countryOfOrigin": "AE",
		// 	  "senderId": "RWD",
		// 	  "senderUserId": "5823XIG",
		// 	  "transactionId": "FBACC0003651563885814959",
		// 	  "transactionDateTime": "2019-07-23T16:42:28Z",
		// 	  "transactionTimeZone": "(GMT+4:00) Asia/Dubai",
		// 	  "language": "EN",
		// 	  "creationDateTime": "2019-07-23T16:42:28Z",
		// 	  "senderLocation": "UAE"
		// 	},
		// 	"dataArea": {
		// 	  "toAddress": "` + request.RecipientEmail + `",
		// 	  "fromAddress": "donotreply@bankfab.com",
		// 	  "fromEntityName": "FAB",
		// 	  "emailSubject": "` + request.EmailSubject + `",
		// 	  "emailBodyContent": "` + request.MessageBody + `",
		// 	  "emailBodyContentType": "text/html"
		// 	}
		//   }`)
		log.Println("Fab client request body email: ", body)
		fabAPIResp, err := httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/email", body, &response, rootCertFile, caChainCertFile, certKeyFile)
		log.Println("success fab email, requestId--> ", request.RequestID, " response--> ", fabAPIResp)
		if err != nil {
			log.Println("fab mail api response", fabAPIResp)
			log.Println("fab mail api err", err)
			updateTracker(&request, &response, "FAILED FAB API(connection error)", fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)), "")
			return nil
		}
		if fabAPIResp.StatusCode != 200 {
			updateTracker(&request, &response, "FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)", err.Error(), "")
			return nil
		}
	} else if request.CommsID == "2" {
		var reqSms FabCommunicationsSMS
		var dSms DataAreaSMS
		reqSms.ApplicationArea = a
		dSms.MobileNumber = request.MobileNumber
		dSms.MessageText = request.SmsMessageBody
		dSms.FromEntityName = "FAB"
		dSms.MessageType = "OTP"
		dSms.OriginatorName = "FAB"
		reqSms.DataArea = dSms
		byteBodySms, errSms := json.Marshal(reqSms)
		if errSms != nil {
			updateTracker(&request, &response, "FAILED", "", errSms.Error())
			return nil
		}

		// fmt.Println("Body:", string(byteBody))
		bodySms := strings.NewReader(string(byteBodySms))
		// body := strings.NewReader(`{
		// 	"applicationArea": {
		// 		"correlationId": "8134481cfe144e1e959cdaa3a569e120",
		// 		"interfaceID": "SMERewards",
		// 		"countryOfOrigin": "AE",
		// 		"senderId": "RWD",
		// 		"senderUserId": "5823XIG",
		// 		"transactionId": "FBACC0003651563885814959",
		// 		"transactionDateTime": "2019-07-23T16:42:28Z",
		// 		"transactionTimeZone": "(GMT+4:00) Asia/Dubai",
		// 		"language": "EN",
		// 		"creationDateTime": "2019-07-23T16:42:28Z",
		// 		"senderLocation": "UAE"
		// 	  },
		// 	"dataArea": {
		// 		"mobileNumber": "` + request.MobileNumber + `",
		// 		"messageText": "` + request.MessageBody + `",
		// 		"fromEntityName": "FAB",
		// 		"messageType": "OTP",
		// 		"originatorName": "FAB"
		// 	}
		// }`)
		log.Println("Fab client request body sms: ", bodySms)
		fabAPIResp, err := httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/sms", bodySms, &response, rootCertFile, caChainCertFile, certKeyFile)
		log.Println("success fab sms, requestId--> ", request.RequestID, " response--> ", fabAPIResp)
		if err != nil {
			log.Println("fab sms api response", fabAPIResp)
			log.Println("fab sms api err", err)
			updateTracker(&request, &response, "FAILED FAB API(connection error)", fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)), "")
			return nil
		}
		if fabAPIResp.StatusCode != 200 {
			updateTracker(&request, &response, "FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)", err.Error(), "")
			return nil
		}
	} else if request.CommsID == "3" {

		log.Println("starting fab mail api")
		var req FabCommunicationsMail
		var d DataAreaMail
		req.ApplicationArea = a
		d.ToAddress = request.RecipientEmail
		d.FromAddress = "donotreply@bankfab.com"
		d.FromEntityName = "FAB"
		d.EmailSubject = request.EmailSubject
		d.EmailBodyContent = request.EmailMessageBody
		d.EmailBodyContentType = "text/html"
		req.DataArea = d
		byteBody, err := json.Marshal(req)
		if err != nil {
			updateTracker(&request, &response, "FAILED", "", err.Error())
		}

		// fmt.Println("Body:", string(byteBody))
		body := strings.NewReader(string(byteBody))
		// body := strings.NewReader(`{
		// 	"applicationArea": {
		// 	  "correlationId": "8134481cfe144e1e959cdaa3a569e120",
		// 	  "interfaceID": "SMERewards",
		// 	  "countryOfOrigin": "AE",
		// 	  "senderId": "RWD",
		// 	  "senderUserId": "5823XIG",
		// 	  "transactionId": "FBACC0003651563885814959",
		// 	  "transactionDateTime": "2019-07-23T16:42:28Z",
		// 	  "transactionTimeZone": "(GMT+4:00) Asia/Dubai",
		// 	  "language": "EN",
		// 	  "creationDateTime": "2019-07-23T16:42:28Z",
		// 	  "senderLocation": "UAE"
		// 	},
		// 	"dataArea": {
		// 	  "toAddress": "` + request.RecipientEmail + `",
		// 	  "fromAddress": "donotreply@bankfab.com",
		// 	  "fromEntityName": "FAB",
		// 	  "emailSubject": "` + request.EmailSubject + `",
		// 	  "emailBodyContent": "` + request.MessageBody + `",
		// 	  "emailBodyContentType": "text/html"
		// 	}
		//   }`)
		log.Println("Fab client request body for mail: ", body)
		fabAPIResp, err := httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/email", body, &response, rootCertFile, caChainCertFile, certKeyFile)
		log.Println("success fab mail, requestId--> ", request.RequestID, " response--> ", fabAPIResp)
		if err != nil {
			log.Println("fab mail api response", fabAPIResp)
			log.Println("fab mail api err", err)
			updateTracker(&request, &response, "FAILED FAB API(connection error)", fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)), "")
		}
		if fabAPIResp.StatusCode != 200 {
			updateTracker(&request, &response, "FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)", err.Error(), "")
		}

		var reqSms FabCommunicationsSMS
		var dSms DataAreaSMS
		reqSms.ApplicationArea = a
		dSms.MobileNumber = request.MobileNumber
		dSms.MessageText = request.SmsMessageBody
		dSms.FromEntityName = "FAB"
		dSms.MessageType = "OTP"
		dSms.OriginatorName = "FAB"
		reqSms.DataArea = dSms
		byteBodySms, errSms := json.Marshal(reqSms)
		if errSms != nil {
			updateTracker(&request, &response, "FAILED", "", errSms.Error())
			return nil
		}

		// fmt.Println("Body:", string(byteBody))
		bodySms := strings.NewReader(string(byteBodySms))
		// body := strings.NewReader(`{
		// 	"applicationArea": {
		// 		"correlationId": "8134481cfe144e1e959cdaa3a569e120",
		// 		"interfaceID": "SMERewards",
		// 		"countryOfOrigin": "AE",
		// 		"senderId": "RWD",
		// 		"senderUserId": "5823XIG",
		// 		"transactionId": "FBACC0003651563885814959",
		// 		"transactionDateTime": "2019-07-23T16:42:28Z",
		// 		"transactionTimeZone": "(GMT+4:00) Asia/Dubai",
		// 		"language": "EN",
		// 		"creationDateTime": "2019-07-23T16:42:28Z",
		// 		"senderLocation": "UAE"
		// 	  },
		// 	"dataArea": {
		// 		"mobileNumber": "` + request.MobileNumber + `",
		// 		"messageText": "` + request.MessageBody + `",
		// 		"fromEntityName": "FAB",
		// 		"messageType": "OTP",
		// 		"originatorName": "FAB"
		// 	}
		// }`)
		log.Println("Fab client request body for sms: ", bodySms)
		fabAPIResp, err = httpClient.NormalClient("POST", "https://services-test.bankfab.com/communication/v1/send/sms", bodySms, &response, rootCertFile, caChainCertFile, certKeyFile)
		log.Println("success fab sms, requestId--> ", request.RequestID, " response--> ", fabAPIResp)
		if err != nil {
			log.Println("fab sms api response", fabAPIResp)
			log.Println("fab sms api err", err)
			updateTrackerSms(&request, &response, "FAILED FAB API(connection error)", fmt.Sprint(strings.Replace(err.Error(), "\"", "'", -1)), "")
			return nil
		}
		if fabAPIResp.StatusCode != 200 {
			updateTrackerSms(&request, &response, "FAILED (FAB API status code:- ` + fmt.Sprint(fabAPIResp.StatusCode) + `)", err.Error(), "")
			return nil
		}

	} else {
		updateTracker(&request, &response, "FAILED", "", `Invalid Comms Id - `+request.CommsID+``)
	}
	updateTracker(&request, &response, "SUCCESS", "", "")
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

	resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectID, body, &response)
	if errParseTracker != nil {
		log.Println("Parse - PUT - requestId--> ", request.RequestID, "Parse - PUT - response--> ", resp, "Parse - PUT - error--> ", errParseTracker)
	}
}
func updateTrackerSms(request *Request, response *Response, status, smsError, lambdaError string) {
	body := strings.NewReader(`{
		"status":"` + status + `"
	}`)
	if smsError != "" && lambdaError == "" {
		body = strings.NewReader(`{
			"status":"` + status + `",
			"smsError": "` + smsError + `"
		}`)
	} else if lambdaError != "" && smsError == "" {
		body = strings.NewReader(`{
			"status":"` + status + `",
			"lambdaError": "` + lambdaError + `"
		}`)
	} else if lambdaError != "" && smsError != "" {
		body = strings.NewReader(`{
			"status":"` + status + `",
			"lambdaError": "` + lambdaError + `",
			"smsError":"` + smsError + `"
		}`)
	}

	resp, errParseTracker := httpClient.ParseClient("PUT", "https://dev-fab-api-gateway.thriwe.com/parse/classes/tracker/"+request.TrackerObjectID, body, &response)
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

// KMSEncryptAPI defines the interface for the Encrypt function.
// We use this interface to test the function using a mocked service.
type KMSEncryptAPI interface {
	Encrypt(ctx context.Context,
		params *kms.EncryptInput,
		optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
}

// EncryptText encrypts some text using an AWS Key Management Service (AWS KMS) key (KMS key).
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, an EncryptOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to Encrypt.
func EncryptText(c context.Context, api KMSEncryptAPI, input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return api.Encrypt(c, input)
}

// KMSDecryptAPI defines the interface for the Decrypt function.
// We use this interface to test the function using a mocked service.
type KMSDecryptAPI interface {
	Decrypt(ctx context.Context,
		params *kms.DecryptInput,
		optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// DecodeData decrypts some text that was encrypted with an AWS Key Management Service (AWS KMS) key (KMS key).
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DecryptOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to Decrypt.
func DecodeData(c context.Context, api KMSDecryptAPI, input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return api.Decrypt(c, input)
}
func KMSEncrypt(request Request, response Response) {
	keyID := "arn:aws:kms:me-south-1:750585140312:key/46653e53-2df4-4a2c-bdab-f5440e062b51"
	text := "Arihant Jain"

	if keyID == "" || text == "" {
		fmt.Println("You must supply the ID of a KMS key and some text")
		fmt.Println("-k KEY-ID -t \"text\"")
	}

	cfg, errLoadDefaultConfig := config.LoadDefaultConfig(context.TODO(), config.WithRegion("me-south-1"))
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
	}

	client := kms.NewFromConfig(cfg)

	input := &kms.EncryptInput{
		KeyId:     &keyID,
		Plaintext: []byte(text),
	}
	log.Println("text to be encrypted-> ", text)
	result, errKMSEncrypt := EncryptText(context.TODO(), client, input)
	if errKMSEncrypt != nil {
		var devMessage string
		var oe *smithy.OperationError
		if errors.As(errKMSEncrypt, &oe) {
			log.Printf("failed to KMS encrypt: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to KMS encrypt: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			devMessage = errKMSEncrypt.Error()
		}
		updateTracker(&request, &response, "FAILED", "", devMessage)
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)

	fmt.Println("encrypted text-> ", blobString)
	KMSDecrypt(blobString, request, response)
}
func KMSDecrypt(data string, request Request, response Response) {

	if data == "" {
		fmt.Println("You must supply the encrypted data as a string")
		fmt.Println("-d DATA")
		return
	}

	cfg, errLoadDefaultConfig := config.LoadDefaultConfig(context.TODO(), config.WithRegion("me-south-1"))
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
	}

	client := kms.NewFromConfig(cfg)
	blob, errDecoding := b64.StdEncoding.DecodeString(data)
	if errDecoding != nil {
		var devMessage string
		var oe *smithy.OperationError
		if errors.As(errDecoding, &oe) {
			log.Printf("failed to b64decode: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to b64decode: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			devMessage = errDecoding.Error()
		}
		updateTracker(&request, &response, "FAILED", "", devMessage)
	}

	input := &kms.DecryptInput{
		CiphertextBlob: blob,
	}
	log.Println("encrypted data before decrypting-> ", data)
	result, errKMSDecrypt := DecodeData(context.TODO(), client, input)
	if errKMSDecrypt != nil {
		var devMessage string
		var oe *smithy.OperationError
		if errors.As(errKMSDecrypt, &oe) {
			log.Printf("failed to KMS decode: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			devMessage = fmt.Sprintf("failed to KMS decode: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			devMessage = errKMSDecrypt.Error()
		}
		updateTracker(&request, &response, "FAILED", "", devMessage)
	}

	fmt.Println("decrypted data-> ", string(result.Plaintext))
}

type FabCommunicationsMail struct {
	ApplicationArea ApplicationArea `json:"applicationArea"`
	DataArea        DataAreaMail    `json:"dataArea"`
}

type FabCommunicationsSMS struct {
	ApplicationArea ApplicationArea `json:"applicationArea"`
	DataArea        DataAreaSMS     `json:"dataArea"`
}
type ApplicationArea struct {
	CorrelationID       string `json:"correlationId"`
	InterfaceID         string `json:"interfaceID"`
	CountryOfOrigin     string `json:"countryOfOrigin"`
	SenderID            string `json:"senderId"`
	SenderUserID        string `json:"senderUserId"`
	TransactionID       string `json:"transactionId"`
	TransactionDateTime string `json:"transactionDateTime"`
	TransactionTimeZone string `json:"transactionTimeZone"`
	Language            string `json:"language"`
	CreationDateTime    string `json:"creationDateTime"`
	SenderLocation      string `json:"senderLocation"`
}

type DataAreaMail struct {
	ToAddress            string      `json:"toAddress"`
	FromAddress          string      `json:"fromAddress"`
	FromEntityName       string      `json:"fromEntityName"`
	EmailSubject         string      `json:"emailSubject"`
	EmailBodyContent     interface{} `json:"emailBodyContent"`
	EmailBodyContentType string      `json:"emailBodyContentType"`
}
type DataAreaSMS struct {
	MobileNumber   string `json:"mobileNumber"`
	MessageText    string `json:"messageText"`
	FromEntityName string `json:"fromEntityName"`
	MessageType    string `json:"messageType"`
	OriginatorName string `json:"originatorName"`
}
