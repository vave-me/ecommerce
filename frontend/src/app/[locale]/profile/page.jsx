"use client";
import React, { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "../../../context/AuthContext";
import { useTranslations } from "next-intl";
import { User as UserIcon } from "@/icons";
import styles from "./ProfilePage.module.css";

/**
 * Profile Page - Handles the base /profile route
 * - If user is authenticated: redirects to /user (their own profile)
 * - If user is not authenticated: shows login prompt
 * - This route should not show a specific user's profile (that's /profile/[username])
 */
export default function ProfilePage() {
    const router = useRouter();
    const { user, isLoading, authChecked } = useAuth();
    const t = useTranslations('UserProfile');

    useEffect(() => {
        // Only redirect after auth has been checked to avoid redirect loops
        if (authChecked && !isLoading) {
            if (user?.userId) {
                // User is authenticated - redirect to their own profile
                router.replace('/user');
            }
            // If not authenticated, stay on this page to show login prompt
        }
    }, [user, isLoading, authChecked, router]);

    // Show loading while checking authentication
    if (isLoading || !authChecked) {
        return (
            <div className={styles.container}>
                <div className={styles.loadingContainer}>
                    <div className={styles.spinner}></div>
                    <h2>{t("loading", { defaultValue: "Loading..." })}</h2>
                    <p>{t("authChecking", { defaultValue: "Checking authentication..." })}</p>
                </div>
            </div>
        );
    }

    // User is authenticated - they should be redirected, but show fallback just in case
    if (user?.userId) {
        return (
            <div className={styles.container}>
                <div className={styles.redirectContainer}>
                    <div className={styles.redirectIcon}>
                        <UserIcon size={48} aria-hidden="true" />
                    </div>
                    <h2>{t("redirectingToProfile", { defaultValue: "Redirecting to your profile..." })}</h2>
                    <p>{t("redirectingDesc", { defaultValue: "Taking you to your personal profile page." })}</p>
                </div>
            </div>
        );
    }

    // User is not authenticated - show login prompt
    return (
        <div className={styles.container}>
            <div className={styles.authPromptContainer}>
                <div className={styles.authPromptIcon}>
                    <UserIcon size={64} aria-hidden="true" />
                </div>
                <h1 className={styles.authPromptTitle}>
                    {t("profileAccessTitle", { defaultValue: "Access Your Profile" })}
                </h1>
                <p className={styles.authPromptDescription}>
                    {t("profileAccessDesc", { 
                        defaultValue: "Sign in to view and manage your profile, listings, and account settings." 
                    })}
                </p>
                <div className={styles.authPromptActions}>
                    <button 
                        className={styles.loginButton}
                        onClick={() => router.push('/login')}
                        type="button"
                    >
                        {t("signIn", { defaultValue: "Sign In" })}
                    </button>
                    <button 
                        className={styles.registerButton}
                        onClick={() => router.push('/register')}
                        type="button"
                    >
                        {t("createAccount", { defaultValue: "Create Account" })}
                    </button>
                </div>
                <div className={styles.helpText}>
                    <p>
                        {t("profileHelpText", { 
                            defaultValue: "Looking for someone else's profile? Use /profile/username" 
                        })}
                    </p>
                </div>
            </div>
        </div>
    );
}
