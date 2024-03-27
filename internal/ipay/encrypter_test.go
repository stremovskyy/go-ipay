package ipay

import (
	"testing"

	"github.com/stremovskyy/go-ipay/internal/log"
)

func TestNewCipher(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string // We expect to test if the key is correctly set
	}{
		{
			name:    "Valid key",
			args:    args{key: "mysecretkey1234567890"},
			wantKey: "mysecretkey1234567890",
		},
		{
			name:    "Empty key",
			args:    args{key: ""},
			wantKey: "",
		},
		{
			name:    "Key with spaces",
			args:    args{key: "my secret key"},
			wantKey: "my secret key",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewCipher(tt.args.key).(*encrypter) // Type assertion to access the key
				if got.key != tt.wantKey {
					t.Errorf("NewCipher() key = %v, want %v", got.key, tt.wantKey)
				}
			},
		)
	}
}

func Test_encrypter_EncryptData(t *testing.T) {
	mockLogger := log.NewLogger("cipher")

	type fields struct {
		key    string
		logger *log.Logger
	}
	type args struct {
		rawData string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		out     string
	}{
		{
			name: "Valid encryption",
			fields: fields{
				key:    "mysecretkey1234567890",
				logger: mockLogger,
			},
			args: args{
				rawData: "Hello, World!",
			},
			wantErr: false,
			out:     "npswqWqNyYfvPH+DXg==.vmQarM/mLFQQG57mIGS2Uw==",
		},
		{
			name: "Empty raw data",
			fields: fields{
				key:    "mysecretkey1234567890",
				logger: mockLogger,
			},
			args: args{
				rawData: "",
			},
			wantErr: true,
		},
		{
			name: "Empty key",
			fields: fields{
				key:    "",
				logger: mockLogger,
			},
			args: args{
				rawData: "Hello, World!",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &encrypter{
					key:    tt.fields.key,
					logger: tt.fields.logger,
				}
				got, err := c.EncryptData(tt.args.rawData)
				if (err != nil) != tt.wantErr {
					t.Errorf("EncryptData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && got == "" {
					t.Errorf("EncryptData() expected non-empty output")
				}

				if !tt.wantErr && got != tt.out {
					t.Errorf("EncryptData() got = %v, want %v", got, tt.out)
				}

			},
		)
	}
}
