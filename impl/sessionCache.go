package impl

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type SessionCache struct {
	cache        *ttlcache.Cache[string, string]
	requestsChan chan SessionCacheRequest
}

type SessionCacheRequest struct {
	sessionKey      string
	responseChannel chan SessionCacheResponse
}

type SessionCacheResponse struct {
	username string
	error    error
}

func (s *SessionCache) runWorker() {
	for {
		request := <-s.requestsChan
		sessionKey := request.sessionKey
		cachedEntry := s.cache.Get(sessionKey)
		if cachedEntry != nil {
			cachedUsername := cachedEntry.Value()
			request.responseChannel <- SessionCacheResponse{
				username: cachedUsername,
				error:    nil,
			}
			continue
		}

		usernameFromValkey, err := valkeyClient.findUsernameBySession(context.Background(), sessionKey)
		if err == nil && usernameFromValkey != "" {
			s.cache.Set(sessionKey, usernameFromValkey, 10*time.Second)
		}
		request.responseChannel <- SessionCacheResponse{
			username: usernameFromValkey,
			error:    err,
		}
	}
}

func (s *SessionCache) findBySession(session string) (string, error) {
	responseChan := make(chan SessionCacheResponse)
	request := SessionCacheRequest{
		sessionKey:      session,
		responseChannel: responseChan,
	}
	s.requestsChan <- request

	response := <-responseChan
	return response.username, response.error
}

func newSessionCache() SessionCache {
	cache := ttlcache.New[string, string]()
	s := SessionCache{
		cache,
		make(chan SessionCacheRequest),
	}
	go s.runWorker()
	return s
}

var sessionCache = newSessionCache()
