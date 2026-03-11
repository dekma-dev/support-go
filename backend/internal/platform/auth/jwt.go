package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	Subject   string
	Role      string
	TokenType string
}

var errInvalidToken = errors.New("invalid token")

func ParseAndValidateHS256(token string, secret string, now time.Time) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, errInvalidToken
	}

	headerPayload := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, errInvalidToken
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(headerPayload))
	expected := mac.Sum(nil)
	if !hmac.Equal(signature, expected) {
		return Claims{}, errInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, errInvalidToken
	}
	var header map[string]any
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, errInvalidToken
	}
	alg, _ := header["alg"].(string)
	if alg != "HS256" {
		return Claims{}, errInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, errInvalidToken
	}
	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return Claims{}, errInvalidToken
	}

	if exp, ok := numericClaim(payload["exp"]); ok && now.Unix() >= int64(exp) {
		return Claims{}, errInvalidToken
	}
	if nbf, ok := numericClaim(payload["nbf"]); ok && now.Unix() < int64(nbf) {
		return Claims{}, errInvalidToken
	}

	subject, _ := payload["sub"].(string)
	role, _ := payload["role"].(string)
	tokenType, _ := payload["token_type"].(string)
	return Claims{
		Subject:   strings.TrimSpace(subject),
		Role:      strings.ToLower(strings.TrimSpace(role)),
		TokenType: strings.ToLower(strings.TrimSpace(tokenType)),
	}, nil
}

func IssueHS256Token(secret string, claims Claims, now time.Time, expiresAt time.Time) (string, error) {
	now = now.UTC()
	expiresAt = expiresAt.UTC()

	headerBytes, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	payload := map[string]any{
		"sub": claims.Subject,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": expiresAt.Unix(),
	}
	if role := strings.ToLower(strings.TrimSpace(claims.Role)); role != "" {
		payload["role"] = role
	}
	if tokenType := strings.ToLower(strings.TrimSpace(claims.TokenType)); tokenType != "" {
		payload["token_type"] = tokenType
	}
	payload["jti"] = newTokenID()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	unsigned := base64.RawURLEncoding.EncodeToString(headerBytes) + "." + base64.RawURLEncoding.EncodeToString(payloadBytes)

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(unsigned))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigned + "." + signature, nil
}

func newTokenID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	}
	return hex.EncodeToString(bytes)
}

func numericClaim(value any) (int, bool) {
	switch parsed := value.(type) {
	case float64:
		return int(parsed), true
	case json.Number:
		intValue, err := parsed.Int64()
		if err != nil {
			return 0, false
		}
		return int(intValue), true
	case string:
		intValue, err := strconv.Atoi(strings.TrimSpace(parsed))
		if err != nil {
			return 0, false
		}
		return int(intValue), true
	default:
		return 0, false
	}
}
