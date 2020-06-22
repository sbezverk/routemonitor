package bmplistener

import (
	"testing"

	"github.com/sbezverk/routemonitor/pkg/classifier"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "empty input",
			input: nil,
		},
	}
	l := &listener{
		classifier: classifier.NewClassifierNLRI(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l.parsingWorker(tt.input)
		})
	}
}
