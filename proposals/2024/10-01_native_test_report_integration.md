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

* Test reporting is a widely used capability across CI/CD platforms and often implemented as a bolt-on component requiring separate management of plugins and infrastructure. This proposal seeks to streamline the experience by integrating test reporting into Vela itself.
* Many platforms like Semaphore, CircleCI, Github Actions, Gitlab CI, Travis CI, Buildkite and Codefresh already offer robust test-reporting mechanisms such as:
  - Enabling workflows to store, view, and manage artifacts, including test logs, for visibility across builds.
  - Parsing JUnit XML, JSON, and other test result formats to display test results, pass/fail rates, and error messages.
  - Visualizing test executions and integrates with advance analytics engines.
* With this integration, Vela can deliver similar capabilities but tailored to its ecosystem while ensuring compatibility with modern development workflows.


**Please briefly answer the following questions:**

1. Why is this required?

* Test reporting is a widely used feature in CI/CD pipelines, providing insights into test results, code quality, and performance metrics. By integrating Projektor's test reporting features natively into Vela, users can easily track and analyze test data across builds.

2. If this is a redesign or refactor, what issues exist in the current implementation?

* Currently, there is no native support for handling, visualizing, and tracking test results across builds.

3. Are there any other workarounds, and if so, what are the drawbacks?

* Yes. Users can set up separate infra and interact with it using a plugin to upload and access test reports. However, this is not a native solution and requires additional setup.

4. Are there any related issues? Please provide them below if any exist.

* https://github.com/go-vela/community/issues/528

## Design

<!--
This section is intended to explain the solution design for the proposal.

NOTE: If there are no current plans for a solution, please leave this section blank.
-->

**Please describe your solution to the proposal. This includes, but is not limited to:**

* Object Storage integration for storing test results.
* Dedicated `test-report` step.
* Backend and UI enhancements to visualize test results.

### 1. Object Store Integration
Vela will integrate with an S3 compatible object storage system to store test results, code coverage data, and other artifacts. This integration will allow users to store large volumes of test data securely and efficiently.

#### Challenges
- **Scalability**: 
  - Support large volumes of test data and artifacts across multiple builds.
  - Bucket lifecycle management to manage data retention and cleanup.
- **Security and Access Control**: 
  - Access control mechanisms to restrict access to test data.
  - Data encryption at rest and in transit.
- **Performance**: Ensure fast and reliable storage and retrieval of test results and artifacts.
- **Cross-Platform Compatibility**: While many object storage systems are S3-compatible, ensuring compatibility with various systems is essential.

#### Implementation Details
- **Backend Configuration**:
- **Typical API Endpoints**:
  - **List Buckets**: GET /
  - **Create Bucket**: PUT /{bucket}
  - **List Objects**: GET /{bucket}
  - **Post Object**: POST /{bucket}/{object}
  - **Get Object**: GET /{bucket}/{object}
  - **Bucket Configuration**: PUT /{bucket}/config


##### Example Configuration in `docker-compose.yml`
```yaml
services:
  vela-worker:
    image: vela-worker:latest
    environment:
      - VELA_STORAGE_TYPE=s3
      - VELA_STORAGE_ENDPOINT=http://minio:9000
      - VELA_STORAGE_BUCKET=test-reports
      - VELA_STORAGE_ACCESS_KEY=minioadmin
      - VELA_STORAGE_SECRET_KEY=minioadmin
  minio:
    image: minio/minio
    command: server /data
    ports:
      - "9000:9000"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio-data:/data
volumes:
  minio-data:
```


### 2. Test Report Step Configuration
  #### Solution 1: Dedicated Docker Image for Test Reporting
  A dedicated `test-report` step will be added at the end of the Vela pipeline. This step uses a specialized Docker image (`vela/test-report-agent:latest`) to handle parsing and reporting tasks. Users can define the format, file path, and retention period for test data within this step, ensuring flexibility for different testing frameworks and workflows.
  ##### Example Configuration in `.vela.yml`
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
  ```
  In this example:
  - **`report_format`** specifies the format of the test results (e.g., JSON, JUnit XML).
  - **`report_path`** defines the path to the test results file generated in the previous steps.

##### Challenges
- Image Maintenance: Ensure the `test-report-agent` image is regularly updated and maintained including dependencies and security patches.
- Latency and Performance: Pulling image for every build may introduce latency, impacting build times.

#### Test Report Step Workflow
- **Execution**:
  - The pipeline’s test steps execute as usual, generating a test results file (e.g., `test-results.json`).
  - The `test-report` step runs afterward, using the `vela/test-report-agent` image to parse and submit test results to the Vela backend.

  #### Solution 2: Leveraging container output
  Alternatively, users can leverage the output of a container to pass test results to the `test-report` step. This approach allows users to generate test results within a container and pass them to the `test-report` step for processing.
  ##### Example Configuration in `.vela.yml`
  ``` yaml
  steps:
    - name: test
      image: golang:latest
      commands:
        - go test ./... -json > vela/outputs/test-results.json
    - name: read test report from outputs
      image: ${IMAGE}       # resolves to `ubuntu`
      pull: on_start
      parameters:
        report_format: json
        report_path: "./test-results.json"
  ```
  In this example:
  - **`report_format`** specifies the format of the test results (e.g., JSON, JUnit XML).
  - **`report_path`** defines the path to the test results file generated in the previous steps.

  ##### Challenges
  - Increased Complexity: Parsing logic must handle various formats and edge cases.
  - Scalability: Ensure the system can handle large volumes of test data efficiently and parsing large files without performance degradation.
  - Security: Implement secure parsing and storage mechanisms to protect sensitive test data.

#### Test Report Step Workflow
- **Execution**:
    - The pipeline’s test steps execute as usual, generating a test results file (e.g., `test-results.json`).
    - The `test-report` step runs afterward, worker will parse the test results file and submit the data to the Vela backend for storage and processing.

### 3. Backend Enhancements
To support this new feature, Vela’s backend will require additional API endpoints and an expanded database schema.
#### Proposed Database Tables for Vela's Test Reporting
Here is a comprehensive list of tables required to store and manage test results, code quality, and related metrics in Vela's backend:
1. **`code_coverage_file`**: Stores file-level details for code coverage.
2. **`code_coverage_run`**: Aggregates coverage data for a specific run.
3. **`performance_results`**: Contains performance metrics, such as request count, average time, and maximum response time.
4. **`code_quality_report`**: Stores code quality report data, including file and group names.
5. **`test_run`**: Represents a single test run, storing counts of passed, failed, and skipped tests along with timing details.
6. **`test_suite`**: Represents high-level test groupings, like test suites, with success and failure counts.
7. **`test_case`**: Stores individual test cases, including their results, duration, and logs.
8. **`test_run_attachment`**: Manages attachments for a test run, like log files and screenshots.
9. **`results_processing`**: Logs the status and errors (if any) from test results processing.

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

**NOTES**: The list above is not a complete list of API endpoints but provides a starting point for implementing test reporting features in Vela. The API endpoints listed above are subject to change based on the final implementation details.

### 4. User Interface Enhancements
Vela’s UI will be enhanced to display test results in an intuitive and user-friendly manner.
- **Build Summary**: Include a **Test Results** section in the build summary view, displaying metrics such as total tests, pass/fail rates, and error messages.
- **Historical Data**: Provide a dashboard view showing trends in test pass/fail rates and test duration over time, allowing users to monitor stability and identify patterns.

### 5. Key Features to Implement
- **Test Reports and Analytics**: Generate and display test results, with pass/fail rates, historical trends, and detailed test case information.
- **Code Coverage Metrics**: Calculate and visualize code coverage data, including line and branch coverage percentages.
- **Performance Metrics**: Track performance metrics like response times, request counts, and error rates.
- **Code Quality Reports**: Display code quality metrics, such as common patterns and potential bugs.
- **Flaky Test Detection**: Identify and flag flaky tests for further investigation.

## Implementation
### Phases
- **Phase 1: Basic Test Reporting**
  - Integrate object storage (ability to hook up Vela to a storage system).
    - Includes backend, API, and database changes.
  - Implement the `test-report` step and backend support for storing test results to a storage system.
  - UI/UX research.
- **Phase 2: UI**
  - Add code coverage, performance metrics, and code quality reporting.
  - Enhance the UI to visualize test results/data, including code coverage, performance metrics, and code quality reporting.
- **Phase 3: Enhanded UI**
  - Enhance UI with visualizations and dashboards/historical data (trends).

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

