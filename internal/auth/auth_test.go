package auth // <-- gleiches Package wie in auth.go, damit du auf nicht‑exportierte Dinge zugreifen kannst

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

// Hinweis: ErrNoAuthHeaderIncluded ist bereits in auth.go definiert,
// du musst sie hier **nicht** erneut deklarieren.

// ------------------------------------------------------------
// TestGetAPIKey
// ------------------------------------------------------------
func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		wantKey       string
		wantErr       error
		wantErrString string
	}{
		{
			name:    "keine Authorization‑Header",
			headers: http.Header{},
			wantKey: "",
			wantErr: ErrNoAuthHeaderIncluded,
		},
		{
			name: "falsches Schema (nicht ApiKey)",
			headers: http.Header{
				"Authorization": []string{"Bearer abcdef"},
			},
			wantKey:       "",
			wantErrString: "malformed authorization header",
		},
		{
			name: "fehlendes Token nach ApiKey",
			headers: http.Header{
				"Authorization": []string{"ApiKey"},
			},
			wantKey:       "",
			wantErrString: "malformed authorization header",
		},
		{
			name: "zu viele Leerzeichen (nur erstes Token wird verwendet)",
			headers: http.Header{
				"Authorization": []string{"ApiKey   my-secret-key"},
			},
			wantKey: "my-secret-key",
			wantErr: nil,
		},
		{
			name: "gültiger Header",
			headers: http.Header{
				"Authorization": []string{"ApiKey my-secret-key"},
			},
			wantKey: "my-secret-key",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := GetAPIKey(tt.headers)

			// Schlüssel prüfen
			if gotKey != tt.wantKey {
				t.Fatalf("unexpected key: got %q, want %q", gotKey, tt.wantKey)
			}

			// Erwarteter Error über Variable (z. B. ErrNoAuthHeaderIncluded)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("unexpected error: got %v, want %v", err, tt.wantErr)
				}
				return
			}

			// Erwarteter Error über Fehlermeldungs‑String (für "malformed …")
			if tt.wantErrString != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrString) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrString, err)
				}
				return
			}

			// Kein Error erwartet
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
