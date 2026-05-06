"use client"
import React, { useState, useEffect } from 'react';
import { 
  getAccessToken, 
  setAccessToken, 
  clearTokens 
} from '../../utils/auth.utils';
import styles from '../auth/tokenHandlingTest.module.css';
/**
 * A simple component to test token storage and handling
 * This doesn't rely on the AuthContext or other components,
 * making it easier to isolate and test the token utilities
 */
export default function MockTokenTest() {
  const [results, setResults] = useState([]);
  const [currentToken, setCurrentToken] = useState(null);
  // Test setting and retrieving a token in memory
  const testSetGetToken = () => {
    const testToken = 'test-token-123';
    setAccessToken(testToken);
    const retrievedToken = getAccessToken();
    const success = retrievedToken === testToken;
    setResults(prev => [...prev, {
      name: 'In-Memory Token Storage',
      passed: success,
      message: success 
        ? 'Token successfully stored and retrieved from memory' 
        : `Token retrieval failed. Expected: ${testToken}, Got: ${retrievedToken}`
    }]);
    // Update display
    setCurrentToken(retrievedToken);
  };
  // Test that localStorage is not being used
  const testLocalStorageNotUsed = () => {
    const testToken = 'test-localStorage-token';
    // First clear any existing token
    clearTokens();
    // Then set a new token
    setAccessToken(testToken);
    // Check if localStorage has the token (it should NOT)
    const inLocalStorage = localStorage.getItem('jwtToken') === testToken;
    const success = !inLocalStorage;
    setResults(prev => [...prev, {
      name: 'localStorage Not Used',
      passed: success,
      message: success 
        ? 'localStorage correctly not being used for token storage'
        : 'Token was found in localStorage, which is incorrect'
    }]);
  };
  // Test clearing tokens
  const testClearTokens = async () => {
    const testToken = 'test-clear-token';
    // Set a token
    setAccessToken(testToken);
    // Clear it
    await clearTokens();
    // Check if cleared
    const retrievedToken = getAccessToken();
    const success = retrievedToken === null;
    setResults(prev => [...prev, {
      name: 'Clear Tokens',
      passed: success,
      message: success 
        ? 'Token successfully cleared from memory' 
        : `Token still exists after clearing. Got: ${retrievedToken}`
    }]);
    // Update display
    setCurrentToken(retrievedToken);
  };
  // Run all tests
  const runAllTests = () => {
    // Clear previous results
    setResults([]);
    // Run tests
    testSetGetToken();
    testLocalStorageNotUsed();
    testClearTokens();
  };
  return (
    <div className={styles.container}>
      <h1 className={styles.heading}>Token Implementation Test</h1>
      <div className={styles.resultsSection}>
        <div className={styles.stateInfo}>
          <p><strong>Current Memory Token:</strong> {currentToken || 'None'}</p>
          <p><strong>localStorage Token:</strong> {localStorage.getItem('jwtToken') || 'None'}</p>
        </div>
        <h2 className={styles.sectionTitle}>Test Results</h2>
        {results.length === 0 ? (
          <p>No tests run yet. Click "Run Tests" to start.</p>
        ) : (
          <ul className={styles.resultsList}>
            {results.map((result, index) => (
              <li 
                key={index} 
                className={`${styles.resultItem} ${
                  result.passed ? styles.resultPassed : styles.resultFailed
                }`}
              >
                <div className={styles.resultName}>{result.name}</div>
                <div className={styles.resultMessage}>{result.message}</div>
                <div className={`${styles.resultStatus} ${
                  result.passed ? styles.statusPassed : styles.statusFailed
                }`}>
                  {result.passed ? 'PASSED' : 'FAILED'}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
      <div className={styles.buttonGroup}>
        <button 
          onClick={runAllTests}
          className={`${styles.button} ${styles.buttonPrimary}`}
        >
          Run All Tests
        </button>
        <button 
          onClick={testSetGetToken}
          className={`${styles.button} ${styles.buttonSuccess}`}
        >
          Test Set/Get Token
        </button>
        <button 
          onClick={testLocalStorageNotUsed}
          className={`${styles.button} ${styles.buttonSuccess}`}
        >
          Test localStorage Not Used
        </button>
        <button 
          onClick={testClearTokens}
          className={`${styles.button} ${styles.buttonDanger}`}
        >
          Test Clear Tokens
        </button>
      </div>
    </div>
  );
} 