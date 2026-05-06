'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import DemoLayout from './components/DemoLayout';
import styles from './VerticalsModule.module.css';

const CarIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M14 16.5V15a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v1.5M10 13V8m4 5V8m-8 5V8m-2 5h12M4 21h16M3 10h18M2 5h20v5H2z"/>
  </svg>
);

const PropertyIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
    <polyline points="9 22 9 12 15 12 15 22"/>
  </svg>
);

const JobsIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <rect x="2" y="7" width="20" height="14" rx="2" ry="2"/>
    <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/>
  </svg>
);

const DealsIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/>
    <line x1="7" y1="7" x2="7.01" y2="7"/>
  </svg>
);

const ServicesIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
  </svg>
);

const StreamIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <polygon points="23 7 16 12 23 17 23 7"/>
    <rect x="1" y="5" width="15" height="14" rx="2" ry="2"/>
  </svg>
);

const PlusIcon = () => (
  <svg className={styles.icon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <line x1="12" y1="5" x2="12" y2="19"/>
    <line x1="5" y1="12" x2="19" y2="12"/>
  </svg>
);

const VerticalsModule = () => {
  const t = useTranslations('Demo');
  const ts = useTranslations('shared');
  const breadcrumbs = [
    { label: t('breadcrumb.verticalsModule') }
  ];

  const verticals = [
    {
      icon: <CarIcon />,
      title: t('verticalsPage.verticals.auto.title'),
      description: t('verticalsPage.verticals.auto.description'),
      fields: t('verticalsPage.verticals.auto.fields'),
      available: true
    },
    {
      icon: <PropertyIcon />,
      title: ts('properties'),
      description: t('verticalsPage.verticals.property.description'),
      fields: t('verticalsPage.verticals.property.fields'),
      available: true
    },
    {
      icon: <JobsIcon />,
      title: ts('jobs'),
      description: t('verticalsPage.verticals.jobs.description'),
      fields: t('verticalsPage.verticals.jobs.fields'),
      available: true
    },
    {
      icon: <DealsIcon />,
      title: ts('deals'),
      description: t('verticalsPage.verticals.deals.description'),
      fields: t('verticalsPage.verticals.deals.fields'),
      available: true
    },
    {
      icon: <ServicesIcon />,
      title: ts('services'),
      description: t('verticalsPage.verticals.services.description'),
      fields: t('verticalsPage.verticals.services.fields'),
      available: true
    },
    {
      icon: <StreamIcon />,
      title: t('verticalsPage.verticals.streams.title'),
      description: t('verticalsPage.verticals.streams.description'),
      fields: t('verticalsPage.verticals.streams.fields'),
      available: false,
      badge: t('verticalsPage.verticals.streams.badge')
    }
  ];

  const advantages = [
    {
      title: t('verticalsPage.advantages.agility.title'),
      description: t('verticalsPage.advantages.agility.description')
    },
    {
      title: t('verticalsPage.advantages.scaling.title'),
      description: t('verticalsPage.advantages.scaling.description')
    },
    {
      title: t('verticalsPage.advantages.contextAware.title'),
      description: t('verticalsPage.advantages.contextAware.description')
    }
  ];

  return (
    <DemoLayout
      title={t('verticalsPage.title')}
      subtitle={t('verticalsPage.subtitle')}
      breadcrumbs={breadcrumbs}
      showAd={true}
    >
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>{t('verticalsPage.available')}</h2>
        
        <div className={styles.verticalsGrid}>
          {verticals.map((vertical, index) => (
            <div key={index} className={styles.verticalCard}>
              <div className={styles.cardHeader}>
                <div className={styles.iconWrapper}>{vertical.icon}</div>
                <div>
                  <h3 className={styles.cardTitle}>
                    {vertical.title}
                    {vertical.badge && <span className={styles.badge}>{vertical.badge}</span>}
                  </h3>
                </div>
              </div>
              <p className={styles.cardDescription}>{vertical.description}</p>
              <div className={styles.cardFields}>
                <strong>{t('verticalsPage.specializedFields')}</strong> {vertical.fields}
              </div>
            </div>
          ))}
          
          <div className={`${styles.verticalCard} ${styles.customCard}`}>
            <div className={styles.cardHeader}>
              <div className={styles.iconWrapper}><PlusIcon /></div>
              <h3 className={styles.cardTitle}>{t('verticalsPage.verticals.custom.title')}</h3>
            </div>
            <p className={styles.cardDescription}>
              {t('verticalsPage.verticals.custom.description')}
            </p>
            <div className={styles.cardFields}>
              <strong>{t('verticalsPage.specializedFields')}</strong> {t('verticalsPage.verticals.custom.fields')}
            </div>
          </div>
        </div>
      </section>

      <section className={styles.advantageSection}>
        <h2 className={styles.sectionTitle}>{t('verticalsPage.advantages.title')}</h2>
        <p className={styles.advantageSubtitle}>
          {t('verticalsPage.advantages.subtitle')}
        </p>
        
        <div className={styles.advantagesGrid}>
          {advantages.map((advantage, index) => (
            <div key={index} className={styles.advantageCard}>
              <h3 className={styles.advantageTitle}>{advantage.title}</h3>
              <p className={styles.advantageText}>{advantage.description}</p>
            </div>
          ))}
        </div>
      </section>
    </DemoLayout>
  );
};

export default VerticalsModule;