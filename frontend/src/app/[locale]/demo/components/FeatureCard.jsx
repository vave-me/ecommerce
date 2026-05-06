import React from 'react';
import styles from './FeatureCard.module.css';

const FeatureCard = ({ icon, title, description, highlight = false, badge = null }) => {
  return (
    <div className={`${styles.featureCard} ${highlight ? styles.highlight : ''}`}>
      {badge && <span className={styles.badge}>{badge}</span>}
      <div className={styles.iconWrapper}>
        {typeof icon === 'string' ? (
          <span className={styles.iconEmoji}>{icon}</span>
        ) : (
          icon
        )}
      </div>
      <h3 className={styles.title}>{title}</h3>
      <p className={styles.description}>{description}</p>
    </div>
  );
};

export default FeatureCard;