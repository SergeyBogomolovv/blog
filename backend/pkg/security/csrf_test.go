// Проверки подписи и привязки CSRF-токенов к сессии.
package security

import (
	"bytes"
	"strings"
	"testing"
)

// Эталон рассчитан Node.js crypto с тестовым ключом из 32 байтов 0x01.
const csrfVector = "0000000000000000000000000000000000000000000000000000000000000000." +
	"f50f082c6345034a207b0b7e945bf896843f962412bd6149a6dfcb5fdceaddc8"

// TestValidCSRF проверяет совместимость подписи и отказ для неверных CSRF-токенов.
func TestValidCSRF(t *testing.T) {
	t.Parallel()
	key := bytes.Repeat([]byte{1}, 32)
	for _, tc := range []struct {
		name, token, session string
		key                  []byte
		want                 bool
	}{
		{"node_signature", csrfVector, "session-a", key, true},
		{"another_session", csrfVector, "session-b", key, false},
		{"logged_out", csrfVector, "", key, false},
		{"another_key", csrfVector, "session-a", bytes.Repeat([]byte{2}, 32), false},
		{"missing", "", "session-a", key, false},
		{"unsigned", strings.Repeat("a", 129), "session-a", key, false},
		{"tampered_nonce", "1" + csrfVector[1:], "session-a", key, false},
		{"tampered_signature", csrfVector[:128] + "0", "session-a", key, false},
		{"non_hex_nonce", "z" + csrfVector[1:], "session-a", key, false},
		{"non_hex_signature", csrfVector[:128] + "z", "session-a", key, false},
		{"extra_data", csrfVector + ".extra", "session-a", key, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := ValidCSRF(tc.key, tc.token, tc.session); got != tc.want {
				t.Errorf("ValidCSRF() = %v, want %v", got, tc.want)
			}
		})
	}
}
