package configs

import (
	"fmt"
	"os"
)

type CatalogConfig struct {
	Title    string
	RootPath string
	Username string
	Password string
}

func (c CatalogConfig) Validate() error {
	if c.RootPath == "" {
		return fmt.Errorf("%s is not set", "CATALOG_ROOT_PATH")
	}

	return nil
}

func (c CatalogConfig) NeedsAuth() bool {
	if c.Username == "" || c.Password == "" {
		return false
	}

	return true
}

func GetCatalogConfig() *CatalogConfig {
	defaultTitle := "OPDS Catalog"

	catalog := CatalogConfig{
		RootPath: os.Getenv("CATALOG_ROOT_PATH"),
		Title:    os.Getenv("CATALOG_TITLE"),
		Username: os.Getenv("CATALOG_USERNAME"),
		Password: os.Getenv("CATALOG_PASSWORD"),
	}

	if catalog.Title == "" {
		catalog.Title = defaultTitle
	}

	return &catalog
}
