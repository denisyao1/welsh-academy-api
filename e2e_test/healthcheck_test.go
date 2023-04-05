package e2etest

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	url := "/api/v1/health"
	req := httptest.NewRequest(GetMethod, url, nil)
	resp, _ := App.Test(req, -1)
	assert.Equal(200, resp.StatusCode)
	msg, _ := io.ReadAll(resp.Body)
	m := map[string]string{}
	json.Unmarshal(msg, &m)
	assert.Equal("Welsh Academy Api is running.", m["Message"], "Health message is not correct")

}
