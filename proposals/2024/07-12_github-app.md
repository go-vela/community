# Vela GitHub App

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :---: |
| **Author(s)** | Win San |
| **Reviewers** |       |
| **Date**      | July 12th, 2024 |
| **Status**    | Reviewed |

<!--
If you're already working with someone, please add them to the proper author/reviewer category.

If not, please leave the reviewer category empty and someone from the Vela team will assign it to themself.

Here is a brief explanation of the different proposal statuses:

1. Reviewed: The proposal is currently under review or has been reviewed.

2. Accepted: The proposal has been accepted and is ready for implementation.

3. In Progress: An accepted proposal is being implemented by actual work.

NOTE: The design is subject to change during this phase.

4. Cancelled: While or before implementation the proposal was cancelled.

NOTE: This can happen for a multitude of reasons.

5. Complete: This feature/change is implemented.
-->

## Background

<!--
This section is intended to describe the new feature, redesign or refactor.
-->

**Please provide a summary of the new feature, redesign or refactor:**

<!--
Provide your description here.
-->

- Create a GitHub App to replace the Vela OAuth App

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

- [GitHub Apps are preferred to OAuth apps](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/differences-between-github-apps-and-oauth-apps) because they use fine-grained permissions, give more control over which repositories the app can access, and use short-lived tokens. 
- These properties can harden security by limiting the damage that could be done if the app's credentials were leaked.
- GitHub Apps can also leverage GitHub's REST API endpoints for checks to create powerful and informative checks.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->
- Permissions and scope could be more granular.

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

- No

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

- No

## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated endpoints or url paths
* new/updated configuration variables (environment, flags, files, etc.)
* performance and user experience tradeoffs
* security concerns or assumptions
* examples or (pseudo) code snippets

<!-- Answer here -->
Replacing the OAuth App with a GitHub App (GHA) is a huge undertaking. Therefore, it might be a easiest to implement the transition in phases.
## Phase 1: GitHub App with advanced checks
GHAs can be installed on either an organization or personal account, each considered a single 'installation.' To encourage teams to install the GHA on their organization, we could develop an app offering advanced checks using [GitHub's Checks API](https://docs.github.com/en/enterprise-server@3.13/rest/checks?apiVersion=2022-11-28), exclusive to GHAs. This API supports features like code annotations and action selection by users. Possible checks we could implement include pipeline config validation, formatting, linting, and status reporting.

GHAs must be configured with the appropriate [permissions](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28). At a minimum, we'll need to include the following:
- [Repository permissions for "Checks"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-checks) to execute check runs
- [Repository permissions for "Contents"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-contents) to read the contents of a repo

Since GHAs use an app ID and private key instead of a client ID and client secret like OAuth Apps, we'll need to update the setup process for the SCM in Vela server.

```diff
// scm/setup.go

type Setup struct {
	// scm Configuration
	...
+	GithubAppID         int64
+	GithubAppPrivateKey string
```

One thing we'll have to consider is how GHA installations will be tracked. 
1. For repos already enabled in Vela, installing the GHA will trigger an installation addition event ([installation event documentation](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads#installation), [installation_repositories event documentation](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads#installation_repositories)). Upon this event, Vela server will search for the repo in the database and update its installation ID field.
2. For repos not yet enabled in Vela but with the GHA installed, enabling the repo will prompt Vela server to retrieve the installation ID associated with the repo. This process requires [authenticating as a GHA](https://docs.github.com/en/enterprise-server@3.13/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app) to access the list of installations across accounts (orgs and users). The installation ID will then be set when creating the repo in Vela.
3. When a user removes the GHA from an organization, an installation removal event ([installation removal event documentation](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads?actionType=deleted#installation), [installation_repositories removal event documentation](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads?actionType=removed#installation_repositories)) will be triggered. Vela server will act on this event by searching the database for the repo and clearing its installation ID field.

```diff
type Repo struct {
	ID           *int64    `json:"id,omitempty"`
	Owner        *User     `json:"owner,omitempty"`
	...
+	InstallID    *int64    `json:"install_id,omitempty"`
}
```

Another thing we'll have to consider is the check run flow. A potential flow could be:
1. Branch with an open PR receives a new commit
2. [Generate a JWT](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app) using the private key and app ID of GHA

```go
// Create a JWT token
func createJWT(appID int64, key *rsa.PrivateKey) (string, error) {
    now := time.Now()
    claims := jwt.MapClaims{
        "iat": now.Unix(),                       // Issued at time
        "exp": now.Add(time.Minute * 10).Unix(), // Expiration time
        "iss": appID,                            // GitHub App ID
    }

    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(key)
}
```

3. [Exchange the JWT](https://docs.github.com/en/enterprise-server@3.13/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation#using-an-installation-access-token-to-authenticate-as-an-app-installation) for an installation access token

```go
// Get installation access token
func getInstallationToken(ctx context.Context, jwtToken string, installationID int64) (string, error) {
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: jwtToken},
    )
    tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)

    token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
    if err != nil {
        return "", err
    }

    return token.GetToken(), nil
}
```

4. Use installation access token to [create check run(s)](https://docs.github.com/en/enterprise-server@3.13/rest/checks/runs?apiVersion=2022-11-28#create-a-check-run)
5. Initiate check run(s)
6. Use installation access token to [update status of check run(s)](https://docs.github.com/en/enterprise-server@3.13/rest/checks/runs?apiVersion=2022-11-28#update-a-check-run)
7. Optional: Take action offered by check run(s)

## Phase 2: Reduce usage of PAT with installation token
By default, Vela clones the repo of a build into a local volume that is mounted into each container. The cloning of the repo happens in the "clone" step. Currently, the repo owner's personal access token (PAT) is used to authenticate the clone request. 

This is not ideal as the PAT allows for unscoped access. More specifically, the PAT cannot be scoped to specific repositories or resources. A GHA installation access token is more favorable as it allows for scoped access. More specifically, during token creation, we are able to [fine tune the access](https://docs.github.com/en/enterprise-server@3.13/apps/creating-github-apps/authenticating-with-a-github-app/generating-an-installation-access-token-for-a-github-app#generating-an-installation-access-token) that the installation access token has, such as specifying individual repositories and permissions. For our purposes, we'd allow the installation access token read-only access to just the repo associated with the build.

When setting up the environment for a pipeline, we'll need to check which token type to use, since we can't assume that a repo has the Vela GHA installed. If the repo has installation ID set, we'll choose the more secure method of [authenticating](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation#about-authentication-as-a-github-app-installation) our clone request with a GHA installation access token.

```go
// compiler/native/environment.go

// helper function that creates the standard set of environment variables for a pipeline.
func environment(b *api.Build, m *internal.Metadata, r *api.Repo, u *api.User, iat string) map[string]string {
	// ...
	env := make(map[string]string)

	// vela specific environment variables
	// ...
	if r.GetInstallID() == 0 {
		env["VELA_NETRC_PASSWORD"] = u.GetToken()
		env["VELA_NETRC_USERNAME"] = "x-oauth-basic"
	} else {
		env["VELA_NETRC_PASSWORD"] = iat
		env["VELA_NETRC_USERNAME"] = "x-access-token"
	}
	// ...
```

We can further utilize installation access tokens to interact with the SCM, replacing a user's PAT for the following functions:
- `Changeset`: Lists the files changed in a commit.
- `GetOrgName`: Fetches the organization name from GitHub.
- `GetRepo`: Retrieves repository information from GitHub.
- `GetOrgAndRepoName`: Returns the names of the organization and repository in the SCM.
- `GetPullRequest`: Retrieves a pull request for a repository.
- `GetBranch`: Retrieves a branch for a repository.

It's important to note that we still need a user's PAT to access user-specific resources, such as the following:
- `ListUserRepos`: Lists all repositories the user has access to.
- `ListUsersTeamsForOrg`: Lists a user's teams for an org.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

Yes

2. What's the estimated time to completion?

<!-- Answer here -->

Several weeks

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

## Questions

**Please list any questions you may have:**

<!-- Answer here -->