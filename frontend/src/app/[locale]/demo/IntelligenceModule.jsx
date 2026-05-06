'use client';

import React from 'react';
import Link from 'next/link';
import { useTranslations } from 'next-intl';
import DemoLayout from './components/DemoLayout';
import FeatureCard from './components/FeatureCard';
import styles from './IntelligenceModule.module.css';

const AiCoreIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
    <path d="M12 22v-6"/>
    <path d="m20.5 15.5-5-3-5 3"/>
  </svg>
);

const VoiceImageIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"/>
    <path d="M19 10v2a7 7 0 0 1-14 0v-2"/>
    <rect x="2" y="16" width="20" height="6" rx="2"/>
    <circle cx="8" cy="19" r="1"/>
    <path d="m14 19 2 2 2-2"/>
  </svg>
);

const AnalyticsIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M3 3v18h18"/>
    <path d="M18.7 8a4 4 0 0 1-4.9 0l-3.2 4.3a4 4 0 0 1-4.9 0l-2.4-3.2"/>
    <circle cx="6" cy="15" r="2"/>
    <circle cx="12" cy="11" r="2"/>
    <circle cx="18" cy="7" r="2"/>
  </svg>
);

const SupportIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M20 16.5v-2.1a2.2 2.2 0 0 0-1.2-2l-1.8-1a2.2 2.2 0 0 1-1.2-2V8a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v1.4a2.2 2.2 0 0 1-1.2 2l-1.8 1a2.2 2.2 0 0 0-1.2 2V17a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2z"/>
    <path d="M8 19v-2.5"/>
    <path d="M16 19v-2.5"/>
  </svg>
);

const IntelligenceModule = () => {
  const t = useTranslations('Demo');
  const breadcrumbs = [
    { label: t('breadcrumb.intelligenceModule') }
  ];

  return (
    <DemoLayout
      title={t('intelligencePage.title')}
      subtitle={t('intelligencePage.subtitle')}
      breadcrumbs={breadcrumbs}
      showAd={true}
    >
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('intelligencePage.capabilities.title')}</h2>
        
        <div className={styles.featuresGrid}>
          <div className={styles.featureItem}>
            <div className={styles.featureIconWrapper}>
              <AiCoreIcon />
            </div>
            <div className={styles.featureContent}>
              <h3 className={styles.featureTitle}>{t('intelligencePage.capabilities.autonomous.title')}</h3>
              <p className={styles.featureText}>
                {t('intelligencePage.capabilities.autonomous.description')}
              </p>
              <div className={styles.codeBlock}>
                <span className={styles.codeComment}>{t('intelligencePage.capabilities.autonomous.example.comment')}</span>
                <span className={styles.codeText}>
                  {t('intelligencePage.capabilities.autonomous.example.command')}
                </span>
              </div>
            </div>
          </div>

          <div className={styles.featureItem}>
            <div className={styles.featureIconWrapper}>
              <VoiceImageIcon />
            </div>
            <div className={styles.featureContent}>
              <h3 className={styles.featureTitle}>{t('intelligencePage.capabilities.multiModal.title')}</h3>
              <p className={styles.featureText}>
                {t('intelligencePage.capabilities.multiModal.description')}
              </p>
              <ul className={styles.featureList}>
                <li>{t('intelligencePage.capabilities.multiModal.features.image')}</li>
                <li>{t('intelligencePage.capabilities.multiModal.features.voice')}</li>
              </ul>
            </div>
          </div>

          <div className={styles.featureItem}>
            <div className={styles.featureIconWrapper}>
              <SupportIcon />
            </div>
            <div className={styles.featureContent}>
              <h3 className={styles.featureTitle}>{t('intelligencePage.capabilities.conversational.title')}</h3>
              <p className={styles.featureText}>
                {t('intelligencePage.capabilities.conversational.description')}
              </p>
            </div>
          </div>

          <div className={styles.featureItem}>
            <div className={styles.featureIconWrapper}>
              <AnalyticsIcon />
            </div>
            <div className={styles.featureContent}>
              <h3 className={styles.featureTitle}>{t('intelligencePage.capabilities.analytics.title')}</h3>
              <p className={styles.featureText}>
                {t('intelligencePage.capabilities.analytics.description')}
              </p>
              <div className={styles.codeBlock}>
                <span className={styles.codeComment}>{t('intelligencePage.capabilities.analytics.example.comment')}</span>
                <span className={styles.codeText}>
                  {t('intelligencePage.capabilities.analytics.example.query')}
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className={styles.differentiatorSection}>
        <div className={styles.differentiatorContent}>
          <h2 className={styles.differentiatorTitle}>{t('intelligencePage.differentiator.title')}</h2>
          <p className={styles.differentiatorText}>
            {t('intelligencePage.differentiator.description')}
          </p>
        </div>
      </section>

      <section className={styles.relatedSection}>
        <h3 className={styles.relatedTitle}>{t('intelligencePage.integration')}</h3>
        <div className={styles.relatedLinks}>
          <Link href="/demo/commerce" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>🛍️</span>
            <div>
              <h4>{t('modules.commerce.title')}</h4>
              <p>{t('intelligencePage.relatedCommerce')}</p>
            </div>
          </Link>
          <Link href="/demo/social" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>💬</span>
            <div>
              <h4>{t('modules.social.title')}</h4>
              <p>{t('intelligencePage.relatedSocial')}</p>
            </div>
          </Link>
          <Link href="/demo/verticals" className={styles.relatedLink}>
            <span className={styles.relatedIcon}>🎯</span>
            <div>
              <h4>{t('modules.verticals.title')}</h4>
              <p>{t('intelligencePage.relatedVerticals')}</p>
            </div>
          </Link>
        </div>
      </section>
    </DemoLayout>
  );
};

export default IntelligenceModule;