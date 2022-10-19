package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/config"
	"github.com/semaphoreci/cli/generators"
	"github.com/spf13/cobra"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

const (
	GithubIntegrationOAuthToken = "github_token"
	GithubIntegrationApp        = "github_app"
)

var flagProjectName string
var flagRepoUrl string
var flagGithubIntegration string

func InitCmd() cobra.Command {
	cmd := cobra.Command{
		Use:   "init",
		Short: "Initialize a project",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			RunInit(cmd, args)
		},
	}

	cmd.Flags().StringVar(&flagRepoUrl, "repo-url", "", "explicitly set the repository url, if not set it is extracted from local git repository")
	cmd.Flags().StringVar(&flagProjectName, "project-name", "", "explicitly set the project name, if not set it is extracted from the repo-url")

	cmd.Flags().StringVar(
		&flagGithubIntegration,
		"github-integration",
		GithubIntegrationOAuthToken,
		fmt.Sprintf("github integration for the project. Possible values are: %s", validIntegrationTypes()),
	)

	return cmd
}

func init() {
	cmd := InitCmd()

	RootCmd.AddCommand(&cmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	var err error
	var name string
	var repoUrl string

	if flagRepoUrl != "" {
		repoUrl = flagRepoUrl
	} else {
		repoUrl, err = getGitOriginUrl()

		utils.Check(err)
	}

	if flagProjectName != "" {
		name = flagProjectName
	} else {
		name, err = ConstructProjectName(repoUrl)

		utils.Check(err)
	}

	if flagGithubIntegration != GithubIntegrationOAuthToken && flagGithubIntegration != GithubIntegrationApp {
		utils.Fail(fmt.Sprintf(
			"Invalid GitHub integration '%s' for project. Possible values are %s",
			flagGithubIntegration,
			validIntegrationTypes(),
		))
	}

	c := client.NewProjectV1AlphaApi()
	projectModel := models.NewProjectV1Alpha(name)
	projectModel.Spec.Repository.Url = repoUrl
	projectModel.Spec.Repository.RunOn = []string{"branches", "tags"}
	projectModel.Spec.Repository.IntegrationType = flagGithubIntegration

	project, err := c.CreateProject(&projectModel)

	utils.Check(err)

	if generators.PipelineFileExists() {
		fmt.Printf("[info] skipping .semaphore/semaphore.yml generation. It is already present in the repository.\n\n")
	} else {
		err = generators.GeneratePipelineYaml()

		utils.Check(err)
	}

	fmt.Printf("Project is created. You can find it at https://%s/projects/%s.\n", config.GetHost(), project.Metadata.Name)
	fmt.Println("")
	fmt.Printf("To run your first pipeline execute:\n")
	fmt.Println("")
	fmt.Printf("  git add .semaphore/semaphore.yml && git commit -m \"First pipeline\" && git push\n")
	fmt.Println("")
}

func ConstructProjectName(repoUrl string) (string, error) {
	formats := []*regexp.Regexp{
		regexp.MustCompile(`.+[:|\/].*\/([^.]*)`),
	}

	for _, r := range formats {
		match := r.FindStringSubmatch(repoUrl)

		if len(match) >= 2 {
			return match[1], nil
		}
	}

	errTemplate := "unsupported git remote format '%s'.\n"
	errTemplate += "\n"
	errTemplate += "Format must be one of the following:\n"
	errTemplate += "  - git@HOST:<owner>/<repo_name>.git\n"
	errTemplate += "  - git@HOST:<owner>/<repo_name>\n"
	errTemplate += "  - https://HOST/<owner>/<repo_name>\n"
	errTemplate += "  - https://HOST/<owner>/<repo_name>.git\n"
	errTemplate += "\n"
	errTemplate += "To add a project with an alternative git url, use the --repo-url flag:\n"
	errTemplate += "  - sem init --repo-url git@github.com:<owner>/<repo_name>.git\n"

	return "", errors.New(fmt.Sprintf(errTemplate, repoUrl))
}

func getGitOriginUrl() (string, error) {
	if flag.Lookup("test.v") == nil {
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			return "", errors.New("not a git repository")
		}

		return gitconfig.OriginURL()
	} else {
		return "git@github.com:/renderedtext/something.git", nil
	}
}

func validIntegrationTypes() string {
	return fmt.Sprintf("\"%s\" (OAuth token), \"%s\"", GithubIntegrationOAuthToken, GithubIntegrationApp)
}
