// internal/finalizer/retry_integration_test.go
package finalizer

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/pubsub"
	"github.com/stretchr/testify/require"
)

type flakyHandler struct {
	calls   int
	succeed int
}

func (f *flakyHandler) Process(_ context.Context, _ *pubsub.Message) error {
	f.calls++
	if f.calls < f.succeed {
		return fmt.Errorf("transient error #%d", f.calls)
	}
	return nil
}

func TestRetryIntegration(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		callCount++
		if callCount < 3 {
			stdhttp.Error(w, "temporary", stdhttp.StatusInternalServerError)
		} else {
			w.WriteHeader(stdhttp.StatusOK)
		}
	}))
	defer ts.Close()

	_, err := commonsHttp.NewHTTPClient("retry-test", ts.URL, &stdhttp.Transport{})
	require.NoError(t, err)

	base := &flakyHandler{calls: 0, succeed: 3}
	rh := NewRetryHandler(base, 3, 10*time.Millisecond)

	err = rh.Process(context.Background(), &pubsub.Message{Value: []byte("{}")})
	require.NoError(t, err)
	require.Equal(t, 3, base.calls)
}
