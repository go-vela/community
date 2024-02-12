## What is build approval?

If attempting to contribute to a repository from a fork, you may encounter the following status:

![status.png](..%2F..%2F..%2F..%2F..%2FDownloads%2F68747470733a2f2f6d656469612e6769742e7461726765742e636f6d2f757365722f31373034332f66696c65732f39303933396364662d656136322d343361662d613562322d633163323362326363343438.png)

On the Vela page for the build, you will see this:

![68747470733a2f2f6d656469612e6769742e7461726765742e636f6d2f757365722f31373034332f66696c65732f61323632333936362d363038352d346435322d393436372d313663633438373231373863.png](..%2F..%2F..%2F..%2F..%2FDownloads%2F68747470733a2f2f6d656469612e6769742e7461726765742e636f6d2f757365722f31373034332f66696c65732f61323632333936362d363038352d346435322d393436372d313663633438373231373863.png)

What this means is that the build is now in a `pending approval` state and must be approved by a repository admin to run.

To run the build, a repository admin need only click the `Approve Build` button or leverage the Vela CLI:

```sh
$ vela approve build --org <ORG> --repo <REPO> --number <NUMBER>
```

To deny the build, simply press the `Cancel Build` button.

This setting is enabled by default on all repositories with the release of `v0.23.0`. It can be edited by repository admins in the settings page:

![68747470733a2f2f6d656469612e6769742e7461726765742e636f6d2f757365722f31373034332f66696c65732f38306437653464662d376532652d346330632d393534652d366534356637376564373331.png](..%2F..%2F..%2F..%2F..%2FDownloads%2F68747470733a2f2f6d656469612e6769742e7461726765742e636f6d2f757365722f31373034332f66696c65732f38306437653464662d376532652d346330632d393534652d366534356637376564373331.png)


## Why keep this setting enabled?

- Approving builds from outside contributors is a vital way of protecting your CI pipeline, especially if you leverage `pull_request` events for status checks.

- Many CI services rely on the authenticity of CI build to imply permissions.

- Any secrets with `pull_request` events enabled will be protected from extraction.