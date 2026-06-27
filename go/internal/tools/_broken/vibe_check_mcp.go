package tools

import (
	"context"
	"fmt"
	"strings"
)

func HandleVibeCheck(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	apiErr := checkArgs(args, "text")
	if apiErr != nil {
		return err(apiErr.Error())
}

	text, _ :=getString(args, "text")
	if strings.TrimSpace(text) == "" {
		return err("text cannot be empty")
}

	vibe := analyzeVibe(text)
	if vibe == "positive" {
		return ok("The vibe is good! 👍")
	} else if vibe == "negative" {
		return ok("The vibe is bad. 👎")
	} else {
		return ok("The vibe is neutral. 🤷")

}

}

func HandleVibeScore(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	apiErr := checkArgs(args, "text")
	if apiErr != nil {
		return err(apiErr.Error())
}

	text, _ :=getString(args, "text")
	if strings.TrimSpace(text) == "" {
		return err("text cannot be empty")
}

	score := calculateVibeScore(text)
	return ok(fmt.Sprintf("Vibe score: %d/100", score))
}

func HandleVibeCompare(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	apiErr := checkArgs(args, "text1", "text2")
	if apiErr != nil {
		return err(apiErr.Error())
}

	text1, _ :=getString(args, "text1")
	text2, _ :=getString(args, "text2")

	if strings.TrimSpace(text1) == "" || strings.TrimSpace(text2) == "" {
		return err("text1 and text2 cannot be empty")
}

	score1 := calculateVibeScore(text1)
	score2 := calculateVibeScore(text2)

	if score1 > score2 {
		return ok(fmt.Sprintf("Text1 (%d) has better vibes than Text2 (%d)", score1, score2))
	} else if score1 < score2 {
		return ok(fmt.Sprintf("Text2 (%d) has better vibes than Text1 (%d)", score2, score1))
	} else {
		return ok("Both texts have the same vibe score")

}

}

func checkArgs(args map[string]interface{}, required ...string) error {
	for _, key := range required {
		if _, exists := args[key]; !exists {
			return fmt.Errorf("missing required argument: %s", key)

	}
	return nil
}

}

func analyzeVibe(text string) string {
	positiveWords := []string{"good", "great", "awesome", "happy", "joy", "love"}
	negativeWords := []string{"bad", "terrible", "awful", "sad", "hate", "angry"}

	textLower := strings.ToLower(text)

	for _, word := range positiveWords {
		if strings.Contains(textLower, word) {
			return "positive"
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(textLower, word) {
			return "negative"
		}
	}

	return "neutral"
}

func calculateVibeScore(text string) int {
	positiveWords := []string{"good", "great", "awesome", "happy", "joy", "love"}
	negativeWords := []string{"bad", "terrible", "awful", "sad", "hate", "angry"}

	textLower := strings.ToLower(text)
	score := 0

	for _, word := range positiveWords {
		if strings.Contains(textLower, word) {
			score += 10
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(textLower, word) {
			score -= 10
		}
	}

	// Ensure score is within 0-100 range
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return score
}