// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import "errors"

// Erreurs spécifiques du package
var (
	// ErrInvalidToken est retourné quand le token d'authentification est invalide
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidModel est retourné quand le modèle spécifié est invalide
	ErrInvalidModel = errors.New("invalid model")

	// ErrRateLimit est retourné quand la limite de requêtes est atteinte
	ErrRateLimit = errors.New("rate limit exceeded")

	// ErrTimeout est retourné quand une requête dépasse le timeout configuré
	ErrTimeout = errors.New("request timeout")

	// ErrStreamCorrupted est retourné quand le stream de réponse est corrompu
	ErrStreamCorrupted = errors.New("stream corrupted")

	// ErrInvalidTemp est retourné quand la température est hors limites (0.0-1.0)
	ErrInvalidTemp = errors.New("temperature must be between 0.0 and 1.0")

	// ErrEmptyMessage est retourné quand le message est vide
	ErrEmptyMessage = errors.New("message cannot be empty")

	// ErrEmptyToken est retourné quand le token est vide
	ErrEmptyToken = errors.New("token cannot be empty")
)
