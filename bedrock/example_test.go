package bedrock_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/openai/openai-go/v3/bedrock"
)

var exampleAWSEnvironment = []string{
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_SESSION_TOKEN",
	"AWS_PROFILE",
	"AWS_DEFAULT_PROFILE",
	"AWS_REGION",
	"AWS_DEFAULT_REGION",
	"AWS_SHARED_CREDENTIALS_FILE",
	"AWS_CONFIG_FILE",
	"AWS_WEB_IDENTITY_TOKEN_FILE",
	"AWS_ROLE_ARN",
	"AWS_ROLE_SESSION_NAME",
	"AWS_CONTAINER_CREDENTIALS_FULL_URI",
	"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI",
	"AWS_BEARER_TOKEN_BEDROCK",
	"AWS_BEDROCK_BASE_URL",
	"AWS_EC2_METADATA_DISABLED",
}

func useExampleAWSConfig() func() {
	dir, err := os.MkdirTemp("", "openai-go-bedrock-example-*")
	if err != nil {
		panic(err)
	}

	credentialsPath := filepath.Join(dir, "credentials")
	if err := os.WriteFile(credentialsPath, []byte(`[default]
aws_access_key_id = example-access-key
aws_secret_access_key = example-secret-key

[production]
aws_access_key_id = example-access-key
aws_secret_access_key = example-secret-key
`), 0o600); err != nil {
		_ = os.RemoveAll(dir)
		panic(err)
	}

	configPath := filepath.Join(dir, "config")
	if err := os.WriteFile(configPath, []byte(`[default]
region = us-west-2

[profile production]
region = us-west-2
`), 0o600); err != nil {
		_ = os.RemoveAll(dir)
		panic(err)
	}

	previous := make(map[string]string, len(exampleAWSEnvironment))
	present := make(map[string]bool, len(exampleAWSEnvironment))
	for _, name := range exampleAWSEnvironment {
		previous[name], present[name] = os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			_ = os.RemoveAll(dir)
			panic(err)
		}
	}

	for _, value := range []struct {
		name  string
		value string
	}{
		{"AWS_SHARED_CREDENTIALS_FILE", credentialsPath},
		{"AWS_CONFIG_FILE", configPath},
		{"AWS_EC2_METADATA_DISABLED", "true"},
	} {
		if err := os.Setenv(value.name, value.value); err != nil {
			_ = os.RemoveAll(dir)
			panic(err)
		}
	}

	return func() {
		for _, name := range exampleAWSEnvironment {
			if present[name] {
				_ = os.Setenv(name, previous[name])
			} else {
				_ = os.Unsetenv(name)
			}
		}
		_ = os.RemoveAll(dir)
	}
}

func ExampleNewClient() {
	restore := useExampleAWSConfig()
	defer restore()

	// The default AWS credential chain supports environment credentials,
	// ~/.aws/credentials, ~/.aws/config, SSO, web identity, and workload roles.
	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSRegion: "us-west-2",
	})
	if err != nil {
		panic(err)
	}
	_ = client

	// Output:
}

func ExampleNewClient_profile() {
	restore := useExampleAWSConfig()
	defer restore()

	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSProfile: "production",
	})
	if err != nil {
		panic(err)
	}
	_ = client

	// Output:
}

func ExampleNewClient_bearer() {
	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSRegion: "us-west-2",
		APIKey:    "bedrock-bearer-token",
		BaseURL:   "https://bedrock-mantle.us-west-2.api.aws/openai/v1",
	})
	if err != nil {
		panic(err)
	}
	_ = client

	// Output:
}
