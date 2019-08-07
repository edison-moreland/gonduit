package jwt

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Light abstraction for token revocation backend
// In the future it can redis, or some other store as using a global hash table is awful
// Dont judge me on this code pls

// TODO: Add to viper config
const revokedTimeToLive = time.Hour * 10

var revocationList revocationStore
var isRevocationListInitialized = false

func initRevocationBackend() error {
	if isRevocationListInitialized {
		return errors.New("Can't reinitialize revocation store")
	}

	revocationList = revocationStore{internal: make(map[string]time.Time)}
	isRevocationListInitialized = true

	return nil
}

// Revoke adds a token to the list of revoked tokens
func Revoke(token string) error {
	if !isRevocationListInitialized {
		err := initRevocationBackend()
		if err != nil {
			return fmt.Errorf("Could not revoke token. Reason: %v", err.Error())
		}
	}

	revocationList.Add(token, revokedTimeToLive)

	return nil
}

// IsRevoked checks if a token has been revoked
func IsRevoked(token string) bool {
	if !isRevocationListInitialized {
		err := initRevocationBackend()
		if err != nil {
			panic(fmt.Sprintf("Could check for revoked token. Reason: %v", err.Error()))
		}
	}

	return revocationList.DoesExists(token)

}

// PruneRevocationList removes token that were added more than 10 hours ago
func PruneRevocationList() error {
	if !isRevocationListInitialized {
		err := initRevocationBackend()
		if err != nil {
			return fmt.Errorf("Could check for revoked token. Reason: %v", err.Error())
		}
	}

	revocationList.Prune()
	return nil
}

type revocationStore struct {
	sync.RWMutex
	internal map[string]time.Time
}

// DoesExists check if a token has been stored
func (rs *revocationStore) DoesExists(token string) bool {
	rs.RLock()
	_, ok := rs.internal[token]
	rs.RUnlock()
	return ok
}

// Add puts a new token in the store
func (rs *revocationStore) Add(token string, timeToLive time.Duration) {
	rs.Lock()
	rs.internal[token] = time.Now().Add(timeToLive)
	rs.Unlock()
}

// Prune removes all expired tokens
func (rs *revocationStore) Prune() {
	rs.Lock()
	now := time.Now()
	for token, expiration := range rs.internal {
		if now.After(expiration) {
			delete(rs.internal, token)
		}
	}
	rs.Unlock()
}
