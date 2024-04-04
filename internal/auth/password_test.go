package auth

import (
	"testing"
)

func TestPassword_Hash(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name         string
		p            Crptr
		args         args
		resultLength int
		wantErr      bool
	}{
		{
			name:         "empty string",
			args:         args{str: ""},
			resultLength: 60,
			wantErr:      false,
		},
		{
			name:         "short string",
			args:         args{str: "jdfhlskajd"},
			resultLength: 60,
			wantErr:      false,
		},
		{
			name:         "extremely long string",
			args:         args{str: "jdlfjsdlksdjflkdsjflkdsjfslkdjfl;sdkfjlskdfjl;sdkfj;dslkfj;dslkfj;dslkfj;sdlkfjd;kslfj;dskfjds;kfdsjf;sdflksdj"},
			resultLength: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCryptor().Encrypt(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Password.Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.resultLength {
				t.Errorf("Password.Hash() = %v, want %v", len(got), tt.resultLength)
			}
		})
	}
}

func TestPassword_Compare(t *testing.T) {
	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name string
		p    Crptr
		args args
		want bool
	}{
		{
			name: "empty password",
			args: args{
				hash:     "$2a$10$tXFSjq1P0MSLBzizJ05HeulzDwbaGwrUYXMnsgiC9Uwbyba1G8hfe",
				password: "",
			},
			want: true,
		},
		{
			name: "short password", args: args{
				hash:     "$2a$10$JOvaQRSoJjEB5zugUS/b4Ow4MIKvb8mwx1pCF1eajghg1s1.6DULK",
				password: "ps",
			},
			want: true,
		},
		{
			name: "extremely long password",
			args: args{
				hash:     "",
				password: "123456789023456hgjfhgkjdfhgjkdfgkljfhgkljfglhsdfjghsljkfl789jkgkjdahfkjdhflksjdfhqwertyuiop[asdfghjkl;lzxcvbnm,./]",
			},
			want: false,
		},
		{
			name: "empty hash",
			args: args{
				hash:     "",
				password: "123",
			},
			want: false,
		},
		{
			name: "invalid hash",
			args: args{
				hash:     "$2a$10$NNFG7QLAijyl5jdupLSWauKlz1/FGZidbhT9gnWK6pAAzbt8ri4xG",
				password: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCryptor().Compare(tt.args.hash, tt.args.password); got != tt.want {
				t.Errorf("Password.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
