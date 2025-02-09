// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import "time"

// Options contient les options configurables du client
type Options struct {
	// BaseURL définit l'URL de base de l'API
	BaseURL string

	// Temperature contrôle la créativité des réponses (0.0-1.0)
	Temperature Temperature

	// Timeout définit le délai maximum pour une requête
	Timeout time.Duration

	// RetryConfig configure la politique de retry
	RetryConfig *RetryConfig

	// PromptSystem définit le prompt système à utiliser
	PromptSystem string

	// Stream indique si le modèle utilise le streaming
	Stream bool

	// Debug active l'affichage des messages de debug
	Debug bool

	// AssistantID spécifie l'ID de l'assistant à utiliser
	AssistantID string
}

// RetryConfig configure le comportement des retries
type RetryConfig struct {
	// MaxRetries est le nombre maximum de tentatives
	MaxRetries int

	// RetryDelay est le délai entre les tentatives
	RetryDelay time.Duration

	// MaxDelay est le délai maximum entre les tentatives
	MaxDelay time.Duration
}

// Constantes par défaut
const (
	defaultTimeout = 30 * time.Second
	minTemperature = 0.0
	maxTemperature = 1.0
	defaultTemp    = 1.0
)

// Option est une fonction qui configure les Options
type Option func(*Options)

// WithTemperature définit la température pour les réponses
func WithTemperature(temp float64) Option {
	return func(o *Options) {
		// Always set the temperature if it's within bounds
		if temp < float64(minTemperature) || temp > float64(maxTemperature) {
			// If out of bounds, keep default value
			return
		}
		o.Temperature = Temperature(temp)
	}
}

// WithTimeout définit le timeout des requêtes
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		if timeout > 0 {
			o.Timeout = timeout
		}
	}
}

// WithRetry configure la politique de retry
func WithRetry(maxRetries int, retryDelay time.Duration) Option {
	return func(o *Options) {
		if maxRetries > 0 && retryDelay > 0 {
			o.RetryConfig = &RetryConfig{
				MaxRetries: maxRetries,
				RetryDelay: retryDelay,
				MaxDelay:   retryDelay * 4, // Exponential backoff max
			}
		}
	}
}

// WithSystemPrompt définit le prompt système à utiliser
func WithSystemPrompt(prompt string) Option {
	return func(o *Options) {
		o.PromptSystem = prompt
	}
}

// WithStream active ou désactive le mode streaming
func WithStream(stream bool) Option {
	return func(o *Options) {
		o.Stream = stream
	}
}

// WithDebug active ou désactive les messages de debug
func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}

const defaultBaseURL = "https://ai.dragonflygroup.fr/api/v1"

// defaultOptions retourne les options par défaut
func defaultOptions() *Options {
	return &Options{
		BaseURL:      defaultBaseURL,
		Temperature:  Temperature(defaultTemp),
		Timeout:      defaultTimeout,
		RetryConfig:  nil,
		PromptSystem: "",
		Stream:       false,
		Debug:        false,
		AssistantID:  "",
	}
}

// WithAssistantID définit l'ID de l'assistant à utiliser
func WithAssistantID(assistantID string) Option {
	return func(o *Options) {
		o.AssistantID = assistantID
	}
}

// WithBaseURL définit l'URL de base de l'API
func WithBaseURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.BaseURL = url
		}
	}
}
