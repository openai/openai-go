package webhooks

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const (
	webhookSecretPrefix = "whsec_"
	// Standard Webhooks symmetric signing keys contain at least 24 bytes.
	minimumWebhookSecretBytes = 24
)

func decodeWebhookSecret(secret string) ([]byte, error) {
	encodedSecret, prefixed := strings.CutPrefix(secret, webhookSecretPrefix)
	if !prefixed {
		return []byte(secret), nil
	}

	decodedSecret, err := base64.StdEncoding.Strict().DecodeString(encodedSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook secret format: %v", err)
	}
	if base64.StdEncoding.EncodeToString(decodedSecret) != encodedSecret {
		return nil, errors.New("invalid webhook secret format: expected canonical base64 encoding")
	}
	if len(decodedSecret) < minimumWebhookSecretBytes {
		return nil, fmt.Errorf("invalid webhook secret: decoded secret must be at least %d bytes", minimumWebhookSecretBytes)
	}

	return decodedSecret, nil
}
