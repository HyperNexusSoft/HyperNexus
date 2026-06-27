package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var http.DefaultClient = http.Client{Timeout: 30 * time.Second}

func HandlePalInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	version, _ :=getString(args, "version")
	if version == "" {
		version = "latest"
	}

	reqURL := fmt.Sprintf("https://api.palworldgame.com/info?version=%s", url.QueryEscape(version))

	req, reqErr := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to fetch pal info: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result map[string]interface{}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON: %v", parseErr))
}

	output, formatErr := json.MarshalIndent(result, "", "  ")
	if formatErr != nil {
		return err(fmt.Sprintf("failed to format output: %v", formatErr))
}

	return ok(string(output))
}

func HandlePalSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query parameter is required")
}

	page, _ :=getInt(args, "page")
	if page == 0 {
		page = 1
	}

	reqURL := fmt.Sprintf("https://api.palworldgame.com/search?q=%s&page=%d",
		url.QueryEscape(query),
		page)

	req, reqErr := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("failed to create request: %v", reqErr))
}

	resp, fetchErr := http.DefaultClient.Do(req)
	if fetchErr != nil {
		return err(fmt.Sprintf("failed to search pals: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("failed to read response body: %v", readErr))
}

	var result struct {
		Results []map[string]interface{} `json:"results"`
		Total   int                     `json:"total"`
	}
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return err(fmt.Sprintf("failed to parse JSON: %v", parseErr))
}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Found %d results (page %d):\n", result.Total, page))
	for i, item := range result.Results {
		output.WriteString(fmt.Sprintf("\nResult %d:\n", i+1))
		for k, v := range item {
			output.WriteString(fmt.Sprintf("  %s: %v\n", k, v))

	}

	return ok(output.String())
}

}

func HandlePalValidateName(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("name parameter is required")
}

	// Palworld naming rules
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,20}$`)
	if !re.MatchString(name) {
		return ok("Name is invalid. Must be 3-20 characters and only contain letters, numbers, underscores, and hyphens.")
}

	return ok("Name is valid according to Palworld naming rules.")
}

func HandlePalCalculateStats(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	level, _ :=getInt(args, "level")
	if level < 1 || level > 50 {
		return err("level must be between 1 and 50")
}

	baseHP, _ :=getInt(args, "base_hp")
	baseAttack, _ :=getInt(args, "base_attack")
	baseDefense, _ :=getInt(args, "base_defense")

	if baseHP == 0 || baseAttack == 0 || baseDefense == 0 {
		return err("base stats (hp, attack, defense) are required")
}

	// Simplified Palworld stat calculation
	hp := baseHP + (level-1)*2
	attack := baseAttack + (level-1)
	defense := baseDefense + (level-1)/2

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Stats for level %d Pal:\n", level))
	output.WriteString(fmt.Sprintf("HP: %d\n", hp))
	output.WriteString(fmt.Sprintf("Attack: %d\n", attack))
	output.WriteString(fmt.Sprintf("Defense: %d\n", defense))

	return ok(output.String())
}

func HandlePalElementAdvantages(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	element := strings.ToLower(getString(args, "element"))
	if element == "" {
		return err("element parameter is required")
}

	advantages := map[string][]string{
		"fire":     {"grass", "ice", "bug"},
		"water":    {"fire", "ground", "rock"},
		"grass":    {"water", "ground", "rock"},
		"electric": {"water", "flying"},
		"ground":   {"fire", "electric", "poison", "rock"},
		"ice":      {"grass", "ground", "flying", "dragon"},
		"dragon":   {"dragon"},
		"dark":     {"ghost", "psychic"},
		"fairy":    {"fighting", "dark", "dragon"},
		"normal":   {},
	}

	weakTo := map[string][]string{
		"fire":     {"water", "ground", "rock"},
		"water":    {"electric", "grass"},
		"grass":    {"fire", "ice", "poison", "flying", "bug"},
		"electric": {"ground"},
		"ground":   {"water", "grass", "ice"},
		"ice":      {"fire", "fighting", "rock", "steel"},
		"dragon":   {"ice", "dragon", "fairy"},
		"dark":     {"fighting", "fairy", "bug"},
		"fairy":    {"poison", "steel"},
		"normal":   {"fighting"},
	}

	if _, exists := advantages[element]; !exists {
		return err(fmt.Sprintf("unknown element: %s", element))
}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Element advantages for %s:\n", element))

	if len(advantages[element]) > 0 {
		output.WriteString("Strong against: ")
		output.WriteString(strings.Join(advantages[element], ", "))
		output.WriteString("\n")
	} else {
		output.WriteString("No particular advantages\n")

	if len(weakTo[element]) > 0 {
		output.WriteString("Weak to: ")
		output.WriteString(strings.Join(weakTo[element], ", "))
		output.WriteString("\n")
	} else {
		output.WriteString("No particular weaknesses\n")

	return ok(output.String())
}
}
}