# Platform Settings

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           | Value |
| :-----------: | :-: |
| **Author(s)** | Easton Crupper |
| **Reviewers** |  |
| **Date**      | July 17th, 2024 |
| **Status**    | In Progress |

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

Consider the following scenario:

A user wishes to leverage authentication generated in a previous step in a plugin. How does that look in Vela today?

```yaml
steps:
    - name: generate credentials
      image: alpine
      commands:
        - ./fetch-credentials.sh
        - echo $FETCHED_CREDENTIALS > /vela/creds.txt

    - name: plugin to use auth
      image: vela/my-plugin:latest
      parameters:
        cred_path: /vela/creds.txt
```

This is the power of the shared Docker volume. All steps can tap into the workspace and "communicate" data. 

However, there are a few shortcomings with this. Some are obvious; some are not.

### Problem 1: environment variable data sourcing

Say the user in the above example does not own or contribute to the `vela/my-plugin` plugin. Further, the plugin does not support `cred_path: <path>`; instead, it expects a token provided as `$PLUGIN_TOKEN`. How can we accommodate this in Vela today?

```yaml
steps:
    - name: generate credentials
      image: alpine
      commands:
        - ./fetch-credentials.sh
        - echo $FETCHED_CREDENTIALS > /vela/creds.txt

    - name: plugin to use auth
      image: vela/my-plugin:latest
      entrypoint:
        - /bin/sh
        - -c
        - export PLUGIN_TOKEN=$(cat /vela/creds.txt); /bin/my-plugin
```

This will work for most cases. However, this seems rather unintuitive. Further, the documentation on such a process is scarce. One must be very familiar with the Vela code base to even know the above is possible.

### Problem 2: masking values that have been generated

Consider the prospect of masking the value of the fetched token like we do for Vela native secrets. This would enable users to enable verbose logging without risking secret exposure. How can we accommodate this in Vela today? 

We can't.

Why? The worker is the application that masks Vela secrets today. It parses container logs and determines where to apply the secret mask. Masking is not a natively supported concept with containers. 

This line of thinking approaches what I believe to be the core problem.

### Problem 3: the worker does not have reasonable access to the Docker volume

The Moby Docker client library does not offer the ability to read the actual content of the volume. It can create, inspect metadata, remove, and update the volume, but it cannot read it. Therefore, any content that the build generates can only be used by the build itself. This is very limiting when considering dynamic environments, value masking, and outputs for storage or rich status updates.

**Please briefly answer the following questions:**

1. Why is this required?

There are several active issues that would be solved / impacted by implementing this:

- https://github.com/go-vela/community/issues/140
- https://github.com/go-vela/community/issues/465
- https://github.com/go-vela/community/issues/448
- https://github.com/go-vela/community/issues/983

It also would help reach feature parity with GitHub Actions's implementation of similar concepts:
https://docs.github.com/en/actions/using-jobs/defining-outputs-for-jobs
https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#masking-a-value-in-a-log

2. If this is a redesign or refactor, what issues exist in the current implementation?

Net new

3. Are there any other workarounds, and if so, what are the drawbacks?

Mentioned in the problem set. There are a couple active workarounds for specific tasks but not comprehensive.

4. Are there any related issues? Please provide them below if any exist.

Linked above

## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* new/updated configuration variables (environment, flags, files, etc.)
* performance and user experience tradeoffs
* security concerns or assumptions
* examples or (pseudo) code snippets

### Proposed Solution: A Sidecar Container

When thinking about the solution for this, one needs not look further than how the worker funnels container logs every single second to the Vela server. This is data that is being generated by the build directed to std out. The worker parses the output and takes action. How can we accomplish the same idea but by reading file information?

If we launch a side car container for each build and "poll" the docker volume using it, we can extract information in between steps.

The code snippet below shows the `outputs.create` function, which would be called at the beginning of each build, setting up the side car container.

```go
    // -- executor/linux/build.go --
    //
	// create outputs container with a timeout equal to the repo timeout
	c.err = c.outputs.create(ctx, c.OutputCtn, (int64(60) * c.build.GetRepo().GetTimeout()))
	if c.err != nil {
		return fmt.Errorf("unable to create outputs container: %w", c.err)
	}
```

In order to collect outputs in the form of env files or (as shown below) a `library.Report` struct, we invoke the runtime client to execute a command on the running outputs container (simply `cat`). The `library.Report` in this example is actually an implementation of GitHub's `output` type for the [Checks API](https://docs.github.com/en/rest/checks/runs), giving users expansive power when designing plugins or step scripts.

```go
    // -- executor/linux/outputs.go
    //
    // poll tails the output for sidecar container.
    func (o *outputSvc) poll(ctx context.Context, ctn *pipeline.Container, stepCtn *pipeline.Container) (map[string]string, map[string]string, *library.Report, error) {
        // exit if outputs container has not been configured
        if len(ctn.Image) == 0 {
            return nil, nil, nil, nil
        }

        // update engine logger with outputs metadata
        //
        // https://pkg.go.dev/github.com/sirupsen/logrus#Entry.WithField
        logger := o.client.Logger.WithField("outputs", ctn.Name)

        logger.Debug("tailing container")

        // grab outputs
        outputBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/.env")
        if err != nil {
            return nil, nil, nil, err
        }

        // grab masked outputs
        maskedBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, "/vela/outputs/masked.env")
        if err != nil {
            return nil, nil, nil, err
        }

        var report library.Report

        // grab report if report path specified
        if stepCtn != nil && len(stepCtn.ReportPath) > 0 {
            logger.Infof("polling report.json file from outputs container %s", ctn.ID)

            reportPath := fmt.Sprintf("/vela/outputs/%s", stepCtn.ReportPath)

            reportBytes, err := o.client.Runtime.PollOutputsContainer(ctx, ctn, reportPath)
            if err != nil {
                return nil, nil, nil, err
            }

            err = json.Unmarshal(reportBytes, &report)
            if err != nil {
                logger.Infof("error unmarshaling report: %s", err)
            }
        }

        return toMap(outputBytes), toMap(maskedBytes), &report, nil
    }
```

For the runtime implementation, the Docker client does have the ability to attach a command to execute to a running container and extract the response of the command. Converting the file output to key-values allows us to have a dynamic environment. Further, we can use a similar idea to read beyond just an env file. This can be seen below:

```go
// PollOutputsContainer
func (c *client) PollOutputsContainer(ctx context.Context, ctn *pipeline.Container, path string) ([]byte, error) {
    // create configuration for command exec
	execConfig := types.ExecConfig{
		Tty:          true,
		Cmd:          []string{"sh", "-c", fmt.Sprintf("cat %s", path)},
		AttachStderr: true,
		AttachStdout: true,
	}

    // establish response 
	responseExec, err := c.Docker.ContainerExecCreate(ctx, ctn.ID, execConfig)
	if err != nil {
		log.Fatal(err)
	}

    // intercept response bytes from `cat` command
	hijackedResponse, err := c.Docker.ContainerExecAttach(ctx, responseExec.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatal(err)
	}

	defer hijackedResponse.Close()

	outputStdout := new(bytes.Buffer)
	outputStderr := new(bytes.Buffer)

	stdcopy.StdCopy(outputStdout, outputStderr, hijackedResponse.Reader)

	if outputStderr.Len() > 0 {
		fmt.Println("Error: ", outputStderr.String())
		return nil, fmt.Errorf("Error: %s", outputStderr.String())
	}

	data := outputStdout.Bytes()

	return data, nil
}
```

So how does this all come together to deliver dynamic environments, dynaminc value masking, and dynamic output reporting? During the build, the executor will do the following:

```go
    // beginning of the build — the external secret image execution prior to the first polling allows for masking of external secrets.
    //
    // -> executor/linux/build.go
	c.err = c.outputs.exec(ctx, c.OutputCtn)
	if c.err != nil {
		return fmt.Errorf("unable to exec outputs container: %w", c.err)
	}

	c.Logger.Info("executing secret images")
	// execute the secret
	c.err = c.secret.exec(ctx, &c.pipeline.Secrets)
	if c.err != nil {
		return fmt.Errorf("unable to execute secret: %w", c.err)
	}

	// poll outputs container for any updates
	opEnv, maskEnv, report, c.err = c.outputs.poll(ctx, c.OutputCtn, nil)
	if c.err != nil {
		return fmt.Errorf("unable to exec outputs container: %w", c.err)
	}

    // .....
    // inside for-loop for steps
    	// plan the step
		c.err = c.PlanStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to plan step: %w", c.err)
		}

        /* 
        In between planning and execution, the executor will perform environment updates and substitution using the values found in the outputs env paths.
        */

		// merge env from outputs
		_step.MergeEnv(opEnv)

		// merge env from masked outputs
		_step.MergeEnv(maskEnv)

		// add masked outputs to secret map so they can be masked in logs
		for key := range maskEnv {
			sec := &pipeline.StepSecret{
				Target: key,
			}
			_step.Secrets = append(_step.Secrets, sec)
		}

		// perform any substitution on dynamic variables
		err = _step.Substitute()
		if err != nil {
			return err
		}

		c.Logger.Infof("executing %s step", _step.Name)
		// execute the step
		c.err = c.ExecStep(ctx, _step)
		if c.err != nil {
			return fmt.Errorf("unable to execute step: %w", c.err)
		}

		// poll outputs
		opEnv, maskEnv, report, c.err = c.outputs.poll(ctx, c.OutputCtn, _step)
		if c.err != nil {
			return fmt.Errorf("unable to exec outputs container: %w", c.err)
		}
        
        // this is POC code that establishes potential for more complicated reporting to be part of the step object
		if _step.ReportStatus {
			libStep, err := step.Load(_step, &c.steps)
			if err != nil {
				return fmt.Errorf("unable to load step %s", _step.Name)
			}

			libStep.SetReport(report)

			_, _, err = c.Vela.Step.Update(c.build.GetRepo().GetOrg(), c.build.GetRepo().GetName(), c.build.GetNumber(), libStep)
			if err != nil {
				return fmt.Errorf("unable to update step %s", _step.Name)
			}
		}
```

**Pros**

* Leverages similar concepts of extracting container logs
* No need for storage options within Vela

**Cons**

* Adding another container to the mix inherently will increase complexity
* Adds another process to the build. A way to make it opt-in?

### Other Solutions

I have done a fair amount of research on different methods of consuming volume data as a running container, but they all seem rather cumbersome. However, I am certainly not going to claim my search as exhaustive. Perhaps one day the Docker SDK client will have the ability to read volume data — in which case that would certainly be the preferred way.

### K8s Considerations

Based on my interpretation of the code base, it appears that the Kubernetes runtime does not support a dynamic environment in the first place. Therefore, it makes sense to approach the implementation as a Docker only concept and potentially open an issue for exploring similar ideas in K8s.

## Implementation

<!--
This section is intended to explain how the solution will be implemented for the proposal.

NOTE: If there are no current plans for implementation, please leave this section blank.
-->

**Please briefly answer the following questions:**

1. Is this something you plan to implement yourself?

<!-- Answer here -->
* Yes

2. What's the estimated time to completion?

<!-- Answer here -->
* 2 weeks or so

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->