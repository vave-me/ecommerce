"use client";

import React from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { ArrowLeft, Settings } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import ErrorBoundary from '@/components/Utils/ErrorBoundary';
import styles from '../MerchantCenter.module.css';

const MerchantSettings = () => {
  let t;
  try {
    t = useTranslations('MerchantSettings');
  } catch (e) {
    t = (key, options) => options?.defaultValue || key;
  }

  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  if (!isAdmin) {
    return (
      <div className={styles.container}>
        <div className={styles.errorState}>
          <h2>{t('accessDenied', { defaultValue: 'Access Denied' })}</h2>
          <p>{t('adminRequired', { defaultValue: 'You need admin privileges to access merchant settings.' })}</p>
        </div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className={styles.container}>
        <div className={styles.header}>
          <div className={styles.headerLeft}>
            <button 
              className={styles.backButton}
              onClick={() => router.push('/admin/merchant')}
            >
              <ArrowLeft size={16} />
              {t('backToMerchant', { defaultValue: 'Back to Merchant Center' })}
            </button>
            <div>
              <h1 className={styles.title}>
                <Settings size={24} />
                {t('title', { defaultValue: 'Merchant Settings' })}
              </h1>
              <p className={styles.subtitle}>
                {t('subtitle', { defaultValue: 'Configure Google Merchant Center integration' })}
              </p>
            </div>
          </div>
        </div>
        
        <div className={styles.contentCard}>
          <div className={styles.emptyState}>
            <Settings size={48} />
            <h3>{t('comingSoon', { defaultValue: 'Coming Soon' })}</h3>
            <p>{t('pageInDevelopment', { defaultValue: 'This page is currently under development.' })}</p>
          </div>
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default MerchantSettings; 