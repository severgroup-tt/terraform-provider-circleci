package circleci

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	circleci "github.com/samanthaq/terraform-provider-circleci/circleci/circleci-go"
)

func init() {
	resource.AddTestSweepers("circleci_project", &resource.Sweeper{
		Name: "circleci_project",
		F:    sweepRepositories,
	})
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sweepRepositories(region string) error {
	org := os.Getenv("GITHUB_ORGANIZATION")
	if org == "" {
		return fmt.Errorf("GITHUB_ORGANIZATION must be set for test sweeper")
	}
	if os.Getenv("GITHUB_TOKEN") == "" {
		return fmt.Errorf("GITHUB_TOKEN must be set for test sweeper")
	}

	client, err := newGithubClient()
	if err != nil {
		return err
	}

	repos, _, err := client.Repositories.List(context.TODO(), org, nil)
	if err != nil {
		return err
	}
	log.Printf("Found %d repos", len(repos))

	for _, r := range repos {
		if name := r.GetName(); strings.HasPrefix(name, "tf-acc-test-") {
			log.Printf("Destroying repository %s", name)

			if _, err := client.Repositories.Delete(context.TODO(), org, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccCircleCIProject_basic(t *testing.T) {
	var project circleci.Project

	rn := "circleci_project.foo"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-test-%s", randString)
	username := os.Getenv("GITHUB_ORGANIZATION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if err := testFixtureGithubRepository(name, username); err != nil {
				t.Fatalf("Could not create github repository: %s", err)
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCircleCIProjectConfig(name, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCircleCIProjectExists(rn, &project),
					resource.TestCheckResourceAttr(rn, "name", name),
					resource.TestCheckResourceAttr(rn, "vcs_type", "github"),
					resource.TestCheckResourceAttr(rn, "username", username),
				),
			},
		},
	})
}

func testAccCircleCIProjectConfig(name string, username string) string {
	return fmt.Sprintf(`
resource "circleci_project" "foo" {
  name     = "%s"
  username = "%s"
  vcs_type = "github"
}
`, name, username)
}

func testAccCheckCircleCIProjectExists(rn string, project *circleci.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}

		projectID := rs.Primary.ID
		if projectID == "" {
			return fmt.Errorf("No project name set")
		}

		client := testAccProvider.Meta().(*circleci.Client)
		gotProject, _, err := client.Projects.Read(projectID)
		if err != nil {
			return err
		}
		*project = *gotProject
		return nil
	}
}
