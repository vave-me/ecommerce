import React from 'react';
import styles from './MetricCard.module.css';

const MetricCard = ({ value, label, suffix = '', prefix = '', trend = null }) => {
  return (
    <div className={styles.metricCard}>
      <div className={styles.value}>
        {prefix}
        <span className={styles.number}>{value}</span>
        {suffix}
      </div>
      <div className={styles.label}>{label}</div>
      {trend && (
        <div className={`${styles.trend} ${trend > 0 ? styles.positive : styles.negative}`}>
          {trend > 0 ? '↑' : '↓'} {Math.abs(trend)}%
        </div>
      )}
    </div>
  );
};

export default MetricCard;