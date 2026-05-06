"use client";
import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { enableUser } from '../../../api/userService';
import styles from './verify.module.css';
export default function VerifyAccount() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const t = useTranslations('Verify');
  const [status, setStatus] = useState('verifying'); // verifying, success, error
  const [errorMessage, setErrorMessage] = useState('');
  useEffect(() => {
    const verifyUser = async () => {
      try {
        // Get userId and verificationToken from URL
        const userId = searchParams.get('id');
        const verificationToken = searchParams.get('token');
        if (!userId || !verificationToken) {
          throw new Error(t('missingInfo'));
        }
        // Make API call to enable user
        await enableUser(userId, verificationToken);
        setStatus('success');
      } catch (error) {
        setStatus('error');
        setErrorMessage(error.response?.data?.message || error.message || t('unexpectedError'));
      }
    };
    verifyUser();
  }, [searchParams, t]);
  const handleGoToLogin = () => {
    router.push('/login');
  };
  return (
    <div className={styles.container}>
      <div className={styles.card}>
        {status === 'verifying' && (
          <>
            <div className={styles.spinner}></div>
            <h1>{t('verifying')}</h1>
            <p>{t('verifyingMessage')}</p>
          </>
        )}
        {status === 'success' && (
          <>
            <div className={styles.successIcon}>✓</div>
            <h1>{t('success')}</h1>
            <p>{t('successMessage')}</p>
            <button className={styles.button} onClick={handleGoToLogin}>
              {t('goToLogin')}
            </button>
          </>
        )}
        {status === 'error' && (
          <>
            <div className={styles.errorIcon}>✗</div>
            <h1>{t('error')}</h1>
            <p>{errorMessage}</p>
            <button className={styles.button} onClick={handleGoToLogin}>
              {t('goToLogin')}
            </button>
          </>
        )}
      </div>
    </div>
  );
} 