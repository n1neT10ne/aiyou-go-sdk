package aiyou

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

const (
	modelNonStream = "az-gpt-4o"
	modelStream    = "neuralmagic/Llama-3.1-Nemotron-70B-Instruct-HF-FP8-dynamic"
)

func getTestToken(t *testing.T) string {
	token := os.Getenv("AIYOU_TEST_TOKEN")
	if token == "" {
		t.Skip("AIYOU_TEST_TOKEN environment variable not set")
	}
	return token
}

func TestCompletion(t *testing.T) {
	testToken := getTestToken(t)

	tests := []struct {
		name        string
		model       string
		token       string
		message     string
		opts        []Option
		mockHandler func(w http.ResponseWriter, r *http.Request)
		want        string
		wantErr     error
	}{
		{
			name:    "simple response",
			model:   modelNonStream,
			token:   testToken,
			message: "Hello",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				resp := `{
					"model": "az-gpt-4o",
					"id": "test-id",
					"created": 1234567890,
					"choices": [
						{
							"index": 0,
							"delta": {
								"role": "assistant",
								"content": "Hello! How can I help you?"
							}
						}
					]
				}`
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, resp)
			},
			want: "Hello! How can I help you?",
		},
		{
			name:    "streaming response",
			model:   modelStream,
			token:   testToken,
			message: "Hello",
			opts:    []Option{WithStream(true)},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				chunks := []string{
					`{"choices":[{"delta":{"role":"assistant","content":"Hello"}}]}`,
					`{"choices":[{"delta":{"content":"!"}}]}`,
					`[DONE]`,
				}
				for _, chunk := range chunks {
					fmt.Fprintf(w, "data: %s\n\n", chunk)
				}
			},
			want: "Hello!",
		},
		{
			name:    "invalid token",
			model:   modelNonStream,
			token:   "invalid-token",
			message: "Hello",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			wantErr: ErrInvalidToken,
		},
		{
			name:    "empty token",
			model:   modelNonStream,
			token:   "",
			message: "Hello",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("handler should not be called")
			},
			wantErr: ErrEmptyToken,
		},
		{
			name:    "empty message",
			model:   modelNonStream,
			token:   testToken,
			message: "",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("handler should not be called")
			},
			wantErr: ErrEmptyMessage,
		},
		{
			name:    "with temperature",
			model:   modelNonStream,
			token:   testToken,
			message: "Hello",
			opts:    []Option{WithTemperature(0.8)},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(body), `"temperature":0.8`) {
					t.Errorf("temperature not set correctly in request: %s", string(body))
				}
				fmt.Fprintln(w, `{"choices":[{"delta":{"content":"Hi"}}]}`)
			},
			want: "Hi",
		},
		{
			name:    "with retry",
			model:   modelNonStream,
			token:   testToken,
			message: "Hello",
			opts:    []Option{WithRetry(2, time.Millisecond)},
			mockHandler: (func() func(http.ResponseWriter, *http.Request) {
				retryCount := 0
				return func(w http.ResponseWriter, r *http.Request) {
					if retryCount == 0 {
						retryCount++
						w.WriteHeader(http.StatusTooManyRequests)
						return
					}
					fmt.Fprintln(w, `{"choices":[{"delta":{"content":"Hi"}}]}`)
				}
			})(),
			want: "Hi",
		},
		{
			name:    "with system prompt",
			model:   modelNonStream,
			token:   testToken,
			message: "Hello",
			opts:    []Option{WithSystemPrompt("You are a French teacher")},
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(body), `"promptSystem":"You are a French teacher"`) {
					t.Errorf("system prompt not set correctly in request: %s", string(body))
				}
				fmt.Fprintln(w, `{"choices":[{"delta":{"content":"Bonjour!"}}]}`)
			},
			want: "Bonjour!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.mockHandler))
			defer server.Close()

			// Override base URL for testing
			opts := append(tt.opts, WithBaseURL(server.URL))

			// Execute test
			got, err := Completion(tt.model, tt.token, tt.message, opts...)

			// Check error
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check result
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestWithTemperature(t *testing.T) {
	tests := []struct {
		name string
		temp float64
		want float64
	}{
		{"valid temperature", 0.8, 0.8},
		{"too low", -0.1, defaultTemp},
		{"too high", 2.1, defaultTemp},
		{"zero", 0.0, 0.0},
		{"max", 2.0, 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTemperature(tt.temp)
			options := defaultOptions()
			opt(options)
			if float64(options.Temperature) != tt.want {
				t.Errorf("expected temperature %v, got %v", tt.want, float64(options.Temperature))
			}
		})
	}
}

func TestWithSystemPrompt(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
		want   string
	}{
		{"empty prompt", "", ""},
		{"valid prompt", "You are a helpful assistant", "You are a helpful assistant"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithSystemPrompt(tt.prompt)
			options := defaultOptions()
			opt(options)
			if options.PromptSystem != tt.want {
				t.Errorf("expected system prompt %q, got %q", tt.want, options.PromptSystem)
			}
		})
	}
}

func TestWithStream(t *testing.T) {
	tests := []struct {
		name   string
		stream bool
		want   bool
	}{
		{"enable streaming", true, true},
		{"disable streaming", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithStream(tt.stream)
			options := defaultOptions()
			opt(options)
			if options.Stream != tt.want {
				t.Errorf("expected stream %v, got %v", tt.want, options.Stream)
			}
		})
	}
}

func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"custom url", "https://custom.api.com", "https://custom.api.com"},
		{"empty url", "", defaultBaseURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBaseURL(tt.url)
			options := defaultOptions()
			opt(options)
			if options.BaseURL != tt.want {
				t.Errorf("expected base URL %q, got %q", tt.want, options.BaseURL)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{"valid timeout", 5 * time.Second, 5 * time.Second},
		{"zero timeout", 0, defaultTimeout},
		{"negative timeout", -1 * time.Second, defaultTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTimeout(tt.timeout)
			options := defaultOptions()
			opt(options)
			if options.Timeout != tt.want {
				t.Errorf("expected timeout %v, got %v", tt.want, options.Timeout)
			}
		})
	}
}
