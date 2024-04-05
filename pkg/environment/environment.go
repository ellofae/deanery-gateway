package environment

import (
	"fmt"
	"os"

	"github.com/ellofae/deanery-gateway/pkg/logger"
)

const (
	SET_STAGE_ENVIRONEMENT_VAR string = "SET_STAGE"
	PRODUCTION_STAGE           string = "prod"
	TEST_STAGE                 string = "test"
)

const (
	PRODUCTION_CONFIGURATION_FILE string = "config_prod.yaml"
	TEST_CONFIGURATION_FILE       string = "config_test.yaml"
)

func ParseEnvironmentVariable() (string, error) {
	logger := logger.GetLogger()

	stage_value, ok := os.LookupEnv(SET_STAGE_ENVIRONEMENT_VAR)
	if !ok {
		logger.Printf("You must specify the project stage. Env variable '%s' is required.\n", SET_STAGE_ENVIRONEMENT_VAR)
		return "", fmt.Errorf("no project stage is specified")
	}

	switch stage_value {
	case PRODUCTION_STAGE:
		logger.Println("Running in production stage.")
		return PRODUCTION_CONFIGURATION_FILE, nil
	case TEST_STAGE:
		logger.Println("Running in test stage.")
		return TEST_CONFIGURATION_FILE, nil
	default:
		logger.Println("Unknown stage passed. Exiting..")
	}

	return "", fmt.Errorf("unknown project stage is specified")
}
