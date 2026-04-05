package deps

import (
	awsSqS "github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Kanishkmittal55/bridgr-api/internal/auth"
	"github.com/Kanishkmittal55/bridgr-api/internal/cloud"
	"github.com/Kanishkmittal55/bridgr-api/internal/httpx"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository"
	"github.com/Kanishkmittal55/bridgr-api/internal/repository/sqlc"
)

// Deps are Bridgr API runtime dependencies (narrow slice of former monolith dependencies).
type Deps struct {
	ResponseWriter  *httpx.ResponseWriter
	HsQuerier       sqlc.Querier
	Repo            *repository.Repo
	SQSClient       *awsSqS.Client
	AccessToApiKeys map[auth.Access][]string
	S3              cloud.Interface
}
