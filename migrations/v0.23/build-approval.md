## What is build approval?

If attempting to contribute to a repository from a fork, you may encounter the following status:

![Screenshot 2024-02-09 at 10 29 36 AM](https://gist.github.com/assets/4662429/73892735-497d-4994-ab7d-7227d8ed5929)

On the Vela page for the build, you will see this:

![Screenshot 2024-02-09 at 10 33 40 AM](https://gist.github.com/assets/4662429/1e90a9bd-e684-49be-9900-38df0a2b61e0)

What this means is that the build is now in a `pending approval` state and must be approved by a repository admin to run.

To run the build, a repository admin need only click the `Approve Build` button or leverage the Vela CLI:

```sh
$ vela approve build --org <ORG> --repo <REPO> --number <NUMBER>
```

To deny the build, simply press the `Cancel Build` button.

This setting is enabled by default on all repositories with the release of `v0.23.0`. It can be edited by repository admins in the settings page:

![Screenshot 2024-02-09 at 10 37 47 AM](https://gist.github.com/assets/4662429/de2196bc-7182-40d3-a966-260297f856f7)


## Why keep this setting enabled?

- Approving builds from outside contributors is a vital way of protecting your CI pipeline, especially if you leverage `pull_request` events for status checks.

- Many CI services rely on the authenticity of CI build to imply permissions.

- Any secrets with `pull_request` events enabled will be protected from extraction.