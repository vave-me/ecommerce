// src/components/NotificationSettings.jsx
"use client"
import React, { useState, memo } from 'react';
import { useTranslations } from 'next-intl';
import { 
    Bell, 
    Mail, 
    Smartphone, 
    Globe, 
    CheckCircle, 
    Shield,
    MessageSquare,
    AlertTriangle,
    Loader
} from '@/icons';
import ToggleSwitch from './ToggleSwitch';
import styles from './NotificationSettings.module.css';
/**
 * Modern NotificationSettings Component
 * Elegant design matching the account settings pattern
 */
const NotificationSettings = memo(({ 
    settings = {
        email: false,
        sms: false,
        push: false,
        marketing: false,
        newsletter: false,
        security: false,
        updates: false
    },
    onUpdateSettings = () => {}
}) => {
    const t = useTranslations('NotificationSettings');
    // Initialize preferences state using defaulted settings
    const [preferences, setPreferences] = useState(settings);
    const [showFeedback, setShowFeedback] = useState(false);
    const [saving, setSaving] = useState(false);
    const handleToggle = async (type) => {
        const updatedPreferences = { ...preferences, [type]: !preferences[type] };
        setPreferences(updatedPreferences);
        // Update parent component
        try {
            await onUpdateSettings(updatedPreferences);
        } catch (error) {
            // Revert the change if API call fails
            setPreferences(preferences);
        }
    };
    const handleSave = async () => {
        setSaving(true);
        try {
            await onUpdateSettings(preferences);
            setShowFeedback(true);
            setTimeout(() => setShowFeedback(false), 4000);
        } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
            setSaving(false);
        }
    };
    return (
        <div className={styles.container}>
            {/* Essential Notifications Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <Shield size={20} className={styles.sectionIcon} />
                        <h3>{t('essentialTitle')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        {t('essentialDescription')}
                    </p>
                </div>
                <div className={styles.card}>
                    <div className={styles.toggleGroup}>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <Mail size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('securityAlertsTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('securityAlertsDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.security}
                                handleToggle={() => handleToggle('security')}
                                label={preferences.security ? t('enabled') : t('disabled')}
                                size="medium"
                            />
                        </div>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <AlertTriangle size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('systemUpdatesTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('systemUpdatesDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.updates}
                                handleToggle={() => handleToggle('updates')}
                                label={preferences.updates ? t('enabled') : t('disabled')}
                                size="medium"
                            />
                        </div>
                    </div>
                </div>
            </section>
            {/* Communication Preferences Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <MessageSquare size={20} className={styles.sectionIcon} />
                        <h3>{t('communicationTitle')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        {t('communicationDescription')}
                    </p>
                </div>
                <div className={styles.card}>
                    <div className={styles.toggleGroup}>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <Mail size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('emailNotificationsTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('emailNotificationsDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.email}
                                handleToggle={() => handleToggle('email')}
                                label={preferences.email ? t('enabled') : t('disabled')}
                                size="medium"
                            />
                        </div>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <Bell size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('pushNotificationsTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('pushNotificationsDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.push}
                                handleToggle={() => handleToggle('push')}
                                label={preferences.push ? t('enabled') : t('disabled')}
                                size="medium"
                            />
                        </div>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <Smartphone size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('smsNotificationsTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('smsNotificationsDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.sms}
                                handleToggle={() => handleToggle('sms')}
                                label={preferences.sms ? t('enabled') : t('disabled')}
                                size="medium"
                            />
                        </div>
                    </div>
                </div>
            </section>
            {/* Marketing & Content Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <Globe size={20} className={styles.sectionIcon} />
                        <h3>{t('marketingTitle')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        {t('marketingDescription')}
                    </p>
                </div>
                <div className={styles.card}>
                    <div className={styles.toggleGroup}>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <Mail size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('marketingEmailsTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('marketingEmailsDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.marketing}
                                handleToggle={() => handleToggle('marketing')}
                                label={preferences.marketing ? t('subscribed') : t('unsubscribed')}
                                size="medium"
                            />
                        </div>
                        <div className={styles.toggleItem}>
                            <div className={styles.toggleInfo}>
                                <div className={styles.toggleHeader}>
                                    <Globe size={18} className={styles.toggleIcon} />
                                    <h4 className={styles.toggleTitle}>{t('newsletterTitle')}</h4>
                                </div>
                                <p className={styles.toggleDescription}>
                                    {t('newsletterDescription')}
                                </p>
                            </div>
                            <ToggleSwitch
                                isOn={preferences.newsletter}
                                handleToggle={() => handleToggle('newsletter')}
                                label={preferences.newsletter ? t('subscribed') : t('unsubscribed')}
                                size="medium"
                            />
                        </div>
                    </div>
                </div>
            </section>
            {/* Save Actions */}
            <div className={styles.actions}>
                <button 
                    className={styles.saveButton} 
                    onClick={handleSave}
                    disabled={saving}
                >
                    {saving ? (
                        <>
                            <Loader size={16} className={styles.buttonSpinner} />
                            {t('savingPreferences')}
                        </>
                    ) : (
                        <>
                            <CheckCircle size={16} />
                            {t('saveAllPreferences')}
                        </>
                    )}
                </button>
            </div>
            {/* Success Feedback */}
            {showFeedback && (
                <div className={styles.feedback}>
                    <div className={styles.feedbackContent}>
                        <CheckCircle size={16} className={styles.feedbackIcon} />
                        <span className={styles.feedbackText}>
                            {t('preferencesUpdated')}
                        </span>
                    </div>
                </div>
            )}
        </div>
    );
});
NotificationSettings.displayName = 'NotificationSettings';
export default NotificationSettings;