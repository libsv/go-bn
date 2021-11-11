package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/libsv/go-bn/models"
	"github.com/stretchr/testify/assert"
)

type closeFunc func()

// nolint: revive // test code
// TestServer creates a test server for testing.
func TestServer(t *testing.T, expReq *models.Request, testFile string) (*httptest.Server, closeFunc) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, *expReq, req)

		// nolint:gosec // test code
		response, err := ioutil.ReadFile(path.Join("./testing/data", testFile+".json"))
		assert.NoError(t, err)

		mm := map[string]interface{}{}
		assert.NoError(t, json.Unmarshal(response, &mm))
		bb, err := json.Marshal(mm)
		assert.NoError(t, err)
		_, _ = w.Write(bb)
	}))

	return svr, svr.Close
}
