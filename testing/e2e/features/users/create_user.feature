Feature: Create Store

  As a visitor
  I should be able to create new user

  Scenario: Creating a user called "Monthy"
    Given a valid user email "redacted-email@example.com"
    And no user with "monty@example" exists
    When I create the user called "Monthy"
    Then a user called "Monthy" exists
