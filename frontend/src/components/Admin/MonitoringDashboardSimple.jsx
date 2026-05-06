import React, { useState, useEffect, useCallback } from 'react';
import { RefreshCw, AlertCircle, CheckCircle, TrendingUp, TrendingDown } from 'lucide-react';
import styles from './MonitoringDashboard.module.css';

/**
 * Simple monitoring dashboard without external chart dependencies
 */
export default function MonitoringDashboardSimple() {
  const [metrics, setMetrics] = useState({
    errors: [],
    performance: {},
    vitals: {},
    system: {}
  });
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState('1h');
  const [autoRefresh, setAutoRefresh] = useState(true);

  // Fetch monitoring data
  const fetchMetrics = useCallback(async () => {
    try {
      // In production, this would fetch from your monitoring API
      // For now, using mock data
      setMetrics({
        errors: [
          { type: 'NETWORK', message: 'API timeout', timestamp: Date.now() - 3600000, severity: 'warning' },
          { type: 'VALIDATION', message: 'Invalid form input', timestamp: Date.now() - 7200000, severity: 'low' }
        ],
        performance: {
          avgLoadTime: 1234,
          p95LoadTime: 2345,
          errorRate: 0.02
        },
        vitals: {
          LCP: { average: 2100, status: 'good' },
          FID: { average: 85, status: 'good' },
          CLS: { average: 0.08, status: 'good' },
          TTFB: { average: 650, status: 'good' }
        },
        system: {
          uptime: 99.95,
          cpu: 45,
          memory: 62
        }
      });
    } catch (error) {
      console.error('Failed to fetch monitoring data:', error);
    } finally {
      setLoading(false);
    }
  }, [timeRange]);

  // Auto-refresh
  useEffect(() => {
    fetchMetrics();
    
    if (autoRefresh) {
      const interval = setInterval(fetchMetrics, 30000);
      return () => clearInterval(interval);
    }
  }, [fetchMetrics, autoRefresh]);

  const calculateErrorRate = () => {
    const total = metrics.errors.length;
    const recent = metrics.errors.filter(e => 
      new Date(e.timestamp) > new Date(Date.now() - 3600000)
    ).length;
    
    return {
      total,
      recent,
      trend: recent > total / 24 ? 'up' : 'down'
    };
  };

  if (loading) {
    return (
      <div className={styles.loading}>
        <RefreshCw className="animate-spin" />
        <p>Loading monitoring data...</p>
      </div>
    );
  }

  const errorRate = calculateErrorRate();

  return (
    <div className={styles.dashboard}>
      {/* Header */}
      <div className={styles.header}>
        <h1>Production Monitoring</h1>
        <div className={styles.controls}>
          <select 
            value={timeRange} 
            onChange={(e) => setTimeRange(e.target.value)}
            className={styles.timeSelect}
          >
            <option value="1h">Last Hour</option>
            <option value="6h">Last 6 Hours</option>
            <option value="24h">Last 24 Hours</option>
            <option value="7d">Last 7 Days</option>
          </select>
          <button 
            className={autoRefresh ? styles.buttonActive : styles.button}
            onClick={() => setAutoRefresh(!autoRefresh)}
          >
            Auto-refresh: {autoRefresh ? 'ON' : 'OFF'}
          </button>
          <button onClick={fetchMetrics} className={styles.button}>
            <RefreshCw size={16} />
            Refresh
          </button>
        </div>
      </div>

      {/* Key Metrics */}
      <div className={styles.metricsGrid}>
        <div className={styles.card}>
          <h3>Error Rate</h3>
          <div className={styles.metric}>
            <div className={styles.metricValue}>
              {errorRate.total}
              {errorRate.trend === 'up' ? 
                <TrendingUp className={styles.trendUp} /> : 
                <TrendingDown className={styles.trendDown} />
              }
            </div>
            <div className={styles.metricLabel}>Total Errors</div>
            <div className={styles.metricSub}>
              {errorRate.recent} in last hour
            </div>
          </div>
        </div>

        <div className={styles.card}>
          <h3>System Health</h3>
          <div className={styles.healthIndicators}>
            <div className={styles.healthItem}>
              <CheckCircle className={styles.healthGood} />
              <span>API</span>
            </div>
            <div className={styles.healthItem}>
              <CheckCircle className={styles.healthGood} />
              <span>Database</span>
            </div>
            <div className={styles.healthItem}>
              <AlertCircle className={styles.healthWarning} />
              <span>Cache</span>
            </div>
          </div>
        </div>

        <div className={styles.card}>
          <h3>Performance</h3>
          <div className={styles.metric}>
            <div className={styles.metricValue}>
              {metrics.performance.avgLoadTime}ms
            </div>
            <div className={styles.metricLabel}>Avg Load Time</div>
            <div className={styles.metricSub}>
              P95: {metrics.performance.p95LoadTime}ms
            </div>
          </div>
        </div>

        <div className={styles.card}>
          <h3>Uptime</h3>
          <div className={styles.metric}>
            <div className={styles.metricValue}>{metrics.system.uptime}%</div>
            <div className={styles.metricLabel}>Last 30 days</div>
            <div className={styles.metricSub}>
              CPU: {metrics.system.cpu}% | Memory: {metrics.system.memory}%
            </div>
          </div>
        </div>
      </div>

      {/* Core Web Vitals */}
      <div className={styles.card}>
        <h3>Core Web Vitals</h3>
        <div className={styles.vitalsGrid}>
          {Object.entries(metrics.vitals).map(([metric, data]) => (
            <div key={metric} className={styles.vitalItem}>
              <span className={styles.vitalMetric}>{metric}</span>
              <span className={`${styles.vitalValue} ${styles[data.status]}`}>
                {data.average}{metric === 'CLS' ? '' : 'ms'}
              </span>
              <span className={`${styles.vitalStatus} ${styles[data.status]}`}>
                {data.status}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Recent Errors */}
      <div className={styles.card}>
        <h3>Recent Errors</h3>
        <div className={styles.errorList}>
          {metrics.errors.length === 0 ? (
            <p className={styles.noErrors}>No recent errors</p>
          ) : (
            metrics.errors.slice(0, 10).map((error, index) => (
              <div key={index} className={styles.errorItem}>
                <div className={styles.errorHeader}>
                  <span className={`${styles.badge} ${styles[error.severity]}`}>
                    {error.type}
                  </span>
                  <span className={styles.errorTime}>
                    {new Date(error.timestamp).toLocaleString()}
                  </span>
                </div>
                <div className={styles.errorMessage}>{error.message}</div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}