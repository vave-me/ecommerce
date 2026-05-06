// src/components/AccountSettings.jsx
"use client"
import React, { useState, useEffect, memo } from "react";
import { useTranslations } from 'next-intl';
import { useAuth } from '../../context/AuthContext';
import { getUserById, updateUser, renameUser } from '../../api/client/userApi';
import { 
    CheckCircle, 
    X, 
    User, 
    Lock, 
    AlertCircle, 
    Loader,
    Mail,
    MapPin,
    FileText,
    Edit3,
    Shield
} from '@/icons';
import styles from "./AccountSettings.module.css";
/**
 * Modern AccountSettings Component
 * Elegant design matching notifications and wishlists pages
 */
const AccountSettings = memo(() => {
    const t = useTranslations('AccountSettings');
    const { user } = useAuth();
    // State for profile information
    const [profile, setProfile] = useState({
        userName: '',
        firstName: '',
        lastName: '',
        email: '',
        bio: '',
        location: '',
        lat: null,
        lng: null,
        thumbnail: '',
    });
    // Loading and error states
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState('');
    // State for password fields
    const [password, setPassword] = useState({
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
    });
    // State for feedback messages
    const [feedbackMessage, setFeedbackMessage] = useState({
        key: "",
        type: "",
        section: "",
        values: {},
    });
    // Fetch user data on component mount
    useEffect(() => {
        const fetchUserData = async () => {
            if (!user?.userId) return;
            setLoading(true);
            setError('');
            try {
                const response = await getUserById(user.userId);
                if (response && response.user) {
                    const userData = response.user;
                    setProfile({
                        userName: userData.userName || '',
                        firstName: userData.firstName || '',
                        lastName: userData.lastName || '',
                        email: userData.email || '',
                        bio: userData.bio || '',
                        location: userData.location || '',
                        lat: userData.lat || null,
                        lng: userData.lng || null,
                        thumbnail: userData.thumbnail || '',
                    });
                }
            } catch (err) {
                setError('fetchError');
            } finally {
                setLoading(false);
            }
        };
        fetchUserData();
    }, [user]);
    // Handle profile input changes
    const handleProfileChange = (e) => {
        const { name, value } = e.target;
        setProfile((prev) => ({ ...prev, [name]: value }));
    };
    // Handle password input changes
    const handlePasswordChange = (e) => {
        const { name, value } = e.target;
        setPassword((prev) => ({ ...prev, [name]: value }));
    };
    // Generic function to set feedback message
    const showFeedback = (key, type, section, values = {}, duration = 3000) => {
        setFeedbackMessage({ key, type, section, values });
        setTimeout(() => {
            setFeedbackMessage({ key: "", type: "", section: "", values: {} });
        }, duration);
    };
    // Handle profile form submission
    const handleSubmitProfile = async (e) => {
        e.preventDefault();
        if (!user?.userId) return;
        setSubmitting(true);
        try {
            // First update username if it has changed
            if (profile.userName !== user.username) {
                await renameUser(user.userId, profile.userName);
            }
            // Then update other profile information
            const updateData = {
                firstName: profile.firstName,
                lastName: profile.lastName,
                bio: profile.bio,
                location: profile.location,
                lat: profile.lat,
                lng: profile.lng,
                thumbnail: profile.thumbnail,
            };
            await updateUser(user.userId, updateData);
            showFeedback("profileUpdateSuccess", "success", "profile");
        } catch (error) {
            showFeedback("profileUpdateError", "error", "profile");
        } finally {
            setSubmitting(false);
        }
    };
    // Handle password form submission
    const handleSubmitPassword = async (e) => {
        e.preventDefault();
        if (!user?.userId) return;
        setFeedbackMessage({ key: "", type: "", section: "", values: {} });
        // Validation
        if (password.newPassword !== password.confirmPassword) {
            showFeedback("errorPasswordMismatch", "error", "password");
            return;
        }
        if (password.newPassword.length < 8) {
            showFeedback("errorPasswordLength", "error", "password", { minLength: 8 });
            return;
        }
        setSubmitting(true);
        try {
            setTimeout(() => {
                showFeedback("passwordChangeSuccess", "success", "password");
                setPassword({ currentPassword: "", newPassword: "", confirmPassword: "" });
                setSubmitting(false);
            }, 1000);
        } catch (error) {
            showFeedback("errorCurrentPassword", "error", "password");
            setSubmitting(false);
        }
    };
    // Render feedback message using translation
    const renderFeedback = (sectionName) => {
        if (feedbackMessage.section === sectionName && feedbackMessage.key) {
            return (
                <div className={`${styles.feedback} ${styles[feedbackMessage.type]}`} role="alert">
                    <div className={styles.feedbackContent}>
                        {feedbackMessage.type === "success" ? (
                            <CheckCircle size={16} className={styles.feedbackIcon} />
                        ) : (
                            <X size={16} className={styles.feedbackIcon} />
                        )}
                        <span className={styles.feedbackText}>
                            {t(feedbackMessage.key, feedbackMessage.values)}
                        </span>
                    </div>
                </div>
            );
        }
        return null;
    };
    // Loading state
    if (loading) {
        return (
            <div className={styles.loadingContainer}>
                <div className={styles.loadingContent}>
                    <Loader size={32} className={styles.loadingIcon} />
                    <h3>{t('loading')}</h3>
                    <p>{t('pageDescription')}</p>
                </div>
            </div>
        );
    }
    // Error state
    if (error) {
        return (
            <div className={styles.errorContainer}>
                <div className={styles.errorContent}>
                    <AlertCircle size={48} className={styles.errorIcon} />
                    <h3>{t('errorTitle')}</h3>
                    <p>{t(error)}</p>
                    <button 
                        className={styles.retryButton}
                        onClick={() => window.location.reload()}
                    >
                        {t('retry')}
                    </button>
                </div>
            </div>
        );
    }
    return (
        <div className={styles.container}>
            {/* Profile Information Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <User size={20} className={styles.sectionIcon} />
                        <h3>{t('profileTitle')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        Update your personal information and profile details
                    </p>
                </div>
                <div className={styles.card}>
                    <form className={styles.form} onSubmit={handleSubmitProfile}>
                        <div className={styles.formGrid}>
                            <div className={styles.formGroup}>
                                <label htmlFor="userName" className={styles.label}>
                                    <Edit3 size={16} className={styles.labelIcon} />
                                    {t('profileUsernameLabel')}
                                </label>
                                <input 
                                    type="text" 
                                    id="userName" 
                                    name="userName" 
                                    className={styles.input} 
                                    value={profile.userName}
                                    onChange={handleProfileChange} 
                                    placeholder="Enter username"
                                    required
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="email" className={styles.label}>
                                    <Mail size={16} className={styles.labelIcon} />
                                    {t('profileEmailLabel')}
                                </label>
                                <input 
                                    type="email" 
                                    id="email" 
                                    name="email" 
                                    className={`${styles.input} ${styles.inputDisabled}`} 
                                    value={profile.email}
                                    disabled
                                    title={t('emailChangeDisabled')}
                                />
                                <span className={styles.inputHelp}>
                                    {t('emailChangeDisabled')}
                                </span>
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="firstName" className={styles.label}>
                                    {t('profileFirstNameLabel')}
                                </label>
                                <input 
                                    type="text" 
                                    id="firstName" 
                                    name="firstName" 
                                    className={styles.input} 
                                    value={profile.firstName}
                                    onChange={handleProfileChange}
                                    placeholder="Enter first name"
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="lastName" className={styles.label}>
                                    {t('profileLastNameLabel')}
                                </label>
                                <input 
                                    type="text" 
                                    id="lastName" 
                                    name="lastName" 
                                    className={styles.input} 
                                    value={profile.lastName}
                                    onChange={handleProfileChange}
                                    placeholder="Enter last name"
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="location" className={styles.label}>
                                    <MapPin size={16} className={styles.labelIcon} />
                                    {t('profileLocationLabel')}
                                </label>
                                <input 
                                    type="text" 
                                    id="location" 
                                    name="location" 
                                    className={styles.input} 
                                    value={profile.location}
                                    onChange={handleProfileChange}
                                    placeholder="Enter your location"
                                />
                            </div>
                            <div className={styles.formGroupFull}>
                                <label htmlFor="bio" className={styles.label}>
                                    <FileText size={16} className={styles.labelIcon} />
                                    {t('profileBioLabel')}
                                </label>
                                <textarea 
                                    id="bio" 
                                    name="bio" 
                                    rows="4" 
                                    className={styles.textarea}
                                    value={profile.bio}
                                    onChange={handleProfileChange}
                                    placeholder="Tell us about yourself..."
                                />
                            </div>
                        </div>
                        {renderFeedback("profile")}
                        <div className={styles.formActions}>
                            <button 
                                type="submit" 
                                className={styles.primaryButton}
                                disabled={submitting}
                            >
                                {submitting ? (
                                    <>
                                        <Loader size={16} className={styles.buttonSpinner} />
                                        {t('profileSaving')}
                                    </>
                                ) : (
                                    <>
                                        <CheckCircle size={16} />
                                        {t('profileSaveButton')}
                                    </>
                                )}
                            </button>
                        </div>
                    </form>
                </div>
            </section>
            {/* Security Section */}
            <section className={styles.section}>
                <div className={styles.sectionHeader}>
                    <div className={styles.sectionTitle}>
                        <Shield size={20} className={styles.sectionIcon} />
                        <h3>{t('passwordTitle')}</h3>
                    </div>
                    <p className={styles.sectionDescription}>
                        Update your password to keep your account secure
                    </p>
                </div>
                <div className={styles.card}>
                    <form className={styles.form} onSubmit={handleSubmitPassword}>
                        <div className={styles.formGrid}>
                            <div className={styles.formGroupFull}>
                                <label htmlFor="currentPassword" className={styles.label}>
                                    <Lock size={16} className={styles.labelIcon} />
                                    {t('passwordCurrentLabel')}
                                </label>
                                <input 
                                    type="password" 
                                    id="currentPassword" 
                                    name="currentPassword" 
                                    className={styles.input}
                                    value={password.currentPassword} 
                                    onChange={handlePasswordChange}
                                    placeholder="Enter current password"
                                    required
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="newPassword" className={styles.label}>
                                    {t('passwordNewLabel')}
                                </label>
                                <input 
                                    type="password" 
                                    id="newPassword" 
                                    name="newPassword" 
                                    className={styles.input}
                                    value={password.newPassword} 
                                    onChange={handlePasswordChange}
                                    placeholder="Enter new password"
                                    required
                                />
                            </div>
                            <div className={styles.formGroup}>
                                <label htmlFor="confirmPassword" className={styles.label}>
                                    {t('passwordConfirmLabel')}
                                </label>
                                <input 
                                    type="password" 
                                    id="confirmPassword" 
                                    name="confirmPassword" 
                                    className={styles.input}
                                    value={password.confirmPassword} 
                                    onChange={handlePasswordChange}
                                    placeholder="Confirm new password"
                                    required
                                />
                            </div>
                        </div>
                        {renderFeedback("password")}
                        <div className={styles.formActions}>
                            <button 
                                type="submit" 
                                className={styles.secondaryButton}
                                disabled={submitting}
                            >
                                {submitting ? (
                                    <>
                                        <Loader size={16} className={styles.buttonSpinner} />
                                        {t('passwordChanging')}
                                    </>
                                ) : (
                                    <>
                                        <Lock size={16} />
                                        {t('passwordChangeButton')}
                                    </>
                                )}
                            </button>
                        </div>
                    </form>
                </div>
            </section>
        </div>
    );
});
AccountSettings.displayName = 'AccountSettings';
export default AccountSettings;