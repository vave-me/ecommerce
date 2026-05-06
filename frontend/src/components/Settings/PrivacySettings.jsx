// src/components/Settings/PrivacySettings.jsx
"use client"
import React, { useState, memo } from 'react';
import { Plus } from '@/icons';
import { FaBan, FaCheckCircle, FaExclamationTriangle, FaShieldAlt, FaUserShield } from '../../utils/iconImports';
import ToggleSwitch from './ToggleSwitch';
import Dropdown from './Dropdown';
import styles from './PrivacySettings.module.css';
/**
 * PrivacySettings Component
 * Allows users to configure privacy settings.
 */
const PrivacySettings = memo(() => {
    // State for privacy settings
    const [privacy, setPrivacy] = useState({
        profileVisibility: 'Public', // Options: Public, Friends Only, Private
        dataSharing: true,
        blockedUsers: ['user456', 'user789'],
    });
    // State for new blocked user
    const [blockedUser, setBlockedUser] = useState('');
    // State for feedback messages
    const [feedbackMessage, setFeedbackMessage] = useState({
        message: "",
        type: "",  // success or error
        section: ""  // visibility, sharing, blocked
    });
    // Handle profile visibility change
    const handleVisibilityChange = (e) => {
        setPrivacy({...privacy, profileVisibility: e.target.value});
        setFeedbackMessage({
            message: `Profile visibility changed to ${e.target.value}.`,
            type: "success",
            section: "visibility"
        });
        setTimeout(() => {
            setFeedbackMessage({ message: "", type: "", section: "" });
        }, 3000);
    };
    // Handle data sharing toggle
    const handleDataSharingToggle = () => {
        setPrivacy({...privacy, dataSharing: !privacy.dataSharing});
        setFeedbackMessage({
            message: !privacy.dataSharing
                ? "Data sharing enabled."
                : "Data sharing disabled.",
            type: "success",
            section: "sharing"
        });
        setTimeout(() => {
            setFeedbackMessage({ message: "", type: "", section: "" });
        }, 3000);
    };
    // Handle adding a blocked user
    const handleAddBlockedUser = () => {
        if (blockedUser.trim() && !privacy.blockedUsers.includes(blockedUser.trim())) {
            setPrivacy({
                ...privacy,
                blockedUsers: [...privacy.blockedUsers, blockedUser.trim()],
            });
            setFeedbackMessage({
                message: `User "${blockedUser.trim()}" has been blocked.`,
                type: "success",
                section: "blocked"
            });
            setBlockedUser('');
            setTimeout(() => {
                setFeedbackMessage({ message: "", type: "", section: "" });
            }, 3000);
        } else if (blockedUser.trim() === '') {
            setFeedbackMessage({
                message: "Please enter a user ID to block.",
                type: "error",
                section: "blocked"
            });
            setTimeout(() => {
                setFeedbackMessage({ message: "", type: "", section: "" });
            }, 3000);
        } else {
            setFeedbackMessage({
                message: "This user is already in your blocked list.",
                type: "error",
                section: "blocked"
            });
            setTimeout(() => {
                setFeedbackMessage({ message: "", type: "", section: "" });
            }, 3000);
        }
    };
    // Handle removing a blocked user
    const handleRemoveBlockedUser = (userId) => {
        const updatedBlockedUsers = privacy.blockedUsers.filter((id) => id !== userId);
        setPrivacy({...privacy, blockedUsers: updatedBlockedUsers});
        setFeedbackMessage({
            message: `User "${userId}" has been unblocked.`,
            type: "success",
            section: "blocked"
        });
        setTimeout(() => {
            setFeedbackMessage({ message: "", type: "", section: "" });
        }, 3000);
    };
    // Handle privacy settings save
    const handleSave = () => {
        // Handle saving privacy settings to the backend
        setFeedbackMessage({
            message: "Privacy settings saved successfully!",
            type: "success",
            section: ""
        });
        setTimeout(() => {
            setFeedbackMessage({ message: "", type: "", section: "" });
        }, 3000);
    };
    return (
        <div className={styles.container}>
            <h2 className={styles.pageTitle}>Privacy Settings</h2>
            <p className={styles.pageDescription}>
                Control who can see your profile and how your information is used.
            </p>
            {/* Profile Visibility Section */}
            <section className={styles.section}>
                <h3 className={styles.sectionTitle}>
                    <FaUserShield className={styles.sectionIcon} />
                    Profile Visibility
                </h3>
                <div className={styles.card}>
                    <div className={styles.cardContent}>
                        <div className={styles.dropdownContainer}>
                            <Dropdown
                                label="Choose who can see your profile:"
                                options={['Public', 'Friends Only', 'Private']}
                                value={privacy.profileVisibility}
                                onChange={handleVisibilityChange}
                            />
                        </div>
                        <div className={styles.visibilityInfo}>
                            {privacy.profileVisibility === 'Public' && (
                                <div className={styles.visibilityExplainer}>
                                    <strong>Public:</strong> Your profile is visible to everyone, including search engines.
                                </div>
                            )}
                            {privacy.profileVisibility === 'Friends Only' && (
                                <div className={styles.visibilityExplainer}>
                                    <strong>Friends Only:</strong> Only your connections can view your full profile.
                                </div>
                            )}
                            {privacy.profileVisibility === 'Private' && (
                                <div className={styles.visibilityExplainer}>
                                    <strong>Private:</strong> Your profile is visible only to you and administrators.
                                </div>
                            )}
                        </div>
                    </div>
                    {feedbackMessage.section === "visibility" && (
                        <div className={`${styles.feedback} ${styles[feedbackMessage.type]}`}>
                            {feedbackMessage.type === "success" ? (
                                <FaCheckCircle className={styles.feedbackIcon} />
                            ) : (
                                <FaExclamationTriangle className={styles.feedbackIcon} />
                            )}
                            <span>{feedbackMessage.message}</span>
                        </div>
                    )}
                </div>
            </section>
            {/* Data Sharing Section */}
            <section className={styles.section}>
                <h3 className={styles.sectionTitle}>
                    <FaShieldAlt className={styles.sectionIcon} />
                    Data Sharing
                </h3>
                <div className={styles.card}>
                    <div className={styles.cardContent}>
                        <div className={styles.dataShareInfo}>
                            <h4 className={styles.dataShareTitle}>Third-party Data Sharing</h4>
                            <p className={styles.dataShareDescription}>
                                Allow us to share your non-personal data with trusted partners to improve
                                services and provide personalized experiences.
                            </p>
                        </div>
                        <div className={styles.toggleContainer}>
                            <ToggleSwitch
                                id="dataSharing"
                                isOn={privacy.dataSharing}
                                handleToggle={handleDataSharingToggle}
                                ariaLabel="Toggle Data Sharing"
                            />
                        </div>
                    </div>
                    {feedbackMessage.section === "sharing" && (
                        <div className={`${styles.feedback} ${styles[feedbackMessage.type]}`}>
                            {feedbackMessage.type === "success" ? (
                                <FaCheckCircle className={styles.feedbackIcon} />
                            ) : (
                                <FaExclamationTriangle className={styles.feedbackIcon} />
                            )}
                            <span>{feedbackMessage.message}</span>
                        </div>
                    )}
                    <div className={styles.infoBox}>
                        <p className={styles.dataShareStatus}>
                            {privacy.dataSharing
                                ? 'Your non-personal data is being shared with trusted third parties.'
                                : 'Your data is not shared with third parties.'}
                        </p>
                    </div>
                </div>
            </section>
            {/* Blocked Users Section */}
            <section className={styles.section}>
                <h3 className={styles.sectionTitle}>
                    <FaBan className={styles.sectionIcon} />
                    Blocked Users
                </h3>
                {feedbackMessage.section === "blocked" && (
                    <div className={`${styles.feedback} ${styles[feedbackMessage.type]}`}>
                        {feedbackMessage.type === "success" ? (
                            <FaCheckCircle className={styles.feedbackIcon} />
                        ) : (
                            <FaExclamationTriangle className={styles.feedbackIcon} />
                        )}
                        <span>{feedbackMessage.message}</span>
                    </div>
                )}
                <div className={styles.blockedList}>
                    {privacy.blockedUsers.length > 0 ? (
                        <ul className={styles.blockedUsersList}>
                            {privacy.blockedUsers.map((userId) => (
                                <li key={userId} className={styles.blockedUserItem}>
                                    <span className={styles.blockedUserName}>{userId}</span>
                                    <button
                                        className={styles.unblockButton}
                                        onClick={() => handleRemoveBlockedUser(userId)}
                                        aria-label={`Unblock ${userId}`}
                                    >
                                        Unblock
                                    </button>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <div className={styles.emptyState}>
                            <p className={styles.emptyMessage}>
                                No users blocked. Blocked users won't be able to contact you or view your content.
                            </p>
                        </div>
                    )}
                </div>
                <div className={styles.addBlockedUserContainer}>
                    <div className={styles.formRow}>
                        <input
                            type="text"
                            placeholder="Enter user ID to block"
                            value={blockedUser}
                            onChange={(e) => setBlockedUser(e.target.value)}
                            className={styles.input}
                        />
                        <button
                            className={styles.addButton}
                            onClick={handleAddBlockedUser}
                            aria-label="Block User"
                        >
                            <Plus className={styles.addIcon} /> Block User
                        </button>
                    </div>
                    <p className={styles.helperText}>
                        Enter the user ID of the person you wish to block. Blocked users cannot see your profile or contact you.
                    </p>
                </div>
            </section>
            <div className={styles.actions}>
                <button className={styles.saveButton} onClick={handleSave}>
                    Save Privacy Settings
                </button>
            </div>
            {/* Global success message */}
            {feedbackMessage.section === "" && feedbackMessage.message && (
                <div className={`${styles.globalFeedback} ${styles[feedbackMessage.type]}`}>
                    <FaShieldAlt className={styles.feedbackIcon} />
                    <span>{feedbackMessage.message}</span>
                </div>
            )}
        </div>
    );
});
PrivacySettings.displayName = 'PrivacySettings';
export default PrivacySettings;