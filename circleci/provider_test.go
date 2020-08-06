package circleci

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	circleci "github.com/samanthaq/terraform-provider-circleci/circleci/circleci-go"
	"golang.org/x/oauth2"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"circleci": testAccProvider,
	}
}

// ?
func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CIRCLECI_BASE_URL"); v == "" {
		t.Fatal("CIRCLECI_BASE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("CIRCLECI_TOKEN"); v == "" {
		t.Fatal("CIRCLECI_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("GITHUB_TOKEN"); v == "" {
		t.Fatal("GITHUB_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("GITHUB_ORGANIZATION"); v == "" {
		t.Fatal("GITHUB_ORGANIZATION must be set for acceptance tests")
	}
}

func newGithubClient() (*github.Client, error) {
	oauthTokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	oauthClient := oauth2.NewClient(context.TODO(), oauthTokenSource)
	if baseURL := os.Getenv("GITHUB_BASE_URL"); baseURL != "" {
		return github.NewEnterpriseClient(baseURL, baseURL, oauthClient)
	} else {
		log.Printf("GITHUB_BASE_URL not set, using github.com")
		return github.NewClient(oauthClient), nil
	}
}

func testFixtureGithubRepository(name, org string) error {
	githubClient, err := newGithubClient()
	if err != nil {
		return err
	}
	repo := &github.Repository{
		Name:     github.String(name),
		AutoInit: github.Bool(true),
	}
	_, _, err = githubClient.Repositories.Create(context.TODO(), org, repo)
	return err
}

func testFixtureCircleCIProject(name, org string) error {
	cciConfig := Config{
		AuthToken: os.Getenv("CIRCLECI_TOKEN"),
		BaseURL:   os.Getenv("CIRCLECI_BASE_URL"),
	}
	cciClient := cciConfig.NewClient()
	// cciClient := testAccProvider.Meta().(*circleci.Client)
	projectInput := &circleci.Project{
		VcsType:  "github",
		Username: org,
		Name:     name,
	}
	_, _, err := cciClient.Projects.Create(projectInput)
	return err
}
