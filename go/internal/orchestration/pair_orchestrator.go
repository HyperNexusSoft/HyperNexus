package orchestration

/**
 * @file pair_orchestrator.go
 * @module go/internal/orchestration
 *
 * WHAT: Multi-model pair programming and rotation system.
 * Coordinates multiple frontier models (Claude, GPT, Gemini) in a shared context.
 *
 * WHY: Cross-model consensus and rotating roles (Planner/Implementer/Tester) 
 * significantly reduce hallucinations and improve code quality.
 */

import (
	"context"
	"fmt"
	"strings"

	"github.com/borghq/borg-go/internal/ai"
)

type PairRole string

const (
	Planner     PairRole = "planner"
	Implementer PairRole = "implementer"
	Tester      PairRole = "tester"
	Critic      PairRole = "critic"
)

type SquadMember struct {
	Name     string   `json:"name"`
	Role     PairRole `json:"role"`
	Provider string   `json:"provider"`
	ModelID  string   `json:"modelId"`
}

type PairSessionResult struct {
	Success     bool     `json:"success"`
	History     []string `json:"history"`
	FinalOutput string   `json:"finalOutput"`
}

type PairOrchestrator struct {
	Squad   []SquadMember
	History []string
}

func NewPairOrchestrator() *PairOrchestrator {
	return &PairOrchestrator{
		History: []string{},
	}
}

func (p *PairOrchestrator) SetupFrontierSquad() {
	p.Squad = []SquadMember{
		{Name: "Claude (Architect)", Role: Planner, Provider: "anthropic", ModelID: "claude-3-5-sonnet-20241022"},
		{Name: "GPT (Engineer)", Role: Implementer, Provider: "openai", ModelID: "gpt-4o"},
		{Name: "Gemini (Reviewer)", Role: Tester, Provider: "google", ModelID: "gemini-1.5-pro"},
	}
}

func (p *PairOrchestrator) RotateRoles() {
	if len(p.Squad) < 2 {
		return
	}

	// [R1, R2, R3] -> [R2, R3, R1]
	firstRole := p.Squad[0].Role
	for i := 0; i < len(p.Squad)-1; i++ {
		p.Squad[i].Role = p.Squad[i+1].Role
	}
	p.Squad[len(p.Squad)-1].Role = firstRole

	fmt.Println("[PairOrchestrator] 🔄 Roles rotated:")
	for _, m := range p.Squad {
		fmt.Printf("  - %s: %s\n", m.Name, m.Role)
	}
}

func (p *PairOrchestrator) RunTask(ctx context.Context, task string) (*PairSessionResult, error) {
	fmt.Printf("[PairOrchestrator] 🚀 Starting Multi-Model Task: \"%s\"\n", task)
	p.History = []string{"USER: " + task}

	// 1. Planning Phase
	plan, err := p.executeTurn(ctx, Planner, "Create a detailed implementation plan for this task: "+task)
	if err != nil {
		return nil, err
	}
	p.History = append(p.History, fmt.Sprintf("PLANNER (%s): %s", p.getMemberName(Planner), plan))

	// 2. Review Phase
	feedback, err := p.executeTurn(ctx, Tester, "Review this plan and identify potential edge cases or bugs: "+plan)
	if err != nil {
		return nil, err
	}
	p.History = append(p.History, fmt.Sprintf("TESTER (%s): %s", p.getMemberName(Tester), feedback))

	// 3. Refinement
	finalPlan, err := p.executeTurn(ctx, Planner, "Refine the plan based on this feedback: "+feedback)
	if err != nil {
		return nil, err
	}
	p.History = append(p.History, fmt.Sprintf("PLANNER (%s): %s", p.getMemberName(Planner), finalPlan))

	// 4. Implementation
	implementation, err := p.executeTurn(ctx, Implementer, "Implement the final plan. Focus on correctness and performance. Plan: "+finalPlan)
	if err != nil {
		return nil, err
	}
	p.History = append(p.History, fmt.Sprintf("IMPLEMENTER (%s): %s", p.getMemberName(Implementer), implementation))

	// 5. Verification
	verification, err := p.executeTurn(ctx, Tester, "Verify the implementation against the plan and task requirements. Implementation: "+implementation)
	if err != nil {
		return nil, err
	}
	p.History = append(p.History, fmt.Sprintf("TESTER (%s): %s", p.getMemberName(Tester), verification))

	success := !strings.Contains(strings.ToLower(verification), "fail")

	return &PairSessionResult{
		Success:     success,
		History:     p.History,
		FinalOutput: implementation,
	}, nil
}

func (p *PairOrchestrator) executeTurn(ctx context.Context, role PairRole, prompt string) (string, error) {
	var member *SquadMember
	for i := range p.Squad {
		if p.Squad[i].Role == role {
			member = &p.Squad[i]
			break
		}
	}

	if member == nil {
		return "", fmt.Errorf("no member assigned to role: %s", role)
	}

	fmt.Printf("[PairOrchestrator] 👤 %s (%s) is thinking...\n", member.Name, member.Role)

	systemPrompt := fmt.Sprintf(`You are part of a multi-agent pair programming squad. 
Your name is %s. Your current role is %s.
Collaborate with your teammates to solve the task perfectly.

SQUAD ROLES:
- PLANNER: Breaks down the task and designs the solution.
- IMPLEMENTER: Writes the actual code and executes tools.
- TESTER: Identifies bugs, edge cases, and verifies correctness.`, member.Name, strings.ToUpper(string(member.Role)))

	fullHistory := strings.Join(p.History, "\n\n")
	turnPrompt := fmt.Sprintf("CONVERSATION HISTORY:\n%s\n\nCURRENT TURN (%s): %s", fullHistory, strings.ToUpper(string(member.Role)), prompt)

	resp, err := ai.AutoRouteWithModel(ctx, member.Provider+"/"+member.ModelID, []ai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: turnPrompt},
	})

	if err != nil {
		fmt.Printf("[PairOrchestrator] ⚠️ Turn failed for %s: %v\n", member.Name, err)
		return "", err
	}

	return resp.Content, nil
}

func (p *PairOrchestrator) getMemberName(role PairRole) string {
	for _, m := range p.Squad {
		if m.Role == role {
			return m.Name
		}
	}
	return "Unknown"
}
