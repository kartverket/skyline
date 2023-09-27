package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSenderType_IsValid(t *testing.T) {
	var tests = []struct {
		name  string
		input SenderType
		want  bool
	}{
		{"msgraph should be valid", 0, true},
		{"Dummy should be valid", 1, true},
		{"random should be invalid", 2, false},
	}
	for _, tt := range tests {
		assert.Equal(t, SenderType.IsValid(tt.input), tt.want, tt.name)
	}
}
