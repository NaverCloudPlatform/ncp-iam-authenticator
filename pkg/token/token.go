package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud/credentials"
	"github.com/NaverCloudPlatform/ncp-iam-authenticator/pkg/constants"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthv1beta1 "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
	"strconv"
	"strings"
	"time"
)

// Generator provides new JWT tokens
type Generator interface {
	// Get a token using credentials in the default credentials chain.
	Get(*credentials.Credentials, string, string) (*Token, error)
	// FormatJSON returns the client auth formatted json for the ExecCredential auth
	FormatJSON(Token) (string, error)
}

// Token is generated and used by Kubernetes client-go to authenticate with a Kubernetes cluster.
type Token struct {
	Token string
}

type Claim struct {
	TimeStamp string `json:"timestamp"`
	AccessKey string `json:"accessKey"`
	HmacSign  string `json:"signature"`
	Path      string `json:"path"`
}

type generator struct {
}

// NewGenerator creates a Generator and returns it.
func NewGenerator() (Generator, error) {
	return generator{}, nil
}

func (g generator) Get(credential *credentials.Credentials, clusterId string, region string) (*Token, error) {
	token, err := makeToken(credential.AccessKey(), credential.SecretKey(), clusterId, region)

	if err != nil {
		return nil, err
	}

	return &Token{
		Token: token,
	}, nil
}

func makeToken(accessKey string, secretKey string, clusterId string, region string) (string, error) {
	timestamp := strconv.FormatInt(makeTimestamp(), 10)
	path := getPathWithParams(clusterId, region)
	hmacSign := makeSignature("GET", path, accessKey, secretKey, timestamp)

	tokenClaim := Claim{
		TimeStamp: timestamp,
		AccessKey: accessKey,
		HmacSign:  hmacSign,
		Path:      path,
	}

	e, err := json.Marshal(tokenClaim)
	if err != nil {
		return "", errors.Wrap(err, "failed to make json")
	}

	token := base64.StdEncoding.EncodeToString(e)
	token = constants.TokenPrefix + token

	return token, nil
}

func getStageFromRegion(region string) string {
	if region == "" {
		return "v1"
	}

	switch region {
	case "FKR":
		return "v1"
	case "KR":
		return "v1"
	case "SGN":
		return "sgn-v1"
	case "KRS":
		return "krs-v1"
	default:
		return strings.ToLower(region) + "-v1"
	}
}

func getPathWithParams(clusterId string, region string) string {
	return "/iam/" + getStageFromRegion(region) + "/user?clusterUuid=" + clusterId
}

func makeSignature(method string, uri string, accessKey string, secretKey string, timestamp string) string {
	space := " "    // one space
	newLine := "\n" // new line

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(method))
	mac.Write([]byte(space))
	mac.Write([]byte(uri))
	mac.Write([]byte(newLine))
	mac.Write([]byte(timestamp))
	mac.Write([]byte(newLine))
	mac.Write([]byte(accessKey))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FormatJSON formats the json to support ExecCredential authentication
func (g generator) FormatJSON(token Token) (string, error) {
	execInput := &clientauthv1beta1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "client.authentication.k8s.io/v1beta1",
			Kind:       "ExecCredential",
		},
		Status: &clientauthv1beta1.ExecCredentialStatus{
			Token: token.Token,
		},
	}
	enc, err := json.Marshal(execInput)

	if err != nil {
		return "", errors.Wrap(err, "failed to make json")
	}

	return string(enc), nil
}
