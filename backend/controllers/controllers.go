package controllers

import (
	"github.com/firebase/genkit/go/genkit"
	"gorm.io/gorm"
)

// This file centralizes the package-level variables that are shared
// across different controller files, like the database connection
// and the Genkit instance.

var (
	// DB is the database connection pool, initialized in main.go.
	DB *gorm.DB
	// G is the Genkit instance, initialized in main.go.
	G *genkit.Genkit
)
