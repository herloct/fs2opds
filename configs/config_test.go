package configs

import (
	"os"
	"testing"
)

func TestGetCatalogConfig(t *testing.T) {
	os.Clearenv()

	catalogConfig := GetCatalogConfig()
	if catalogConfig.Title == "" {
		t.Error("CatalogConfig.Title, should have default value")
	}
}

func TestCatalogConfigValidate(t *testing.T) {
	os.Clearenv()
	os.Setenv("CATALOG_ROOT_PATH", "/books")

	catalogConfig := GetCatalogConfig()
	err := catalogConfig.Validate()
	if err != nil {
		t.Error("CatalogConfig.Validate(), should not throws error")
	}
}

func TestCatalogConfigValidateWithEmptyEnv(t *testing.T) {
	os.Clearenv()

	catalogConfig := GetCatalogConfig()
	err := catalogConfig.Validate()
	if err == nil {
		t.Error("CatalogConfig.Validate(), should throws error")
	}
}

func TestCatalogConfigNeedsAuth(t *testing.T) {
	var tests = []struct {
		message  string
		username string
		password string
		want     bool
	}{
		{"both username and password is empty", "", "", false},
		{"username is empty", "", "password", false},
		{"password is empty", "username", "", false},
		{"both username and password is not empty", "username", "password", true},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			catalogConfig := CatalogConfig{
				Username: tt.username,
				Password: tt.password,
			}
			if got := catalogConfig.NeedsAuth(); got != tt.want {
				t.Errorf("CatalogConfig.NeedsAuth() with %s, should be %v", tt.message, tt.want)
			}
		})
	}
}
