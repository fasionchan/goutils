package httpx

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1234567890",
		"name": "John Doe",
		"admin": true,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	fmt.Println(tokenString)
}

func TestNewClientJwtPathTokenConfigured(t *testing.T) {
	client, err := NewClient("http://browserd.example/api?_jwtPathTokenSecret=sec&_jwtPathTokenSeconds=120&keep=1")
	require.NoError(t, err)
	require.NotNil(t, client)

	auth, ok := client.requestAuthenticator.(*JwtPathTokenAuthenticator)
	require.True(t, ok)
	assert.Equal(t, 120*time.Second, auth.GetExpireDuration())

	query := client.GetBaseUrl().Query()
	assert.False(t, query.Has(QueryParameterNameJwtPathTokenSecret))
	assert.False(t, query.Has(QueryParameterNameJwtPathTokenSeconds))
	assert.Equal(t, "1", query.Get("keep"))
}

func TestNewClientJwtPathTokenDefaultSeconds(t *testing.T) {
	client, err := NewClient("http://browserd.example/?_jwtPathTokenSecret=sec")
	require.NoError(t, err)

	auth, ok := client.requestAuthenticator.(*JwtPathTokenAuthenticator)
	require.True(t, ok)
	assert.Equal(t, DefaultJwtPathTokenExpireDuration, auth.GetExpireDuration())
}

func TestNewClientJwtPathTokenNoSecret(t *testing.T) {
	client, err := NewClient("http://browserd.example/?keep=1")
	require.NoError(t, err)
	assert.Nil(t, client.requestAuthenticator)

	client, err = NewClient("http://browserd.example/?_jwtPathTokenSecret=&_jwtPathTokenSeconds=120")
	require.NoError(t, err)
	assert.Nil(t, client.requestAuthenticator)
}

func TestNewClientJwtPathTokenInvalidSeconds(t *testing.T) {
	for _, seconds := range []string{"abc", "0", "-1", ""} {
		t.Run(seconds, func(t *testing.T) {
			_, err := NewClient("http://browserd.example/?_jwtPathTokenSecret=sec&_jwtPathTokenSeconds=" + seconds)
			require.Error(t, err)
			assert.Contains(t, err.Error(), QueryParameterNameJwtPathTokenSeconds)
		})
	}
}

func TestBuildRequestWithBodyJwtPathToken(t *testing.T) {
	client, err := NewClient("http://browserd.example/?_jwtPathTokenSecret=sec&_jwtPathTokenSeconds=120")
	require.NoError(t, err)

	request, err := client.BuildRequestWithBody(
		context.Background(),
		http.MethodGet,
		"sessions/abc",
		nil,
		http.Header{"Authorization": []string{"Bearer old"}},
		nil,
		nil,
		"",
	)
	require.NoError(t, err)

	token := request.Header.Get("Authorization")
	require.NotEmpty(t, token)
	assert.False(t, strings.HasPrefix(token, "Bearer "))

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("sec"), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, "sessions/abc", claims["path"])

	exp, ok := claims["exp"].(float64)
	require.True(t, ok)
	assert.InDelta(t, float64(time.Now().Add(120*time.Second).Unix()), exp, 5)
}
