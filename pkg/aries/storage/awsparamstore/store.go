package awsparamstore

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hyperledger/aries-framework-go/spi/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	appconfig "github.com/medibloc/vc-service/pkg/config"
)

type store struct {
	ssmClient *ssm.Client
	keyPrefix string
}

func newStore(keyPrefix string) (*store, error) {
	// Load AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     appconfig.Cfg.AWSAccessKey,
				SecretAccessKey: appconfig.Cfg.AWSSecretAccessKey,
			},
		}),
		config.WithRegion(appconfig.Cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS default config: %w", err)
	}

	return &store{
		ssmClient: ssm.NewFromConfig(cfg),
		keyPrefix: keyPrefix,
	}, nil
}

func (s *store) wrapKey(key string) string {
	return fmt.Sprintf("%s-%s", s.keyPrefix, key)
}

func (s *store) Put(key string, value []byte, tags ...storage.Tag) error {
	if key == "" {
		return errors.New("key cannot be blank")
	}
	if value == nil {
		return errors.New("value cannot be nil")
	}

	_, err := s.ssmClient.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:           aws.String(s.wrapKey(key)),
		Value:          aws.String(base64.RawURLEncoding.EncodeToString(value)),
		AllowedPattern: nil,
		DataType:       aws.String("text"),
		Description:    aws.String("BBS+ private key for VC selective disclosure"),
		KeyId:          aws.String(appconfig.Cfg.AWSKmsID),
		Overwrite:      false,
		Policies:       nil,
		Tags:           nil,
		Tier:           types.ParameterTierStandard,
		Type:           types.ParameterTypeSecureString,
	})
	if err != nil {
		return fmt.Errorf("failed to put parameter: %w", err)
	}

	return nil
}

func (s *store) Get(key string) ([]byte, error) {
	out, err := s.ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(s.wrapKey(key)),
		WithDecryption: true,
	})
	if err != nil {
		var notFoundErr *types.ParameterNotFound
		if errors.As(err, &notFoundErr) {
			return nil, storage.ErrDataNotFound
		}
		return nil, fmt.Errorf("failed to get parameter: %w", err)
	}

	value, err := base64.RawURLEncoding.DecodeString(aws.ToString(out.Parameter.Value))
	if err != nil {
		return nil, fmt.Errorf("failed to decode a base64-encoded value: %w", err)
	}

	return value, nil
}

func (s *store) GetTags(key string) ([]storage.Tag, error) {
	return nil, errors.New("not supported")
}

func (s *store) GetBulk(keys ...string) ([][]byte, error) {
	return nil, errors.New("not supported")
}

func (s *store) Query(expression string, options ...storage.QueryOption) (storage.Iterator, error) {
	return nil, errors.New("not supported")
}

func (s *store) Delete(key string) error {
	_, err := s.ssmClient.DeleteParameter(context.TODO(), &ssm.DeleteParameterInput{
		Name: aws.String(s.wrapKey(key)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete parameter: %w", err)
	}
	return nil
}

func (s *store) Batch(operations []storage.Operation) error {
	return errors.New("not supported")
}

func (s *store) Flush() error {
	return nil
}

func (s *store) Close() error {
	return nil
}
