'use client';
import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import styles from './ResetPassword.module.css';
import {resetPassword} from "../../../api/userApi";
export default function ResetPasswordPage() {
  const t = useTranslations('ResetPassword');
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get('token');
  const email = searchParams.get('email');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  useEffect(() => {
    // Check if token and email are present
    if (!token || !email) {
      setError(t('invalidLink'));
    }
  }, [token, email, t]);
  const checkPasswordStrength = (password) => {
    // Calculate password strength (simple implementation)
    let strength = 0;
    if (password.length >= 8) strength += 1;
    if (password.match(/[A-Z]/)) strength += 1;
    if (password.match(/[0-9]/)) strength += 1;
    if (password.match(/[^A-Za-z0-9]/)) strength += 1;
    setPasswordStrength(strength);
  };
  const handlePasswordChange = (e) => {
    const newPassword = e.target.value;
    setPassword(newPassword);
    checkPasswordStrength(newPassword);
  };
  const handleSubmit = async (e) => {
    e.preventDefault();
    // Validation
    if (password !== confirmPassword) {
      setError(t('passwordsDoNotMatch'));
      return;
    }
    if (password.length < 8) {
      setError(t('passwordTooShort'));
      return;
    }
    setError('');
    setIsLoading(true);
    try {
      await resetPassword(token, email, password);
      setSuccess(true);
    } catch (err) {
      setError(err.response?.data?.message || t('genericError'));
    } finally {
      setIsLoading(false);
    }
  };
  return (
    <div className={styles.container}>
      <div className={styles.contentWrapper}>
        <div className={styles.formContainer}>
          <h1 className={styles.title}>
            {t('title')}
          </h1>
          <p className={styles.description}>
            {t('description')}
          </p>
          {error && (
            <div className={styles.errorMessage}>
              <p>{error}</p>
            </div>
          )}
          {!token || !email ? (
            <div className={styles.warningMessage}>
              <p>{t('invalidLink')}</p>
              <Link 
                href="/forgot-password" 
                className={styles.actionLink}
              >
                {t('requestNewLink')}
              </Link>
            </div>
          ) : success ? (
            <div className={styles.successMessage}>
              <p>{t('passwordResetSuccess')}</p>
              <p className={styles.successSubtext}>{t('loginWithNewPassword')}</p>
              <Link 
                href="/login"
                className={styles.loginButton}
              >
                {t('goToLogin')}
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit}>
              {/* New Password */}
              <div className={styles.formGroup}>
                <label htmlFor="password" className={styles.inputLabel}>
                  {t('newPassword')}
                </label>
                <div className={styles.inputWrapper}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                  </svg>
                  <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={handlePasswordChange}
                    required
                    placeholder={t('newPasswordPlaceholder')}
                    className={styles.input}
                  />
                </div>
                {/* Password strength meter */}
                {password && (
                  <div className={styles.strengthMeter}>
                    <div className={`${styles.strengthIndicator} ${passwordStrength >= 1 ? styles.strengthIndicatorWeak : ''}`}></div>
                    <div className={`${styles.strengthIndicator} ${passwordStrength >= 2 ? styles.strengthIndicatorFair : ''}`}></div>
                    <div className={`${styles.strengthIndicator} ${passwordStrength >= 3 ? styles.strengthIndicatorGood : ''}`}></div>
                    <div className={`${styles.strengthIndicator} ${passwordStrength >= 4 ? styles.strengthIndicatorStrong : ''}`}></div>
                    <p className={styles.strengthText}>
                      {t('passwordStrength')}: {
                        passwordStrength === 0 ? t('veryWeak') :
                        passwordStrength === 1 ? t('weak') :
                        passwordStrength === 2 ? t('medium') :
                        passwordStrength === 3 ? t('strong') :
                        t('veryStrong')
                      }
                    </p>
                  </div>
                )}
              </div>
              {/* Confirm Password */}
              <div className={styles.formGroup}>
                <label htmlFor="confirmPassword" className={styles.inputLabel}>
                  {t('confirmPassword')}
                </label>
                <div className={styles.inputWrapper}>
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                  </svg>
                  <input
                    id="confirmPassword"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    placeholder={t('confirmPasswordPlaceholder')}
                    className={styles.input}
                  />
                </div>
                {password && confirmPassword && password !== confirmPassword && (
                  <p className={styles.errorText}>{t('passwordsDoNotMatch')}</p>
                )}
              </div>
              <button
                type="submit"
                disabled={isLoading}
                className={styles.submitButton}
              >
                {isLoading ? (
                  <>
                    <svg className={styles.spinnerIcon} xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <circle cx="12" cy="12" r="10" stroke="currentColor" className="opacity-25"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    {t('processing')}
                  </>
                ) : t('resetButton')}
              </button>
            </form>
          )}
          {(token && email && !success) && (
            <div className={styles.footer}>
              <Link 
                href="/login" 
                className={styles.backLink}
              >
                &larr; {t('backToLogin')}
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 