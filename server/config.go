package server

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func LoadLastFMAPIKey() string {
	name := os.Getenv("LASTFM_API_KEY_SSM_PARAMETER")
	if name == "" {
		log.Fatal("LASTFM_API_KEY_SSM_PARAMETER is not set")
	}

	ctx := context.Background()
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	out, err := ssm.NewFromConfig(cfg).GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("Failed to read LastFM API key from SSM parameter %s: %v", name, err)
	}

	return strings.TrimSpace(aws.ToString(out.Parameter.Value))
}
