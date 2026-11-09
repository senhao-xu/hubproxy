package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-containerregistry/pkg/v1"
)

func TestCollectPlatforms(t *testing.T) {
	tests := []struct {
		name      string
		manifests []v1.Descriptor
		want      []string
	}{
		{
			name: "variant kept, deduped and sorted",
			manifests: []v1.Descriptor{
				{Platform: &v1.Platform{OS: "linux", Architecture: "arm", Variant: "v7"}},
				{Platform: &v1.Platform{OS: "linux", Architecture: "amd64"}},
				{Platform: &v1.Platform{OS: "linux", Architecture: "arm"}},
				{Platform: &v1.Platform{OS: "linux", Architecture: "arm", Variant: "v7"}},
				{Platform: &v1.Platform{OS: "linux", Architecture: "arm64", Variant: "v8"}},
				{Platform: &v1.Platform{OS: "windows", Architecture: "amd64"}},
			},
			want: []string{"linux/amd64", "linux/arm", "linux/arm/v7", "linux/arm64/v8", "windows/amd64"},
		},
		{
			name: "nil platform skipped",
			manifests: []v1.Descriptor{
				{Platform: nil},
				{Platform: &v1.Platform{OS: "linux", Architecture: "amd64"}},
			},
			want: []string{"linux/amd64"},
		},
		{
			name:      "no manifests yields empty list",
			manifests: nil,
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectPlatforms(&v1.IndexManifest{Manifests: tt.manifests})
			if !slices.Equal(got, tt.want) {
				t.Fatalf("platforms = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadDebouncer(t *testing.T) {
	d := NewDownloadDebouncer(time.Minute)
	if !d.ShouldAllow("user", "content") {
		t.Fatal("first request denied")
	}
	if d.ShouldAllow("user", "content") {
		t.Fatal("duplicate request allowed")
	}
	if !d.ShouldAllow("other", "content") {
		t.Fatal("different user denied")
	}
}

func TestTokenStoreCreateConsume(t *testing.T) {
	store := newTokenStore[SingleDownloadRequest]()
	req := SingleDownloadRequest{Image: "nginx:latest", Platform: "linux/amd64", UseCompressedLayers: true}

	token, err := store.create(req, "127.0.0.1", "ua")
	if err != nil {
		t.Fatal(err)
	}

	got, ok := store.consume(token, "127.0.0.1", "ua")
	if !ok {
		t.Fatal("token not consumed")
	}
	if got != req {
		t.Fatalf("request = %#v, want %#v", got, req)
	}
	if _, ok := store.consume(token, "127.0.0.1", "ua"); ok {
		t.Fatal("token consumed twice")
	}
}

func TestTokenStoreRejectsDifferentClient(t *testing.T) {
	store := newTokenStore[SingleDownloadRequest]()
	token, err := store.create(SingleDownloadRequest{Image: "nginx:latest"}, "127.0.0.1", "ua")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.consume(token, "127.0.0.2", "ua"); ok {
		t.Fatal("token accepted for different IP")
	}
}

func TestGenerateContentFingerprintStable(t *testing.T) {
	a := generateContentFingerprint([]string{"b:1", "a:1"}, "linux/amd64")
	b := generateContentFingerprint([]string{"a:1", "b:1"}, "linux/amd64")
	c := generateContentFingerprint([]string{"a:1", "b:1"}, "linux/arm64")
	if a != b || a == c {
		t.Fatalf("unexpected fingerprints: %q %q %q", a, b, c)
	}
}

func TestResolveImageRef(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("query preserves underscores", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/api/image/download?image=user/my_app:v1", nil)
		if got := resolveImageRef(c); got != "user/my_app:v1" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("missing image is empty", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/api/image/download", nil)
		if got := resolveImageRef(c); got != "" {
			t.Fatalf("got %q", got)
		}
	})
}

func TestWriteDownloadErrorSkipsJSONAfterBodyStarted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("before write returns json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		writeDownloadError(c, errors.New("boom"), "镜像下载失败")
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "镜像下载失败") || !strings.Contains(body, "boom") {
			t.Fatalf("body = %q", body)
		}
	})

	t.Run("after write skips json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		if _, err := c.Writer.Write([]byte("tar-bytes")); err != nil {
			t.Fatal(err)
		}
		writeDownloadError(c, errors.New("boom"), "镜像下载失败")
		if got := w.Body.String(); got != "tar-bytes" {
			t.Fatalf("body corrupted: %q", got)
		}
	})
}
