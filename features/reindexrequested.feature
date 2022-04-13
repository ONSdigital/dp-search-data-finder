Feature: ReindexRequested


  Scenario: Posting and checking a response
    When these reindex-requested events are consumed:
            | JobID |   SearchIndex |   TraceID |
            | 1234  |   ons         |   ed12323eddd3232 |
    Then I should receive a reindex-requested response