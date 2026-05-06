'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import DemoLayout from './components/DemoLayout';
import FeatureCard from './components/FeatureCard';
import styles from './TechnicalModule.module.css';

const GoIcon = () => (
  <svg className={styles.techIcon} viewBox="0 0 35 35">
    <path fill="#00ADD8" d="M17.5 35a17.5 17.5 0 1 1 0-35 17.5 17.5 0 0 1 0 35zm-5.7-9.5s-2.1.2-2.1 2c0 2.3 2.1 2 2.1 2s2.2.3 2.2-2c0-2.2-2.2-2-2.2-2zm13.1 0s-2.1.2-2.1 2c0 2.3 2.1 2 2.1 2s2.2.3 2.2-2c0-2.2-2.1-2-2.1-2z"/>
    <path fill="#fff" d="M15.8 17.3h3.3v-3h-3.3v-2h5.3v7h-5.3z"/>
    <path fill="#00ADD8" d="M12.3 20.8a1.6 1.6 0 0 0 1.3-1.3l.3-1.2h-3.3v-2h-.2l-2 7h2.2l.3-1.2c.2-.7 1-2.3 1.4-2.3zm8.3-4.5c-1.8 0-3.3 1.5-3.3 3.3s1.5 3.3 3.3 3.3 3.3-1.5 3.3-3.3-1.5-3.3-3.3-3.3zm0 4.5c-.7 0-1.2-.5-1.2-1.2s.5-1.2 1.2-1.2 1.2.5 1.2 1.2-.5 1.2-1.2 1.2z"/>
  </svg>
);

const K8sIcon = () => (
  <svg className={styles.techIcon} viewBox="0 0 35 35">
    <path fill="#326CE5" d="M17.5 35a17.5 17.5 0 1 1 0-35 17.5 17.5 0 0 1 0 35z"/>
    <path fill="#fff" d="m17.5 17.5 6.4-3.7v-3.1l-6.4 3.7-6.4-3.7v3.1zM11.1 21.2l6.4 3.7 6.4-3.7v-3.1l-6.4 3.7-6.4-3.7z"/>
    <path fill="#326CE5" d="m16.8 17.5 7.1-4.1v-3.1l-7.1 4.1-7.2-4.1v3.1zM10.4 21.2l7.2 4.1 7.1-4.1v-3.1l-7.1 4.1-7.2-4.1z"/>
    <path fill="#fff" d="M17.5 2.5A15 15 0 1 0 32.5 17.5 15 15 0 0 0 17.5 2.5zm0 27.5a12.5 12.5 0 1 1 12.5-12.5A12.5 12.5 0 0 1 17.5 30z"/>
  </svg>
);

const GrpcIcon = () => (
  <svg className={styles.techIcon} viewBox="0 0 35 35">
    <path fill="#4285F4" d="M17.5 35a17.5 17.5 0 1 1 0-35 17.5 17.5 0 0 1 0 35z"/>
    <path fill="#fff" d="M17.5 30a12.5 12.5 0 1 1 0-25 12.5 12.5 0 0 1 0 25z"/>
    <path fill="#4285F4" d="M17.5 27.5a10 10 0 1 1 0-20 10 10 0 0 1 0 20z"/>
    <g fill="#fff">
      <path d="M17.5 7.5c-4.1 0-7.5 2.5-9.1 6l1.8.8c1.2-2.7 3.8-4.6 7.3-4.6s6.1 1.8 7.3 4.6l1.8-.8c-1.6-3.5-5-6-9.1-6zm0 20c4.1 0 7.5-2.5 9.1-6l-1.8-.8c-1.2 2.7-3.8 4.6-7.3 4.6s-6.1-1.8-7.3-4.6l-1.8.8c1.6 3.5 5 6 9.1 6z"/>
      <path d="M7.5 17.5a10 10 0 0 0 1.2 4.8l1.7-.9a8 8 0 0 1-1-3.9 8 8 0 0 1 1-3.9l-1.7-.9a10 10 0 0 0-1.2 4.8zm20 0a10 10 0 0 1-1.2 4.8l-1.7-.9a8 8 0 0 0 1-3.9 8 8 0 0 0-1-3.9l1.7-.9a10 10 0 0 1 1.2 4.8z"/>
    </g>
  </svg>
);

const TechnicalModule = () => {
  const t = useTranslations('Demo');
  const breadcrumbs = [
    { label: t('breadcrumb.technicalOverview') }
  ];

  return (
    <DemoLayout
      title={t('technicalPage.title')}
      subtitle={t('technicalPage.subtitle')}
      breadcrumbs={breadcrumbs}
      showAd={true}
    >
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('technicalPage.architecture.title')}</h2>
        <p className={styles.sectionSubtitle}>
          {t('technicalPage.architecture.subtitle')}
        </p>
        
        <div className={styles.featureGrid}>
          <FeatureCard
            icon="⚡"
            title={t('technicalPage.architecture.performance.title')}
            description={t('technicalPage.architecture.performance.description')}
            highlight={true}
          />
          
          <FeatureCard
            icon="🧠"
            title={t('technicalPage.architecture.aiNative.title')}
            description={t('technicalPage.architecture.aiNative.description')}
          />
          
          <FeatureCard
            icon="☁️"
            title={t('technicalPage.architecture.cloudNative.title')}
            description={t('technicalPage.architecture.cloudNative.description')}
          />
        </div>

        <div className={styles.techStack}>
          <div className={styles.techStackItem}>
            <GoIcon />
            <span>{t('technicalPage.techStack.go')}</span>
          </div>
          <div className={styles.techStackItem}>
            <K8sIcon />
            <span>{t('technicalPage.techStack.kubernetes')}</span>
          </div>
          <div className={styles.techStackItem}>
            <GrpcIcon />
            <span>{t('technicalPage.techStack.grpc')}</span>
          </div>
          <div className={styles.techStackItem}>
            <span className={styles.techStackIcon}>⚛️</span>
            <span>{t('technicalPage.techStack.react')}</span>
          </div>
          <div className={styles.techStackItem}>
            <span className={styles.techStackIcon}>💳</span>
            <span>{t('technicalPage.techStack.stripe')}</span>
          </div>
        </div>
      </section>

      <section className={`${styles.section} ${styles.modulesSection}`}>
        <h2 className={styles.sectionTitle}>{t('technicalPage.platform.title')}</h2>
        
        <div className={styles.modulesGrid}>
          <div className={styles.moduleCard}>
            <h3 className={styles.moduleTitle}>{t('technicalPage.platform.modules.commerce.title')}</h3>
            <p className={styles.moduleDescription}>
              {t('technicalPage.platform.modules.commerce.description')}
            </p>
            <a href="/features/commerce" className={styles.moduleLink}>{t('technicalPage.platform.modules.commerce.link')}</a>
          </div>

          <div className={styles.moduleCard}>
            <h3 className={styles.moduleTitle}>{t('technicalPage.platform.modules.social.title')}</h3>
            <p className={styles.moduleDescription}>
              {t('technicalPage.platform.modules.social.description')}
            </p>
            <a href="/features/social" className={styles.moduleLink}>{t('technicalPage.platform.modules.social.link')}</a>
          </div>

          <div className={styles.moduleCard}>
            <h3 className={styles.moduleTitle}>{t('technicalPage.platform.modules.verticals.title')}</h3>
            <p className={styles.moduleDescription}>
              {t('technicalPage.platform.modules.verticals.description')}
            </p>
            <a href="/features/verticals" className={styles.moduleLink}>{t('technicalPage.platform.modules.verticals.link')}</a>
          </div>

          <div className={styles.moduleCard}>
            <h3 className={styles.moduleTitle}>{t('technicalPage.platform.modules.intelligence.title')}</h3>
            <p className={styles.moduleDescription}>
              {t('technicalPage.platform.modules.intelligence.description')}
            </p>
            <a href="/features/intelligence" className={styles.moduleLink}>{t('technicalPage.platform.modules.intelligence.link')}</a>
          </div>

          <div className={styles.moduleCard}>
            <h3 className={styles.moduleTitle}>{t('technicalPage.platform.modules.infrastructure.title')}</h3>
            <p className={styles.moduleDescription}>
              {t('technicalPage.platform.modules.infrastructure.description')}
            </p>
            <a href="/features/infrastructure" className={styles.moduleLink}>{t('technicalPage.platform.modules.infrastructure.link')}</a>
          </div>
        </div>
      </section>

      <section className={styles.ctaSection}>
        <h2 className={styles.ctaTitle}>{t('technicalPage.ready.title')}</h2>
        <p className={styles.ctaSubtitle}>
          {t('technicalPage.ready.subtitle')}
        </p>
        <a href="#contact-sales" className={styles.ctaButton}>
          {t('technicalPage.cta.engineer')}
        </a>
      </section>
    </DemoLayout>
  );
};

export default TechnicalModule;