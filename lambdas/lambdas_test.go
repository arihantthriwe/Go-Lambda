package lambdas

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

func TestHandler(t *testing.T) {
	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Handler(tt.args.ctx, tt.args.sqsEvent); (err != nil) != tt.wantErr {
				t.Errorf("Handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_create(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := create(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateTracker(t *testing.T) {
	type args struct {
		request     *Request
		response    *Response
		status      string
		mailError   string
		lambdaError string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			updateTracker(tt.args.request, tt.args.response, tt.args.status, tt.args.mailError, tt.args.lambdaError)
		})
	}
}

func Test_downloadFileFromS3(t *testing.T) {
	type args struct {
		downloadManagerAWS *manager.Downloader
		file               *os.File
		bucket             string
		key                string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if gotErr := downloadFileFromS3(tt.args.downloadManagerAWS, tt.args.file, tt.args.bucket, tt.args.key); gotErr != tt.wantErr {
				t.Errorf("downloadFileFromS3() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
