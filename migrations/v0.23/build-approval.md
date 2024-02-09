## What is build approval?

If attempting to contribute to a repository from a fork, you may encounter the following status:

![Screenshot 2024-02-09 at 10 29 36 AM](https://media.git.target.com/user/17043/files/90939cdf-ea62-43af-a5b2-c1c23b2cc448)

On the Vela page for the build, you will see this:

![Screenshot 2024-02-09 at 10 33 40 AM](https://media.git.target.com/user/17043/files/a2623966-6085-4d52-9467-16cc4872178c)

What this means is that the build is now in a `pending approval` state and must be approved by a repository admin to run.

To run the build, a repository admin need only click the `Approve Build` button or leverage the Vela CLI:

```sh
$ vela approve build --org <ORG> --repo <REPO> --number <NUMBER>
```

To deny the build, simply press the `Cancel Build` button.

This setting is enabled by default on all repositories with the release of `v0.23.0`. It can be edited by repository admins in the settings page:

![Screenshot 2024-02-09 at 10 37 47 AM](https://media.git.target.com/user/17043/files/80d7e4df-7e2e-4c0c-954e-6e45f77ed731)


## Why keep this setting enabled?

- Approving builds from outside contributors is a vital way of protecting your CI pipeline, especially if you leverage `pull_request` events for status checks.

- Many CI services rely on the authenticity of CI build to imply permissions.

- Any secrets with `pull_request` events enabled will be protected from extraction.