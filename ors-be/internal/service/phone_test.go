package service

import "testing"

func TestIsPhoneValid(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"digits", "13800138000", true},
		{"spaces and hyphen", "138 0013-8000", true},
		{"letter", "138abc", false},
		{"plus", "+8613800138000", false},
		{"continuous hyphen", "138--0013", false},
		{"spaces between hyphen", "138- -0013", true},
		{"only spaces", "   ", false},
		{"only hyphen", "-", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPhoneValid(tt.phone); got != tt.want {
				t.Errorf("isPhoneValid(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}
