package client

import (
	"fmt"
	"net/http"
)

type AuthKey struct {
	Id          string `json:"id"`
	Uuid        string `json:"uuid"`
	UserId      string `json:"user_id"`
	Raw         string `json:"authkey_raw"`
	Comment     string `json:"comment"`
	ReadOnly    bool   `json:"read_only"`
	DateCreated string `json:"created"`
	Expiration  string `json:"expiration"`
	// LastUsed    string `json:"last_used,omitempty"`
}

type authKeyEnvelope struct {
	AuthKey AuthKey `json:"AuthKey"`
}

type AuthKeyCreatePayload struct {
	Uuid     string `json:"uuid,omitempty"`
	UserId   string `json:"user_id"`
	Raw      string `json:"authkey_raw"`
	Comment  string `json:"comment"`
	ReadOnly bool   `json:"read_only"`
}

func (k AuthKey) String() string {
	return fmt.Sprintf("auth key %q", k.Uuid)
}

func (k AuthKey) CreatePayload() AuthKeyCreatePayload {
	return AuthKeyCreatePayload{
		Uuid:     k.Uuid,
		UserId:   k.UserId,
		Raw:      k.Raw,
		Comment:  k.Comment,
		ReadOnly: k.ReadOnly,
	}
}

func (c *ApiClient) GetAuthKeys() (*[]AuthKey, error) {
	return getResourceList(c, "/auth_keys", func(e authKeyEnvelope) AuthKey {
		return e.AuthKey
	})
}

func (c *ApiClient) GetAuthKeyByComment(comment string) (*AuthKey, error) {
	// get all authKeys
	authKeys, err := c.GetAuthKeys()
	if err != nil {
		return nil, err
	}

	// find first by comment
	authKey, found := findFirst(*authKeys, func(k AuthKey) bool {
		return k.Comment == comment
	})
	if !found {
		return nil, &ApiError{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("auth key %q not found", comment)}
	}

	return authKey, nil
}

func (c *ApiClient) CreateAuthKey(authKey *AuthKey) (*AuthKey, error) {
	return create(c, fmt.Sprintf("/auth_keys/add/%s", authKey.UserId), *authKey, func(e authKeyEnvelope) AuthKey {
		return e.AuthKey
	})
}

func (c *ApiClient) EnsureAuthKey(authKey *AuthKey) (*AuthKey, error) {
	return ensure(*authKey, func() (*AuthKey, error) {
		return c.GetAuthKeyByComment(authKey.Comment)
	}, func() (*AuthKey, error) {
		return c.CreateAuthKey(authKey)
	})
}
