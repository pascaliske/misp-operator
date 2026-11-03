package client

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	AuthKey     string `json:"authkey"`
	OrgId       string `json:"org_id"`
	Push        bool   `json:"push"`
	Pull        bool   `json:"pull"`
	RemoteOrgId string `json:"remote_org_id"`
}

type serverEnvelope struct {
	Server Server `json:"Server"`
}

type ServerCreatePayload struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	AuthKey     string `json:"authkey"`
	OrgId       string `json:"org_id"`
	Push        bool   `json:"push"`
	Pull        bool   `json:"pull"`
	RemoteOrgId string `json:"remote_org_id"`
}

func (s Server) String() string {
	return fmt.Sprintf("server %q", s.Name)
}

func (s Server) CreatePayload() ServerCreatePayload {
	return ServerCreatePayload{
		Name:        s.Name,
		Url:         s.Url,
		AuthKey:     s.AuthKey,
		OrgId:       s.OrgId,
		Push:        s.Push,
		Pull:        s.Pull,
		RemoteOrgId: s.RemoteOrgId,
	}
}

func (c *ApiClient) GetServers() (*[]Server, error) {
	return getResourceList(c, "/servers", func(e serverEnvelope) Server {
		return e.Server
	})
}

func (c *ApiClient) GetServerByUrl(url string) (*Server, error) {
	// get all servers
	servers, err := c.GetServers()
	if err != nil {
		return nil, err
	}

	// find first by url
	server, found := findFirst(*servers, func(o Server) bool {
		return o.Url == url
	})
	if !found {
		return nil, &ApiError{StatusCode: http.StatusNotFound, Message: fmt.Sprintf("server %q not found", url)}
	}

	return server, nil
}

func (c *ApiClient) CreateServer(server *Server) (*Server, error) {
	return create(c, "/servers/add", *server, func(e serverEnvelope) Server {
		return e.Server
	})
}

func (c *ApiClient) EnsureServer(server *Server) (*Server, error) {
	return ensure(*server, func() (*Server, error) {
		return c.GetServerByUrl(server.Url)
	}, func() (*Server, error) {
		return c.CreateServer(server)
	})
}

func (c *ApiClient) TestConnection(server *Server) error {
	fmt.Printf("==> test connection for server %q\n", server.Url)

	var endpoint = fmt.Sprintf("/servers/testConnection/%s", server.Id)
	var result any

	if err := c.do(context.Background(), http.MethodGet, endpoint, nil, &result); err != nil {
		return fmt.Errorf("testing server connection: %w", err)
	}

	return nil
}
