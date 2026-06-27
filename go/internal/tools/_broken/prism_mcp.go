package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	http.DefaultClient = http.DefaultClient
	uuidRegex  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

func HandlePrismStatus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ :=getString(args, "host")
	if host == "" {
		return err("host parameter is required")
}

	u := fmt.Sprintf("http://%s:9000/api/v1/status", host)
	req, reqErr := http.NewRequestWithContext(ctx, "GET", u, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, apiErr := http.DefaultClient.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	version, _ :=getString(result, "version")
	if version == "" {
		return err("version not found in response")
}

	return ok(fmt.Sprintf("Prism MCP running version %s", version))
}

func HandlePrismClusterList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ :=getString(args, "host")
	if host == "" {
		return err("host parameter is required")
}

	u := fmt.Sprintf("http://%s:9000/api/nutanix/v3/clusters/list", host)
	reqBody := map[string]interface{}{"kind": "cluster"}
	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(string(jsonBody)))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := http.DefaultClient.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	entities := getSlice(result, "entities")
	if entities == nil {
		return ok("No clusters found")
}

	var clusterNames []string
	for _, entity := range entities {
		if cluster, found := entity.(map[string]interface{}); found {
			if name := getString(cluster, "name"); name != "" {
				clusterNames = append(clusterNames, name)

		}
	}

	sort.Strings(clusterNames)
	return ok(fmt.Sprintf("Clusters: %s", strings.Join(clusterNames, ", ")))
}

}

func HandlePrismVmList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ :=getString(args, "host")
	if host == "" {
		return err("host parameter is required")
}

	clusterUUID, _ :=getString(args, "cluster_uuid")
	if clusterUUID == "" || !uuidRegex.MatchString(clusterUUID) {
		return err("valid cluster_uuid parameter is required")
}

	u := fmt.Sprintf("http://%s:9000/api/nutanix/v3/vms/list", host)
	reqBody := map[string]interface{}{
		"kind": "vm",
		"filter": fmt.Sprintf("cluster_reference==%s", clusterUUID),
	}
	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(string(jsonBody)))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := http.DefaultClient.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	entities := getSlice(result, "entities")
	if entities == nil {
		return ok("No VMs found")
}

	var vmNames []string
	for _, entity := range entities {
		if vm, found := entity.(map[string]interface{}); found {
			if name := getString(vm, "name"); name != "" {
				vmNames = append(vmNames, name)

		}
	}

	sort.Strings(vmNames)
	return ok(fmt.Sprintf("VMs: %s", strings.Join(vmNames, ", ")))
}

}

func HandlePrismTaskList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ :=getString(args, "host")
	if host == "" {
		return err("host parameter is required")
}

	limit, _ :=getInt(args, "limit")
	if limit <= 0 {
		limit = 10
	}

	u := fmt.Sprintf("http://%s:9000/api/nutanix/v3/tasks/list", host)
	reqBody := map[string]interface{}{
		"kind":  "task",
		"limit": limit,
	}
	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(string(jsonBody)))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := http.DefaultClient.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	entities := getSlice(result, "entities")
	if entities == nil {
		return ok("No tasks found")
}

	var taskInfo []string
	for _, entity := range entities {
		if task, found := entity.(map[string]interface{}); found {
			taskUUID, _ :=getString(task, "metadata.uuid")
			status, _ :=getString(task, "status")
			operation, _ :=getString(task, "operation_type")
			if taskUUID != "" && status != "" {
				taskInfo = append(taskInfo, fmt.Sprintf("%s (%s): %s", taskUUID, status, operation))

		}
	}

	return ok(fmt.Sprintf("Recent tasks:\n%s", strings.Join(taskInfo, "\n")))
}

}

func HandlePrismAlertList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ :=getString(args, "host")
	if host == "" {
		return err("host parameter is required")
}

	severity, _ :=getString(args, "severity")
	if severity == "" {
		severity = "WARNING,ERROR,CRITICAL"
	}

	u := fmt.Sprintf("http://%s:9000/api/nutanix/v3/alerts/list", host)
	reqBody := map[string]interface{}{
		"kind": "alert",
		"filter": fmt.Sprintf("severity==%s", severity),
	}
	jsonBody, jsonErr := json.Marshal(reqBody)
	if jsonErr != nil {
		return err(fmt.Sprintf("failed to marshal request body: %v", jsonErr))
}

	req, reqErr := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(string(jsonBody)))
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	req.Header.Set("Content-Type", "application/json")

	resp, apiErr := http.DefaultClient.Do(req)
	if apiErr != nil {
		return err(fmt.Sprintf("request failed: %v", apiErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response: %v", readErr))
}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("failed to parse response: %v", parseErr))
}

	entities := getSlice(result, "entities")
	if entities == nil {
		return ok("No alerts found")
}

	var alertInfo []string
	for _, entity := range entities {
		if alert, found := entity.(map[string]interface{}); found {
			alertUUID, _ :=getString(alert, "metadata.uuid")
			message, _ :=getString(alert, "message")
			severity, _ :=getString(alert, "severity")
			if alertUUID != "" && message != "" {
				alertInfo = append(alertInfo, fmt.Sprintf("%s [%s]: %s", alertUUID, severity, message))

		}
	}

	return ok(fmt.Sprintf("Active alerts:\n%s", strings.Join(alertInfo, "\n")))
}

}

// Helper function to safely get a slice from map
func getSlice(m map[string]interface{}, key string) []interface{} {
	if val, found := m[key]; found {
		if slice, found := val.([]interface{}); found {
			return slice
		}
	}
	return nil
}