package healthlake

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"moul.io/http2curl"
)

const (
	jsonExtension         = ".json"
	healthLakeServiceName = "healthlake"
	awsRegion             = "ap-south-1"
	acceptHeaderName      = "Accept"
	contentTypeHeaderName = "Content-Type"
	bundleResourceType    = "Bundle"
)

var jsonContentType = mime.TypeByExtension(jsonExtension)

type Repository interface {
	CreateResource(ctx context.Context, resource []byte,
		resourceType string, optFns ...func(*Options)) error
}

type HealthLake struct {
	url     string
	options Options
	awsCfg  aws.Config
	signer  *v4.Signer
}

const ahlURL = "https://healthlake.%s.amazonaws.com/datastore/%s/r4/"

func NewHealthLakeRepository(config config) Repository {
	url := fmt.Sprintf(ahlURL, config.awsRegion, config.dataStoreID)
	awsConfig, _ := GetAWSConfig()
	signer := v4.NewSigner()

	return &HealthLake{
		url:    url,
		awsCfg: awsConfig,
		signer: signer,
	}
}

func (hl *HealthLake) CreateResource(ctx context.Context, resourceData []byte,
	resourceType string, optFns ...func(*Options)) error {
	var url string
	switch resourceType {
	case bundleResourceType:
		url = hl.url
	default:
		url = fmt.Sprintf("%s%s", hl.url, resourceType)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url,
		bytes.NewBuffer(resourceData))
	if err != nil {
		return err
	}

	hl.addStandardHeadersForPost(req)
	if err := hl.signRequest(ctx, resourceData, req); err != nil {
		return err
	}

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	var authProfileResponse map[string]interface{}
	options := hl.options.Copy()
	for _, fn := range optFns {
		fn(&options)
	}
	if err := Call(ctx, options.GetHTTPClient(), req, &authProfileResponse); err != nil {
		return err
	}

	return nil
}

func (hl *HealthLake) getPayloadHash(bv []byte) string {
	h := sha256.New()
	h.Write(bv)
	return hex.EncodeToString(h.Sum(nil))
}

func (hl *HealthLake) signRequest(ctx context.Context, requestBody []byte, req *http.Request) error {
	creds, err := hl.awsCfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	return hl.signer.SignHTTP(ctx, creds, req, hl.getPayloadHash(requestBody),
		healthLakeServiceName, awsRegion, time.Now())
}

func (hl *HealthLake) addStandardHeadersForPost(req *http.Request) {
	req.Header.Add(acceptHeaderName, jsonContentType)
	req.Header.Add(contentTypeHeaderName, jsonContentType)
}
