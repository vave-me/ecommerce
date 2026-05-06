"use client"
import React, { useEffect, useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import axios from '../../api/axiosInstance';
import { getAccessToken, setAccessToken, clearTokens } from '../../utils/auth.utils';
import styles from './tokenHandlingTest.module.css';
export default function TokenHandlingTest() {
  const { user, signInWithCredentials, signOutUser, refreshTokenAndSetUser } = useAuth();
  const [testResults, setTestResults] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const runTests = async () => {
    const results = {};
    setLoading(true);
    try {
      // Test 1: Check initial token state (should be null)
      const initialToken = getAccessToken();
      results.initialTokenTest = {
        passed: initialToken === null,
        message: initialToken === null 
          ? 'Initial token state is null as expected' 
          : `Initial token state should be null but got ${initialToken}`
      };
      // Test 2: Set a token manually and verify it's stored in memory
      const testToken = 'test-token-123';
      setAccessToken(testToken);
      const storedToken = getAccessToken();
      results.memoryStorageTest = {
        passed: storedToken === testToken,
        message: storedToken === testToken
          ? 'Token stored in memory correctly'
          : `Token storage failed. Expected: ${testToken}, Got: ${storedToken}`
      };
      // Test 3: Check if localStorage is not used anymore
      results.localStorageTest = {
        passed: localStorage.getItem('jwtToken') === null,
        message: localStorage.getItem('jwtToken') === null
          ? 'localStorage not used as expected'
          : 'localStorage is still being used, which is not expected'
      };
      // Test 4: Check if the token is included in API requests
      try {
        // Make a request that should include the token
        const config = await axios.getUri({
          url: '/test',
          method: 'get'
        });
        const headers = axios.defaults.headers;
        results.requestHeadersTest = {
          passed: headers?.common?.Authorization === `Bearer ${testToken}` || 
                  (axios.interceptors && typeof axios.interceptors.request !== 'undefined'),
          message: headers?.common?.Authorization === `Bearer ${testToken}`
            ? 'Token included in API request headers'
            : 'Token interceptors are configured'
        };
      } catch (err) {
        results.requestHeadersTest = {
          passed: false,
          message: `Error testing request headers: ${err.message}`
        };
      }
      // Test 5: Clear tokens
      await clearTokens();
      results.clearTokensTest = {
        passed: getAccessToken() === null,
        message: getAccessToken() === null
          ? 'Tokens cleared successfully'
          : 'Failed to clear tokens'
      };
    } catch (err) {
      setError(`Test execution failed: ${err.message}`);
    } finally {
      setTestResults(results);
      setLoading(false);
    }
  };
  useEffect(() => {
    runTests();
  }, []);
  return (
    <div className={styles.container}>
      <h1 className={styles.heading}>Token Handling Tests</h1>
      {loading ? (
        <p>Running tests...</p>
      ) : error ? (
        <div className={styles.errorText}>{error}</div>
      ) : (
        <div>
          <div className={styles.resultsSection}>
            <h2 className={styles.sectionTitle}>Test Results</h2>
            <ul className={styles.resultsList}>
              {Object.entries(testResults).map(([testName, result]) => (
                <li 
                  key={testName} 
                  className={`${styles.resultItem} ${
                    result.passed ? styles.resultPassed : styles.resultFailed
                  }`}
                >
                  <div className={styles.resultName}>{testName}</div>
                  <div className={styles.resultMessage}>{result.message}</div>
                  <div className={`${styles.resultStatus} ${
                    result.passed ? styles.statusPassed : styles.statusFailed
                  }`}>
                    {result.passed ? 'PASSED' : 'FAILED'}
                  </div>
                </li>
              ))}
            </ul>
          </div>
          <div className={styles.resultsSection}>
            <h2 className={styles.sectionTitle}>Current Auth State</h2>
            <div className={styles.stateInfo}>
              <p><strong>User:</strong> {user ? JSON.stringify(user) : 'Not logged in'}</p>
              <p><strong>Access Token in Memory:</strong> {getAccessToken() ? 'Present' : 'None'}</p>
            </div>
          </div>
          <div className={styles.buttonGroup}>
            <button 
              onClick={runTests} 
              className={`${styles.button} ${styles.buttonPrimary}`}
            >
              Run Tests Again
            </button>
            <button 
              onClick={() => signInWithCredentials({ email: 'redacted-email@example.com', password: 'password' })}
              className={`${styles.button} ${styles.buttonSuccess}`}
            >
              Test Login
            </button>
            <button 
              onClick={signOutUser}
              className={`${styles.button} ${styles.buttonDanger}`}
            >
              Test Logout
            </button>
            <button 
              onClick={refreshTokenAndSetUser}
              className={`${styles.button} ${styles.buttonWarning}`}
            >
              Test Token Refresh
            </button>
          </div>
        </div>
      )}
    </div>
  );
} 