package core

import "strings"

type Env string

const (
	EnvDevelopment Env = "dev"
	EnvStaging     Env = "stage"
	EnvProduction  Env = "prod"
)

func parseEnv(env string) Env {
	switch strings.ToLower(env) {
	case "dev", "development":
		return EnvDevelopment
	case "stage", "staging":
		return EnvStaging
	case "prod", "production":
		return EnvProduction
	default:
		return EnvDevelopment
	}
}
