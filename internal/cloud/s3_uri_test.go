package cloud

import "testing"

func TestParseS3URI(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		b, k, err := ParseS3URI("s3://my-bucket/inputs/1/uuid/file.pdf")
		if err != nil {
			t.Fatal(err)
		}
		if b != "my-bucket" || k != "inputs/1/uuid/file.pdf" {
			t.Fatalf("got %q %q", b, k)
		}
	})
	t.Run("reject non s3", func(t *testing.T) {
		_, _, err := ParseS3URI("https://example.com/a")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestContentTypeByExtension(t *testing.T) {
	if ct := ContentTypeByExtension("x.pdf"); ct != "application/pdf" {
		t.Fatalf("pdf: %q", ct)
	}
}
