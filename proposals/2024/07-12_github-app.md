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

## Rollout Phases
### Phase 1: Advanced checks
- Create a GitHub App (GHA) that provides advanced checks
- This will incentivize teams to install the app onto their organizational repos
- Checks could include pipeline validation, formatting, linting, and status reporting.
- Make use of [GitHub's REST API endpoints for checks](https://docs.github.com/en/enterprise-server@3.13/rest/checks?apiVersion=2022-11-28), exclusive to GHAs
	- Allows for code annotations
	- Enables users to select an action offered by the app

**GHA Installation Scenarios**
1. User installs GHA onto an organization account and enables organization repo in Vela
	1. Vela server grabs installation id for that organization. Requires [authenticating as a GHA](https://docs.github.com/en/enterprise-server@3.13/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app#authentication-as-a-github-app) in order to grab the [list of installations](https://docs.github.com/en/enterprise-server@3.13/rest/apps/installations?apiVersion=2022-11-28#list-repositories-accessible-to-the-app-installation) across accounts (orgs and users).
	2. Vela server creates Repo in DB with installation id set
3. User installs GHA onto an organization repo but repo has already been enabled in Vela
	1. Acting on an installation addition event ([1](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads#installation), [2](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads#installation_repositories)), Vela server searches for Repo in DB
	2. Vela server sets installation id for that Repo
4. User removes GHA from an organization repo
	1. Acting on an installation removal event ([1](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads?actionType=deleted#installation), [2](https://docs.github.com/en/enterprise-server@3.13/webhooks/webhook-events-and-payloads?actionType=removed#installation_repositories)), Vela server searches for Repo in DB
	2. Vela server zeroes installation id for that Repo??

**Check Run Flow**
1. Branch with an open PR receives a new commit
2. [Generate a JWT](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app) using the private key and app ID of GHA
3. [Exchange the JWT](https://docs.github.com/en/enterprise-server@3.13/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation#using-an-installation-access-token-to-authenticate-as-an-app-installation) for an installation access token
5. Use installation access token to [create check run(s)](https://docs.github.com/en/enterprise-server@3.13/rest/checks/runs?apiVersion=2022-11-28#create-a-check-run)
6. Initiate check run(s)
7. Use installation access token to [update status of check run(s)](https://docs.github.com/en/enterprise-server@3.13/rest/checks/runs?apiVersion=2022-11-28#update-a-check-run)
8. Optional: Take action offered by check run(s)

**Code Integration**
- Add GHA's private key and app ID to SCM setup
```diff
// scm/setup.go

type Setup struct {
	// scm Configuration
	...
+	GithubAppID         int64
+	GithubAppPrivateKey string
```
- Add GHA installation ID to Repo
```diff
type Repo struct {
	ID           *int64    `json:"id,omitempty"`
	Owner        *User     `json:"owner,omitempty"`
	...
+	InstallID    *int64    `json:"install_id,omitempty"`
}
```
- Add check suite ID to Build
```diff
type Build struct {
	ID            *int64              `json:"id,omitempty"`
	Repo          *Repo               `json:"repo,omitempty"`
	...
+	CheckID       *int64              `json:"check_id,omitempty"`
}
```
- Generate a JWT using the GHA's private key and app ID
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
- Exchange the JWT for an installation access token
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

**Questions**
- What [permissions](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28) should the GitHub App have?
    - [Repository permissions for "Checks"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-checks)
    - [Repository permissions for "Commit statuses"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-commit-statuses)
    - [Repository permissions for "Contents"](https://docs.github.com/en/enterprise-server@3.13/rest/authentication/permissions-required-for-github-apps?apiVersion=2022-11-28#repository-permissions-for-contents)
- What checks should the GHA run?

### Phase 2: Replacing NETRC

### Phase 3: Moving off of OAuth App

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