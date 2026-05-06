"use client";
import React, { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useParams, useRouter } from "next/navigation";
import { useAuth } from "../../../../context/AuthContext";
import { getBaseUserById } from "../../../../api/userApi";
// Child Components
import ProfileHeader from "../../../../components/Profile/ProfileHeader";
import Tabs from "../../../../components/Profile/Tabs";
import CommentList from "../../../../components/Profile/CommentsList";
import ReviewList from "../../../../components/Profile/ReviewList";
import ContactInfo from "../../../../components/Profile/ContactInfo";
import SocialLinks from "../../../../components/Profile/SocialLinks";
import ItemList from "../../../../components/Profile/ItemList";
// Styles
import styles from "./Profile.module.css";

/**
 * Public Profile Page - Shows another user's public profile
 * Route: /profile/[username]
 * - Shows public information for the specified username
 * - If username matches current user, redirects to /user (private profile)
 * - Handles authentication properly
 */
export default function PublicProfile() {
    const t = useTranslations('SellerProfile');
    const params = useParams();
    const router = useRouter();
    const { user: currentUser, authChecked, isLoading: authLoading } = useAuth();
    
    // Get username from route params
    const username = params?.username;
    
    // Component state
    const [profileUser, setProfileUser] = useState(null);
    const [activeTab, setActiveTab] = useState("Items");
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    // Check if viewing own profile and redirect if needed
    useEffect(() => {
        if (authChecked && currentUser?.username === username) {
            // User is trying to view their own profile via public route
            // Redirect to private profile page
            router.replace('/user');
            return;
        }
    }, [authChecked, currentUser, username, router]);

    // Fetch profile user data
    useEffect(() => {
        const fetchProfileUser = async () => {
            if (!username) {
                setError(t('errorNoSellerId', { defaultValue: 'No user specified' }));
                setLoading(false);
                return;
            }

            // Don't fetch if redirecting to own profile
            if (currentUser?.username === username) {
                return;
            }

            setLoading(true);
            setError(null);

            try {
                // Use the public API to get user data by username
                const result = await getBaseUserById(username);
                
                if (result && result.user) {
                    setProfileUser(result.user);
                } else {
                    throw new Error(t('errorSellerNotFound', { defaultValue: 'User not found' }));
                }
            } catch (err) {
                // Error: 'Error fetching profile user:', err...
                setError(
                    err.response?.status === 404 
                        ? t('errorSellerNotFound', { defaultValue: 'User not found' })
                        : t('errorLoadingFailed', { 
                            defaultValue: 'Failed to load profile',
                            message: err.message || 'Unknown error'
                        })
                );
            } finally {
                setLoading(false);
            }
        };

        // Only fetch after auth has been checked to avoid unnecessary requests
        if (authChecked) {
            fetchProfileUser();
        }
    }, [username, currentUser, authChecked, t]);

    // Show loading while checking auth or fetching profile
    if (authLoading || !authChecked || loading) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingContainer}>
                    <div className={styles.spinner}></div>
                    <h2>{t("loading", { defaultValue: "Loading profile..." })}</h2>
                    <p>{t("loadingDesc", { defaultValue: "Please wait while we fetch the profile information." })}</p>
                </div>
            </div>
        );
    }

    // Show error state
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.errorContainer} role="alert">
                    <div className={styles.errorIcon}>⚠️</div>
                    <h2 className={styles.errorTitle}>
                        {t("errorTitle", { defaultValue: "Profile Not Found" })}
                    </h2>
                    <p className={styles.errorMessage}>{error}</p>
                    <div className={styles.errorActions}>
                        <button 
                            onClick={() => router.back()}
                            className={styles.backButton}
                            type="button"
                        >
                            {t("goBack", { defaultValue: "Go Back" })}
                        </button>
                        <button 
                            onClick={() => window.location.reload()}
                            className={styles.retryButton}
                            type="button"
                        >
                            {t("tryAgain", { defaultValue: "Try Again" })}
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    // Show empty state if no profile user found
    if (!profileUser) {
        return (
            <div className={styles.container}>
                <div className={styles.emptyContainer}>
                    <div className={styles.emptyIcon}>👤</div>
                    <h2 className={styles.emptyTitle}>
                        {t('errorSellerNotFound', { defaultValue: 'User not found' })}
                    </h2>
                    <p className={styles.emptyMessage}>
                        {t('userNotFoundDesc', { 
                            defaultValue: 'The user you are looking for does not exist or has been removed.',
                            username 
                        })}
                    </p>
                    <button 
                        onClick={() => router.push('/')}
                        className={styles.homeButton}
                        type="button"
                    >
                        {t("goHome", { defaultValue: "Go Home" })}
                    </button>
                </div>
            </div>
        );
    }

    // Render the public profile
    return (
        <div className={styles.container}>
            {/* Profile Header */}
            <ProfileHeader 
                username={profileUser.userName || profileUser.username} 
                userId={profileUser.id || profileUser.userId} 
                isPublicProfile={true}
                userData={profileUser}
            />

            {/* Tabs Navigation */}
            <Tabs 
                activeTab={activeTab} 
                setActiveTab={setActiveTab}
                isPublicProfile={true}
            />

            {/* Tab Content */}
            <div className={styles.content}>
                {/* Items/Products Tab */}
                {activeTab === "Items" && (
                    <div id="panel-Items" role="tabpanel" tabIndex={0} aria-labelledby="tab-Items">
                        <section className={styles.section} aria-label={t("section.listedItems", { defaultValue: "Listed Items" })}>
                            <ItemList userId={profileUser.id || profileUser.userId} />
                        </section>
                    </div>
                )}

                {/* Comments Tab */}
                {activeTab === "Comments" && (
                    <div id="panel-Comments" role="tabpanel" tabIndex={0} aria-labelledby="tab-Comments">
                        <section className={styles.section} aria-label={t("section.comments", { defaultValue: "Comments" })}>
                            <CommentList 
                                username={profileUser.userName || profileUser.username} 
                                userId={profileUser.id || profileUser.userId} 
                            />
                        </section>
                    </div>
                )}

                {/* Reviews Tab */}
                {activeTab === "Reviews" && (
                    <div id="panel-Reviews" role="tabpanel" tabIndex={0} aria-labelledby="tab-Reviews">
                        <section className={styles.section} aria-label={t("section.reviews", { defaultValue: "Reviews" })}>
                            <ReviewList userId={profileUser.id || profileUser.userId} />
                        </section>
                    </div>
                )}

                {/* Contact Tab */}
                {activeTab === "Contact" && (
                    <div id="panel-Contact" role="tabpanel" tabIndex={0} aria-labelledby="tab-Contact">
                        <section className={styles.section} aria-label={t("section.contactInfo", { defaultValue: "Contact Information" })}>
                            <ContactInfo 
                                userId={profileUser.id || profileUser.userId} 
                                user={profileUser}
                            />
                        </section>
                    </div>
                )}

                {/* Gallery/Social Links Tab */}
                {activeTab === "Gallery" && (
                    <div id="panel-Gallery" role="tabpanel" tabIndex={0} aria-labelledby="tab-Gallery">
                        <section className={styles.section} aria-label={t("section.gallery", { defaultValue: "Gallery" })}>
                            <SocialLinks 
                                username={profileUser.userName || profileUser.username} 
                                userId={profileUser.id || profileUser.userId} 
                            />
                        </section>
                    </div>
                )}
            </div>
        </div>
    );
}