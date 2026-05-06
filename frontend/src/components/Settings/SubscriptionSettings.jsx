// src/components/SubscriptionSettings.jsx
"use client"
import React, { useState, memo } from 'react';
import { useTranslations } from 'next-intl';
import { 
    Plus, 
    CreditCard, 
    Mail, 
    Calendar, 
    CheckCircle, 
    ExternalLink, 
    Info, 
    History, 
    Shield,
    Globe,
    Loader,
    Bell
} from '@/icons';
import ToggleSwitch from './ToggleSwitch';
import styles from './SubscriptionSettings.module.css';
/**
 * Modern SubscriptionSettings Component
 * Elegant design matching the account and notification settings pattern
 */
const SubscriptionSettings = memo(({
    availableNewsletters = [
        { 
            id: 'weeklyDigest', 
            title: 'Weekly Digest', 
            description: 'Get a weekly roundup of the most important updates and featured content.',
            category: 'content'
        },
        { 
            id: 'productUpdates', 
            title: 'Product Updates', 
            description: 'Be the first to know about new features, improvements, and releases.',
            category: 'product'
        },
        { 
            id: 'securityAlerts', 
            title: 'Security Alerts', 
            description: 'Important security notifications, tips, and best practices for your account.',
            category: 'security'
        }
    ],
    userSubscriptions = ['weeklyDigest'],
    currentPlan = { 
        name: 'Pro Plan', 
        price: '$15.99', 
        billingCycle: 'monthly', 
        nextBillingDate: '2024-06-15',
        paymentMethod: '•••• 4242'
    },
    onSave = () => {}
}) => {
    const t = useTranslations('SubscriptionSettings');
    // Initialize subscription state based on user's current subscriptions
    const [subscriptions, setSubscriptions] = useState(() => {
        const initialState = {};
        availableNewsletters.forEach((newsletter) => {
            initialState[newsletter.id] = userSubscriptions.includes(newsletter.id);
        });
        return initialState;
    });
    const [showFeedback, setShowFeedback] = useState(false);
    const [saving, setSaving] = useState(false);
    const [loadingAction, setLoadingAction] = useState('');
    const handleToggle = async (newsletterId) => {
        const updatedSubscriptions = {
            ...subscriptions,
            [newsletterId]: !subscriptions[newsletterId]
        };
        setSubscriptions(updatedSubscriptions);
        // Optionally update parent immediately
        try {
            await onSave(updatedSubscriptions);
        } catch (error) {
            // Revert change on error
            setSubscriptions(subscriptions);
        }
    };
    const handleSave = async () => {
        setSaving(true);
        try {
            await onSave(subscriptions);
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
    const handlePlanAction = async (action) => {
        setLoadingAction(action);
        // Simulate action
        setTimeout(() => {
            setLoadingAction('');
        }, 1500);
    };
    const getCategoryIcon = (category) => {
        switch (category) {
            case 'security': return Shield;
            case 'product': return Bell;
            case 'content': return Globe;
            default: return Mail;
        }
    };
    return (
        <div className={styles.container}>
            {/* Current Plan Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <CreditCard size={20} className={styles.sectionIcon} />
                        <h3>{t('currentPlan')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        Manage your subscription plan and billing information
                    </p>
                </div>
                <div className={styles.card}>
                    <div className={styles.planContainer}>
                        <div className={styles.planInfo}>
                            <div className={styles.planHeader}>
                                <h4 className={styles.planName}>{t('proPlan')}</h4>
                                <span className={styles.planBadge}>Active</span>
                            </div>
                            <div className={styles.planPrice}>
                                {t('monthlyPrice', { price: currentPlan.price })}
                            </div>
                        </div>
                        <div className={styles.planDetails}>
                            <div className={styles.planDetail}>
                                <Calendar size={16} className={styles.detailIcon} />
                                <span className={styles.detailLabel}>{t('nextBillingDate')}</span>
                                <span className={styles.detailValue}>{currentPlan.nextBillingDate}</span>
                            </div>
                            <div className={styles.planDetail}>
                                <CreditCard size={16} className={styles.detailIcon} />
                                <span className={styles.detailLabel}>{t('paymentMethod')}</span>
                                <span className={styles.detailValue}>{currentPlan.paymentMethod}</span>
                            </div>
                        </div>
                        <div className={styles.planActions}>
                            <button 
                                className={styles.secondaryButton}
                                onClick={() => handlePlanAction('history')}
                                disabled={loadingAction === 'history'}
                            >
                                {loadingAction === 'history' ? (
                                    <>
                                        <Loader size={16} className={styles.buttonSpinner} />
                                        Loading...
                                    </>
                                ) : (
                                    <>
                                        <History size={16} />
                                        {t('billingHistory')}
                                    </>
                                )}
                            </button>
                            <button 
                                className={styles.primaryButton}
                                onClick={() => handlePlanAction('update')}
                                disabled={loadingAction === 'update'}
                            >
                                {loadingAction === 'update' ? (
                                    <>
                                        <Loader size={16} className={styles.buttonSpinner} />
                                        Loading...
                                    </>
                                ) : (
                                    <>
                                        <CreditCard size={16} />
                                        {t('updatePlan')}
                                    </>
                                )}
                            </button>
                        </div>
                    </div>
                </div>
            </section>
            {/* Newsletter Subscriptions Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <Mail size={20} className={styles.sectionIcon} />
                        <h3>{t('newsletterSubscriptions')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        Subscribe to newsletters and content updates that interest you
                    </p>
                </div>
                <div className={styles.card}>
                    <div className={styles.newslettersList}>
                        {availableNewsletters.map((newsletter) => {
                            const IconComponent = getCategoryIcon(newsletter.category);
                            return (
                                <div key={newsletter.id} className={styles.newsletterItem}>
                                    <div className={styles.newsletterInfo}>
                                        <div className={styles.newsletterHeader}>
                                            <IconComponent size={18} className={styles.newsletterIcon} />
                                            <h4 className={styles.newsletterTitle}>
                                                {t(newsletter.id)}
                                            </h4>
                                            {subscriptions[newsletter.id] && (
                                                <CheckCircle size={16} className={styles.subscribedIcon} />
                                            )}
                                        </div>
                                        <p className={styles.newsletterDescription}>
                                            {t(`${newsletter.id}Desc`)}
                                        </p>
                                    </div>
                                    <div className={styles.newsletterToggle}>
                                        <ToggleSwitch
                                            isOn={subscriptions[newsletter.id]}
                                            handleToggle={() => handleToggle(newsletter.id)}
                                            label={subscriptions[newsletter.id] ? t('subscribed') : t('notSubscribed')}
                                            size="medium"
                                        />
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                    <div className={styles.addNewsletterSection}>
                        <button className={styles.addButton}>
                            <Plus size={16} />
                            {t('addCustomNewsletter')}
                        </button>
                    </div>
                </div>
            </section>
            {/* Additional Information Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <Info size={20} className={styles.sectionIcon} />
                        <h3>{t('additionalInfo')}</h3>
                    </div>
                </div>
                <div className={styles.card}>
                    <div className={styles.infoBox}>
                        <Info size={20} className={styles.infoIcon} />
                        <div className={styles.infoContent}>
                            <p className={styles.infoText}>
                                {t('communicationCenterInfo')}
                            </p>
                            <a href="#" className={styles.infoLink}>
                                Communication Center 
                                <ExternalLink size={14} />
                            </a>
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
                            Saving...
                        </>
                    ) : (
                        <>
                            <CheckCircle size={16} />
                            {t('savePreferences')}
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
                            Subscription preferences saved successfully!
                        </span>
                    </div>
                </div>
            )}
        </div>
    );
});
SubscriptionSettings.displayName = 'SubscriptionSettings';
export default SubscriptionSettings;