package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func newOpenAIClient() (*openai.Client, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &client, nil
}

func SuggestFromReport(reportJSON []byte, specBytes []byte) (string, error) {
	client, err := newOpenAIClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	prompt := fmt.Sprintf(
		`You are specgate suggest.

Input is specgate's JSON report describing missing documentation in an OpenAPI spec.

Task:
- For each operation:
  - If the operation appears in errors.missing_summaries, propose a Summary.
  - If the operation appears in warnings.missing_descriptions, propose a Description.
- If an operation is missing only one of these, DO NOT propose the other.
- Do NOT include any other operations.
- Do NOT explain what you did.

Output format (conditional):
- <METHOD> <PATH>
  - If proposing a summary:
    Summary: <text>
  - If proposing a description:
    Description: <text>

Rules:
- Summary: max 80 characters, imperative verb.
- Description: 1–2 sentences, max 240 characters.
- Do not invent features, headers, or status codes not present in the input.

Specgate JSON report:
%s

OpenAPI spec:
%s`,
		string(reportJSON),
		string(specBytes),
	)

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT4oMini,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
	})
	if err != nil {
		return "", err
	}

	out := strings.TrimSpace(resp.OutputText())
	out = strings.TrimPrefix(out, "```")
	out = strings.TrimSuffix(out, "```")
	out = strings.TrimPrefix(out, "plaintext\n")
	out = strings.TrimPrefix(out, "text\n")
	out = strings.TrimPrefix(out, "markdown\n")

	const maxOut = 12_000
	if len(out) > maxOut {
		out = out[:maxOut] + "\n\n*(truncated)*"
	}

	return out, nil
}