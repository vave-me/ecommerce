/**
 * REAL-TIME PERFORMANCE DASHBOARD
 * Production-ready monitoring component for Core Web Vitals and optimization metrics
 * Displays real-time performance data with actionable insights
 */
import React, { useState, useEffect, useCallback, useMemo, memo } from 'react';
import { Activity, Zap, TrendingUp, AlertTriangle, CheckCircle, Info } from 'lucide-react';
import webVitalsOptimizer from '../../utils/coreWebVitalsOptimizer';
import { useOptimizedState } from '../Enhanced/UltimatePerformanceWrapper';
const PerformanceDashboard = memo(({ 
  isVisible = false, 
  position = 'bottom-right',
  enableAlerts = true,
  autoRefresh = true,
  refreshInterval = 5000 
}) => {
  const [metrics, updateMetrics] = useOptimizedState({
    LCP: { value: 0, rating: 'good', threshold: 2500 },
    FID: { value: 0, rating: 'good', threshold: 100 },
    CLS: { value: 0, rating: 'good', threshold: 0.1 },
    FCP: { value: 0, rating: 'good', threshold: 1800 },
    TTFB: { value: 0, rating: 'good', threshold: 800 }
  }, { throttleMs: 1000 });
  const [resourceMetrics, setResourceMetrics] = useState({
    jsSize: 0,
    cssSize: 0,
    imageSize: 0,
    totalSize: 0,
    loadTime: 0,
    renderTime: 0
  });
  const [alerts, setAlerts] = useState([]);
  const [isExpanded, setIsExpanded] = useState(false);
  // Update metrics from Core Web Vitals optimizer
  const updatePerformanceMetrics = useCallback(() => {
    try {
      const vitalsMetrics = webVitalsOptimizer.getMetrics();
      if (vitalsMetrics && Object.keys(vitalsMetrics).length > 0) {
        updateMetrics(vitalsMetrics);
        // Check for performance alerts
        if (enableAlerts) {
          checkPerformanceAlerts(vitalsMetrics);
        }
      }
    } catch (error) {
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', error);
        }
        // Could set error state here if available
        throw error;
    }
  }, [updateMetrics, enableAlerts]);
  // Calculate resource metrics
  const updateResourceMetrics = useCallback(() => {
    if (typeof window === 'undefined' || !window.performance) return;
    try {
      const resources = performance.getEntriesByType('resource');
      const navigation = performance.getEntriesByType('navigation')[0];
      const resourceSizes = resources.reduce((acc, resource) => {
        const size = resource.transferSize || resource.encodedBodySize || 0;
        if (resource.name.includes('.js')) {
          acc.jsSize += size;
        } else if (resource.name.includes('.css')) {
          acc.cssSize += size;
        } else if (resource.name.match(/\.(jpg|jpeg|png|gif|webp|avif|svg)$/i)) {
          acc.imageSize += size;
        }
        acc.totalSize += size;
        return acc;
      }, { jsSize: 0, cssSize: 0, imageSize: 0, totalSize: 0 });
      setResourceMetrics({
        ...resourceSizes,
        loadTime: navigation?.loadEventEnd - navigation?.loadEventStart || 0,
        renderTime: navigation?.domContentLoadedEventEnd - navigation?.domContentLoadedEventStart || 0
      });
    } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    }
  }, []);
  // Check for performance alerts
  const checkPerformanceAlerts = useCallback((currentMetrics) => {
    const newAlerts = [];
    Object.entries(currentMetrics).forEach(([metric, data]) => {
      if (data.rating === 'poor') {
        newAlerts.push({
          id: `${metric}-${Date.now()}`,
          type: 'error',
          metric,
          message: `${metric} is poor (${Math.round(data.value)}${metric === 'CLS' ? '' : 'ms'})`,
          recommendation: getPerformanceRecommendation(metric)
        });
      } else if (data.rating === 'needs-improvement') {
        newAlerts.push({
          id: `${metric}-${Date.now()}`,
          type: 'warning',
          metric,
          message: `${metric} needs improvement (${Math.round(data.value)}${metric === 'CLS' ? '' : 'ms'})`,
          recommendation: getPerformanceRecommendation(metric)
        });
      }
    });
    // Check resource budget
    if (resourceMetrics.totalSize > 1000 * 1024) { // 1MB
      newAlerts.push({
        id: `bundle-size-${Date.now()}`,
        type: 'warning',
        metric: 'Bundle Size',
        message: `Large bundle size: ${(resourceMetrics.totalSize / 1024).toFixed(1)}KB`,
        recommendation: 'Consider code splitting and lazy loading'
      });
    }
    setAlerts(newAlerts);
  }, [resourceMetrics.totalSize]);
  // Get performance recommendations
  const getPerformanceRecommendation = (metric) => {
    const recommendations = {
      LCP: 'Optimize largest contentful paint by preloading critical resources and optimizing images',
      FID: 'Reduce first input delay by minimizing JavaScript execution time and using web workers',
      CLS: 'Improve cumulative layout shift by specifying image dimensions and avoiding dynamic content',
      FCP: 'Speed up first contentful paint by optimizing critical rendering path',
      TTFB: 'Reduce time to first byte by optimizing server response time and CDN usage'
    };
    return recommendations[metric] || 'Optimize this metric for better performance';
  };
  // Auto-refresh metrics
  useEffect(() => {
    if (!autoRefresh || !isVisible) return;
    updatePerformanceMetrics();
    updateResourceMetrics();
    const interval = setInterval(() => {
      updatePerformanceMetrics();
      updateResourceMetrics();
    }, refreshInterval);
    return () => clearInterval(interval);
  }, [autoRefresh, isVisible, refreshInterval, updatePerformanceMetrics, updateResourceMetrics]);
  // Performance score calculation
  const performanceScore = useMemo(() => {
    const scores = Object.values(metrics).map(metric => {
      switch (metric.rating) {
        case 'good': return 100;
        case 'needs-improvement': return 75;
        case 'poor': return 50;
        default: return 50;
      }
    });
    return Math.round(scores.reduce((a, b) => a + b, 0) / scores.length);
  }, [metrics]);
  // Position styles
  const getPositionStyles = () => {
    const positions = {
      'top-left': { top: '20px', left: '20px' },
      'top-right': { top: '20px', right: '20px' },
      'bottom-left': { bottom: '20px', left: '20px' },
      'bottom-right': { bottom: '20px', right: '20px' }
    };
    return positions[position] || positions['bottom-right'];
  };
  // Metric card component
  const MetricCard = ({ label, value, rating, threshold, unit = 'ms' }) => {
    const getColorByRating = (rating) => {
      switch (rating) {
        case 'good': return '#10b981';
        case 'needs-improvement': return '#f59e0b';
        case 'poor': return '#ef4444';
        default: return '#6b7280';
      }
    };
    return (
      <div className="metric-card" style={{
        padding: '12px',
        borderRadius: '8px',
        backgroundColor: '#f9fafb',
        border: `2px solid ${getColorByRating(rating)}`,
        marginBottom: '8px'
      }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center' 
        }}>
          <span style={{ fontWeight: '600', fontSize: '14px' }}>{label}</span>
          <span style={{ 
            color: getColorByRating(rating), 
            fontWeight: 'bold',
            fontSize: '16px'
          }}>
            {Math.round(value)}{unit}
          </span>
        </div>
        <div style={{ 
          fontSize: '12px', 
          color: '#6b7280', 
          marginTop: '4px' 
        }}>
          Threshold: {threshold}{unit}
        </div>
      </div>
    );
  };
  // Resource summary component
  const ResourceSummary = () => (
    <div style={{ 
      padding: '12px', 
      backgroundColor: '#f3f4f6', 
      borderRadius: '8px',
      marginTop: '12px' 
    }}>
      <h4 style={{ margin: '0 0 8px 0', fontSize: '14px', fontWeight: '600' }}>
        Resource Usage
      </h4>
      <div style={{ fontSize: '12px', color: '#4b5563' }}>
        <div>JS: {(resourceMetrics.jsSize / 1024).toFixed(1)}KB</div>
        <div>CSS: {(resourceMetrics.cssSize / 1024).toFixed(1)}KB</div>
        <div>Images: {(resourceMetrics.imageSize / 1024).toFixed(1)}KB</div>
        <div style={{ fontWeight: '600', marginTop: '4px' }}>
          Total: {(resourceMetrics.totalSize / 1024).toFixed(1)}KB
        </div>
      </div>
    </div>
  );
  // Alert component
  const AlertItem = ({ alert, onDismiss }) => {
    const getAlertIcon = (type) => {
      switch (type) {
        case 'error': return <AlertTriangle size={16} color="#ef4444" />;
        case 'warning': return <Info size={16} color="#f59e0b" />;
        default: return <CheckCircle size={16} color="#10b981" />;
      }
    };
    return (
      <div style={{
        padding: '8px 12px',
        backgroundColor: alert.type === 'error' ? '#fef2f2' : '#fffbeb',
        border: `1px solid ${alert.type === 'error' ? '#fecaca' : '#fed7aa'}`,
        borderRadius: '6px',
        marginBottom: '6px',
        fontSize: '12px'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
          {getAlertIcon(alert.type)}
          <span style={{ fontWeight: '600' }}>{alert.message}</span>
          <button 
            onClick={() => onDismiss(alert.id)}
            style={{
              marginLeft: 'auto',
              background: 'none',
              border: 'none',
              cursor: 'pointer',
              fontSize: '14px'
            }}
          >
            ×
          </button>
        </div>
        <div style={{ marginTop: '4px', color: '#6b7280' }}>
          {alert.recommendation}
        </div>
      </div>
    );
  };
  if (!isVisible && process.env.NODE_ENV === 'production') {
    return null;
  }
  return (
    <div 
      style={{
        position: 'fixed',
        ...getPositionStyles(),
        zIndex: 9999,
        backgroundColor: 'white',
        borderRadius: '12px',
        boxShadow: '0 10px 25px rgba(0,0,0,0.15)',
        border: '1px solid #e5e7eb',
        width: isExpanded ? '320px' : '200px',
        maxHeight: '80vh',
        overflow: 'auto',
        fontFamily: 'system-ui, -apple-system, sans-serif'
      }}
    >
      {/* Header */}
      <div style={{
        padding: '12px 16px',
        borderBottom: '1px solid #e5e7eb',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        cursor: 'pointer'
      }} onClick={() => setIsExpanded(!isExpanded)}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <Activity size={18} color="#3b82f6" />
          <span style={{ fontWeight: '600', fontSize: '14px' }}>
            Performance
          </span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <div style={{
            backgroundColor: performanceScore >= 90 ? '#10b981' : 
                           performanceScore >= 70 ? '#f59e0b' : '#ef4444',
            color: 'white',
            padding: '2px 8px',
            borderRadius: '12px',
            fontSize: '12px',
            fontWeight: '600'
          }}>
            {performanceScore}
          </div>
          <span style={{ fontSize: '18px', color: '#9ca3af' }}>
            {isExpanded ? '−' : '+'}
          </span>
        </div>
      </div>
      {/* Expanded content */}
      {isExpanded && (
        <div style={{ padding: '16px' }}>
          {/* Core Web Vitals */}
          <div style={{ marginBottom: '16px' }}>
            <h3 style={{ 
              margin: '0 0 12px 0', 
              fontSize: '16px', 
              fontWeight: '600',
              display: 'flex',
              alignItems: 'center',
              gap: '6px'
            }}>
              <Zap size={16} color="#3b82f6" />
              Core Web Vitals
            </h3>
            {Object.entries(metrics).map(([key, metric]) => (
              <MetricCard
                key={key}
                label={key}
                value={metric.value}
                rating={metric.rating}
                threshold={metric.threshold}
                unit={key === 'CLS' ? '' : 'ms'}
              />
            ))}
          </div>
          {/* Resource Usage */}
          <ResourceSummary />
          {/* Performance Alerts */}
          {alerts.length > 0 && (
            <div style={{ marginTop: '16px' }}>
              <h4 style={{ 
                margin: '0 0 8px 0', 
                fontSize: '14px', 
                fontWeight: '600',
                display: 'flex',
                alignItems: 'center',
                gap: '6px'
              }}>
                <TrendingUp size={14} color="#f59e0b" />
                Alerts ({alerts.length})
              </h4>
              {alerts.map(alert => (
                <AlertItem
                  key={alert.id}
                  alert={alert}
                  onDismiss={(id) => setAlerts(alerts.filter(a => a.id !== id))}
                />
              ))}
            </div>
          )}
          {/* Quick Actions */}
          <div style={{ 
            marginTop: '16px', 
            paddingTop: '12px', 
            borderTop: '1px solid #e5e7eb' 
          }}>
            <button
              onClick={() => {
                updatePerformanceMetrics();
                updateResourceMetrics();
              }}
              style={{
                width: '100%',
                padding: '8px 12px',
                backgroundColor: '#3b82f6',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                fontSize: '12px',
                fontWeight: '600',
                cursor: 'pointer'
              }}
            >
              Refresh Metrics
            </button>
          </div>
        </div>
      )}
    </div>
  );
});
export default PerformanceDashboard; 