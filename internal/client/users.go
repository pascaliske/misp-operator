package client

import (
	"fmt"
	"net/http"
)

type User struct {
	Id           string `json:"id"`
	OrgId        string `json:"org_id"`
	RoleId       string `json:"role_id"`
	Email        string `json:"email"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

type userEnvelope struct {
	User User `json:"User"`
}

type UserCreatePayload struct {
	OrgId  string `json:"org_id"`
	RoleId string `json:"role_id"`
	Email  string `json:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("user %q", u.Email)
}

func (u User) CreatePayload() UserCreatePayload {
	return UserCreatePayload{
		OrgId:  u.OrgId,
		RoleId: u.RoleId,
		Email:  u.Email,
	}
}

func (c *ApiClient) GetUsers() (*[]User, error) {
	return getResourceList(c, "/admin/users", func(e userEnvelope) User {
		return e.User
	})
}

func (c *ApiClient) GetUserByEmail(email string) (*User, error) {
	// get all users
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}

	// find first by name
	user, found := findFirst(*users, func(o User) bool {
		return o.Email == email
	})
	if !found {
		return nil, &ApiError{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("user %q not found", email)}
	}

	return user, nil
}

func (c *ApiClient) CreateUser(user *User) (*User, error) {
	return create(c, "/admin/users/add", *user, func(e userEnvelope) User {
		return e.User
	})
}

func (c *ApiClient) EnsureUser(user *User) (*User, error) {
	return ensure(*user, func() (*User, error) {
		return c.GetUserByEmail(user.Email)
	}, func() (*User, error) {
		return c.CreateUser(user)
	})
}
