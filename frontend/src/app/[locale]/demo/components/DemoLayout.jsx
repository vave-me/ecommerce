'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import styles from './DemoLayout.module.css';
import CompactVaveAd from '@/components/ads/CompactVaveAd';

const DemoLayout = ({ children, title, subtitle, showAd = true, breadcrumbs = [] }) => {
  const t = useTranslations('Demo');

  return (
    <div className={styles.demoLayout}>
      <nav className={styles.breadcrumbs} aria-label="Breadcrumb">
        <Link href="/demo" className={styles.breadcrumbLink}>
          {t('breadcrumb.demo')}
        </Link>
        {breadcrumbs.map((crumb, index) => (
          <React.Fragment key={index}>
            <span className={styles.breadcrumbSeparator}>/</span>
            {crumb.href ? (
              <Link href={crumb.href} className={styles.breadcrumbLink}>
                {crumb.label}
              </Link>
            ) : (
              <span className={styles.breadcrumbCurrent}>{crumb.label}</span>
            )}
          </React.Fragment>
        ))}
      </nav>

      <header className={styles.header}>
        <h1 className={styles.title}>{title}</h1>
        {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
      </header>

      <main className={styles.content}>
        {children}
      </main>

      {showAd && (
        <aside className={styles.adSection}>
          <CompactVaveAd />
        </aside>
      )}
    </div>
  );
};

export default DemoLayout;