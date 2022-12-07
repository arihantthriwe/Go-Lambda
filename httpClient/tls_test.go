package httpClient

import (
	"crypto"
	"crypto/tls"
	"reflect"
	"testing"
)

func TestCustomLoadX509KeyPair(t *testing.T) {
	type args struct {
		certFile string
		keyFile  string
	}
	tests := []struct {
		name     string
		args     args
		wantCert tls.Certificate
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotCert, err := CustomLoadX509KeyPair(tt.args.certFile, tt.args.keyFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomLoadX509KeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCert, tt.wantCert) {
				t.Errorf("CustomLoadX509KeyPair() = %v, want %v", gotCert, tt.wantCert)
			}
		})
	}
}

func TestX509KeyPair(t *testing.T) {
	type args struct {
		certPEMBlock []byte
		keyPEMBlock  []byte
		pw           []byte
	}
	tests := []struct {
		name     string
		args     args
		wantCert tls.Certificate
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotCert, err := X509KeyPair(tt.args.certPEMBlock, tt.args.keyPEMBlock, tt.args.pw)
			if (err != nil) != tt.wantErr {
				t.Errorf("X509KeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCert, tt.wantCert) {
				t.Errorf("X509KeyPair() = %v, want %v", gotCert, tt.wantCert)
			}
		})
	}
}

func Test_parsePrivateKey(t *testing.T) {
	type args struct {
		der []byte
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.PrivateKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parsePrivateKey(tt.args.der)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
