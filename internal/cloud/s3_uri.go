package cloud

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

// ParseS3URI splits s3://bucket/key into bucket and key (key may contain '/').
func ParseS3URI(uri string) (bucket, key string, err error) {
	u := strings.TrimSpace(uri)
	if !strings.HasPrefix(u, "s3://") {
		return "", "", fmt.Errorf("not an s3 URI")
	}
	rest := strings.TrimPrefix(u, "s3://")
	if rest == "" {
		return "", "", fmt.Errorf("empty s3 URI")
	}
	slash := strings.IndexByte(rest, '/')
	if slash <= 0 || slash == len(rest)-1 {
		return "", "", fmt.Errorf("invalid s3 URI: need bucket and key")
	}
	return rest[:slash], rest[slash+1:], nil
}

// ContentTypeByExtension returns a MIME type from the object key extension, or application/octet-stream.
func ContentTypeByExtension(key string) string {
	ext := strings.ToLower(filepath.Ext(key))
	if ext == "" {
		return "application/octet-stream"
	}
	ct := mime.TypeByExtension(ext)
	if ct == "" {
		return "application/octet-stream"
	}
	return ct
}
