# Acceptance tests against Folly

[Folly](https://github.com/equinix-labs/folly) is a mock of the Equinix APIs. It reads the API descriptions of the
[Equinix API catalog](https://docs.equinix.com/.well-known/api-catalog), validates every request and response against
them, and answers with scenarios that [Arazzo](https://spec.openapis.org/arazzo/latest.html) workflows describe.

The workflow `.github/workflows/folly_acctest.yml` runs the acceptance tests in `tests.txt` against Folly. The tests need
no Equinix account and no secrets, so they also run on pull requests from forks.

## What is here

| File | Content |
|---|---|
| `tests.txt` | The acceptance tests that run against Folly. |
| `folly.yaml` | The Folly configuration: the scenarios, and the polling loops of each scenario. |
| `flows/*.arazzo.yaml` | The Arazzo workflows. Each one describes the calls that the provider makes for one resource. |

| Test | Scenario | What it shows |
|---|---|---|
| `TestAccCloudRouterCreateOnlyRequiredParameters_PFCR` | `cloud-router` | Create, wait for `PROVISIONED`, rename (JSON Patch), delete, wait for `DEPROVISIONED`. |
| `TestAccFabricNetworkCreateOnlyRequiredParameters_PFCR` | `network` | The same lifecycle with the status in `operation.equinixStatus`, and values that the API computes (the change log). |
| `TestAccFabricRouteFilterPolicy_PFCR` | `route-filter` | Create and delete. The provider treats a 4xx response as deleted. |
| `TestAccDataSourceNetworkAccount_basic` | `ne-account` | A Network Edge data source with one fixed answer. |

## How a run works

1. The workflow checks out Folly, builds it, and runs `folly sync` to download the OpenAPI descriptions.
2. It starts Folly with `folly.yaml` and waits for `GET /__folly/health`.
3. It runs `go test` with `EQUINIX_API_ENDPOINT` pointing at Folly. This is how the other acceptance workflows run, with
   a different endpoint and dummy credentials. Folly issues the OAuth token itself (`flows/authenticate.arazzo.yaml`).

Folly answers a request from three sources, in this order:

1. **The scenario.** The workflow steps say which request comes next and what the response must say. A polling step
   (`onFailure` with `type: retry`) answers "keep waiting" for the number of requests that `loops` in `folly.yaml` sets,
   and then answers "done". The flow can also set values that the real API computes, for example
   `$response.body#/changeLog/createdBy == 'folly'`. Folly sets the field in the response.
2. **The resource store.** Folly keeps the resources that the client created during the run. A read returns the values
   that the create call sent, and a patch changes them. The state that a scenario sets stays on the resource.
3. **The spec.** Any other value comes from the example of the API description.

`tolerantReads: true` makes Folly answer a `GET` that the flow does not list from the store. The test framework
(`terraform-plugin-testing`) reads a resource between the steps of a test, and the provider does not send these reads.

## Add a test

1. Run the test against a plain Folly first: `folly serve --out <specs>`, with `EQUINIX_API_ENDPOINT` set. Folly logs
   every request. If the test waits ten minutes, the provider waits for a state that the spec example does not have: write a flow.
2. Write the flow in `flows/`. Use the step names of an existing flow, and the `operationId` values of the API
   description. Folly refuses to start, and names the problem, if the flow does not match the specs.
3. Add a scenario to `folly.yaml`, with an entry in `loops` for every polling step.
4. Add the test name to `tests.txt`. A scenario serves one test, so do not add a second test for the same resource type.

## Known limits

- A scenario keeps its state for one run of a workflow. Two tests for one resource type need two Folly processes.
- Search responses (`POST .../search`) are not built from the store yet.
- Folly validates requests against the API description. A test that sends a value that the description does not list
  (for example a device type code that is missing from an enum) fails with the documented 400. That is a gap in the
  description, and Folly shows it.
- Tests that need values from outside the request (the region of a metro, for example) need a literal value in the flow.
