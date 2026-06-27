package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getChromaBaseURL() string {
	base := osGetEnv("CHROMA_BASE_URL", "http://localhost:8000")
	return strings.TrimRight(base, "/")
}

func osGetEnv(key, fallback string) string {
	if v, found := osLookupEnv(key); found {
		return v
	}
	return fallback
}

func osLookupEnv(key string) (string, bool) {
	return "", false
}

func chromaRequest(method, path string, body interface{}) (map[string]interface{}, error) {
	baseURL := getChromaBaseURL()
	fullURL := baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, jsonErr := json.Marshal(body)
		if jsonErr != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", jsonErr)
}

		reqBody = bytes.NewReader(jsonData)

	req, reqErr := http.NewRequest(method, fullURL, reqBody)
	if reqErr != nil {
		return nil, fmt.Errorf("failed to create request: %w", reqErr)
}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, doErr := client.Do(req)
	if doErr != nil {
		return nil, fmt.Errorf("request failed: %w", doErr)
}

	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response: %w", readErr)
}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("chroma API error (status %d): %s", resp.StatusCode, string(respBody))
}

	var result map[string]interface{}
	if len(respBody) > 0 {
		if unmarshalErr := json.Unmarshal(respBody, &result); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to parse response: %w", unmarshalErr)

	}

	return result, nil
}

}
}

// HandleListCollections lists all collections in the Chroma knowledge base
func HandleListCollections(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	limit, _ :=getInt(args, "limit")
	offset, _ :=getInt(args, "offset")

	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))

	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))

	path := "/api/v1/collections"
	if len(params) > 0 {
		path += "?" + params.Encode()

	result, apiErr := chromaRequest("GET", path, nil)
	if apiErr != nil {
		return err(apiErr.Error())
}

	jsonBytes, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format result: %v", marshalErr))
}

	return ok(string(jsonBytes))
}

}
}
}

// HandleCreateCollection creates a new collection in the Chroma knowledge base
func HandleCreateCollection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("name is required")
}

	metadataRaw, _ :=getString(args, "metadata")
	getOrCreate, _ :=getBool(args, "get_or_create")

	body := map[string]interface{}{
		"name": name,
	}

	if metadataRaw != "" {
		var metadata map[string]interface{}
		if parseErr := json.Unmarshal([]byte(metadataRaw), &metadata); parseErr != nil {
			return err(fmt.Sprintf("invalid metadata JSON: %v", parseErr))
}

		body["metadata"] = metadata
	}

	if getOrCreate {
		body["get_or_create"] = true
	}

	result, apiErr := chromaRequest("POST", "/api/v1/collections", body)
	if apiErr != nil {
		return err(apiErr.Error())
}

	jsonBytes, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format result: %v", marshalErr))
}

	return ok(string(jsonBytes))
}

// HandleAddDocuments adds documents with embeddings to a collection
func HandleAddDocuments(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	collectionID, _ :=getString(args, "collection_id")
	if collectionID == "" {
		return err("collection_id is required")
}

	documentsRaw, _ :=getString(args, "documents")
	if documentsRaw == "" {
		return err("documents is required (JSON array of strings)")
}

	var documents []string
	if parseErr := json.Unmarshal([]byte(documentsRaw), &documents); parseErr != nil {
		return err(fmt.Sprintf("invalid documents JSON: %v", parseErr))
}

	body := map[string]interface{}{
		"documents": documents,
	}

	idsRaw, _ :=getString(args, "ids")
	if idsRaw != "" {
		var ids []string
		if parseErr := json.Unmarshal([]byte(idsRaw), &ids); parseErr != nil {
			return err(fmt.Sprintf("invalid ids JSON: %v", parseErr))
}

		body["ids"] = ids
	} else {
		ids := make([]string, len(documents))
		for i := range documents {
			ids[i] = fmt.Sprintf("doc_%d_%d", time.Now().UnixNano(), i)

		body["ids"] = ids
	}

	embeddingsRaw, _ :=getString(args, "embeddings")
	if embeddingsRaw != "" {
		var embeddings [][]float64
		if parseErr := json.Unmarshal([]byte(embeddingsRaw), &embeddings); parseErr != nil {
			return err(fmt.Sprintf("invalid embeddings JSON: %v", parseErr))
}

		body["embeddings"] = embeddings
	}

	metadatasRaw, _ :=getString(args, "metadatas")
	if metadatasRaw != "" {
		var metadatas []map[string]interface{}
		if parseErr := json.Unmarshal([]byte(metadatasRaw), &metadatas); parseErr != nil {
			return err(fmt.Sprintf("invalid metadatas JSON: %v", parseErr))
}

		body["metadatas"] = metadatas
	}

	path := "/api/v1/collections/" + url.PathEscape(collectionID) + "/add"
	result, apiErr := chromaRequest("POST", path, body)
	if apiErr != nil {
		return err(apiErr.Error())
}

	jsonBytes, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format result: %v", marshalErr))
}

	return ok(string(jsonBytes))
}

}

// HandleQueryCollection queries a collection with query texts or embeddings
func HandleQueryCollection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	collectionID, _ :=getString(args, "collection_id")
	if collectionID == "" {
		return err("collection_id is required")
}

	queryTextsRaw, _ :=getString(args, "query_texts")
	queryEmbeddingsRaw, _ :=getString(args, "query_embeddings")

	if queryTextsRaw == "" && queryEmbeddingsRaw == "" {
		return err("either query_texts or query_embeddings is required")
}

	body := map[string]interface{}{}

	if queryTextsRaw != "" {
		var queryTexts []string
		if parseErr := json.Unmarshal([]byte(queryTextsRaw), &queryTexts); parseErr != nil {
			return err(fmt.Sprintf("invalid query_texts JSON: %v", parseErr))
}

		body["query_texts"] = queryTexts
	}

	if queryEmbeddingsRaw != "" {
		var queryEmbeddings [][]float64
		if parseErr := json.Unmarshal([]byte(queryEmbeddingsRaw), &queryEmbeddings); parseErr != nil {
			return err(fmt.Sprintf("invalid query_embeddings JSON: %v", parseErr))
}

		body["query_embeddings"] = queryEmbeddings
	}

	nResults, _ :=getInt(args, "n_results")
	if nResults > 0 {
		body["n_results"] = nResults
	}

	whereRaw, _ :=getString(args, "where")
	if whereRaw != "" {
		var where map[string]interface{}
		if parseErr := json.Unmarshal([]byte(whereRaw), &where); parseErr != nil {
			return err(fmt.Sprintf("invalid where JSON: %v", parseErr))
}

		body["where"] = where
	}

	whereDocumentRaw, _ :=getString(args, "where_document")
	if whereDocumentRaw != "" {
		var whereDocument map[string]interface{}
		if parseErr := json.Unmarshal([]byte(whereDocumentRaw), &whereDocument); parseErr != nil {
			return err(fmt.Sprintf("invalid where_document JSON: %v", parseErr))
}

		body["where_document"] = whereDocument
	}

	includeRaw, _ :=getString(args, "include")
	if includeRaw != "" {
		var include []string
		if parseErr := json.Unmarshal([]byte(includeRaw), &include); parseErr != nil {
			return err(fmt.Sprintf("invalid include JSON: %v", parseErr))
}

		body["include"] = include
	}

	path := "/api/v1/collections/" + url.PathEscape(collectionID) + "/query"
	result, apiErr := chromaRequest("POST", path, body)
	if apiErr != nil {
		return err(apiErr.Error())
}

	jsonBytes, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format result: %v", marshalErr))
}

	return ok(string(jsonBytes))
}

// HandleDeleteCollection deletes a collection from the Chroma knowledge base
func HandleDeleteCollection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	collectionName, _ :=getString(args, "collection_name")
	if collectionName == "" {
		return err("collection_name is required")
}

	path := "/api/v1/collections/" + url.PathEscape(collectionName)
	result, apiErr := chromaRequest("DELETE", path, nil)
	if apiErr != nil {
		return err(apiErr.Error())
}

	jsonBytes, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		return err(fmt.Sprintf("failed to format result: %v", marshalErr))
}

	return ok(string(jsonBytes))
}