// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ListModels récupère la liste des modèles disponibles
func ListModels(token string, opts ...Option) ([]Model, error) {
	// Validation du token
	if token == "" {
		return nil, ErrEmptyToken
	}

	// Configuration
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Création du client HTTP avec timeout
	client := &http.Client{
		Timeout: options.Timeout,
	}

	// Création de la requête HTTP
	httpReq, err := http.NewRequest(
		"POST",
		options.BaseURL+"/models",
		bytes.NewReader([]byte("{}")), // Corps vide requis
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	// Fonction pour exécuter la requête avec retry
	var lastErr error
	maxRetries := getMaxRetries(options.RetryConfig)
	debugPrint(options, "Starting models request with max retries: %d", maxRetries)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := getRetryDelay(attempt, options.RetryConfig)
			debugPrint(options, "Retry attempt %d/%d, waiting %v", attempt, maxRetries, delay)
			time.Sleep(delay)
		}

		// Exécution de la requête
		resp, err := client.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("error executing request: %w", err)
			continue
		}
		defer resp.Body.Close()

		// Gestion des erreurs HTTP
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			lastErr = handleHTTPError(resp)
			debugPrint(options, "HTTP error: %v", lastErr)
			if !shouldRetry(resp.StatusCode) {
				return nil, lastErr
			}
			continue
		}

		debugPrint(options, "Models request successful: %s", resp.Status)

		// Utilisation de TeeReader pour le debug et le décodage
		var buf bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &buf)

		// Lecture de la réponse
		var modelsResp ModelsResponse
		if err := json.NewDecoder(teeReader).Decode(&modelsResp); err != nil {
			lastErr = fmt.Errorf("error decoding response: %w", err)
			if options.Debug {
				debugPrint(options, "Raw response body: %s", buf.String())
			}
			debugPrint(options, "Error decoding response: %v", err)
			continue
		}

		if options.Debug {
			debugPrint(options, "Raw response body: %s", buf.String())
		}

		debugJSON(options, "Models Response", modelsResp)

		// Extraction des modèles de la structure imbriquée
		var models []Model
		for _, provider := range modelsResp {
			for _, m := range provider.Models {
				models = append(models, Model{
					Name: m.Name,
				})
			}
		}

		return models, nil
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// Completion envoie une requête à l'API AI.You et retourne la réponse
func Completion(
	model string,
	token string,
	message string,
	opts ...Option,
) (string, error) {
	// Validation des entrées
	if token == "" {
		return "", ErrEmptyToken
	}
	if message == "" {
		return "", ErrEmptyMessage
	}

	// Configuration
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	debugJSON(options, "Options", options)
	debugPrint(options, "Stream mode: %v", options.Stream)

	// Préparation de la requête
	req := apiRequest{
		Messages: []apiMessage{
			{
				Role: "user",
				Content: []content{
					{
						Type: "text",
						Text: message,
					},
				},
			},
		},
		Model:        model,
		Temperature:  options.Temperature,
		Stream:       options.Stream,
		PromptSystem: options.PromptSystem,
		AssistantID:  options.AssistantID,
	}

	// Création du client HTTP avec timeout
	client := &http.Client{
		Timeout: options.Timeout,
	}

	// Encodage de la requête
	body, err := json.Marshal(req)
	if err != nil {
		debugPrint(options, "Error marshaling request: %v", err)
		return "", fmt.Errorf("error marshaling request: %w", err)
	}
	debugJSON(options, "Request", req)
	debugPrint(options, "Request stream mode: %v", req.Stream)

	// Création de la requête HTTP
	httpReq, err := http.NewRequest(
		"POST",
		options.BaseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	// Fonction pour exécuter la requête avec retry
	var lastErr error
	maxRetries := getMaxRetries(options.RetryConfig)
	debugPrint(options, "Starting request with max retries: %d", maxRetries)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := getRetryDelay(attempt, options.RetryConfig)
			debugPrint(options, "Retry attempt %d/%d, waiting %v", attempt, maxRetries, delay)
			time.Sleep(delay)
		}

		// Exécution de la requête
		resp, err := client.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("error executing request: %w", err)
			continue
		}
		defer resp.Body.Close()

		// Gestion des erreurs HTTP
		if resp.StatusCode != http.StatusOK {
			var buf bytes.Buffer
			teeReader := io.TeeReader(resp.Body, &buf)
			body, _ := io.ReadAll(teeReader)
			lastErr = handleHTTPError(resp)
			if options.Debug {
				debugPrint(options, "Raw error response body: %s", string(body))
			}
			debugPrint(options, "HTTP error: %v", lastErr)
			if !shouldRetry(resp.StatusCode) {
				return "", lastErr
			}
			continue
		}

		debugPrint(options, "Request successful: %s", resp.Status)

		// Utilisation de TeeReader pour le debug et le décodage
		var buf bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &buf)

		// Traitement de la réponse
		if options.Stream {
			if options.Debug {
				// Pour le streaming, on capture le début de la réponse pour le debug
				body := make([]byte, 1024)
				n, _ := teeReader.Read(body)
				if n > 0 {
					debugPrint(options, "Start of stream response: %s", string(body[:n]))
				}
				// Réinitialisation du body pour le streaming
				resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body[:n]), resp.Body))
			}
			return processStream(resp.Body, options)
		}

		// Lecture de la réponse non-streaming
		var apiResp apiResponse
		if err := json.NewDecoder(teeReader).Decode(&apiResp); err != nil {
			lastErr = fmt.Errorf("error decoding response: %w", err)
			if options.Debug {
				debugPrint(options, "Raw response body: %s", buf.String())
			}
			debugPrint(options, "Error decoding response: %v", err)
			continue
		}

		if options.Debug {
			debugPrint(options, "Raw response body: %s", buf.String())
		}
		debugJSON(options, "Response", apiResp)

		// Debug détaillé de la structure de la réponse
		debugPrint(options, "Response choices length: %d", len(apiResp.Response.Choices))
		if len(apiResp.Response.Choices) > 0 {
			debugPrint(options, "First choice details: index=%d, finish_reason=%s, role=%s",
				apiResp.Response.Choices[0].Index,
				apiResp.Response.Choices[0].FinishReason,
				apiResp.Response.Choices[0].Message.Role)
		}

		// Extraction du contenu
		if len(apiResp.Response.Choices) > 0 {
			return apiResp.Response.Choices[0].Message.Content, nil
		}
		return "", fmt.Errorf("no content in response")
	}

	return "", fmt.Errorf("max retries exceeded: %w", lastErr)
}

// handleHTTPError convertit les erreurs HTTP en erreurs typées
func handleHTTPError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return ErrInvalidToken
	case http.StatusTooManyRequests:
		return ErrRateLimit
	case http.StatusBadRequest:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad request: %s", string(body))
	default:
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}
}

// shouldRetry détermine si une erreur HTTP doit être retentée
func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode >= 500
}

// getMaxRetries retourne le nombre maximum de retries
func getMaxRetries(config *RetryConfig) int {
	if config == nil {
		return 0
	}
	return config.MaxRetries
}

// getRetryDelay calcule le délai avant la prochaine tentative
func getRetryDelay(attempt int, config *RetryConfig) time.Duration {
	if config == nil {
		return 0
	}

	delay := config.RetryDelay
	// Exponential backoff
	for i := 1; i < attempt; i++ {
		delay *= 2
	}

	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	// Ajouter un jitter de ±30%
	jitterRange := float64(delay) * 0.3
	jitter := time.Duration(rand.Float64()*jitterRange*2 - jitterRange)
	delay += jitter

	return delay
}
