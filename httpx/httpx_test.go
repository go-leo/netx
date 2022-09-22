package httpx

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	request, err := new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/get").
		Query("name", "jack").
		AddQuery("name", "zhaoyi").
		QueryString("class=322&id=1").
		Queries(url.Values{"money": []string{"12345"}}).
		RemoveQuery("money").
		Build(context.Background())
	assert.NoError(t, err)
	resp, err := PooledClient().Do(request)
	assert.NoError(t, err)
	body, err := ResponseHelper{resp: resp}.TextBody()
	assert.NoError(t, err)
	t.Log(body)
	t.Log(resp.Header)
}

func TestHeader(t *testing.T) {
	request, err := new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/headers").
		QueryString("class=322&id=1").
		Query("name", "jack").
		Header("id", "12345").
		AddHeader("id", "45678").
		Headers(http.Header{"token": []string{"this is token"}}).
		RemoveHeader("token").
		UserAgent("leo-httpx").
		Build(context.Background())
	assert.NoError(t, err)
	resp, err := PooledClient().Do(request)
	assert.NoError(t, err)
	body, err := ResponseHelper{resp: resp}.TextBody()
	assert.NoError(t, err)
	t.Log(body)
	t.Log(resp.Header)
}

func TestBasicAuth(t *testing.T) {
	request, err := new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/basic-auth").
		QueryString("class=322&id=1").
		Query("name", "jack").
		Header("id", "12345").
		AddHeader("id", "45678").
		Headers(http.Header{"token": []string{"3344555"}}).
		BasicAuth("basic", "auth").
		Build(context.Background())
	assert.NoError(t, err)
	resp, err := PooledClient().Do(request)
	assert.NoError(t, err)
	body, err := ResponseHelper{resp: resp}.TextBody()
	assert.NoError(t, err)
	t.Log(body)
	t.Log(resp.Header)
}

func TestBearerAuth(t *testing.T) {
	request, err := new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/bearer").
		QueryString("class=322&id=1").
		Query("name", "jack").
		Header("id", "12345").
		AddHeader("id", "45678").
		Headers(http.Header{"token": []string{"3344555"}}).
		BearerAuth("this is bearer auth").
		Build(context.Background())
	assert.NoError(t, err)
	resp, err := PooledClient().Do(request)
	assert.NoError(t, err)
	body, err := ResponseHelper{resp: resp}.TextBody()
	assert.NoError(t, err)
	t.Log(body)
	t.Log(resp.Header)
}

func TestCache(t *testing.T) {
	request, err := new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/cache").
		Build(context.Background())
	assert.NoError(t, err)
	resp, err := PooledClient().Do(request)
	assert.NoError(t, err)
	helper := ResponseHelper{resp: resp}
	statusCode := helper.StatusCode()
	t.Log(statusCode)
	lastModified := helper.LastModified()
	t.Log(lastModified)
	etag := helper.Etag()
	t.Log(etag)

	request, err = new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/cache").
		IfModifiedSince(lastModified).
		Build(context.Background())
	assert.NoError(t, err)
	resp, err = PooledClient().Do(request)
	assert.NoError(t, err)
	helper = ResponseHelper{resp: resp}
	statusCode = helper.StatusCode()
	t.Log(statusCode)

	request, err = new(RequestBuilder).
		Get().
		URLString("http://httpbin.org/cache").
		IfNoneMatch(etag).
		Build(context.Background())
	assert.NoError(t, err)
	resp, err = PooledClient().Do(request)
	assert.NoError(t, err)
	helper = ResponseHelper{resp: resp}
	statusCode = helper.StatusCode()
	t.Log(statusCode)

}
