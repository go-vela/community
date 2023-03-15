# Build Tokens

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value                                                                                |
| :-----------: | :----------------------------------------------------------------------------------: |
| **Author(s)** | Easton.Crupper                                                                     |
| **Reviewers** |  |
| **Date**      | February 23rd, 2023                                                                      |
| **Status**    | In Progress                                                                             |

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

Instead of using the symmetric secret shared between the server and the worker, an executor will request a build token for each build it pulls from the queue. This token will be structured in the same way as our user tokens (JWT), and minted with claims such as `BuildID` and `Repo`. 

The permissions for API endpoints that are used to update build resources as well as pull secrets will be updated to only accept build tokens with the proper claims. 

In order to centralize token creation/validation, a new `internal/token` package will be added to the server code, wherein user tokens + build tokens will be minted and parsed.

Requests from workers will no longer be treated as a "user" with admin rights everywhere. Therefore JWT claims are being placed into the context prior to establishing a user.

Platform admins will be required to supply a `VELA_SERVER_PRIVATE_KEY` which will be used to sign and validate both token types: user and build.

**Please briefly answer the following questions:**

1. Why is this required?

<!-- Answer here -->

Using a symmetric token poses a significant security risk, as arbitrary code is executed worker-side. It is far more secure to limit the access of an executor to Vela resources. More specifically, it makes sense to treat the `vela-worker` as an `admin` of the repo / build it is presently executing rather than a `platform admin`.

2. If this is a redesign or refactor, what issues exist in the current implementation?

<!-- Answer here -->

Stated in the answer to question 1. 

3. Are there any other workarounds, and if so, what are the drawbacks?

<!-- Answer here -->

No

4. Are there any related issues? Please provide them below if any exist.

<!-- Answer here -->

N/A

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

### Private Key

Platform admins will need to supply a `VELA_SERVER_PRIVATE_KEY`, which will sign all tokens. This replaces the `user.Hash` as a signing mechanism (this field will be deprecated as a result of this change, since signing keys is its only use). 

Due to the relatively quick expirations of all tokens, the risk of collisions is extremely low. Each token will be signed with a specified `subject`, `issued_at`, and `expires_at`, which will ensure uniqueness.


### Token Manager

The changes for build tokens will include the creation of an `internal/token` package. Rather than minting, composing, parsing tokens in the middleware, the token manager will handle all of it and will be placed in the context. The manager will be defined as: 

```
type Manager struct {
	// PrivateKey key used to sign tokens
	PrivateKey string

	// SignMethod method to sign tokens
	SignMethod jwt.SigningMethod

	// UserAccessTokenDuration specifies the token duration to use for users
	UserAccessTokenDuration time.Duration

	// UserRefreshTokenDuration specifies the token duration for user refresh
	UserRefreshTokenDuration time.Duration

	// BuildTokenBufferDuration specifies the additional token duration of build tokens beyond repo timeout
	BuildTokenBufferDuration time.Duration
}
```

Explanations of the fields are in the comments. The process of minting and parsing will be very similar to how it is done today with user tokens. The primary difference is the addition of `TokenType`, `BuildID`, and `Repo` in the claims.


### API Changes

In order for the worker to be issued a build token, it must request one. This is where the only additional endpoint comes into play. The path `api/v1/repos/:org/:repo/build/:build/token` will, given that the caller is a worker (this will be determined using the symmetric token — more on that in [Permissions](02-23_build-tokens.md#permissions)), mint a build token and return it in the response. 

### Permissions

The bulk of the changes come in how the server handles permissions from workers. Instead of declaring a temporary admin user, as we do [today](https://github.com/go-vela/server/blob/74d2a68088ba20d6a91646277d98fbccb3caf641/router/middleware/user/user.go#L36-L48), the server will determine if the incoming token has correct access to the build using its build token. This is where the `BuildID` claim is useful.

Current
```
		// special handling for workers
		secret := c.MustGet("secret").(string)
		if strings.EqualFold(at, secret) {
			u := new(library.User)
			u.SetName("vela-worker")
			u.SetActive(true)
			u.SetAdmin(true)

			ToContext(c, u)
			c.Next()

			return
		}
```

Build Token Implementation (full picture [here](https://github.com/go-vela/server/pull/765))

`claims.go` 
```
		// parse and validate the token and return the associated the user
		claims, err = tm.ParseToken(at)
		if err != nil {
			util.HandleError(c, http.StatusUnauthorized, err)
			return
		}

		ToContext(c, claims)
		c.Next()
```
`user.go`
```
		// if token is not a user token, establish empty user to better handle nil checks
		if !strings.EqualFold(cl.TokenType, constants.UserAccessTokenType) {
			u := new(library.User)

			ToContext(c, u)
			c.Next()
        }
```

Further, `MustSecretAdmin()` will be adjusted to check the `Repo` claim. This will ensure that the caller is not attempting to get a secret to which they do not have access. 

This begs the question: how does the server know a worker is asking for a build token? For now, the symmetric token will continue to be used. In `claims.Establish()`: 

```
		// special handling for workers
		secret := c.MustGet("secret").(string)
		if strings.EqualFold(at, secret) {
			claims.Subject = "vela-worker"
			claims.TokenType = constants.ServerWorkerTokenType
			ToContext(c, claims)
			c.Next()

			return
		}
```

This is why the `VELA_SECRET` and `VELA_SERVER_SECRET` will still remain in the new implementation. Additional guard rails, such as build status, will be set in place to make sure that the symmetric token cannot be abused to request build tokens for builds that have already happened, yet to happen, etc.



### Worker Changes

The changes to the worker would be minimal. The check in routine will continue to use the symmetric token. However, each executor's Vela SDK client will be setup with a build token as the auth header rather than inheriting the symmetric token from the worker. For example, if a worker were to have 3 executors all running builds, each executor would have a different build token that it uses for its operations: pulling secrets, updating builds, logs, etc.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->

Yes, with help. Much of the work has been done.

2. What's the estimated time to completion?

<!-- Answer here -->

Several weeks.

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

Open PRs:

- https://github.com/go-vela/server/pull/765
- https://github.com/go-vela/sdk-go/pull/201
- https://github.com/go-vela/worker/pull/427

Merged PRs:

- https://github.com/go-vela/types/pull/276

## Questions
