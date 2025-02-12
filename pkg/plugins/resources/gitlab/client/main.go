package client

import (
	"net/http"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// Spec defines a specification for a "Gitlab" resource
// parsed from an updatecli manifest file
type Spec struct {
	// [S][C][T] URL specifies the default Gitlab url in case of Gitlab enterprise
	URL string `yaml:",omitempty" jsonschema:"required"`
	// [S][C][T] Username specifies the username used to authenticate with Gitlab API
	Username string `yaml:",omitempty"`
	// [S][C][T] Token specifies the credential used to authenticate with
	Token string `yaml:",omitempty"`
}

const (

	// GITLABDOMAIN defines the default gitlab url
	GITLABDOMAIN string = "gitlab.com"
)

type Client *scm.Client

func New(s Spec) (Client, error) {

	url := EnsureValidURL(s.URL)

	client, err := gitlab.New(url)

	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{}

	if len(s.Token) >= 0 {
		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.StaticTokenSource(
					&scm.Token{
						Token: s.Token,
					},
				),
			},
		}
	}

	return client, nil

}

func EnsureValidURL(u string) string {
	if u == "" {
		u = GITLABDOMAIN
	}

	if !strings.HasPrefix(u, "https://") && !strings.HasPrefix(u, "http://") {
		u = "https://" + u
	}

	return u
}
