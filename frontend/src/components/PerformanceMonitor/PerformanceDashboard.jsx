import React, { useState, useEffect, memo } from 'react';
import styles from './PerformanceDashboard.module.css';
const PerformanceDashboard = memo(() => {
  const [metrics, setMetrics] = useState({
    apiCalls: {},
    renderTimes: {},
    errors: []
  });
  useEffect(() => {
    // Update metrics every second
    const interval = setInterval(() => {
      const currentMetrics = performanceMetrics.getMetrics();
      setMetrics(currentMetrics);
    }, 1000);
    return () => clearInterval(interval);
  }, []);
  const calculateAverage = (times) => {
    if (!times || times.length === 0) return 0;
    return times.reduce((a, b) => a + b, 0) / times.length;
  };
  const formatTime = (ms) => {
    return `${ms.toFixed(2)}ms`;
  };
  return (
    <div className={styles.dashboard}>
      <h2>Performance Monitor</h2>
      {/* API Calls Section */}
      <section className={styles.section}>
        <h3>API Calls</h3>
        <div className={styles.grid}>
          {Object.entries(metrics.apiCalls).map(([endpoint, times]) => (
            <div key={endpoint} className={styles.card}>
              <h4>{endpoint}</h4>
              <p>Calls: {times.length}</p>
              <p>Avg Time: {formatTime(calculateAverage(times))}</p>
            </div>
          ))}
        </div>
      </section>
      {/* Render Times Section */}
      <section className={styles.section}>
        <h3>Component Render Times</h3>
        <div className={styles.grid}>
          {Object.entries(metrics.renderTimes).map(([component, times]) => (
            <div key={component} className={styles.card}>
              <h4>{component}</h4>
              <p>Renders: {times.length}</p>
              <p>Avg Time: {formatTime(calculateAverage(times))}</p>
            </div>
          ))}
        </div>
      </section>
      {/* Errors Section */}
      <section className={styles.section}>
        <h3>Errors</h3>
        <div className={styles.errorList}>
          {metrics.errors.map((error, index) => (
            <div key={index} className={styles.errorCard}>
              <p className={styles.errorTime}>{new Date(error.timestamp).toLocaleTimeString()}</p>
              <p className={styles.errorMessage}>{error.message}</p>
              {error.stack && (
                <pre className={styles.errorStack}>{error.stack}</pre>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
});
PerformanceDashboard.displayName = 'PerformanceDashboard';
export default PerformanceDashboard; 