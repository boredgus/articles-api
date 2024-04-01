package auth

import (
	mocks "a-article/internal/mocks/utils"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT_Generate(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		payload JWTPayload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "success",
			fields:  fields{now: mocks.MockTimeNow(time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC))},
			args:    args{payload: JWTPayload{Username: "user", UserOId: "oid", Role: "user"}},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDIyNTI4MDAsImlhdCI6MTcwMjE2NjQwMCwidXNlcm5hbWUiOiJ1c2VyIiwidXNlcl9vaWQiOiJvaWQiLCJyb2xlIjoidXNlciJ9",
			wantErr: nil,
		},
	}
	secretKey := []byte("secret")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JWT{
				secretKey: secretKey,
				now:       tt.fields.now,
			}.Generate(tt.args.payload)

			assert.Equal(t, strings.Join(strings.Split(got, ".")[:2], "."), tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestJWT_Decode(t *testing.T) {
	type fields struct {
		now func() time.Time
	}
	type args struct {
		token string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    JWTPayload
		wantErr error
	}{
		{
			name:    "invalid signing method",
			fields:  fields{now: mocks.MockTimeNow(time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC))},
			args:    args{token: "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.AlAC2kGm2UZwvgEQS-sbkmbPmRzDuDpTd2WacERgzH4pXs2qU_04w0UcFcTh0DQwy24IMLqwnngrxFhRJuUtdGQBgfCOQjUmhFOwHnfteWlkU2y3yWyZgROUdXirHlgP3zBgLfOpqi_0VDDgupdq-ysBDuciD3ibdmx5AKF_Y9yypHKQQmzHa-JpD66sUeiYq2BzsnKiV9JyEclux9cdXc1dL4jfZlKEOL0wL75aCib3lznPNOAikMzME4SunCD9dQoDKwCQ75smhfVdINUssc_9MMacPcXaaZ5qTIboxZSgsynon8cUXSFH63UBFrl_AixJpJ6Bn-AC8akqBgjl7g"},
			want:    JWTPayload{},
			wantErr: UnexpectedSigningMethod,
		},
		{
			name:    "token is expired",
			args:    args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDIyNTI4MDAsImlhdCI6MTcwMjE2NjQwMCwidXNlcm5hbWUiOiJ1ZGF2IiwidXNlcl9pZCI6ImRlODZjN2FmLTRjMzEtNGYzYi04NWYyLTFmMDQ0ZDQ5NDczOCIsInJvbGUiOiJ1c2VyIn0.Bz3gotNj0qU2xyiFcMPvCO2Ve29qRGP13b2GAaF1mBE"},
			want:    JWTPayload{},
			wantErr: jwt.ErrTokenExpired,
		},
		{
			name:    "successfully decoded",
			fields:  fields{now: mocks.MockTimeNow(time.Date(2033, 12, 10, 0, 0, 0, 0, time.UTC))},
			args:    args{token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIwMTc4NzIwMDAsImlhdCI6MjAxNzc4NTYwMCwidXNlcm5hbWUiOiJ1ZGF2IiwidXNlcl9vaWQiOiJkZTg2YzdhZi00YzMxLTRmM2ItODVmMi0xZjA0NGQ0OTQ3MzgiLCJyb2xlIjoidXNlciJ9.T0t3N3bLkW9NaK26GA1cRzzm7UAKhge21jJNNIaXHzk"},
			want:    JWTPayload{Username: "udav", UserOId: "de86c7af-4c31-4f3b-85f2-1f044d494738", Role: "user"},
			wantErr: nil,
		},
	}
	secretKey := []byte("secret")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := JWT{
				secretKey: secretKey,
				now:       tt.fields.now,
			}
			got, err := tr.Decode(tt.args.token)
			assert.Equal(t, got, tt.want)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
