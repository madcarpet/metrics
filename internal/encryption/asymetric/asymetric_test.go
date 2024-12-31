package asymetric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsymetricEncryption(t *testing.T) {
	testString := "this is a test string for encryption"
	enc, err := EncryptWithPubKey("test_keys/pubkey.pem", []byte(testString))
	require.NoError(t, err)
	dec, err := DecryptWithPrivateKey("test_keys/privkey.pem", enc)
	require.NoError(t, err)
	assert.Equal(t, testString, string(dec))
}
