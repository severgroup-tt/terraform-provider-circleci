package circleci

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	circleci "github.com/samanthaq/terraform-provider-circleci/circleci/circleci-go"
)

func TestAccCircleCIEnvironmentVariable_basic(t *testing.T) {
	var envvar circleci.EnvironmentVariable

	rn := "circleci_environment_variable.foo"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	projectName := fmt.Sprintf("tf-acc-test-%s", randString)
	username := os.Getenv("GITHUB_ORGANIZATION")
	name := "OHNO"
	value := "fruitypebbles"

	project := circleci.Project{
		VcsType: "github",
		Username: username,
		Name: projectName,
	}
	projectID := circleci.ProjectIdFromProject(project)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if err := testFixtureGithubRepository(projectName, username); err != nil {
				t.Fatalf("Could not create github repository: %s", err)
			}
			if err := testFixtureCircleCIProject(projectName, username); err != nil {
				t.Fatalf("Could not create circleci project: %s", err)
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCircleCIEnvironmentVariableConfig(projectID, name, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCircleCIEnvironmentVariableExists(rn, &envvar),
					resource.TestCheckResourceAttr(rn, "project_id", projectID),
					resource.TestCheckResourceAttr(rn, "name", name),
					resource.TestCheckResourceAttr(rn, "value", value),
				),
			},
		},
	})
}

func testAccCircleCIEnvironmentVariableConfig(projectID, name, value string) string {
	return fmt.Sprintf(`
resource "circleci_environment_variable" "foo" {
  project_id = "%s"
  name       = "%s"
  value      = "%s"
}
`, projectID, name, value)
}

func testAccCheckCircleCIEnvironmentVariableExists(rn string, envvar *circleci.EnvironmentVariable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set")
		}

		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("No project ID set")
		}
		name := rs.Primary.Attributes["name"]
		if name == "" {
			return fmt.Errorf("No name set")
		}

		envvarInput := &circleci.EnvironmentVariable{
			ProjectId: projectID,
			Name: name,
		}

		client := testAccProvider.Meta().(*circleci.Client)
		gotEnvvar, _, err := client.EnvironmentVariables.Read(envvarInput)
		if err != nil {
			return err
		}
		fmt.Printf("got: %s\n", gotEnvvar)
		*envvar = *gotEnvvar
		return nil
	}
}
