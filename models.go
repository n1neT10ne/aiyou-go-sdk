// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// apiMessage représente le format du message envoyé à l'API
type apiMessage struct {
	Role    string    `json:"role"`
	Content []content `json:"content"`
}

// content représente le contenu d'un message
type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ContextWindow est un type personnalisé pour gérer les valeurs de context_window
// qui peuvent être soit des entiers soit des chaînes dans le JSON
type ContextWindow int

func (c *ContextWindow) UnmarshalJSON(data []byte) error {
	// Essayer d'abord de décoder comme une chaîne
	var strValue string
	if err := json.Unmarshal(data, &strValue); err == nil {
		// Si c'est une chaîne, la convertir en entier
		value, err := strconv.Atoi(strValue)
		if err != nil {
			return fmt.Errorf("invalid context window string value: %s", strValue)
		}
		*c = ContextWindow(value)
		return nil
	}

	// Si ce n'est pas une chaîne, essayer de décoder comme un entier
	var intValue int
	if err := json.Unmarshal(data, &intValue); err != nil {
		return fmt.Errorf("context window must be a string or integer")
	}
	*c = ContextWindow(intValue)
	return nil
}

func (c ContextWindow) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(c))
}

// Temperature est un type personnalisé pour s'assurer que la température est toujours
// formatée comme un float dans le JSON
type Temperature float64

func (t Temperature) MarshalJSON() ([]byte, error) {
	// Force le format float avec .0 pour les nombres entiers
	return []byte(fmt.Sprintf("%.1f", float64(t))), nil
}

// apiRequest représente la requête complète envoyée à l'API
type apiRequest struct {
	Messages     []apiMessage `json:"messages"`
	Model        string       `json:"model,omitempty"`
	AssistantID  string       `json:"assistantId,omitempty"`
	Temperature  Temperature  `json:"temperature"`
	Stream       bool         `json:"stream"`
	PromptSystem string       `json:"promptSystem,omitempty"`
}

// apiResponse représente la réponse de l'API en mode non-streaming
type apiResponse struct {
	Response struct {
		Model   string   `json:"model"`
		ID      string   `json:"id"`
		Created int64    `json:"created"`
		Choices []choice `json:"choices"`
		Usage   usage    `json:"usage"`
	} `json:"response"`
}

// choice représente un choix dans la réponse
type choice struct {
	Index        int     `json:"index"`
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// message représente le message dans la réponse
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Refusal string `json:"refusal"`
}

// ModelsResponse représente la réponse de l'API pour la liste des modèles
type ModelsResponse []struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

// Model représente un modèle disponible
type Model struct {
	Name string
}

// usage représente les statistiques d'utilisation des tokens
type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// streamResponse représente un chunk de réponse en mode streaming
type streamResponse struct {
	Model   string         `json:"model"`
	ID      string         `json:"id"`
	Created int64          `json:"created"`
	Choices []streamChoice `json:"choices"`
	Usage   *usage         `json:"usage,omitempty"`
}

// streamChoice représente un choix dans la réponse streaming
type streamChoice struct {
	Index int   `json:"index"`
	Delta delta `json:"delta"`
}

// delta représente le contenu incrémental dans la réponse streaming
type delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
