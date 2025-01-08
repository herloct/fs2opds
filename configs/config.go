package configs

import "os"

type CatalogConfig struct {
	Title    string
	RootPath string
	Username string
	Password string
}

func (c CatalogConfig) IsValid() bool {
	if c.Title == "" || c.RootPath == "" {
		return false
	}

	return true
}

func (c CatalogConfig) NeedsAuth() bool {
	if c.Username == "" || c.Password == "" {
		return false
	}

	return true
}

func GetCatalogConfig() CatalogConfig {
	catalog := CatalogConfig{
		Title:    os.Getenv("CATALOG_TITLE"),
		RootPath: os.Getenv("CATALOG_ROOT_PATH"),
		Username: os.Getenv("CATALOG_USERNAME"),
		Password: os.Getenv("CATALOG_PASSWORD"),
	}

	return catalog
}
