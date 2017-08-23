package webhook

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookRequest(t *testing.T) {
	realWebHookStr := "{\"update_id\":878863831,\n\"message\":{\"message_id\":4,\"from\":{\"id\":89922360,\"first_name\":\"Denis\",\"last_name\":\"Isaev\",\"language_code\":\"ru\"},\"chat\":{\"id\":89922360,\"first_name\":\"Denis\",\"last_name\":\"Isaev\",\"type\":\"private\"},\"date\":1503427617,\"text\":\"\\u043f\\u0440\\u0438\\u0432\"}}"

	req, err := http.NewRequest(http.MethodGet, apiPath,
		strings.NewReader(realWebHookStr))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handleHTTPRequest).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Zero(t, rr.Body.Len())
}
