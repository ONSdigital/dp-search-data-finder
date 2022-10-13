Feature: Consuming reindex-requested kafka messages

  Scenario: Get published-index request to zebedee is successful
    Given all of the downstream services are healthy
    And get published-index request to zebedee is successful
    And get requests to dataset-api is successful
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | ons              | ed12323eddd3232  |
    Then the URLs of zebedee and dataset documents are retrieved successfully

  Scenario: Error received from get published-index request to zebedee
    Given all of the downstream services are healthy
    And get published-index request to zebedee is unsuccessful
    And get requests to dataset-api is successful
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | ons              | ed12323eddd3232  |
    Then the URLs of zebedee documents are not retrieved

  Scenario: Error received from request to dataset-api
    Given all of the downstream services are healthy
    And get published-index request to zebedee is successful
    And get requests to dataset-api is unsuccessful
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | ons              | ed12323eddd3232  |
    Then the URLs of dataset documents are not retrieved
