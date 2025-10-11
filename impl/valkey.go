package impl

import (
	"context"

	"github.com/valkey-io/valkey-go"
)

const SESSION_VALKEY_PREFIX = "session_"

type ValkeyClient struct {
	client valkey.Client
}

func (c *ValkeyClient) createSession(context context.Context, username string, sessionKey string) error {
	ttlSeconds := int64(config.Sessions.Ttl.Seconds())
	return c.client.Do(context, c.client.B().Setex().Key(SESSION_VALKEY_PREFIX+sessionKey).Seconds(ttlSeconds).Value(username).Build()).Error()
}

func (c *ValkeyClient) findUsernameBySession(context context.Context, sessionKey string) (string, error) {
	message, err := c.client.Do(context, c.client.B().Getex().Key(SESSION_VALKEY_PREFIX+sessionKey).Ex(config.Sessions.Ttl).Build()).ToMessage()
	if message.IsNil() {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return message.ToString()
}

var valkeyClient = func() ValkeyClient {
	clientImpl, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{config.Valkey.Address}})
	if err != nil {
		panic(err)
	}
	return ValkeyClient{
		client: clientImpl,
	}
}()
