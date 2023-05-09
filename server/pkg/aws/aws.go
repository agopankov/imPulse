package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secrets struct {
	BinanceAPIKey    string `json:"BINANCE_API_KEY"`
	BinanceSecretKey string `json:"BINANCE_SECRET_KEY"`
}

func GetSecrets(ctx context.Context, secretID string) (*Secrets, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)

	secretValue, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return nil, err
	}

	var secrets Secrets
	err = json.Unmarshal([]byte(*secretValue.SecretString), &secrets)
	if err != nil {
		return nil, err
	}

	return &secrets, nil
}
