// Copyright (c) 2024 Cyrille BARTHELEMY
//
// This software is released under the MIT License.
// https://github.com/n1neT10ne/aiyou-go-sdk/blob/main/LICENSE

package aiyou

import (
	"encoding/json"
	"fmt"
	"os"
)

// debugPrint affiche un message de debug si le mode debug est activé
func debugPrint(options *Options, format string, args ...interface{}) {
	if !options.Debug {
		return
	}
	fmt.Fprintf(os.Stderr, "[AIYOU DEBUG] "+format+"\n", args...)
}

// debugJSON affiche une structure en JSON indenté si le mode debug est activé
func debugJSON(options *Options, prefix string, v interface{}) {
	if !options.Debug {
		return
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		debugPrint(options, "Error marshaling JSON for %s: %v", prefix, err)
		return
	}
	fmt.Fprintf(os.Stderr, "[AIYOU DEBUG] %s:\n%s\n", prefix, string(data))
}
