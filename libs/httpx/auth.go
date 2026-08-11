/*
 * Author: fasion
 * Created time: 2024-12-12 14:22:42
 * Last Modified by: fasion
 * Last Modified time: 2025-06-27 15:50:18
 */

package httpx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/form3tech-oss/jwt-go"
)

const (
	HmacHashMethodSha1    = "sha1"
	HmacHashMethodSha256  = "sha256"
	HmacHashMethodSha512  = "sha512"
	HmacHashMethodDefault = HmacHashMethodSha256

	HeaderHmacAccessKey     = "X-HMAC-ACCESS-KEY"
	HeaderHmacAlgorithm     = "X-HMAC-ALGORITHM"
	HeaderHmacSignedHeaders = "X-HMAC-SIGNED-HEADERS"
	HeaderHmacSignature     = "X-HMAC-SIGNATURE"
)

var hmacHashMethods = map[string]func() hash.Hash{
	HmacHashMethodSha1:   sha1.New,
	HmacHashMethodSha256: sha256.New,
	HmacHashMethodSha512: sha512.New,
}

type RequestAuthenticator interface {
	Authenticate(request *http.Request) error
}

type RequestAuthenticatorFunc func(*http.Request) error

func (f RequestAuthenticatorFunc) Authenticate(request *http.Request) error {
	if f == nil {
		return nil
	}
	return f(request)
}

type JwtPathTokenAuthenticator struct {
	secret         string
	expireDuration time.Duration
}

func NewJwtPathTokenAuthenticator(secret string, expireDuration time.Duration) *JwtPathTokenAuthenticator {
	return &JwtPathTokenAuthenticator{
		secret:         secret,
		expireDuration: expireDuration,
	}
}

func NewJwtPathTokenAuthenticatorFromUrlValues(values url.Values) (*JwtPathTokenAuthenticator, error) {
	secret := values.Get(QueryParameterNameJwtPathTokenSecret)
	if secret == "" {
		return nil, nil
	}

	expireDuration := DefaultJwtPathTokenExpireDuration
	if values.Has(QueryParameterNameJwtPathTokenSeconds) {
		raw := values.Get(QueryParameterNameJwtPathTokenSeconds)
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds <= 0 {
			return nil, fmt.Errorf("invalid %s: %q", QueryParameterNameJwtPathTokenSeconds, raw)
		}
		expireDuration = time.Duration(seconds) * time.Second
	}

	return NewJwtPathTokenAuthenticator(secret, expireDuration), nil
}

func GenerateJwtPathToken(path, secret string, expireDuration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"path": path,
		"exp":  time.Now().Add(expireDuration).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (authenticator *JwtPathTokenAuthenticator) GetExpireDuration() time.Duration {
	if authenticator == nil {
		return 0
	}
	return authenticator.expireDuration
}

func (authenticator *JwtPathTokenAuthenticator) Authenticate(request *http.Request) error {
	if authenticator == nil {
		return baseutils.NewNilError("JwtPathTokenAuthenticator")
	}
	if request == nil || request.URL == nil {
		return baseutils.NewNilError("http.Request")
	}

	path := strings.Trim(request.URL.Path, "/")
	token, err := GenerateJwtPathToken(path, authenticator.secret, authenticator.expireDuration)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", token)
	return nil
}

type HmacAuthenticatorPtr = *HmacAuthenticator

type HmacAuthenticator struct {
	hashMethod string
	accessKey  string
	secretKey  string
	paths      types.Strings
}

func NewHmacAuthenticator(accessKey, secretKey string) *HmacAuthenticator {
	return &HmacAuthenticator{
		hashMethod: HmacHashMethodDefault,
		accessKey:  accessKey,
		secretKey:  secretKey,
	}
}

func NewHmacAuthenticatorFromUrlValues(values url.Values) *HmacAuthenticator {
	if !values.Has(QueryParameterNameHmacAuthAccessKey) {
		return nil
	}

	return NewHmacAuthenticator(
		values.Get(QueryParameterNameHmacAuthAccessKey),
		values.Get(QueryParameterNameHmacAuthSecretKey),
	).WithHashMethod(values.Get(QueryParameterNameHmacAuthHashMethod))
}

func (hmacAuthenticator *HmacAuthenticator) GetPaths() types.Strings {
	if hmacAuthenticator == nil {
		return nil
	}
	return hmacAuthenticator.paths
}

func (hmacAuthenticator *HmacAuthenticator) WithHashMethod(hashMethod string) *HmacAuthenticator {
	if hmacAuthenticator == nil {
		return nil
	}

	if hashMethod == "" {
		hashMethod = HmacHashMethodDefault
	}

	hmacAuthenticator.hashMethod = hashMethod
	return hmacAuthenticator
}

func (hmacAuthenticator *HmacAuthenticator) WithPaths(paths types.Strings) *HmacAuthenticator {
	if hmacAuthenticator == nil {
		return nil
	}

	if paths.Empty() {
		paths = types.NewStrings("*")
	}

	hmacAuthenticator.paths = paths
	return hmacAuthenticator
}

func (hmacAuthenticator *HmacAuthenticator) WithCommaSeparatedPaths(pathPrefixes string) *HmacAuthenticator {
	return hmacAuthenticator.WithPaths(types.NewStrings(pathPrefixes).Split(",").TrimSpace().PurgeZero())
}

func (hmacAuthenticator *HmacAuthenticator) Authenticate(request *http.Request) error {
	if hmacAuthenticator == nil {
		return baseutils.NewNilError("HmacAuthenticator")
	}

	hashMethod, ok := hmacHashMethods[hmacAuthenticator.hashMethod]
	if !ok {
		return baseutils.NewNotImplementedError(fmt.Sprintf("HmacAuthenticator-%s", hmacAuthenticator.hashMethod))
	}

	hmacHash := hmac.New(hashMethod, []byte(hmacAuthenticator.secretKey))

	now := time.Now()
	date := now.UTC().Format(time.RFC1123)

	requestUrl := request.URL
	buffer := bytes.NewBuffer(nil)
	fmt.Fprintf(io.MultiWriter(hmacHash, buffer), "%s\n%s\n%s\n%s\n%s\n%s:%s\n",
		request.Method,
		requestUrl.Path,
		requestUrl.RawQuery,
		hmacAuthenticator.accessKey,
		date,
		HeaderHmacAccessKey,
		hmacAuthenticator.accessKey,
	)

	request.Header.Set("Date", date)
	request.Header.Set(HeaderHmacAccessKey, hmacAuthenticator.accessKey)
	request.Header.Set(HeaderHmacAlgorithm, fmt.Sprintf("hmac-%s", hmacAuthenticator.hashMethod))
	request.Header.Set(HeaderHmacSignedHeaders, HeaderHmacAccessKey)
	request.Header.Set(HeaderHmacSignature, base64.StdEncoding.EncodeToString(hmacHash.Sum(nil)))

	return nil
}

type HmacAuthenticators []*HmacAuthenticator

func NewHmacAuthenticatorsFromUrlValues(values url.Values) HmacAuthenticators {
	accessKeys := values[QueryParameterNameHmacAuthAccessKey]
	secretKeys := values[QueryParameterNameHmacAuthSecretKey]
	hashMethods := values[QueryParameterNameHmacAuthHashMethod]
	pathPrefixes := values[QueryParameterNameHmacAuthPaths]

	nSecretKeys := len(secretKeys)
	nHashMethods := len(hashMethods)
	nPaths := len(pathPrefixes)

	// 默认密钥为最后一个（密钥相同的key写在最后面，只写一次）
	var defaultSecretKey string
	if nSecretKeys > 0 {
		defaultSecretKey = secretKeys[nSecretKeys-1]
	}

	// 默认哈希方法为最后一个
	var defaultHashMethod string
	if nHashMethods > 0 {
		defaultHashMethod = hashMethods[nHashMethods-1]
	}

	return stl.MapPro(accessKeys, func(index int, accessKey string, _ []string) *HmacAuthenticator {
		secretKey := defaultSecretKey
		if index < nSecretKeys {
			secretKey = secretKeys[index]
		}

		hashMethod := defaultHashMethod
		if index < nHashMethods {
			hashMethod = hashMethods[index]
		}

		var pathPrefix string
		if index < nPaths {
			pathPrefix = pathPrefixes[index]
		}

		return NewHmacAuthenticator(accessKey, secretKey).
			WithHashMethod(hashMethod).
			WithCommaSeparatedPaths(pathPrefix)
	})
}

func (authenticators HmacAuthenticators) Len() int {
	return len(authenticators)
}

func (authenticators HmacAuthenticators) MappingByPath() HmacAuthenticatorByString {
	return stl.MappingByKeys(authenticators, HmacAuthenticatorPtr.GetPaths)
}

func (authenticators HmacAuthenticators) RequestAuthenticatorByPath() RequestAuthenticatorByPath {
	return RequestAuthenticatorByPath(authenticators.MappingByPath().AuthenticatorMapping())
}

func (authenticators HmacAuthenticators) RequestAuthenticator() RequestAuthenticator {
	n := authenticators.Len()

	switch {
	case n == 0:
		return nil
	case n == 1:
		authenticator := authenticators[0]
		paths := authenticator.GetPaths()
		if paths.Empty() || paths.Contain("*") {
			return authenticator
		}
	}

	return authenticators.RequestAuthenticatorByPath()
}

type HmacAuthenticatorByString map[string]*HmacAuthenticator

func (m HmacAuthenticatorByString) AuthenticatorMapping() RequestAuthenticatorByString {
	return stl.MapMap[RequestAuthenticatorByString](m, func(key string, value *HmacAuthenticator, _ HmacAuthenticatorByString) (string, RequestAuthenticator) {
		return key, value
	})
}

type BasicAuthenticator struct {
	user     string
	password string
}

func NewBasicAuthenticator(user, password string) *BasicAuthenticator {
	return &BasicAuthenticator{user: user, password: password}
}

func NewBasicAuthenticatorFromUrlUserinfo(user *url.Userinfo) *BasicAuthenticator {
	password, _ := user.Password()
	return NewBasicAuthenticator(user.Username(), password)
}

func (authenticator *BasicAuthenticator) Authenticate(request *http.Request) error {
	if authenticator == nil {
		return baseutils.NewNilError("BasicAuthenticator")
	}

	if request != nil {
		request.SetBasicAuth(authenticator.user, authenticator.password)
	}

	return nil
}

type RequestAuthenticatorByString map[string]RequestAuthenticator

// todo：后续需要支持前缀匹配，用trie树？
type RequestAuthenticatorByPath map[string]RequestAuthenticator

func (authenticators RequestAuthenticatorByPath) Empty() bool {
	return len(authenticators) == 0
}

func (authenticators RequestAuthenticatorByPath) Authenticate(request *http.Request) error {
	if authenticators.Empty() {
		return nil
	}

	path := request.URL.Path
	authenticator, ok := authenticators[path]
	if !ok {
		// 尝试通配认证器
		authenticator, ok = authenticators["*"]
	}
	if !ok {
		return baseutils.NewGenericNotFoundError("authenticator", path)
	}

	return authenticator.Authenticate(request)
}
