// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import (
	"bytes"
	"strings"
	"testing"
)

func TestStreamReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "simple event",
			input: "data: {\"content\":\"Hello\"}\n\n",
			want:  []string{"{\"content\":\"Hello\"}"},
		},
		{
			name: "multiple events",
			input: "data: {\"content\":\"Hello\"}\n\n" +
				"data: {\"content\":\"World\"}\n\n",
			want: []string{
				"{\"content\":\"Hello\"}",
				"{\"content\":\"World\"}",
			},
		},
		{
			name: "with comments",
			input: ": keep-alive\n" +
				"data: {\"content\":\"Hello\"}\n\n",
			want: []string{"{\"content\":\"Hello\"}"},
		},
		{
			name: "with empty lines",
			input: "\n" +
				"data: {\"content\":\"Hello\"}\n\n" +
				"\n",
			want: []string{"{\"content\":\"Hello\"}"},
		},
		{
			name:  "empty input",
			input: "",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := newStreamReader(strings.NewReader(tt.input))
			var got []string

			for {
				data, err := reader.readEvent()
				if err != nil {
					if len(data) > 0 {
						got = append(got, string(data))
					}
					break
				}
				got = append(got, string(data))
			}

			if len(got) != len(tt.want) {
				t.Errorf("got %d events, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("event %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestProcessStream(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		debug   bool
		want    string
		wantErr bool
	}{
		{
			name: "simple message",
			input: "data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n" +
				"data: {\"choices\":[{\"delta\":{\"content\":\" World\"}}]}\n\n" +
				"data: [DONE]\n\n",
			want: "Hello World",
		},
		{
			name: "with role prefix",
			input: "data: {\"choices\":[{\"delta\":{\"role\":\"assistant\",\"content\":\"Hello\"}}]}\n\n" +
				"data: {\"choices\":[{\"delta\":{\"content\":\"!\"}}]}\n\n" +
				"data: [DONE]\n\n",
			want: "Hello!",
		},
		{
			name:    "invalid json",
			input:   "data: {invalid json}\n\n",
			wantErr: true,
		},
		{
			name:  "empty choices",
			input: "data: {\"choices\":[]}\n\n",
			want:  "",
		},
		{
			name: "with debug enabled",
			input: "data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n" +
				"data: {\"choices\":[{\"delta\":{\"content\":\" World\"}}]}\n\n" +
				"data: [DONE]\n\n",
			debug: true,
			want:  "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := defaultOptions()
			options.Debug = tt.debug

			got, err := processStream(bytes.NewReader([]byte(tt.input)), options)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWithDebug(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
		want  bool
	}{
		{"enable debug", true, true},
		{"disable debug", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithDebug(tt.debug)
			options := defaultOptions()
			opt(options)
			if options.Debug != tt.want {
				t.Errorf("expected debug %v, got %v", tt.want, options.Debug)
			}
		})
	}
}
