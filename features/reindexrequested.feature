Feature: Consuming reindex-requested kafka messages

  Scenario: Get published-index request to zebedee is successful
    Given all of the downstream services are healthy
    And get published-index request to zebedee is successful
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | ons              | ed12323eddd3232  |
    Then the URLs of zebedee documents are retrieved successfully

  Scenario: Error received from get published-index request to zebedee
    Given all of the downstream services are healthy
    And get published-index request to zebedee is unsuccessful
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | ons              | ed12323eddd3232  |
    Then the URLs of zebedee documents are not retrieved
