## What's changing in `v0.23.2`?

The main change with `v0.23.2` is the introduction of a new setting for internal secrets: `allow substitution`. 

**IMPORTANT POSSIBLY BREAKING CHANGES**: 

With the release of `v0.23.2`, any _existing_ secret with `allow command = false` automatically set `allow substitution = false` as well. Further, _all_ `shared` secrets were manually updated to have both these values be `false` too.*

If your pipeline was impacted by this change, first consider further securing your pipeline by following the best practices section below before flipping the setting back. 

If you aren't sure if this impacts your pipeline, read all the context below.

*Plaftorm Administrators have the option to not perform this DML migration when upgrading.

### Background (`allow command`)

You may be familiar with the existing `allow command` setting, which determines whether or not a secret will be injected into the step's environment depending on the existence of `commands`:

```yaml
version: "1"

secrets:
  - name: no_commands # 'allow command' set to false
    key: vela/no_commands
    type: org
    engine: native
  
  - name: yes_commands # 'allow command' set to true
    key: vela/yes_commands
    type: org
    engine: native
    
steps:
  - name: print secret mask
    image: alpine:latest
    secrets: [ no_commands, yes_commands ]
    commands:
      - echo $NO_COMMANDS # will print nothing (not injected)
      - echo $YES_COMMANDS # will print secret mask ***
```

Users elect to set `allow command = false` as an additional security measure to prevent malicious pipelines from extracting the secret using a set of commands.

### Where does `allow substitution` come in?

During runtime, Vela substitutes any keys in the form ${KEY} with their corresponding value in the step environment. This often comes in the form of `parameters` for various plugins:

```yaml
steps:
  - name: docker build
    image: target/vela-kaniko:latest
    parameters:
      dry_run: true
      registry: docker.company.com
      repo: docker.company.com/some/repo
      tags:
        - ${VELA_BUILD_COMMIT:0:8}
```

This substitution process, prior to `v0.23.2`, extended to _all_ secrets specified in the step. For example:

```yaml
steps:
  - name: docker build
    image: target/vela-kaniko:latest
    secrets: [ api_token ]
    parameters:
      dry_run: true
      registry: docker.company.com
      repo: docker.company.com/some/repo
      build_args:
        - TKN=${API_TOKEN}
      tags:
        - ${VELA_BUILD_COMMIT:0:8}
```

While this type of runtime substitution can be very useful, it does open your pipeline up to attempts at exposure via manipulating those substitution requests.

This is why with `v0.23.2`, Vela is offering additional protection in the form of the `allow substitution` setting. By setting this value to `false`, secrets will _only_ be injected into the environment for the step and bypass any substitution.

**NOTE**: Substitution in the `commands` block is _not_ impacted by this change. All commands are converted to a shell script that pulls from the container environment during runtime.

### Mitigation / Best Practices

If you have a secret which a plugin expects as a specific environment variable, you can leverage `target` rather than substitution:

```diff
steps:
  - name: custom plugin
    image: docker.company.com/my-org/my-plugin
    secrets:
-      - github_token
+      - source: github_token
+        target: PARAMETER_API_TOKEN
    parameters:
-      api_token: ${GITHUB_TOKEN}
```

If you have need of leveraging substitution for a secret and environment injection will not do, _please_ protect that secret by using other Vela security measures:

- Properly set [build approval](https://go-vela.github.io/docs/usage/repo_settings/#outside-contributor-permissions) setting for your repository
- Limit your secret to certain webhook events, being cautious about PR events
- If the secret is `shared`, consider migrating it to be an `org` or `repo` secret. Having a shared secret without command+substitution prevention is _not a secure practice_.


### Updating Settings

You can update both `allow command` and `allow substitution` in the UI by editing the secret and selecting the corresponding boxes.

You can also, as always, leverage our CLI (provided you have upgraded to `v0.23.2` — `brew upgrade vela`).
```sh
vela update secret --secret.type org --org vela --name example --commands false --substitution false
```

