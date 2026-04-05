package deps

import (
	awsSqS "github.com/aws/aws-sdk-go-v2/service/sqs"
	hsHttp "github.com/hassleskip/hassle-go/pkg/http"
	"github.com/hassleskip/hassle-go/pkg/middleware/auth"

	"github.com/hassleskip/bridgr-api/internal/cloud"
	"github.com/hassleskip/bridgr-api/internal/repository"
	"github.com/hassleskip/bridgr-api/internal/repository/sqlc"
)

// Deps are Bridgr API runtime dependencies (narrow slice of former monolith dependencies).
type Deps struct {
	ResponseWriter  *hsHttp.ResponseWriter
	HsQuerier       sqlc.Querier
	Repo            *repository.Repo
	SQSClient       *awsSqS.Client
	AccessToApiKeys map[auth.Access][]string
	S3              cloud.Interface
}
