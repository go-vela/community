## What is build approval?

If attempting to contribute to a repository from a fork, you may encounter the following status:

![status](https://github.com/go-vela/community/assets/4662429/5ca3e0a6-c087-4493-9a1e-a2acfa150e4d)

On the Vela page for the build, you will see this:

![velaUI](https://github.com/go-vela/community/assets/4662429/2515943e-7bd0-4f3a-9e96-14f28f9c4e15)

What this means is that the build is now in a `pending approval` state and must be approved by a repository admin to run.

To run the build, a repository admin need only click the `Approve Build` button or leverage the Vela CLI:

```sh
$ vela approve build --org <ORG> --repo <REPO> --number <NUMBER>
```

To deny the build, simply press the `Cancel Build` button.

This setting is enabled by default on all repositories with the release of `v0.23.0`. It can be edited by repository admins in the settings page:

![permission](https://github.com/go-vela/community/assets/4662429/cd4bab4e-149b-49b0-b7ff-244818ec582e)


## Why keep this setting enabled?

- Approving builds from outside contributors is a vital way of protecting your CI pipeline, especially if you leverage `pull_request` events for status checks.

- Many CI services rely on the authenticity of CI build to imply permissions.

- Any secrets with `pull_request` events enabled will be protected from extraction.
