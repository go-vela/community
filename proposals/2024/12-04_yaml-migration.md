# YAML Library Migration

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
| **Date**      | December 4th, 2024 |
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

Vela uses github.com/buildkite/yaml as the YAML parsing library for pipeline and template files.

The buildkite project has been [archived](https://github.com/buildkite/yaml?tab=readme-ov-file#yaml-support-for-the-go-language) and has [lingering security vulnerabilities](https://www.mend.io/vulnerability-database/CVE-2022-3064). 

In `v0.26.0`, Vela should will the official `gopkg.in/yaml.v3` instead of `github.com/buildkite/yaml`.

The Vela compiler will _also_ do a small amount of custom unmarshaling in order to preserve one of the most common patterns that buildkite allowed: collapsing map keys (so long as the key is `<<:`).

For exactly one (1) release, Vela will have backwards compatibility such that `version: "legacy"` declared at the top of the Vela file will use the old buildkite library. This will let users with established and potentially complex pipelines to fix their YAML oddities at their own pace while also having a firm deadline.

**Please briefly answer the following questions:**

1. Why is this required?

Vela should no longer use an archived and vulnerable library as the primary parser for pipeline and template files.

2. If this is a redesign or refactor, what issues exist in the current implementation?

N/A

3. Are there any other workarounds, and if so, what are the drawbacks?

N/A

4. Are there any related issues? Please provide them below if any exist.

This has been attempted a number of times before, and it always has failed due to the lack of backwards compatibility.

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

### Backwards Compatibility

Vela will have the ability to compile every pipeline file (templates included) using the buildkite library for _**one**_ release cycle. This will be determined by `version: "legacy"` at the top level of the pipeline config. 

This quick band-aid will give users the ability to debug this change on their own timeline. 

This can be accomplished by having _two_ YAML types packages:

* go-vela/server/compiler/types/yaml/yaml —> `go-yaml`
* go-vela/server/compiler/types/yaml/buildkite —> `buildkite`

The types will be identical (except for parsing StageSlice due to the small differences). However, `buildkite` will have a `ToYAML()` function that converts already-parsed YAML objects into the new standard `go-yaml` type.

Determing _which_ library to use will be a simple helper function:
```go
func ParseYAML(data []byte) (*types.Build, []string, error) {
	var (
		rootNode yaml.Node
		warnings []string
		version  string
	)

    /*
        The below block grabs the root YAML node, traverses all map keys at the
        top of the pipeline, and determines which version to use. This will work
        for both pipeline files and templates.
    */
	err := yaml.Unmarshal(data, &rootNode)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal pipeline version yaml: %w", err)
	}

	if len(rootNode.Content) == 0 || rootNode.Content[0].Kind != yaml.MappingNode {
		return nil, nil, fmt.Errorf("unable to find pipeline version in yaml")
	}

	for i, subNode := range rootNode.Content[0].Content {
		if subNode.Kind == yaml.ScalarNode && subNode.Value == "version" {
			if len(rootNode.Content[0].Content) > i {
				version = rootNode.Content[0].Content[i+1].Value

				break
			}
		}
	}

    /*
        Once the version is determined, we can choose the correct library.
    */
	config := new(types.Build)

	switch version {
	case "legacy":
		legacyConfig := new(legacyTypes.Build)

		err := bkYaml.Unmarshal(data, legacyConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to unmarshal legacy yaml: %w", err)
		}

		config = legacyConfig.ToYAML()

		warnings = append(warnings, "using legacy version. Upgrade to go-yaml v3")

	default:
		// unmarshal the bytes into the yaml configuration
		err := yaml.Unmarshal(data, config)
		if err != nil {
			// if error is related to duplicate `<<` keys, attempt to fix
			if strings.Contains(err.Error(), "mapping key \"<<\" already defined") {
				root := new(yaml.Node)

				if err := yaml.Unmarshal(data, root); err != nil {
					fmt.Println("error unmarshalling YAML:", err)

					return nil, nil, err
				}

				warnings = collapseMergeAnchors(root.Content[0], warnings)

				newData, err := yaml.Marshal(root)
				if err != nil {
					return nil, nil, err
				}

				err = yaml.Unmarshal(newData, config)
				if err != nil {
					return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
				}
			} else {
				return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
			}
		}
	}

	return config, warnings, nil
}
```

This code will be removed in `v0.27.0`.

### Custom Node Handling for `<<:` keys

Since go-yaml v3 gives devs the ability to get at the nitty gritty details of the YAML file, we can "solve" for one of the most common issues we are set to face with this migration.

Below is a simple recursive function that will merge anchor keys throughout a pipeline:

```go
func collapseMergeAnchors(node *yaml.Node, warnings []string) []string {
	// only replace on maps
	if node.Kind == yaml.MappingNode {
		var (
			anchors      []*yaml.Node
			keysToRemove []int
			firstIndex   int
			firstFound   bool
		)

		// traverse mapping node content
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]

			// anchor found
			if keyNode.Value == "<<" {
				if (i+1) < len(node.Content) && node.Content[i+1].Kind == yaml.AliasNode {
					anchors = append(anchors, node.Content[i+1])
				}

				if !firstFound {
					firstIndex = i
					firstFound = true
				} else {
					keysToRemove = append(keysToRemove, i)
				}
			}
		}

		// only replace if there were duplicates
		if len(anchors) > 1 && firstFound {
			seqNode := &yaml.Node{
				Kind:    yaml.SequenceNode,
				Content: anchors,
			}

			node.Content[firstIndex] = &yaml.Node{Kind: yaml.ScalarNode, Value: "<<"}
			node.Content[firstIndex+1] = seqNode

			for i := len(keysToRemove) - 1; i >= 0; i-- {
				index := keysToRemove[i]

				warnings = append(warnings, fmt.Sprintf("%d:duplicate << keys in single YAML map", node.Content[index].Line))
				node.Content = append(node.Content[:index], node.Content[index+2:]...)
			}
		}

		// go to next level
		for _, content := range node.Content {
			warnings = collapseMergeAnchors(content, warnings)
		}
	} else if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			warnings = collapseMergeAnchors(item, warnings)
		}
	}

	return warnings
}
```

This function will only be invoked if the original pass at unmarshaling using `v3` contains the specific error related to using duplicate map keys.

### Warnings

Lastly, as one might have seen in the example code snippets, Vela will introduce pipeline warnings with this migration. Not only will it help identify issues today related to the migration, it lays the groundwork for future efforts.

The warnings will, to start, just be a string slice in the `pipelines` table.


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
* 2-3 weeks development, 2-3 months implementation of release, monitoring, education, etc.

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

More details can be found on the [discussion page](https://github.com/go-vela/community/discussions/1023) for this migration.