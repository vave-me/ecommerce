"use client";
import React, { useEffect, useState, useCallback } from 'react';
import { useTranslations } from "next-intl";
import { useAuth } from "../../../context/AuthContext";
import {
    FileText as FileTextIcon,
    Image as ImageIcon,
    MessageSquare as MessageSquareIcon,
    Package as PackageIcon,
    Star as StarIcon,
    User as UserIcon,
    Loader
} from "@/icons";
// Components
import ProfileHeader from "../../../components/Profile/ProfileHeader";
import ItemList from "../../../components/Profile/ItemList";
import CommentList from "../../../components/Profile/CommentsList";
import ReviewList from "../../../components/Profile/ReviewList";
import SocialLinks from "../../../components/Profile/SocialLinks";
// Styles
import styles from "./UserProfile.module.css";

/* ---------------------------------------------------------------------- */
/* Tab Configuration - Optimized for User Profile                        */
/* ---------------------------------------------------------------------- */
const TABS = [
    { id: "Items", key: "tabs.listings", icon: PackageIcon },
    { id: "Comments", key: "tabs.comments", icon: MessageSquareIcon },
    { id: "Reviews", key: "tabs.reviews", icon: StarIcon },
    { id: "Gallery", key: "tabs.gallery", icon: ImageIcon },
    // Additional tabs can be enabled as needed
    // { id: "Posts", key: "tabs.posts", icon: FileTextIcon },
    // { id: "Jobs", key: "tabs.jobs", icon: FileTextIcon },
    // { id: "Vehicles", key: "tabs.vehicles", icon: FileTextIcon },
    // { id: "Services", key: "tabs.services", icon: FileTextIcon },
    // { id: "Videos", key: "tabs.videos", icon: FileTextIcon },
    // { id: "Deals", key: "tabs.deals", icon: FileTextIcon },
];

/**
 * UserProfile Component - Modern, Clean Implementation
 * Features:
 * - Professional design aligned with design system
 * - Proper error boundaries and loading states
 * - Accessibility compliant (WCAG 2.1 AA)
 * - Responsive design with sticky navigation
 * - Smooth animations and transitions
 */
export default function UserProfile() {
    const t = useTranslations('UserProfile');
    const { user, isLoading: authLoading, authChecked } = useAuth();
    
    // Component state
    const [activeTab, setActiveTab] = useState("Items");
    const [isScrolled, setIsScrolled] = useState(false);
    const [isClient, setIsClient] = useState(false);

    // Handle client-side hydration
    useEffect(() => {
        setIsClient(true);
    }, []);

    // Optimized scroll handler with sticky tabs
    const handleScroll = useCallback(() => {
        const scrollY = window.scrollY;
        setIsScrolled(scrollY > 200);
    }, []);

    useEffect(() => {
        if (!isClient) return;
        
        window.addEventListener("scroll", handleScroll, { passive: true });
        return () => window.removeEventListener("scroll", handleScroll);
    }, [isClient, handleScroll]);

    // Tab change handler with accessibility
    const handleTabChange = useCallback((tabId) => {
        setActiveTab(tabId);
        
        // Focus management for accessibility
        const panelElement = document.getElementById(`panel-${tabId}`);
        if (panelElement) {
            panelElement.focus();
        }
    }, []);

    // Show loading while checking authentication
    if (authLoading || !authChecked || !isClient) {
        return (
            <div className={styles.container}>
                <div className={styles.emptyState}>
                    <div className={styles.emptyStateContent}>
                        <Loader className={styles.emptyStateIcon} size={48} />
                        <h2>{t("loading", { defaultValue: "Loading profile..." })}</h2>
                        <p>{t("loadingDesc", { defaultValue: "Please wait while we load your profile information." })}</p>
                    </div>
                </div>
            </div>
        );
    }

    // Handle unauthenticated state
    if (!user) {
        return (
            <div className={styles.container}>
                <div className={styles.emptyState}>
                    <div className={styles.emptyStateContent}>
                        <UserIcon className={styles.emptyStateIcon} size={64} />
                        <h2>{t("noProfileTitle", { defaultValue: "Sign in to view your profile" })}</h2>
                        <p>{t("noProfileMsg", { defaultValue: "You need to be signed in to access your profile and manage your content." })}</p>
                    </div>
                </div>
            </div>
        );
    }

    // Handle missing user data
    if (!user.userId) {
        return (
            <div className={styles.container}>
                <div className={styles.errorState}>
                    <h2>{t("errorTitle", { defaultValue: "Profile Error" })}</h2>
                    <p>{t("errorUserData", { defaultValue: "Unable to load profile data. Please try refreshing the page." })}</p>
                    <button 
                        onClick={() => window.location.reload()}
                        style={{
                            padding: '10px 20px',
                            background: '#3b82f6',
                            color: 'white',
                            border: 'none',
                            borderRadius: '8px',
                            cursor: 'pointer',
                            fontSize: '14px',
                            fontWeight: '500'
                        }}
                    >
                        {t("refresh", { defaultValue: "Refresh Page" })}
                    </button>
                </div>
            </div>
        );
    }

    // Render the profile interface
    return (
        <div className={styles.container}>
            {/* Profile Header */}
            <ProfileHeader 
                username={user.username || user.userName} 
                userId={user.userId}
                userData={user}
                isPublicProfile={false}
            />

            {/* Tabs Navigation */}
            <div className={`${styles.tabsContainer} ${isScrolled ? styles.stickyTabs : ""}`}>
                <nav 
                    className={styles.tabs} 
                    role="tablist" 
                    aria-label={t('tabsAriaLabel', { defaultValue: "Profile navigation" })}
                >
                    {TABS.map((tab) => {
                        const IconComponent = tab.icon;
                        return (
                            <button
                                key={tab.id}
                                id={`tab-${tab.id}`}
                                className={`${styles.tabButton} ${activeTab === tab.id ? styles.activeTab : ""}`}
                                onClick={() => handleTabChange(tab.id)}
                                aria-selected={activeTab === tab.id}
                                role="tab"
                                aria-controls={`panel-${tab.id}`}
                                type="button"
                                tabIndex={activeTab === tab.id ? 0 : -1}
                            >
                                <IconComponent className={styles.tabIcon} aria-hidden="true" />
                                <span className={styles.tabLabel}>
                                    {t(tab.key, { defaultValue: tab.id })}
                                </span>
                            </button>
                        );
                    })}
                </nav>
            </div>

            {/* Tab Content */}
            <div className={styles.content}>
                {/* Items Tab Panel */}
                {activeTab === "Items" && (
                    <div 
                        id="panel-Items" 
                        role="tabpanel" 
                        tabIndex={0} 
                        aria-labelledby="tab-Items"
                        className={styles.section}
                    >
                        <header className={styles.sectionHeader}>
                            <h2 className={styles.sectionTitle}>
                                {t("section.listedItems", { defaultValue: "Listed Items" })}
                            </h2>
                            <span className={styles.sectionCount} aria-label={t("itemsCount", { defaultValue: "Items count" })}>
                                -
                            </span>
                        </header>
                        <ItemList userId={user.userId} />
                    </div>
                )}

                {/* Comments Tab Panel */}
                {activeTab === "Comments" && (
                    <div 
                        id="panel-Comments" 
                        role="tabpanel" 
                        tabIndex={0} 
                        aria-labelledby="tab-Comments"
                        className={styles.section}
                    >
                        <header className={styles.sectionHeader}>
                            <h2 className={styles.sectionTitle}>
                                {t("section.comments", { defaultValue: "Comments" })}
                            </h2>
                            <span className={styles.sectionCount} aria-label={t("commentsCount", { defaultValue: "Comments count" })}>
                                -
                            </span>
                        </header>
                        <CommentList username={user.username || user.userName} userId={user.userId} />
                    </div>
                )}

                {/* Reviews Tab Panel */}
                {activeTab === "Reviews" && (
                    <div 
                        id="panel-Reviews" 
                        role="tabpanel" 
                        tabIndex={0} 
                        aria-labelledby="tab-Reviews"
                        className={styles.section}
                    >
                        <header className={styles.sectionHeader}>
                            <h2 className={styles.sectionTitle}>
                                {t("section.reviews", { defaultValue: "Reviews" })}
                            </h2>
                            <span className={styles.sectionCount} aria-label={t("reviewsCount", { defaultValue: "Reviews count" })}>
                                -
                            </span>
                        </header>
                        <ReviewList username={user.username || user.userName} userId={user.userId} />
                    </div>
                )}

                {/* Gallery Tab Panel */}
                {activeTab === "Gallery" && (
                    <div 
                        id="panel-Gallery" 
                        role="tabpanel" 
                        tabIndex={0} 
                        aria-labelledby="tab-Gallery"
                        className={styles.section}
                    >
                        <header className={styles.sectionHeader}>
                            <h2 className={styles.sectionTitle}>
                                {t("section.gallery", { defaultValue: "Gallery" })}
                            </h2>
                        </header>
                        <SocialLinks username={user.username || user.userName} userId={user.userId} />
                    </div>
                )}
            </div>
        </div>
    );
}