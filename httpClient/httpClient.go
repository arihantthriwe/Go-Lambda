package httpClient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
)

func ParseClient(method, url string, payload *strings.Reader, v interface{}) (*http.Response, error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Minute,
		}).Dial,
		TLSHandshakeTimeout: 1 * time.Minute,
		MaxIdleConns:        10,
		IdleConnTimeout:     1 * time.Minute,
	}
	http_client := &http.Client{
		Timeout:   time.Minute * 1,
		Transport: netTransport,
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Parse-Master-Key", "DEV_MASTER_KEY")
	req.Header.Add("X-Parse-Application-Id", "DEV_APPLICATION_ID")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Server", "1")

	resp, err := http_client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if resp.StatusCode == 400 {
			type temp struct {
				Code  int    `json:"code"`
				Error string `json:"error"`
			}
			var a temp
			_ = json.NewDecoder(resp.Body).Decode(&a)
			fmt.Printf("response error %v", a)
			errR := errors.New(fmt.Sprint(a))
			return nil, errR
		}

		err = fmt.Errorf(fmt.Sprintf("response error from parse client - %v", resp))
		return nil, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func NormalClient(method, url string, payload *strings.Reader, v interface{}, rootCertFile *os.File, caChainCertFile *os.File, certKey *os.File) (*http.Response, error) {
	// tl, err := tls.LoadX509KeyPair(cert.Name(), key.Name())
	// if err != nil {
	// 	return nil, err
	// }
	// Load CA cert
	caCert, err := os.ReadFile(caChainCertFile.Name())
	if err != nil {
		log.Println("caCert os.ReadFile error--> ", err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	rootCert, err := CustomLoadX509KeyPair(rootCertFile.Name(), certKey.Name())
	if err != nil {
		log.Println("rootCert LoadX509KeyPair error--> ", err)
		return nil, err
	}
	fmt.Println(rootCert)
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{rootCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:  1 * time.Minute,
			Deadline: time.Now().Add(1 * time.Minute),
		}).Dial,
		TLSHandshakeTimeout:   1 * time.Minute,
		ResponseHeaderTimeout: 1 * time.Minute,
		MaxIdleConns:          10,
		IdleConnTimeout:       1 * time.Minute,
		TLSClientConfig:       tlsConfig,
	}
	http_client := &http.Client{
		Timeout:   time.Minute * 1,
		Transport: netTransport,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		log.Println("NewRequest req--> ", req, " -error--> ", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	// resp, err := http_client.Post(url, "application/json", payload)
	resp, err := http_client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("client: deadline exceeded")
		} else {
			log.Printf("client: request error: %v", err)
		}
		log.Println("http_client.Do req, ", req, " -resp--> ", resp, " -error--> ", err)
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if resp.StatusCode == 400 {
			var a interface{}
			_ = json.NewDecoder(resp.Body).Decode(&a)
			fmt.Printf("response error %v", a)
			errR := errors.New(fmt.Sprint(a))
			return nil, errR
		}
		if resp.StatusCode == 404 {
			b, _ := io.ReadAll(resp.Body)
			if len(b) > 0 {
				log.Println("404 response body ---> ", string(b))
			}
		}

		err = fmt.Errorf(fmt.Sprintf("response error from fab client - %v", resp))
		return nil, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
func NewNormalClient(method, url string, payload *strings.Reader, v interface{}, rootCertFile *os.File, caChainCertFile *os.File, certKey *os.File) (*http.Response, error) {
	// tl, err := tls.LoadX509KeyPair(cert.Name(), key.Name())
	// if err != nil {
	// 	return nil, err
	// }
	// Load CA cert
	caCert, err := os.ReadFile(caChainCertFile.Name())
	if err != nil {
		log.Println("caCert os.ReadFile error--> ", err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	rootCert, err := CustomLoadX509KeyPair(rootCertFile.Name(), certKey.Name())
	if err != nil {
		log.Println("rootCert LoadX509KeyPair error--> ", err)
		return nil, err
	}
	fmt.Println(rootCert)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{rootCert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:  1 * time.Minute,
			Deadline: time.Now().Add(1 * time.Minute),
		}).Dial,
		TLSHandshakeTimeout:   1 * time.Minute,
		ResponseHeaderTimeout: 1 * time.Minute,
		MaxIdleConns:          10,
		IdleConnTimeout:       1 * time.Minute,
		TLSClientConfig:       tlsConfig,
	}
	http_client := &http.Client{
		Timeout:   time.Minute * 1,
		Transport: netTransport,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		log.Println("NewRequest req--> ", req, " -error--> ", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http_client.Post(url, "application/json", payload)
	// resp, err := http_client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("client: deadline exceeded")
		} else {
			log.Printf("client: request error: %v", err)
		}
		log.Println("http_client.Do req, ", req, " -resp--> ", resp, " -error--> ", err)
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if resp.StatusCode == 400 {
			var a interface{}
			_ = json.NewDecoder(resp.Body).Decode(&a)
			fmt.Printf("response error %v", a)
			errR := errors.New(fmt.Sprint(a))
			return nil, errR
		}

		err = fmt.Errorf(fmt.Sprintf("response error from parse client - %v", resp))
		return nil, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
func ProfanityCheck(messageBody, requestId, projectCode string) (bool, error) {
	body := TrackerRequestType{MessageBody: messageBody, RequestId: requestId, Status: "INITIATED", CommsId: "1", ProjectCode: projectCode}
	var responseTracker TrackerResponseType
	errTrackerRequest := requests.
		URL("https://dzh3x5bk4x5wy3v5mdzrznamiy0ywojf.lambda-url.ap-south-1.on.aws/").
		Header("X-Parse-Application-Id", "DEV_APPLICATION_ID").
		Header("X-Parse-Master-Key", "DEV_MASTER_KEY").
		Header("X-Auth-Server", "1").
		BodyJSON(&body).
		ToJSON(&responseTracker).
		Fetch(context.Background())
	if errTrackerRequest != nil {
		log.Println(errTrackerRequest)
		return false, fmt.Errorf("got an error creating initializer tracker for ThriweComms request")
	}
	if responseTracker.StatusCode != 200 {
		log.Println(responseTracker.Status)
		return false, fmt.Errorf(responseTracker.Status)
	}
	return true, nil
}

type TrackerRequestType struct {
	MessageBody string `json:"messageBody"`
	RequestId   string `json:"requestId"`
	Status      string `json:"status"`
	CommsId     string `json:"commsId"`
	ProjectCode string `json:"projectCode"`
}
type TrackerResponseType struct {
	StatusCode int       `json:"statusCode"`
	Status     string    `json:"status"`
	ObjectId   string    `json:"objectId"`
	CreatedAt  time.Time `json:"createdAt"`
}
