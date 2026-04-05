package env

import (
	"fmt"
	"os"
)

type Env string

const (
	Development  Env = "development"
	Sandbox      Env = "sandbox"
	Staging      Env = "staging"
	Uat          Env = "uat"
	Production   Env = "production"
	Organization Env = "organization"
)

func IsNonDevelopment(e Env) bool {
	return e != Development
}

func ResolveEnv() (Env, error) {
	switch Env(os.Getenv("ENV")) {
	case Development:
		return Development, nil
	case Sandbox:
		return Sandbox, nil
	case Staging:
		return Staging, nil
	case Uat:
		return Uat, nil
	case Production:
		return Production, nil
	case Organization:
		return Organization, nil
	default:
		return "", fmt.Errorf("unknown environment: %q", os.Getenv("ENV"))
	}
}

func ResolveEnvOrDie() Env {
	e, err := ResolveEnv()
	if err != nil {
		panic(err.Error())
	}
	return e
}
