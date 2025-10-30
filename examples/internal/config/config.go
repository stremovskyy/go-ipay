package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config aggregates credentials that examples require.
type Config struct {
	MerchantName         string
	MerchantNameWithdraw string
	MerchantID           string
	MerchantIDWithdraw   string
	Login                string
	MerchantKey          string
	MerchantKeyWithdraw  string
	SystemKey            string
	SuccessRedirect      string
	FailRedirect         string
	SubMerchantID        int
	IpayPaymentID        int
	CardToken            string
	GoogleToken          string
	WebhookURL           string
	AppleContainer       string
}

var defaultEnvPaths = []string{
	".env.local",
	".env",
	"examples/.env.local",
	"examples/.env",
}

// MustLoad loads configuration and panics if required values are missing.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Load populates configuration from environment variables. It optionally pulls
// values from a .env-compatible file to simplify local development.
func Load() (*Config, error) {
	if err := hydrateEnv(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	var err error

	if cfg.MerchantName, err = requireString("IPAY_MERCHANT_NAME"); err != nil {
		return nil, err
	}
	if cfg.MerchantNameWithdraw, err = requireString("IPAY_MERCHANT_NAME_WITHDRAW"); err != nil {
		return nil, err
	}
	if cfg.MerchantID, err = requireString("IPAY_MERCHANT_ID"); err != nil {
		return nil, err
	}
	if cfg.MerchantIDWithdraw, err = requireString("IPAY_MERCHANT_ID_WITHDRAW"); err != nil {
		return nil, err
	}
	if cfg.Login, err = requireString("IPAY_LOGIN"); err != nil {
		return nil, err
	}
	if cfg.MerchantKey, err = requireString("IPAY_MERCHANT_KEY"); err != nil {
		return nil, err
	}
	if cfg.MerchantKeyWithdraw, err = requireString("IPAY_MERCHANT_KEY_WITHDRAW"); err != nil {
		return nil, err
	}
	if cfg.SystemKey, err = requireString("IPAY_SYSTEM_KEY"); err != nil {
		return nil, err
	}
	if cfg.SuccessRedirect, err = requireString("IPAY_SUCCESS_REDIRECT"); err != nil {
		return nil, err
	}
	if cfg.FailRedirect, err = requireString("IPAY_FAIL_REDIRECT"); err != nil {
		return nil, err
	}
	if cfg.CardToken, err = requireString("IPAY_CARD_TOKEN"); err != nil {
		return nil, err
	}
	if cfg.GoogleToken, err = requireString("IPAY_GOOGLE_TOKEN"); err != nil {
		return nil, err
	}
	if cfg.WebhookURL, err = requireString("IPAY_WEBHOOK_URL"); err != nil {
		return nil, err
	}
	if cfg.AppleContainer, err = requireString("IPAY_APPLE_CONTAINER"); err != nil {
		return nil, err
	}
	if cfg.SubMerchantID, err = requireInt("IPAY_SUB_MERCHANT_ID"); err != nil {
		return nil, err
	}
	if cfg.IpayPaymentID, err = requireInt("IPAY_IPAY_PAYMENT_ID"); err != nil {
		return nil, err
	}

	return cfg, nil
}

func hydrateEnv() error {
	if custom := strings.TrimSpace(os.Getenv("IPAY_EXAMPLES_ENV_FILE")); custom != "" {
		if err := loadEnvFile(custom); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("env file %q (referenced by IPAY_EXAMPLES_ENV_FILE) not found", custom)
			}
			return err
		}
		return nil
	}

	var attempted bool
	for _, path := range defaultEnvPaths {
		attempted = true
		if err := loadEnvFile(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		// stop after successfully loading the first available file
		return nil
	}

	if attempted {
		// none of the default files existed; that's acceptable if env vars are pre-set
	}

	return nil
}

func loadEnvFile(path string) error {
	// #nosec G304 -- configuration files are explicitly chosen by the developer.
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			val = strings.Trim(val, `"'`)

			if key == "" {
				continue
			}

			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
	}

	return scanner.Err()
}

func requireString(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is empty", key)
	}
	return val, nil
}

func requireInt(key string) (int, error) {
	raw, err := requireString(key)
	if err != nil {
		return 0, err
	}
	value, convErr := strconv.Atoi(raw)
	if convErr != nil {
		return 0, fmt.Errorf("environment variable %s must be an integer: %w", key, convErr)
	}
	return value, nil
}
