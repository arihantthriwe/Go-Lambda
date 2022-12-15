package main

import (
	"golambda/lambdas"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// e := echo.New()
	// e.POST("/test",Handler)
	// e.Logger.Fatal(e.Start(":8000"))
	lambda.Start(lambdas.Handler)

}
