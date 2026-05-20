package validator

import (
	"testing"
)

type signUpFixture struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func TestValidateStruct_SignUpRequest(t *testing.T) {
	v := New()

	tests := []struct {
		name   string
		req    signUpFixture
		fields []string
	}{
		{
			name:   "valid",
			req:    signUpFixture{Email: "a@b.com", Password: "password1"},
			fields: nil,
		},
		{
			name:   "invalid email",
			req:    signUpFixture{Email: "not-email", Password: "password1"},
			fields: []string{"email"},
		},
		{
			name:   "short password",
			req:    signUpFixture{Email: "a@b.com", Password: "short"},
			fields: []string{"password"},
		},
		{
			name:   "missing fields",
			req:    signUpFixture{},
			fields: []string{"email", "password"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := v.ValidateStruct(&tc.req)
			if tc.fields == nil {
				if len(errs) != 0 {
					t.Fatalf("expected no errors, got %v", errs)
				}
				return
			}
			for _, f := range tc.fields {
				if _, ok := errs[f]; !ok {
					t.Fatalf("expected error for field %q, got %v", f, errs)
				}
			}
		})
	}
}
