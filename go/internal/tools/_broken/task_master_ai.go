package tools

import (
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

type task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	ParentID    string `json:"parent_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var taskStore = struct {
	tasks map[string]task
}{
	tasks: make(map[string]task),
}

func init() {
	taskStore.tasks["1"] = task{
		ID: "1", Title: "Initialize Project", Description: "Set up the project structure and dependencies",
		Status: "done", Priority: "high", ParentID: "", CreatedAt: time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
		UpdatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
	}
	taskStore.tasks["2"] = task{
		ID: "2", Title: "Design Architecture", Description: "Design the system architecture and data models",
		Status: "in-progress", Priority: "high", ParentID: "1", CreatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	taskStore.tasks["3"] = task{
		ID: "3", Title: "Implement Core Features", Description: "Build the core application features",
		Status: "pending", Priority: "medium", ParentID: "2", CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func nextTaskID() string {
	maxID := 0
	for k := range taskStore.tasks {
		idNum, _ := strconv.Atoi(k)
		if idNum > maxID {
			maxID = idNum
		}
	}
	return strconv.Itoa(maxID + 1)
}

func HandleCreateTask(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	title, _ :=getString(args, "title")
	description, _ :=getString(args, "description")
	priority, _ :=getString(args, "priority")
	parentID, _ :=getString(args, "parent_id")

	if title == "" {
		return err("title is required")
}

	if priority == "" {
		priority = "medium"
	}
	validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
	if !validPriorities[priority] {
		return err("invalid priority: must be low, medium, high, or critical")
}

	if parentID != "" {
		if _, exists := taskStore.tasks[parentID]; !exists {
			return err("parent task not found: " + parentID)

	}

	now := time.Now().Format(time.RFC3339)
	newTask := task{
		ID: nextTaskID(), Title: title, Description: description,
		Status: "pending", Priority: priority, ParentID: parentID,
		CreatedAt: now, UpdatedAt: now,
	}
	taskStore.tasks[newTask.ID] = newTask

	result, marshalErr := json.MarshalIndent(newTask, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal task: " + marshalErr.Error())
}

	return ok("Task created successfully:\n" + string(result))
}

}

func HandleListTasks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	statusFilter, _ :=getString(args, "status")

	var filtered []task
	for _, t := range taskStore.tasks {
		if statusFilter == "" || strings.EqualFold(t.Status, statusFilter) {
			filtered = append(filtered, t)

	}

	if len(filtered) == 0 {
		return ok("No tasks found" + func() string {
}
			if statusFilter != "" {
				return " with status: " + statusFilter
			}
			return ""
		}())

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d task(s):\n\n", len(filtered)))
	for _, t := range filtered {
		sb.WriteString(fmt.Sprintf("ID: %s | Title: %s | Status: %s | Priority: %s\n", t.ID, t.Title, t.Status, t.Priority))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("  Description: %s\n", t.Description))

		if t.ParentID != "" {
			sb.WriteString(fmt.Sprintf("  Parent: %s\n", t.ParentID))

		sb.WriteString("\n")

	return ok(sb.String())
}

}
}
}
}

func HandleUpdateTask(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ :=getString(args, "id")
	status, _ :=getString(args, "status")
	priority, _ :=getString(args, "priority")

	if id == "" {
		return err("task id is required")
}

	existing, exists := taskStore.tasks[id]
	if !exists {
		return err("task not found: " + id)
}

	updated := existing
	changed := false

	if status != "" {
		validStatuses := map[string]bool{"pending": true, "in-progress": true, "done": true, "blocked": true, "cancelled": true}
		if !validStatuses[status] {
			return err("invalid status: must be pending, in-progress, done, blocked, or cancelled")
}

		updated.Status = status
		changed = true
	}

	if priority != "" {
		validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
		if !validPriorities[priority] {
			return err("invalid priority: must be low, medium, high, or critical")
}

		updated.Priority = priority
		changed = true
	}

	if !changed {
		return err("no updates provided; specify at least status or priority")
}

	updated.UpdatedAt = time.Now().Format(time.RFC3339)
	taskStore.tasks[id] = updated

	result, marshalErr := json.MarshalIndent(updated, "", "  ")
	if marshalErr != nil {
		return err("failed to marshal task: " + marshalErr.Error())
}

	return ok("Task updated successfully:\n" + string(result))
}

func HandleAnalyzeProject(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	projectName, _ :=getString(args, "project_name")
	complexity, _ :=getString(args, "complexity")

	if projectName == "" {
		return err("project_name is required")
}

	if complexity == "" {
		complexity = "medium"
	}

	validComplexities := map[string]bool{"low": true, "medium": true, "high": true}
	if !validComplexities[complexity] {
		return err("invalid complexity: must be low, medium, or high")
}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Project Analysis: %s\n", projectName))
	sb.WriteString(fmt.Sprintf("Complexity Level: %s\n\n", strings.Title(complexity)))
	sb.WriteString("Recommended Task Breakdown:\n")
	sb.WriteString("============================\n\n")

	phases := []struct {
		name   string
		tasks  []string
	}{
		{"Planning & Requirements", []string{"Define project scope", "Gather requirements", "Create technical specification"}},
		{"Design", []string{"System architecture design", "Database schema design", "API contract definition"}},
		{"Implementation", []string{"Core module development", "Integration layer", "Business logic implementation"}},
		{"Testing & QA", []string{"Unit test suite", "Integration testing", "Performance benchmarking"}},
		{"Deployment", []string{"CI/CD pipeline setup", "Production deployment", "Monitoring setup"}},
	}

	if complexity == "high" {
		phases[2].tasks = append(phases[2].tasks, "Advanced feature development", "Security hardening", "Scalability implementation")
		phases[3].tasks = append(phases[3].tasks, "Load testing", "Security audit", "Chaos engineering")
	} else if complexity == "low" {
		phases[2].tasks = phases[2].tasks[:2]
		phases[3].tasks = phases[3].tasks[:1]
	}

	totalTasks := 0
	for i, phase := range phases {
		sb.WriteString(fmt.Sprintf("Phase %d: %s\n", i+1, phase.name))
		for j, t := range phase.tasks {
			sb.WriteString(fmt.Sprintf("  %d.%d %s\n", i+1, j+1, t))
			totalTasks++
		}
		sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("Total recommended tasks: %d\n", totalTasks))
	sb.WriteString(fmt.Sprintf("Estimated timeline: %s\n", func() string {
		switch complexity {
		case "low":
			return "2-4 weeks"
}
		case "high":
			return "8-16 weeks"
}
		default:
			return "4-8 weeks"
		}
	}()))

	return ok(sb.String())

func HandleGetTask(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ :=getString(args, "id")
	if id == "" {
		return err("task id is required")
}

	t, exists := taskStore.tasks[id]
	if !exists {
		return err("task not found: " + id)
}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Task ID: %s\n", t.ID))
	sb.WriteString(fmt.Sprintf("Title: %s\n", t.Title))
	sb.WriteString(fmt.Sprintf("Description: %s\n", t.Description))
	sb.WriteString(fmt.Sprintf("Status: %s\n", t.Status))
	sb.WriteString(fmt.Sprintf("Priority: %s\n", t.Priority))
	if t.ParentID != "" {
		sb.WriteString(fmt.Sprintf("Parent Task: %s\n", t.ParentID))

	sb.WriteString(fmt.Sprintf("Created: %s\n", t.CreatedAt))
	sb.WriteString(fmt.Sprintf("Updated: %s\n", t.UpdatedAt))

	var subtasks []string
	for _, st := range taskStore.tasks {
		if st.ParentID == id {
			subtasks = append(subtasks, fmt.Sprintf("- [%s] %s (%s)", st.ID, st.Title, st.Status))

	}
	if len(subtasks) > 0 {
		sb.WriteString("\nSubtasks:\n")
		for _, st := range subtasks {
			sb.WriteString(st + "\n")

	}

	return ok(sb.String())
}

}
}
}

func HandleExpandTask(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ :=getString(args, "id")
	numSubtasks, _ :=getInt(args, "num_subtasks")

	if id == "" {
		return err("task id is required")
}

	parent, exists := taskStore.tasks[id]
	if !exists {
		return err("task not found: " + id)
}

	if numSubtasks <= 0 {
		numSubtasks = 3
	}
	if numSubtasks > 10 {
		numSubtasks = 10
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expanding task '%s' (ID: %s) into %d subtask(s):\n\n", parent.Title, id, numSubtasks))

	now := time.Now().Format(time.RFC3339)
	templates := []struct {
		suffix string
		desc   string
	}{
		{"Research and Analysis", "Research requirements and analyze dependencies for: " + parent.Title},
		{"Design and Planning", "Create detailed design and plan for: " + parent.Title},
		{"Core Implementation", "Implement the core functionality of: " + parent.Title},
		{"Integration", "Integrate with existing systems for: " + parent.Title},
		{"Testing", "Write and execute tests for: " + parent.Title},
		{"Documentation", "Document the implementation of: " + parent.Title},
		{"Review and Refinement", "Review and refine the implementation of: " + parent.Title},
		{"Performance Optimization", "Optimize performance for: " + parent.Title},
		{"Security Review", "Conduct security review for: " + parent.Title},
		{"Deployment Preparation", "Prepare deployment for: " + parent.Title},
	}

	createdCount := 0
	for i := 0; i < numSubtasks && i < len(templates); i++ {
		subID := nextTaskID()
		subTitle := fmt.Sprintf("%d. %s %s", i+1, parent.Title, templates[i].suffix)
		subTask := task{
			ID: subID, Title: subTitle, Description: templates[i].desc,
			Status: "pending", Priority: parent.Priority, ParentID: id,
			CreatedAt: now, UpdatedAt: now,
		}
		taskStore.tasks[subID] = subTask
		sb.WriteString(fmt.Sprintf("  Created subtask %s: %s\n", subID, subTitle))
		createdCount++
	}

	sb.WriteString(fmt.Sprintf("\nSuccessfully created %d subtask(s) for task %s.\n", createdCount, id))
	return ok(sb.String())
}

func HandleFetchTaskContext(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	urlStr, _ :=getString(args, "url")
	taskID, _ :=getString(args, "task_id")

	if urlStr == "" {
		return err("url is required")
}

	parsedURL, parseErr := url.Parse(urlStr)
	if parseErr != nil {
		return err("invalid URL: " + parseErr.Error())
}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return err("only http and https URLs are supported")
}

	client := http.DefaultClient
	req, reqErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if reqErr != nil {
		return err("failed to create request: " + reqErr.Error())
}

	resp, fetchErr := client.Do(req)
	if fetchErr != nil {
		return err("failed to fetch URL: " + fetchErr.Error())
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("HTTP request failed with status: %d", resp.StatusCode))
}

	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1024*100))
	if readErr != nil {
		return err("failed to read response: " + readErr.Error())
}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Fetched context from: %s\n", urlStr))
	if taskID != "" {
		if t, exists := taskStore.tasks[taskID]; exists {
			sb.WriteString(fmt.Sprintf("For task: %s (ID: %s)\n", t.Title, taskID))
		} else {
			sb.WriteString(fmt.Sprintf("Warning: task %s not found\n", taskID))

	}
	sb.WriteString(fmt.Sprintf("Content length: %d bytes\n\n", len(body)))

	content := string(body)
	if len(content) > 2000 {
		content = content[:2000] + "\n... (truncated)"
	}
	sb.WriteString(content)

	return ok(sb.String())
}
}