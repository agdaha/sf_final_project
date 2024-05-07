package reader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRss(t *testing.T) {
	// Dummy JSON response
	expected := `<rss version="2.0">
    <channel>
        <title>Golang Weekly</title>
        <description>
			A weekly newsletter about the Go programming language
        </description>
        <link>https://golangweekly.com/</link>
        <item>
            <title>Going supersonic</title>
            <link>https://golangweekly.com/issues/505</link>
            <description>
				D1
            </description>
            <pubDate>Tue, 30 Apr 2024 00:00:00 +0000</pubDate>
            <guid>https://golangweekly.com/issues/505</guid>
        </item>

    </channel>
</rss>`
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	defer svr.Close()
	c, err := GetRss(svr.URL)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(c.Chanel.Items) != 1 {
		t.Errorf(" wrong len items got:%v want:1", len(c.Chanel.Items))
	}
}
