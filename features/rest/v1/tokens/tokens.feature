Feature: get tokens
  In order to know the tokens in the pools
  As an API user
  I need to be able to request the tokens endpoint

  Scenario: does not allow POST method
    When I send "POST" request to "/version"
    Then the response code should be 405
    And the response should match json:
      """
      {
        "error": "Method not allowed"
      }
      """

  Scenario: should get tokens array
    When I send "GET" request to "/v1/tokens"
    Then the response code should be 200
    And the response should match json:
      """
      [
        "BNB",
        "ERD-D85",
        "FSN-F1B",
        "FTM-585",
        "LOK-3C0",
        "TCAN-014",
        "TOMOB-1E1"
      ]
      """

