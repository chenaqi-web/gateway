package tests

import (
	"context"
	"testing"

	httpclient "backend/gateway/internal/client/http"
	"backend/gateway/internal/config"
	"backend/gateway/internal/model/dto"
)

func TestSearchVectors(t *testing.T) {
	client := httpclient.NewHTTPClient(&config.Config{
		HTTP: config.HTTPConfig{
			AgentServerAddr: "http://127.0.0.1:8080",
			RequestTimeout:  30,
		},
	})

	resp, err := client.SearchVectors(context.Background(), "test_knowledge", &dto.VectorSearchRequest{
		Content: "如何重置密码",
		TopK:    new(5),
	})
	if err != nil {
		t.Fatalf("SearchVectors failed: %v", err)
	}

	t.Logf("code=%d msg=%s count=%d", resp.Code, resp.Msg, len(resp.Data))
	for i, item := range resp.Data {
		t.Logf("[%d] chunk_id=%s score=%v title=%s", i, item.ChunkID, item.Score, item.Title)
	}

	if len(resp.Data) == 0 {
		t.Fatal("expected search results, got empty data")
	}
}
