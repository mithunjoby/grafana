package githuboauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v82/github"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"

	provisioning "github.com/grafana/grafana/apps/provisioning/pkg/apis/provisioning/v0alpha1"
	"github.com/grafana/grafana/apps/provisioning/pkg/connection"
	"github.com/grafana/grafana/apps/provisioning/pkg/connection/oauth"
)

// provider implements the GitHub-specific parts of an OAuth app connection.
type provider struct {
	client *http.Client
}

// HTTPClient returns the overriding HTTP client, if any. It exists primarily
// for testing: the oauth exchange and GitHub API calls honor it in place of
// the default transport.
func (p *provider) HTTPClient() *http.Client {
	return p.client
}

func (p *provider) Endpoint() oauth2.Endpoint {
	return oauth2github.Endpoint
}

func (p *provider) ListRepositories(ctx context.Context, accessToken string) ([]provisioning.ExternalRepository, error) {
	if p.client != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, p.client)
	}
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
	client := github.NewClient(httpClient)

	opts := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var result []provisioning.ExternalRepository
	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			if resp != nil && (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) {
				return nil, connection.ErrAuthentication
			}
			return nil, fmt.Errorf("list repositories: %w", err)
		}

		for _, r := range repos {
			result = append(result, provisioning.ExternalRepository{
				Name:  r.GetName(),
				Owner: r.GetOwner().GetLogin(),
				URL:   r.GetHTMLURL(),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return result, nil
}

var _ oauth.Provider = (*provider)(nil)
