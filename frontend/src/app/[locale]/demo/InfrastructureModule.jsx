'use client';

import React from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import DemoLayout from './components/DemoLayout';
import FeatureCard from './components/FeatureCard';
import MetricCard from './components/MetricCard';
import styles from './InfrastructureModule.module.css';

const MicroserviceIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2z"/>
    <path d="M12 16.5A4.5 4.5 0 1 0 7.5 12"/>
    <path d="M12 12h4.5"/>
    <path d="M16.5 12V7.5"/>
  </svg>
);

const ApiIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M18 8V6a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v2"/>
    <path d="M6 8v8a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V8"/>
    <path d="M10 12h4"/>
  </svg>
);

const ScaleIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M21 12H3M3 12l4-4M3 12l4 4M21 12l-4-4M21 12l-4 4"/>
    <path d="M12 21V3M12 3l4 4M12 3L8 7M12 21l4-4M12 21l-4 4"/>
  </svg>
);

const InfrastructureModule = () => {
  const t = useTranslations('Demo');
  const breadcrumbs = [
    { label: t('breadcrumb.infrastructureModule') }
  ];

  return (
    <DemoLayout
      title={t('infrastructurePage.title')}
      subtitle={t('infrastructurePage.subtitle')}
      breadcrumbs={breadcrumbs}
      showAd={true}
    >
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('infrastructurePage.pillars.title')}</h2>
        
        <div className={styles.featureGrid}>
          <FeatureCard
            icon={<MicroserviceIcon />}
            title={t('infrastructurePage.pillars.microservices.title')}
            description={t('infrastructurePage.pillars.microservices.description')}
          />
          
          <FeatureCard
            icon={<ApiIcon />}
            title={t('infrastructurePage.pillars.apiFirst.title')}
            description={t('infrastructurePage.pillars.apiFirst.description')}
          />
          
          <FeatureCard
            icon={<ScaleIcon />}
            title={t('infrastructurePage.pillars.kubernetes.title')}
            description={t('infrastructurePage.pillars.kubernetes.description')}
          />
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('infrastructurePage.costSavings.title')}</h2>
        
        <div className={styles.costSavingContainer}>
          <div className={styles.costSavingItem}>
            <div className={styles.costSavingContent}>
              <h3 className={styles.costSavingTitle}>{t('infrastructurePage.costSavings.media.title')}</h3>
              <p className={styles.costSavingText}>
                {t('infrastructurePage.costSavings.media.description')}
              </p>
            </div>
            <div className={styles.savingsCard}>
              <h4 className={styles.savingsLabel}>{t('infrastructurePage.costSavings.media.savings')}</h4>
              <div className={styles.savingsAmount}>€160 - €270+</div>
              <small className={styles.savingsNote}>{t('infrastructurePage.costSavings.media.note')}</small>
            </div>
          </div>

          <div className={styles.costSavingItem}>
            <div className={styles.costSavingContent}>
              <h3 className={styles.costSavingTitle}>{t('infrastructurePage.costSavings.geocoding.title')}</h3>
              <p className={styles.costSavingText}>
                {t('infrastructurePage.costSavings.geocoding.description')}
              </p>
            </div>
            <div className={styles.savingsCard}>
              <h4 className={styles.savingsLabel}>{t('infrastructurePage.costSavings.geocoding.savings')}</h4>
              <div className={styles.savingsAmount}>€50 - €75+</div>
              <small className={styles.savingsNote}>{t('infrastructurePage.costSavings.geocoding.note')}</small>
            </div>
          </div>
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('infrastructurePage.performance.title')}</h2>
        
        <div className={styles.metricsGrid}>
          <MetricCard
            value="10"
            suffix="x"
            label={t('infrastructurePage.performance.metrics.speed')}
          />
          <MetricCard
            value="99.9"
            suffix="%"
            label={t('infrastructurePage.performance.metrics.uptime')}
          />
          <MetricCard
            value="<100"
            suffix="ms"
            label={t('infrastructurePage.performance.metrics.response')}
          />
        </div>

        <p className={styles.performanceText}>
          {t('infrastructurePage.performance.description')}
        </p>
      </section>

      <section className={styles.relatedSection}>
        <h3 className={styles.relatedTitle}>{t('infrastructurePage.powers')}</h3>
        <div className={styles.relatedLinks}>
          <Link href="/demo/technical" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>🔧</span>
            <div>
              <h4>{t('breadcrumb.technicalOverview')}</h4>
              <p>{t('infrastructurePage.relatedTechnical')}</p>
            </div>
          </Link>
          <Link href="/demo/verticals" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>📊</span>
            <div>
              <h4>{t('modules.verticals.title')}</h4>
              <p>{t('infrastructurePage.relatedVerticals')}</p>
            </div>
          </Link>
        </div>
      </section>
    </DemoLayout>
  );
};

export default InfrastructureModule;