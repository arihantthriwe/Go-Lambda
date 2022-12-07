package httpClient

import (
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseClient(t *testing.T) {
	type args struct {
		method  string
		url     string
		payload *strings.Reader
		v       interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseClient(tt.args.method, tt.args.url, tt.args.payload, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalClient(t *testing.T) {
	type args struct {
		method          string
		url             string
		payload         *strings.Reader
		v               interface{}
		rootCertFile    *os.File
		caChainCertFile *os.File
		certKey         *os.File
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalClient(tt.args.method, tt.args.url, tt.args.payload, tt.args.v, tt.args.rootCertFile, tt.args.caChainCertFile, tt.args.certKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NormalClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNormalClient(t *testing.T) {
	type args struct {
		method          string
		url             string
		payload         *strings.Reader
		v               interface{}
		rootCertFile    *os.File
		caChainCertFile *os.File
		certKey         *os.File
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewNormalClient(tt.args.method, tt.args.url, tt.args.payload, tt.args.v, tt.args.rootCertFile, tt.args.caChainCertFile, tt.args.certKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNormalClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNormalClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfanityCheck(t *testing.T) {
	type args struct {
		messageBody string
		requestId   string
		projectCode string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ProfanityCheck(tt.args.messageBody, tt.args.requestId, tt.args.projectCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfanityCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProfanityCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
