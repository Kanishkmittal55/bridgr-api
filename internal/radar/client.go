package radar

import (
	"context"
	"encoding/json"
	"io"
	"time"

	discoveryv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/discovery/v1"
	pdfv1 "github.com/Kanishkmittal55/bridgr-api/internal/gen/radar/services/pdf/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FetchUrlResult holds URL content fetched via generic_url spider.
type FetchUrlResult struct {
	URL     string
	Title   string
	Content string
}

// Client is a gRPC client for the Python radar service (DiscoveryService, PdfExtractionService).
type Client struct {
	conn         *grpc.ClientConn
	client       discoveryv1.DiscoveryServiceClient
	pdfClient    pdfv1.PdfExtractionServiceClient
	crawlTimeout time.Duration
}

// Config holds client configuration.
type Config struct {
	Addr         string
	CrawlTimeout time.Duration
}

// DefaultCrawlTimeout is the default timeout for crawl operations.
const DefaultCrawlTimeout = 30 * time.Minute

// NewClient creates a new radar gRPC client.
// In Docker, use RADAR_ADDR=radar:50051 (service name). Locally, use localhost:50051.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	addr := cfg.Addr
	if addr == "" {
		addr = "localhost:50051"
	}
	timeout := cfg.CrawlTimeout
	if timeout <= 0 {
		timeout = DefaultCrawlTimeout
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:         conn,
		client:       discoveryv1.NewDiscoveryServiceClient(conn),
		pdfClient:    pdfv1.NewPdfExtractionServiceClient(conn),
		crawlTimeout: timeout,
	}, nil
}

// CrawlSourceStream calls the radar DiscoveryService.CrawlSource RPC (server-streaming).
// Returns a stream; call Recv() in a loop to receive CrawlSourceResponse chunks (e.g. per employer for Workday).
// Caller must call cancel when done to release the timeout context.
func (c *Client) CrawlSourceStream(ctx context.Context, req *discoveryv1.CrawlSourceRequest) (discoveryv1.DiscoveryService_CrawlSourceClient, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(ctx, c.crawlTimeout)
	stream, err := c.client.CrawlSource(ctx, req)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	return stream, cancel, nil
}

// CancelCrawl calls the radar DiscoveryService.CancelCrawl RPC.
// Signals the Python radar to stop the crawl for the given run_uuid.
func (c *Client) CancelCrawl(ctx context.Context, runUUID string) (bool, error) {
	if c.client == nil {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := c.client.CancelCrawl(ctx, &discoveryv1.CancelCrawlRequest{RunUuid: runUUID})
	if err != nil {
		return false, err
	}
	return resp.Cancelled, nil
}

// FetchUrlsTimeout is the timeout for FetchUrls (typically a few URLs).
const FetchUrlsTimeout = 5 * time.Minute

// FetchUrls fetches content from the given URLs via the radar generic_url spider.
// Calls CrawlSource with source_site=generic_url and params.urls. Returns parsed
// url/title/content for each discovery. Does not ingest to DB.
func (c *Client) FetchUrls(ctx context.Context, urls []string) ([]FetchUrlResult, error) {
	if c.client == nil {
		return nil, nil
	}
	if len(urls) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(ctx, FetchUrlsTimeout)
	defer cancel()

	req := &discoveryv1.CrawlSourceRequest{
		SourceSite:    "generic_url",
		DiscoveryType: "idea",
		Params: &discoveryv1.CrawlParams{
			Urls: urls,
		},
	}

	stream, cancel, err := c.CrawlSourceStream(ctx, req)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var results []FetchUrlResult
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return results, err
		}
		if resp == nil || len(resp.Discoveries) == 0 {
			continue
		}
		for _, dwp := range resp.Discoveries {
			if dwp == nil || dwp.Raw == nil || len(dwp.Raw.Payload) == 0 {
				continue
			}
			var payload struct {
				URL     string `json:"url"`
				Title   string `json:"title"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(dwp.Raw.Payload, &payload); err != nil {
				continue
			}
			results = append(results, FetchUrlResult{
				URL:     payload.URL,
				Title:   payload.Title,
				Content: payload.Content,
			})
		}
	}
	return results, nil
}

// ExtractText calls the radar PdfExtractionService.ExtractText RPC.
func (c *Client) ExtractText(ctx context.Context, pdfContent []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	resp, err := c.pdfClient.ExtractText(ctx, &pdfv1.ExtractTextRequest{
		PdfContent: pdfContent,
		Filename:   filename,
	})
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", nil
	}
	return resp.ExtractedText, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
