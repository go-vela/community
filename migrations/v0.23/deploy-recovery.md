## Deployments prior to `v0.23.0`

- Deployment data in Vela was being sourced from GitHub
- Any GET request for deployments was merely a pass-thru to GitHub's database
- Incredibly slow response time
- Tying deployments to their respective builds was cumbersome, unreliable, and not represented in the UI


## Deployments in `v0.23.0`

- Deployments processed by Vela are now stored in Vela
- Any deployments from other sources will not show up in the Vela deployments page
- Much more responsive
- Builds from deployments are linked to their deployment record


## Cutover Process

The Vela admins determined it would be a clean cutover from old deployments to new deployments with the release of `v0.23.0`.

### What does this mean?

- _All prior deployments can be found in GitHub. They will **NOT** be in Vela._
- All new deployments will be available in both GitHub and Vela, but the Vela deployments will be more informative.
- The option to `redeploy` old deployments will be possible, but the process will be different. More info on that below.

### Why was this clean cutover process chosen?

- Prior to `v0.23.0`, deployments were carbon copies of the data in GitHub. There was no augmentation. As such, no data loss
- The complications of deciphering which GitHub deployments actually map back to Vela was not in the spirit of supporting our own data store.
- Enhancements were made to the CLI that make redeployments of legacy deploys easier.

## Process of Redploying Old Deployments

### Why redeploy an old deployment?

- Vela's `redeploy` feature was a UI concept that took the previous deployment's form and created a new deployment in GitHub.
- The benefit of this was the saving of `parameters`, `ref`, and `target`.
- If a team did not save these values in another location, they are unfortunately not simple to find for old deploys.

### Option 1: Webhook Payload Method

If the deployment was very recent, chances are the payload still exists in the webhook payloads section in the repo settings.

1. Navigate to https://github.com/ORG/REPO/settings/hooks and click the vela webhook
2. Scroll to the bottom of the page and find the webhook for a deployment. This webhook will have a header `X-GitHub-Event: deployment`
3. The values of this payload map back to Vela in the following way:

```
Environment --> Target
Ref --> Ref
Payload --> Parameters
```

4. You can either fill these values into the UI, or you can use the CLI, which now has a file input option for parameters:

```sh
$ vela --version
vela version v0.23.0

$ vela add deployment --org ORG --repo REPO --target <ENVIRONMENT> --ref <REF> --parameters-file <PATH_TO_JSON_PAYLOAD>
```

### Option 2: GitHub API Method

Perhaps an easier way is just to request the data from GitHub in the form of an API request.

GitHub's UI for deployments leaves a lot to be desired, so a simple cURL or Insomnia request is the best approach:

1. Generate a PAT that can read repo deployments
2. Perform the following request:

```sh
curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/ORG/REPO/deployments | jq '.[0]'
```

This will grab the last deployment for the repo. To fetch deployments from further in the past, change the `0` in the `jq` argument.

3. Just like the webhook method, parse the fields of the deployment object as so:

```
Environment --> Target
Ref --> Ref
Payload --> Parameters
```

4. You can either fill these values into the UI, or you can use the CLI, which now has a file input option for parameters:

```sh
$ vela --version
vela version v0.23.0

$ vela add deployment --org ORG --repo REPO --target <ENVIRONMENT> --ref <REF> --parameters-file <PATH_TO_JSON_PAYLOAD>
```