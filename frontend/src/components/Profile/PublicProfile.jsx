"use client";
import React, { useEffect, useState, memo } from "react";
import { useTranslations } from "next-intl";
import { getBaseUserById } from "../../api/client/userApi";
// Child Components
import ProfileHeader from "./ProfileHeader";
import Tabs from "./Tabs";
import ReviewList from "./ReviewList";
import SocialLinks from "./SocialLinks";
import styles from "./PublicProfile.module.css";
const PublicProfile = memo(function PublicProfile({ userId }) {
    const t = useTranslations('SellerProfile');
    const [seller, setSeller] = useState(null);
    const [activeTab, setActiveTab] = useState("Items");
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    useEffect(() => {
        if (!userId) {
            setError(t('errorNoSellerId'));
            setLoading(false);
            return;
        }
        setLoading(true);
        setError(null);
        getBaseUserById(userId)
            .then(result => {
                if (result && result.user) {
                    setSeller(result.user);
                } else {
                    throw new Error(t('errorSellerNotFound'));
                }
            })
            .catch((err) => {
                setError(t('errorLoadingFailed', { message: err.message || 'Unknown error' }));
            })
            .finally(() => {
                setLoading(false);
            });
    }, [userId, t]);
    if (loading) {
        return <div className={styles.loadingContainer}>{t("loading")}</div>;
    }
    if (error) {
        return <div className={styles.errorContainer} role="alert">{error}</div>;
    }
    if (!seller) {
        return <div className={styles.emptyContainer}>{t('errorSellerNotFound')}</div>;
    }
    return (
        <div className={styles.container}>
            <ProfileHeader 
                username={seller.userName} 
                userId={seller.id} 
                isPublicProfile={true}
            />
            <Tabs activeTab={activeTab} setActiveTab={setActiveTab} />
            <div className={styles.content}>
                {/* Items Tab Panel */}
                <div id={`panel-Items`} role="tabpanel" tabIndex={0} aria-labelledby={`tab-Items`}
                     hidden={activeTab !== "Items"}>
                    {activeTab === "Items" && (
                        <section className={styles.section} aria-label={t("section.listedItems")}>
                            <h3 className={styles.sectionTitle}>{t("section.listedItems")}</h3>
                        </section>
                    )}
                </div>
                {/* Reviews Tab Panel */}
                <div id={`panel-Reviews`} role="tabpanel" tabIndex={0} aria-labelledby={`tab-Reviews`}
                     hidden={activeTab !== "Reviews"}>
                    {activeTab === "Reviews" && (
                        <section className={styles.section} aria-label={t("section.reviews")}>
                            <h3 className={styles.sectionTitle}>{t("section.reviews")}</h3>
                            <ReviewList userId={seller.id} />
                        </section>
                    )}
                </div>
                {/* Gallery Tab Panel */}
                <div id={`panel-Gallery`} role="tabpanel" tabIndex={0} aria-labelledby={`tab-Gallery`}
                     hidden={activeTab !== "Gallery"}>
                    {activeTab === "Gallery" && (
                        <section className={styles.section} aria-label={t("section.gallery")}>
                            <h3 className={styles.sectionTitle}>{t("section.gallery")}</h3>
                            <SocialLinks userId={seller.id} />
                        </section>
                    )}
                </div>
            </div>
        </div>
    );
}, (prevProps, nextProps) => {
    // Only re-render if userId changes
    return prevProps.userId === nextProps.userId;
});
export default PublicProfile; 