package deploy

import (
	"log"
	"os"

	"github.com/samber/lo"
)

// ENV loads the main environment variable at the start of the program. Is uses DevEnv as a default value (when no
// environment variable is set).
//
// Supported values are:
//   - DevENV
//   - StagingEnv
//   - ProdENV
//
// Any unsupported value, other than empty, will result in a fatal error and exit the program.
var ENV, _ = lo.Coalesce(os.Getenv("ENV"), DevENV)

const (
	DevENV     = "dev"
	ProdENV    = "prod"
	StagingEnv = "staging"
)

// IsReleaseEnv returns true if the current environment is different from the default DevENV.
func IsReleaseEnv() bool {
	return ENV == ProdENV || ENV == StagingEnv
}

func init() {
	// Prevent invalid environment values on deployment. Local dev should never set the environment variable, except
	// if set to DevENV.
	if ENV != DevENV && ENV != ProdENV && ENV != StagingEnv {
		log.Fatalf("unrecognized value for variable 'ENV': '%s'", ENV)
	}
}
