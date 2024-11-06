# Projektor Native Integration

<!--
The name of this markdown file should:

1. Short and contain no more then 30 characters

2. Contain the date of submission in MM-DD format

3. Clearly state what the proposal is being submitted for
-->

| Key           |       Value        |
| :-----------: |:------------------:|
| **Author(s)** |     Tim.Huynh      |
| **Reviewers** |                    |
| **Date**      | October 21st, 2024 |
| **Status**    |    In Progress     |

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

* Currently, there is a lack of native support for handling, visualizing, and tracking test results across builds. Inspired by the feature set in Projektor.dev, this proposal aims to add a dedicated, native support `test-report` feature to Vela.This feature will allow users to parse, store, and visualize test results in a more user-friendly manner.

**Please briefly answer the following questions:**

1. Why is this required?

* Current Projektor is an open-source project and not actively maintained. 
* This allows users to access test reports without needing use special test plugins that interact with other infrastructure.
* It may set the groundwork for feature flags, meaning faster feature rollout and easier feature rollback.

2. If this is a redesign or refactor, what issues exist in the current implementation?

* Currently, there is no native support for handling, visualizing, and tracking test results across builds.

3. Are there any other workarounds, and if so, what are the drawbacks?

* Yes. Users can set up separate infra and interact with it using a plugin to upload and access test reports. However, this is not a native solution and requires additional setup.

4. Are there any related issues? Please provide them below if any exist.

* This proposal will replace the current projektor-vela plugin. But it will not be replacing projektor-tap plugin.

## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* Dedicated `test-report` step.
* Backend and UI enhancements to visualize test results.
* Object Storage integration for storing test results.
* Slack integration for real-time notifications of test results.


### 1. Test Report Step Configuration
A dedicated `test-report` step will be added at the end of the Vela pipeline. This step uses a specialized Docker image (`vela/test-report-agent:latest`) to handle parsing and reporting tasks. Users can define the format, file path, and retention period for test data within this step, ensuring flexibility for different testing frameworks and workflows.
#### Example Configuration in `.vela.yml`
```yaml
steps:
  - name: test
    image: golang:latest
    commands:
      - go test ./... -json > test-results.json
  - name: test-report
    image: vela/test-report-agent:latest
    parameters:
      report_format: json
      report_path: "./test-results.json"
      retention_days: 30
      notify_slack: true
```
In this example:
- **`report_format`** specifies the format of the test results (e.g., JSON, JUnit XML).
- **`report_path`** defines the path to the test results file generated in the previous steps.
- **`retention_days`** allows users to set a retention policy for test data.
- **`notify_slack`** indicates whether to send notifications to Slack regarding the test results.
### 2. Test Report Step Workflow
- **Execution**:
    - The pipeline’s test steps execute as usual, generating a test results file (e.g., `test-results.json`).
    - The `test-report` step runs afterward, using the `vela/test-report-agent` image to parse and submit test results to the Vela backend.
### 3. Backend Enhancements
To support this new feature, Vela’s backend will require additional API endpoints and an expanded database schema.
#### Proposed Database Tables for Vela's Test Reporting
Based on the Projektor.dev architecture, here is a comprehensive list of tables required to store and manage test results, code quality, and related metrics in Vela's backend:
1. **`code_coverage_file`**: Stores file-level details for code coverage.
2. **`code_coverage_group`**: Represents groups or categories of coverage data.
3. **`code_coverage_stats`**: Holds detailed coverage statistics like statements covered, lines missed, and branches missed.
4. **`code_coverage_run`**: Aggregates coverage data for a specific run.
5. **`performance_results`**: Contains performance metrics, such as request count, average time, and maximum response time.
6. **`code_quality_report`**: Stores code quality report data, including file and group names.
7. **`test_run`**: Represents a single test run, storing counts of passed, failed, and skipped tests along with timing details.
8. **`results_metadata`**: Contains metadata related to each test run, such as CI information and group labels.
9. **`test_run_system_attributes`**: Stores system attributes for test runs, such as pinned status.
10. **`test_suite`**: Represents high-level test groupings, like test suites, with success and failure counts.
11. **`test_suite_group`**: Groups related test suites together for organized reporting.
12. **`test_case`**: Stores individual test cases, including their results, duration, and logs.
13. **`test_failure`**: Captures detailed information on test failures, including failure messages and stack traces.
14. **`test_run_attachment`**: Manages attachments for a test run, like log files and screenshots.
15. **`results_processing`**: Logs the status and errors (if any) from test results processing.
16. **`results_processing_failure`**: Tracks specific failure cases encountered during test processing.
17. **`shedlock`**: Used for distributed locking to prevent concurrent processing issues.

Tables related to git metadata are omitted as Vela already has a robust git integration system.


#### Proposed API Endpoints
Vela’s backend will expose new API endpoints to support test reporting and analytics features.
1. **`/api/v1/test-report`**: POST endpoint to submit test results for processing and storage.
2. **`/api/v1/test-report/{run_id}`**: GET endpoint to retrieve test results for a specific test run.
3. **`/api/v1/code-coverage`**: POST endpoint to submit code coverage data for processing and storage.
4. **`/api/v1/code-coverage/{run_id}`**: GET endpoint to fetch code coverage data for a specific run.
5. **`/api/v1/performance-metrics`**: POST endpoint to submit performance metrics for processing and storage.
6. **`/api/v1/performance-metrics/{run_id}`**: GET endpoint to retrieve performance metrics for a specific run.
7. **`/api/v1/code-quality`**: POST endpoint to submit code quality reports for processing and storage.
8. **`/api/v1/code-quality/{run_id}`**: GET endpoint to fetch code quality data for a specific run.
9. **`/api/v1/flaky-tests`**: GET endpoint to list flaky tests detected in the system.
10. **`/api/v1/test-notifications`**: POST endpoint to send test result notifications to Slack channels.


### 4. Object Store Integration
To effectively manage test artifacts and large volumes of test result data, Go-Vela will incorporate an object store as part of the test reporting solution.
#### Key Features of the Object Store
- **Artifact Storage**: Store test artifacts such as logs, screenshots, and detailed reports generated during test runs.
- **Access and Retrieval**: Provide a straightforward API for storing and retrieving artifacts, ensuring easy access from Vela’s user interface and other tools.
- **Scalability**: Enable seamless scalability to accommodate growing amounts of test data over time.
- **Data Retention Policies**: Implement retention policies to manage the lifecycle of stored artifacts, ensuring that outdated data is archived or deleted as necessary.
### 5. Slack Integration (Feature Flag)
To keep teams informed about test results in real-time, Go-Vela will include integration with Slack as a feature flag. This feature will notify designated channels or users about the outcomes of test runs and any significant changes in test performance.
#### Notification Features
- **Test Result Notifications**: Send messages to a specified Slack channel with a summary of test results, including the number of tests run, passed, failed, and any relevant error messages.
- **Flaky Test Alerts**: Notify teams when flaky tests are detected, prompting investigation and resolution.
#### Example Configuration for Slack
Users can specify Slack settings in their `.vela.yml` configuration:
```yaml
slack:
    - name: test-results
      image: vela/test-report-agent:latest
      ruleset:
        status: [ failure ]
      secrets: [ slack_webhook ]
      parameters:
        results: test-results/*.xml
        project: hey-vela # Name of your project to include in the Slack message
        filepath: heyvela_failure_message.json 
```
### 6. User Interface Enhancements
Vela’s UI will be enhanced to display test results in an intuitive and user-friendly manner.
- **Build Summary**: Include a **Test Results** section in the build summary view, displaying metrics such as total tests, pass/fail rates, and error messages.
- **Historical Data**: Provide a dashboard view showing trends in test pass/fail rates and test duration over time, allowing users to monitor stability and identify patterns.

### 7. Key Features to Implement
- **Test Reports and Analytics**: Generate and display test results, with pass/fail rates, historical trends, and detailed test case information.
- **Code Coverage Metrics**: Calculate and visualize code coverage data, including line and branch coverage percentages.
- **Performance Metrics**: Track performance metrics like response times, request counts, and error rates.
- **Code Quality Reports**: Display code quality metrics, such as common patterns and potential bugs.
- **Flaky Test Detection**: Identify and flag flaky tests for further investigation.
- **Real-time Notifications**: Send Slack notifications for test results, flaky tests, and other critical events.

**NOTES**: The list above is not a complete list of API endpoints but provides a starting point for implementing test reporting features in Vela. The API endpoints listed above are subject to change based on the final implementation details.
## Implementation
### Phases
- **Phase 1: Basic Test Reporting**
  - Implement the `test-report` step and backend support for storing test results.
  - Simple UI to display test results.
- **Phase 2: Advanced Features**
  - Add code coverage, performance metrics, and code quality reporting.
  - Integrate object storage for test artifacts.
  - Enhance the UI to visualize test data and metrics.
- **Phase 3: Slack Integration**
  - Implement Slack notifications for test results and flaky tests.
  - Allow users to configure Slack settings in `.vela.yml`.
- **Phase 4: Historical Data and Analytics**
    - Develop a dashboard to show historical test data and trends.
    - Add analytics features to track test performance over time.

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
* Multi-release

**Please provide all tasks (gists, issues, pull requests, etc.) completed to implement the design:**

<!-- Answer here -->

