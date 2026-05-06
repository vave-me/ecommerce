import React, { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Alert, AlertDescription } from '../ui/alert';
import { Badge } from '../ui/badge';
import { Button } from '../ui/button';
import { 
  BarChart, Bar, LineChart, Line, PieChart, Pie, Cell,
  XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer 
} from 'recharts';
import { RefreshCw, AlertCircle, CheckCircle, XCircle, TrendingUp, TrendingDown } from 'lucide-react';
import styles from './MonitoringDashboard.module.css';

/**
 * Real-time monitoring dashboard for production metrics
 */
export default function MonitoringDashboard() {
  const [metrics, setMetrics] = useState({
    errors: [],
    performance: {},
    vitals: {},
    system: {}
  });
  const [loading, setLoading] = useState(true);
  const [timeRange, setTimeRange] = useState('1h'); // 1h, 6h, 24h, 7d
  const [autoRefresh, setAutoRefresh] = useState(true);

  // Fetch monitoring data
  const fetchMetrics = useCallback(async () => {
    try {
      // In a real app, this would fetch from your monitoring API
      const response = await fetch(`/api/admin/monitoring?range=${timeRange}`);
      const data = await response.json();
      
      setMetrics({
        errors: data.errors || [],
        performance: data.performance || {},
        vitals: data.vitals || {},
        system: data.system || {}
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
      const interval = setInterval(fetchMetrics, 30000); // 30 seconds
      return () => clearInterval(interval);
    }
  }, [fetchMetrics, autoRefresh]);

  // Error rate calculation
  const calculateErrorRate = () => {
    const total = metrics.errors.length;
    const last24h = metrics.errors.filter(e => 
      new Date(e.timestamp) > new Date(Date.now() - 24 * 60 * 60 * 1000)
    ).length;
    
    return {
      total,
      last24h,
      trend: last24h > total / 7 ? 'up' : 'down'
    };
  };

  // Core Web Vitals status
  const getVitalStatus = (metric, value) => {
    const thresholds = {
      LCP: { good: 2500, poor: 4000 },
      FID: { good: 100, poor: 300 },
      CLS: { good: 0.1, poor: 0.25 },
      FCP: { good: 1800, poor: 3000 },
      TTFB: { good: 800, poor: 1800 }
    };
    
    const threshold = thresholds[metric];
    if (!threshold) return 'unknown';
    
    if (value <= threshold.good) return 'good';
    if (value <= threshold.poor) return 'needs-improvement';
    return 'poor';
  };

  // Chart data preparation
  const prepareErrorChart = () => {
    const hourly = {};
    metrics.errors.forEach(error => {
      const hour = new Date(error.timestamp).getHours();
      hourly[hour] = (hourly[hour] || 0) + 1;
    });
    
    return Object.entries(hourly).map(([hour, count]) => ({
      hour: `${hour}:00`,
      errors: count
    }));
  };

  const prepareVitalsChart = () => {
    return Object.entries(metrics.vitals).map(([metric, data]) => ({
      metric,
      value: data.average || 0,
      status: getVitalStatus(metric, data.average || 0)
    }));
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
        <h1>Production Monitoring Dashboard</h1>
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
          <Button 
            variant={autoRefresh ? "default" : "outline"}
            onClick={() => setAutoRefresh(!autoRefresh)}
          >
            Auto-refresh: {autoRefresh ? 'ON' : 'OFF'}
          </Button>
          <Button onClick={fetchMetrics} variant="outline">
            <RefreshCw size={16} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Key Metrics */}
      <div className={styles.metricsGrid}>
        <Card>
          <CardHeader>
            <CardTitle>Error Rate</CardTitle>
          </CardHeader>
          <CardContent>
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
                {errorRate.last24h} in last 24h
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Health</CardTitle>
          </CardHeader>
          <CardContent>
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
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Performance Score</CardTitle>
          </CardHeader>
          <CardContent>
            <div className={styles.performanceScore}>
              <div className={styles.scoreCircle}>
                <span className={styles.scoreValue}>87</span>
                <span className={styles.scoreLabel}>/100</span>
              </div>
              <Badge variant="success">Good</Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Uptime</CardTitle>
          </CardHeader>
          <CardContent>
            <div className={styles.metric}>
              <div className={styles.metricValue}>99.95%</div>
              <div className={styles.metricLabel}>Last 30 days</div>
              <div className={styles.metricSub}>
                2 incidents resolved
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts */}
      <div className={styles.chartsGrid}>
        {/* Error Timeline */}
        <Card className={styles.chartCard}>
          <CardHeader>
            <CardTitle>Error Timeline</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={prepareErrorChart()}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="hour" />
                <YAxis />
                <Tooltip />
                <Line 
                  type="monotone" 
                  dataKey="errors" 
                  stroke="#ef4444" 
                  strokeWidth={2}
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Core Web Vitals */}
        <Card className={styles.chartCard}>
          <CardHeader>
            <CardTitle>Core Web Vitals</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={prepareVitalsChart()}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="metric" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="value" fill="#3b82f6">
                  {prepareVitalsChart().map((entry, index) => (
                    <Cell 
                      key={`cell-${index}`} 
                      fill={
                        entry.status === 'good' ? '#10b981' :
                        entry.status === 'needs-improvement' ? '#f59e0b' :
                        '#ef4444'
                      }
                    />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Recent Errors */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Errors</CardTitle>
        </CardHeader>
        <CardContent>
          <div className={styles.errorList}>
            {metrics.errors.slice(0, 10).map((error, index) => (
              <div key={index} className={styles.errorItem}>
                <div className={styles.errorHeader}>
                  <Badge variant={error.severity === 'critical' ? 'destructive' : 'secondary'}>
                    {error.type}
                  </Badge>
                  <span className={styles.errorTime}>
                    {new Date(error.timestamp).toLocaleString()}
                  </span>
                </div>
                <div className={styles.errorMessage}>{error.message}</div>
                <div className={styles.errorMeta}>
                  <span>User: {error.userId || 'Anonymous'}</span>
                  <span>Context: {error.context}</span>
                  {error.count > 1 && (
                    <Badge variant="outline">{error.count} occurrences</Badge>
                  )}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* System Alerts */}
      {metrics.alerts && metrics.alerts.length > 0 && (
        <Alert variant="warning">
          <AlertCircle />
          <AlertDescription>
            <strong>Active Alerts:</strong>
            <ul>
              {metrics.alerts.map((alert, index) => (
                <li key={index}>{alert.message}</li>
              ))}
            </ul>
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}