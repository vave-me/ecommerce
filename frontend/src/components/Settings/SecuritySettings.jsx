// src/components/Settings/SecuritySettings.jsx
"use client"
import React, {useState, memo} from 'react';
import {useTranslations} from 'next-intl';
import {useAuth} from '../../context/AuthContext';
import {Shield, AlertTriangle, CheckCircle, X, Loader, Clock, LogOut} from '@/icons';
import {forgotPassword} from '../../api/client/userApi';
import styles from './SecuritySettings.module.css';
/**
 * SecuritySettings Component
 * Allows users to configure security settings.
 */
const SecuritySettings = memo(() => {
    const t = useTranslations('SecuritySettings');
    const {user, signOutUser} = useAuth();
    const [email, setEmail] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [sessionActivity] = useState([
        {
            id: '1',
            device: 'Chrome on Windows',
            location: 'Warsaw, Poland',
            time: '2023-06-15T10:30:00',
            current: true
        },
        {
            id: '2',
            device: 'Safari on iPhone',
            location: 'Kraków, Poland',
            time: '2023-06-14T18:45:00',
            current: false
        }
    ]);
    const [feedback, setFeedback] = useState({
        key: '',
        type: '',
        section: ''
    });
    const showFeedback = (key, type, section, duration = 5000) => {
        setFeedback({key, type, section});
        setTimeout(() => {
            setFeedback({key: '', type: '', section: ''});
        }, duration);
    };
    const handleSendResetLink = async (e) => {
        e.preventDefault();
        if (!email) return;
        setSubmitting(true);
        try {
            await forgotPassword(email);
            showFeedback('resetLinkSent', 'success', 'passwordReset');
            setEmail('');
        } catch (error) {
            showFeedback('resetLinkError', 'error', 'passwordReset');
        } finally {
            setSubmitting(false);
        }
    };
    const handleLogoutOtherSessions = async () => {
        setSubmitting(true);
        try {
            // API call would go here
            // In a real implementation, we would call an API endpoint to revoke other sessions
            setTimeout(() => {
                showFeedback('logoutSessionsSuccess', 'success', 'sessions');
                setSubmitting(false);
            }, 1000);
        } catch (error) {
            showFeedback('logoutSessionsError', 'error', 'sessions');
            setSubmitting(false);
        }
    };
    const handleLogoutAllSessions = async () => {
        if (window.confirm(t('confirmLogoutAll'))) {
            try {
                await signOutUser();
            } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
        }
    };
    const renderFeedback = (section) => {
        if (feedback.section === section && feedback.key) {
            return (
                <div className={`${styles.feedback} ${styles[feedback.type]}`} role="alert">
                    {feedback.type === 'success' ? (
                        <CheckCircle className={styles.feedbackIcon} aria-hidden="true"/>
                    ) : (
                        <X className={styles.feedbackIcon} aria-hidden="true"/>
                    )}
                    <span>{t(feedback.key)}</span>
                </div>
            );
        }
        return null;
    };
    const formatDate = (dateString) => {
        try {
            const date = new Date(dateString);
            return new Intl.DateTimeFormat('en-US', {
                dateStyle: 'medium',
                timeStyle: 'short'
            }).format(date);
        } catch (error) {
            return dateString;
        }
    };
    return (
        <div className={styles.container}>
            <h2 className={styles.pageTitle}>{t('pageTitle')}</h2>
            <p className={styles.pageDescription}>{t('pageDescription')}</p>
            {/* Password Reset Section */}
            <section className={styles.section} aria-labelledby="password-reset-title">
                <h3 id="password-reset-title" className={styles.sectionTitle}>
                    <Shield className={styles.sectionIcon} aria-hidden="true"/>
                    {t('passwordResetTitle')}
                </h3>
                <p className={styles.sectionDescription}>{t('passwordResetDescription')}</p>
                <form className={styles.form} onSubmit={handleSendResetLink}>
                    <div className={styles.formGroup}>
                        <label htmlFor="email" className={styles.label}>
                            {t('emailLabel')}
                        </label>
                        <input
                            type="email"
                            id="email"
                            className={styles.input}
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder={t('emailPlaceholder')}
                            required
                        />
                    </div>
                    {renderFeedback('passwordReset')}
                    <button 
                        type="submit" 
                        className={styles.primaryButton}
                        disabled={submitting || !email}
                    >
                        {submitting ? (
                            <>
                                <Loader className={styles.buttonIcon} />
                                {t('sendingLink')}
                            </>
                        ) : t('sendResetLink')}
                    </button>
                </form>
            </section>
            {/* Active Sessions Section */}
            <section className={styles.section} aria-labelledby="sessions-title">
                <h3 id="sessions-title" className={styles.sectionTitle}>
                    <Clock className={styles.sectionIcon} aria-hidden="true"/>
                    {t('activeSessions')}
                </h3>
                <p className={styles.sectionDescription}>{t('activeSessionsDescription')}</p>
                {renderFeedback('sessions')}
                <div className={styles.sessionsList}>
                    {sessionActivity.map((session) => (
                        <div 
                            key={session.id} 
                            className={`${styles.sessionItem} ${session.current ? styles.currentSession : ''}`}
                        >
                            <div className={styles.sessionInfo}>
                                <div className={styles.sessionDevice}>
                                    <strong>{session.device}</strong>
                                    {session.current && (
                                        <span className={styles.currentBadge}>{t('currentSession')}</span>
                                    )}
                                </div>
                                <div className={styles.sessionMeta}>
                                    <span>{session.location}</span>
                                    <span className={styles.sessionTime}>
                                        {formatDate(session.time)}
                                    </span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
                <div className={styles.actionButtons}>
                    <button 
                        className={styles.secondaryButton} 
                        onClick={handleLogoutOtherSessions}
                        disabled={submitting}
                    >
                        {submitting ? (
                            <>
                                <Loader className={styles.buttonIcon} />
                                {t('loggingOut')}
                            </>
                        ) : t('logoutOtherSessions')}
                    </button>
                    <button 
                        className={styles.dangerButton} 
                        onClick={handleLogoutAllSessions}
                    >
                        <LogOut className={styles.buttonIcon} />
                        {t('logoutAllSessions')}
                    </button>
                </div>
            </section>
            {/* Account Security Section */}
            <section className={styles.section} aria-labelledby="account-security-title">
                <h3 id="account-security-title" className={styles.sectionTitle}>
                    <AlertTriangle className={styles.sectionIcon} aria-hidden="true"/>
                    {t('dangerZoneTitle')}
                </h3>
                <p className={styles.sectionDescription}>{t('dangerZoneDescription')}</p>
                <div className={styles.dangerZone}>
                    <div className={styles.dangerAction}>
                        <div>
                            <h4 className={styles.dangerTitle}>{t('disableAccount')}</h4>
                            <p className={styles.dangerDescription}>{t('disableAccountDescription')}</p>
                        </div>
                        <button className={styles.outlineDangerButton}>
                            {t('disableAccountButton')}
                        </button>
                    </div>
                    <div className={styles.dangerAction}>
                        <div>
                            <h4 className={styles.dangerTitle}>{t('deleteAccount')}</h4>
                            <p className={styles.dangerDescription}>{t('deleteAccountDescription')}</p>
                        </div>
                        <button className={styles.dangerButton}>
                            {t('deleteAccountButton')}
                        </button>
                    </div>
                </div>
            </section>
        </div>
    );
});
SecuritySettings.displayName = 'SecuritySettings';
export default SecuritySettings;