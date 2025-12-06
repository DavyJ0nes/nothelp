package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davyj0nes/nothelp/internal/config"
	"github.com/spf13/cobra"
)

const (
	defaultLLMModel = "gpt-4o-mini"
	baseLLMURL      = "https://api.openai.com/v1"
)

type dailyNote struct {
	Date string
	Body string
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float32       `json:"temperature,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func ReviewCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "review",
		Short: "generate a 7-day review using an LLM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return reviewRun(cmd, outputPath)
		},
	}

	cmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"",
		"destination path for the review JSON (default: <data>/reviews/review-<start>-to-<end>.json)",
	)

	return cmd
}

func reviewRun(cmd *cobra.Command, outputPath string) error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	notes, missingDates, err := collectNotes(conf, 7)
	if err != nil {
		return err
	}

	for _, date := range missingDates {
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: no note found for %s, skipping\n", date)
	}

	if len(notes) == 0 {
		return errors.New("no notes found for the requested period")
	}

	startDate := notes[0].Date
	endDate := notes[len(notes)-1].Date

	if outputPath == "" {
		outputPath = filepath.Join(conf.DataLocation, "reviews", fmt.Sprintf("review-%s-to-%s.json", startDate, endDate))
	}

	prompt := buildReviewPrompt(notes, startDate, endDate)

	client, err := newLLMClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
	defer cancel()

	raw, err := client.Complete(ctx, prompt)
	if err != nil {
		return err
	}

	jsonBody, err := normalizeLLMJSON(raw)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, jsonBody, 0o600); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "review saved to %s\n", outputPath)

	return nil
}

func collectNotes(conf config.Config, days int) ([]dailyNote, []string, error) {
	var (
		notes        []dailyNote
		missingDates []string
	)

	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.Add(-time.Duration(i) * 24 * time.Hour).Format(dateFormat)

		body, found, err := readNoteForDate(conf, date)
		if err != nil {
			return nil, nil, err
		}

		if !found {
			missingDates = append(missingDates, date)
			continue
		}

		notes = append(notes, dailyNote{
			Date: date,
			Body: body,
		})
	}

	return notes, missingDates, nil
}

func readNoteForDate(conf config.Config, date string) (string, bool, error) {
	dataPath := conf.GetDataFilePath(date)
	if exists(dataPath) {
		body, err := os.ReadFile(dataPath)
		return string(body), true, err
	}

	if archivePath, ok := fileInArchive(conf, date); ok {
		body, err := os.ReadFile(archivePath)
		return string(body), true, err
	}

	return "", false, nil
}

func buildReviewPrompt(notes []dailyNote, startDate, endDate string) string {
	var b strings.Builder

	b.WriteString("You are an expert personal-notes analyst. ")
	b.WriteString("Using the provided daily notes, produce a concise JSON summary for the week.\n\n")

	b.WriteString("Template fields to extract:\n")
	b.WriteString("- Sleep (0-10) and Energy (0-10).\n")
	b.WriteString("- Morning Checklist, Focus for Today, and Evening Checklist contain checkboxes like [ ] or [x].\n\n")

	b.WriteString("Tasks:\n")
	b.WriteString("1) For each day, capture Sleep and Energy (null if missing).\n")
	b.WriteString("2) Perform sentiment analysis of the author's mood using that day's writing (label positive/neutral/negative plus a brief reason).\n")
	b.WriteString("3) Mark a checklist as completed only if all items in that section are checked (treat [x], [X], or similar as checked). Track morning_completed, focus_completed, evening_completed per day.\n")
	b.WriteString("4) Compute averages for Sleep and Energy across available days.\n")
	b.WriteString("5) Count how many days had each checklist fully completed across the period.\n\n")

	b.WriteString("Return a single JSON object with this shape (numbers may be floats):\n")
	b.WriteString("{\n")
	b.WriteString(`  "period": {"start": "YYYY-MM-DD", "end": "YYYY-MM-DD"},` + "\n")
	b.WriteString(`  "averages": {"sleep": <number|null>, "energy": <number|null>},` + "\n")
	b.WriteString(`  "checklists_completed_totals": {"morning": <int>, "evening": <int>, "focus": <int>},` + "\n")
	b.WriteString(`  "days": [` + "\n")
	b.WriteString(`    {` + "\n")
	b.WriteString(`      "date": "YYYY-MM-DD",` + "\n")
	b.WriteString(`      "sleep": <number|null>,` + "\n")
	b.WriteString(`      "energy": <number|null>,` + "\n")
	b.WriteString(`      "sentiment": {"label": "positive|neutral|negative", "summary": "<short reason>"},` + "\n")
	b.WriteString(`      "checklists": {"morning_completed": <bool>, "focus_completed": <bool>, "evening_completed": <bool>}` + "\n")
	b.WriteString("    }\n")
	b.WriteString("  ]\n")
	b.WriteString("}\n\n")

	b.WriteString("Rules: Use only the provided notes. If data is missing, set the field to null. ")
	b.WriteString("Respond with JSON only, no markdown or commentary.\n\n")

	fmt.Fprintf(&b, "Analyze the period from %s to %s.\n\n", startDate, endDate)
	b.WriteString("Notes:\n")

	for idx, note := range notes {
		fmt.Fprintf(&b, "DAY: %s\n", note.Date)
		b.WriteString(note.Body)

		if idx < len(notes)-1 {
			b.WriteString("\n---\n")
		}

		b.WriteString("\n\n")
	}

	return b.String()
}

type llmClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func newLLMClient() (*llmClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}

	model := os.Getenv("NOTHELP_OPENAI_MODEL")
	if model == "" {
		model = defaultLLMModel
	}

	baseURL := os.Getenv("NOTHELP_OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = baseLLMURL
	}

	return &llmClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}, nil
}

func (c *llmClient) Complete(ctx context.Context, prompt string) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "You are a structured data analyst. Always return valid JSON.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm request failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}

	if parsed.Error != nil {
		return "", fmt.Errorf("llm error: %s", parsed.Error.Message)
	}

	if len(parsed.Choices) == 0 {
		return "", errors.New("llm response contained no choices")
	}

	return parsed.Choices[0].Message.Content, nil
}

func normalizeLLMJSON(raw string) ([]byte, error) {
	clean := stripCodeFences(raw)

	var payload any
	if err := json.Unmarshal([]byte(clean), &payload); err != nil {
		return nil, fmt.Errorf("llm did not return valid JSON: %w", err)
	}

	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, err
	}

	return pretty, nil
}

func stripCodeFences(s string) string {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```JSON")
		s = strings.TrimPrefix(s, "```")
	}

	s = strings.TrimSuffix(s, "```")

	return strings.TrimSpace(s)
}
