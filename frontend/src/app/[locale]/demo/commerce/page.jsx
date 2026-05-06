'use client';

import React from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import DemoLayout from '../components/DemoLayout';
import FeatureCard from '../components/FeatureCard';
import MetricCard from '../components/MetricCard';
import styles from './CommercePage.module.css';

const CartIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="9" cy="21" r="1"/>
    <circle cx="20" cy="21" r="1"/>
    <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
  </svg>
);

const PaymentIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
    <line x1="1" y1="10" x2="23" y2="10"/>
  </svg>
);

const OrderIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M12 2L2 7l10 5 10-5-10-5z"/>
    <path d="M2 17l10 5 10-5"/>
    <path d="M2 12l10 5 10-5"/>
  </svg>
);

const CommercePage = () => {
  const t = useTranslations('Demo');
  const breadcrumbs = [
    { label: t('breadcrumb.commerceEngine') }
  ];

  return (
    <DemoLayout
      title={t('commercePage.title')}
      subtitle={t('commercePage.subtitle')}
      breadcrumbs={breadcrumbs}
      showAd={true}
    >
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('commercePage.coreCapabilities')}</h2>
        
        <div className={styles.featureGrid}>
          <FeatureCard
            icon={<CartIcon />}
            title={t('commercePage.features.multiVendor.title')}
            description={t('commercePage.features.multiVendor.description')}
            highlight={true}
          />
          
          <FeatureCard
            icon={<PaymentIcon />}
            title={t('commercePage.features.payment.title')}
            description={t('commercePage.features.payment.description')}
          />
          
          <FeatureCard
            icon={<OrderIcon />}
            title={t('commercePage.features.order.title')}
            description={t('commercePage.features.order.description')}
          />
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('commercePage.metrics.title')}</h2>
        
        <div className={styles.metricsGrid}>
          <MetricCard
            value="<50"
            suffix="ms"
            label={t('commercePage.metrics.checkoutTime')}
          />
          <MetricCard
            value="99.99"
            suffix="%"
            label={t('commercePage.metrics.successRate')}
          />
          <MetricCard
            value="10K+"
            label={t('commercePage.metrics.transactions')}
          />
          <MetricCard
            value="0"
            label={t('commercePage.metrics.pciCompliance')}
          />
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('commercePage.models.title')}</h2>
        
        <div className={styles.modelsGrid}>
          <div className={styles.modelCard}>
            <h3 className={styles.modelTitle}>{t('commercePage.models.leasing.title')}</h3>
            <p className={styles.modelDescription}>
              {t('commercePage.models.leasing.description')}
            </p>
          </div>
          
          <div className={styles.modelCard}>
            <h3 className={styles.modelTitle}>{t('commercePage.models.reservations.title')}</h3>
            <p className={styles.modelDescription}>
              {t('commercePage.models.reservations.description')}
            </p>
          </div>
          
          <div className={styles.modelCard}>
            <h3 className={styles.modelTitle}>{t('commercePage.models.subscription.title')}</h3>
            <p className={styles.modelDescription}>
              {t('commercePage.models.subscription.description')}
            </p>
          </div>
          
          <div className={styles.modelCard}>
            <h3 className={styles.modelTitle}>{t('commercePage.models.auction.title')}</h3>
            <p className={styles.modelDescription}>
              {t('commercePage.models.auction.description')}
            </p>
          </div>
        </div>
      </section>

      <section className={styles.relatedSection}>
        <h3 className={styles.relatedTitle}>{t('commercePage.relatedModules')}</h3>
        <div className={styles.relatedLinks}>
          <Link href="/demo/social" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>💬</span>
            <div>
              <h4>{t('modules.social.title')}</h4>
              <p>{t('commercePage.relatedSocial')}</p>
            </div>
          </Link>
          <Link href="/demo/intelligence" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>🤖</span>
            <div>
              <h4>{t('modules.intelligence.title')}</h4>
              <p>{t('commercePage.relatedIntelligence')}</p>
            </div>
          </Link>
        </div>
      </section>
    </DemoLayout>
  );
};

export default CommercePage;