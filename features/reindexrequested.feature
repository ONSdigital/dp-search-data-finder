Feature: ReindexRequested



  Scenario: Posting and checking a response
    When these reindex-requested events are consumed:
            | JobID |
            | 1234  |
    Then I should receive a reindex-requested response