'use client';

import React from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import DemoLayout from '../components/DemoLayout';
import FeatureCard from '../components/FeatureCard';
import MetricCard from '../components/MetricCard';
import styles from './SocialPage.module.css';

const MessagingIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
  </svg>
);

const FeedIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
    <line x1="9" y1="9" x2="15" y2="9"/>
    <line x1="9" y1="13" x2="15" y2="13"/>
    <line x1="9" y1="17" x2="15" y2="17"/>
  </svg>
);

const FollowersIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
    <circle cx="8.5" cy="7" r="4"/>
    <line x1="20" y1="8" x2="20" y2="14"/>
    <line x1="23" y1="11" x2="17" y2="11"/>
  </svg>
);

const ReviewIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
  </svg>
);

const SocialPage = () => {
  const t = useTranslations('Demo');
  const ts = useTranslations('shared');
  const breadcrumbs = [
    { label: t('breadcrumb.socialLayer') }
  ];

  return (
    <DemoLayout
      title={t('socialPage.title')}
      subtitle={t('socialPage.subtitle')}
      breadcrumbs={breadcrumbs}
      showAd={true}
    >
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('socialPage.features.title')}</h2>
        
        <div className={styles.featureGrid}>
          <FeatureCard
            icon={<MessagingIcon />}
            title={t('socialPage.features.messaging.title')}
            description={t('socialPage.features.messaging.description')}
            highlight={true}
          />
          
          <FeatureCard
            icon={<FeedIcon />}
            title={t('socialPage.features.feeds.title')}
            description={t('socialPage.features.feeds.description')}
          />
          
          <FeatureCard
            icon={<FollowersIcon />}
            title={t('socialPage.features.following.title')}
            description={t('socialPage.features.following.description')}
          />
          
          <FeatureCard
            icon={<ReviewIcon />}
            title={t('socialPage.features.reviews.title')}
            description={t('socialPage.features.reviews.description')}
          />
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('socialPage.metrics.title')}</h2>
        
        <div className={styles.metricsGrid}>
          <MetricCard
            value="<100"
            suffix="ms"
            label={t('socialPage.metrics.messageDelivery')}
          />
          <MetricCard
            value="50K+"
            label={t('socialPage.metrics.connections')}
          />
          <MetricCard
            value="3x"
            label={t('socialPage.metrics.engagement')}
            trend={25}
          />
          <MetricCard
            value="92"
            suffix="%"
            label={t('socialPage.metrics.retention')}
            trend={15}
          />
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('socialPage.tools.title')}</h2>
        
        <div className={styles.toolsGrid}>
          <div className={styles.toolCard}>
            <h3 className={styles.toolTitle}>{t('socialPage.tools.comments.title')}</h3>
            <p className={styles.toolDescription}>
              {t('socialPage.tools.comments.description')}
            </p>
            <ul className={styles.toolFeatures}>
              <li>{t('socialPage.tools.comments.features.nested')}</li>
              <li>{t('socialPage.tools.comments.features.mentions')}</li>
              <li>{t('socialPage.tools.comments.features.spam')}</li>
              <li>{t('socialPage.tools.comments.features.moderation')}</li>
            </ul>
          </div>
          
          <div className={styles.toolCard}>
            <h3 className={styles.toolTitle}>{t('socialPage.tools.profiles.title')}</h3>
            <p className={styles.toolDescription}>
              {t('socialPage.tools.profiles.description')}
            </p>
            <ul className={styles.toolFeatures}>
              <li>{t('socialPage.tools.profiles.features.timeline')}</li>
              <li>{t('socialPage.tools.profiles.features.portfolio')}</li>
              <li>{t('socialPage.tools.profiles.features.verification')}</li>
              <li>{t('socialPage.tools.profiles.features.bio')}</li>
            </ul>
          </div>
          
          <div className={styles.toolCard}>
            <h3 className={styles.toolTitle}>{ts('notifications')}</h3>
            <p className={styles.toolDescription}>
              {t('socialPage.tools.notifications.description')}
            </p>
            <ul className={styles.toolFeatures}>
              <li>{t('socialPage.tools.notifications.features.inApp')}</li>
              <li>{t('socialPage.tools.notifications.features.email')}</li>
              <li>{t('socialPage.tools.notifications.features.push')}</li>
              <li>{t('socialPage.tools.notifications.features.preferences')}</li>
            </ul>
          </div>
        </div>
      </section>

      <section className={styles.relatedSection}>
        <h3 className={styles.relatedTitle}>{t('socialPage.enhance')}</h3>
        <div className={styles.relatedLinks}>
          <Link href="/demo/commerce" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>🛒</span>
            <div>
              <h4>{t('modules.commerce.title')}</h4>
              <p>{t('socialPage.relatedCommerce')}</p>
            </div>
          </Link>
          <Link href="/demo/intelligence" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>🤖</span>
            <div>
              <h4>{t('modules.intelligence.title')}</h4>
              <p>{t('socialPage.relatedIntelligence')}</p>
            </div>
          </Link>
          <Link href="/demo/verticals" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>📱</span>
            <div>
              <h4>{t('modules.verticals.title')}</h4>
              <p>{t('socialPage.relatedVerticals')}</p>
            </div>
          </Link>
        </div>
      </section>
    </DemoLayout>
  );
};

export default SocialPage;