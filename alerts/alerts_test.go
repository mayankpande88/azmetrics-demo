package alerts

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// fakeCred satisfies azcore.TokenCredential without contacting Azure AD.
type fakeCred struct{}

func (fakeCred) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

// mockTransport returns a canned 200 for any request. Like most unit tests, it
// never reaches real ARM, so it CANNOT observe the api-version the client sends.
// This is exactly why the v0.13.0 api-version regression passes CI: the test that
// "covers" ListAlertNames is blind to the query parameter that breaks in production.
type mockTransport struct{}

func (mockTransport) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"value":[]}`)),
		Request:    req,
	}, nil
}

func TestListAlertNames_GreenCI(t *testing.T) {
	opts := &arm.ClientOptions{}
	opts.Transport = mockTransport{}

	names, err := ListAlertNames(context.Background(), "sub-123", fakeCred{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(names))
	}
}
