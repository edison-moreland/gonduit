package jwt

import (
	"time"
)

// Light abstraction for token revocation backend
// In the future it can redis, or some other store as using a global hash table is awful
// Dont judge me on this code pls

// TODO: Add to viper config
const revokedTimeToLive = time.Hour * 10

var _revocationStore map[string]time.Time
var _isRevocationStoreInitialized = false

func getRevocationStore() map[string]time.Time {
	if !_isRevocationStoreInitialized {
		_revocationStore = make(map[string]time.Time)
		_isRevocationStoreInitialized = true
	}

	return _revocationStore
}

// Revoke adds a token to the list of revoked tokens
func Revoke(token string) {
	rs := getRevocationStore()
	rs[token] = time.Now().Add(revokedTimeToLive)
}

// IsRevoked checks if a token has been revoked
func IsRevoked(token string) bool {
	rs := getRevocationStore()

	_, ok := rs[token]
	return ok
}

// PruneRevocationList removes token that were added more than 10 hours ago
func PruneRevocationList() {
	rs := getRevocationStore()
	now := time.Now()

	for token, expiration := range rs {
		if now.After(expiration) {
			delete(rs, token)
		}
	}
}
