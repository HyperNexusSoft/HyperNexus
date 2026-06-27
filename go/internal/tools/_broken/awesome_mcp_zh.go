package tools

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

func HandleSubmitResource(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	description, _ :=getString(args, "description")
	link, _ :=getString(args, "link")
	category, _ :=getString(args, "category")

	if name == "" || description == "" || link == "" || category == "" {
		return err("missing required parameters: name, description, link, or category")
}

	// Validate URL format
	_, parseErr := url.ParseRequestURI(link)
	if parseErr != nil {
		return err("invalid URL format")
}

	// Simulate submission to Awesome-MCP-ZH
	// In a real implementation, this would be an API call to the repository
	response := fmt.Sprintf("Submitted resource:\nName: %s\nDescription: %s\nLink: %s\nCategory: %s",
		name, description, link, category)

	return ok(response)
}

func HandleCheckResource(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	link, _ :=getString(args, "link")
	if link == "" {
		return err("missing required parameter: link")
}

	// Validate URL format
	_, parseErr := url.ParseRequestURI(link)
	if parseErr != nil {
		return err("invalid URL format")
}

	// Simulate checking if resource exists in Awesome-MCP-ZH
	// In a real implementation, this would check the repository contents
	exists := strings.Contains(link, "github.com") // Simple check for example

	if exists {
		return ok("Resource already exists in Awesome-MCP-ZH")
	}

	return ok("Resource does not exist in Awesome-MCP-ZH")
}

func HandleGetCategories(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Simulate getting available categories from Awesome-MCP-ZH
	// In a real implementation, this would fetch from the repository
	categories := []string{
		"🔍 搜索",
		"📚 知识库",
		"🤖 工具",
		"📊 数据",
		"📝 文档",
		"🛠️ 开发",
	}

	response := "Available categories:\n" + strings.Join(categories, "\n")
	return ok(response)
}

func HandleValidateContribution(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	description, _ :=getString(args, "description")
	link, _ :=getString(args, "link")
	category, _ :=getString(args, "category")

	if name == "" || description == "" || link == "" || category == "" {
		return err("missing required parameters: name, description, link, or category")
}

	// Validate URL format
	_, parseErr := url.ParseRequestURI(link)
	if parseErr != nil {
		return err("invalid URL format")
}

	// Simulate validation against Awesome-MCP-ZH standards
	// In a real implementation, this would check against the contributing guidelines
	valid := true
	reasons := []string{}

	if len(name) < 3 {
		valid = false
		reasons = append(reasons, "name too short")

	if len(description) < 10 {
		valid = false
		reasons = append(reasons, "description too short")

	if !strings.Contains(link, "github.com") && !strings.Contains(link, "gitlab.com") {
		valid = false
		reasons = append(reasons, "link should be to a Git repository")

	if valid {
		return ok("Contribution meets Awesome-MCP-ZH standards")
	}

	return ok("Contribution does not meet standards:\n" + strings.Join(reasons, "\n"))
}

}
}
}

func HandleGetContributionGuidelines(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Simulate getting contribution guidelines from Awesome-MCP-ZH
	// In a real implementation, this would fetch from the repository
	guidelines := `Awesome-MCP-ZH Contribution Guidelines:
1. Resources must be real and verifiable MCP implementations
2. Provide clear documentation in Chinese
3. Include installation or usage instructions
4. Follow the format standards in CONTRIBUTING.md
5. Avoid duplicate entries
6. Ensure the resource is valuable for Chinese users`

	return ok(guidelines)
}

func HandleCheckResourceFormat(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	markdown, _ :=getString(args, "markdown")
	if markdown == "" {
		return err("missing required parameter: markdown")
}

	// Simulate checking if markdown follows Awesome-MCP-ZH format
	// In a real implementation, this would parse the markdown structure
	lines := strings.Split(markdown, "\n")
	if len(lines) < 3 {
		return ok("Markdown too short to validate format")
	}

	// Check for table format
	if !strings.Contains(lines[0], "|") || !strings.Contains(lines[1], "|") {
		return ok("Markdown does not follow table format")
	}

	return ok("Markdown follows Awesome-MCP-ZH format")
}