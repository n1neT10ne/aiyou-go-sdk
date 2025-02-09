// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// streamReader gère la lecture d'un stream SSE
type streamReader struct {
	reader *bufio.Reader
	buffer bytes.Buffer
}

// newStreamReader crée un nouveau lecteur de stream
func newStreamReader(r io.Reader) *streamReader {
	return &streamReader{
		reader: bufio.NewReader(r),
	}
}

// readEvent lit le prochain événement SSE
func (s *streamReader) readEvent() ([]byte, error) {
	s.buffer.Reset()
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if s.buffer.Len() > 0 {
					return s.buffer.Bytes(), nil
				}
			}
			return nil, err
		}

		// Ignore les commentaires
		if bytes.HasPrefix(line, []byte(":")) {
			continue
		}

		// Ligne vide marque la fin d'un événement
		if len(bytes.TrimSpace(line)) == 0 {
			if s.buffer.Len() > 0 {
				return s.buffer.Bytes(), nil
			}
			continue
		}

		// Vérifie si c'est une ligne de données
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))
			data = bytes.TrimSpace(data)
			s.buffer.Write(data)
		}
	}
}

// processStream traite le stream et reconstruit la réponse complète
func processStream(r io.Reader, options *Options) (string, error) {
	debugPrint(options, "Starting stream processing")
	reader := newStreamReader(r)
	var result strings.Builder

	for {
		data, err := reader.readEvent()
		if len(data) > 0 {
			debugPrint(options, "Received event: %s", string(data))
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// Ignore les événements vides ou [DONE]
		if len(data) == 0 {
			continue
		}
		if string(data) == "[DONE]" {
			debugPrint(options, "Received [DONE] event")
			continue
		}

		// Parse la réponse
		var resp streamResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			debugPrint(options, "Error parsing stream response: %v", err)
			return "", ErrStreamCorrupted
		}
		debugJSON(options, "Parsed stream response", resp)

		// Ajoute le contenu au résultat
		for _, choice := range resp.Choices {
			if choice.Delta.Content != "" {
				result.WriteString(choice.Delta.Content)
				debugPrint(options, "Added content: %q", choice.Delta.Content)
			}
		}
	}

	finalResult := result.String()
	debugPrint(options, "Stream processing completed, final result: %q", finalResult)
	return finalResult, nil
}
