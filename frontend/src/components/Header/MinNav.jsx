"use client";
import React, {useMemo, memo} from "react";
import PropTypes from "prop-types";
import {useRouter} from "next/navigation";
import {useTranslations} from "next-intl"; //  Import hook
import {
    ArrowLeftIcon,
    Bell,
    FilmIcon,
    Heart,
    Home,
    MessageCircle,
    NewspaperIcon,
    Percent,
    ShoppingBag,
    TruckIcon
} from "@/icons";
import styles from "./MinNav.module.css";
import MinNavItem from "./MinNavItem"; // Assume MinNavItem handles its own translations (e.g., badge aria-label)
// Base navigation items configuration (without text)
const baseNavItemsData = [
    {type: "back", icon: ArrowLeftIcon, id: 'back'},
    {to: "/home", icon: Home, id: 'home'},
    {to: "/news", icon: NewspaperIcon, id: 'news'},
    {to: "/deals", icon: Percent, id: 'deals'},
    {to: "/products", icon: ShoppingBag, id: 'shop'}, // Note: ID 'shop' used for label/tooltip keys
    {to: "/videos", icon: FilmIcon, id: 'videos'},
    {to: "/cars", icon: TruckIcon, id: 'vehicles'}, // Note: ID 'vehicles' used for label/tooltip keys
    {to: "/wishlist", icon: Heart, id: 'saved', badgeCount: 2}, // Note: ID 'saved' used for label/tooltip keys
    {to: "/messages", icon: MessageCircle, id: 'messages', badgeCount: 5},
    {to: "/notifications", icon: Bell, id: 'alerts', badgeCount: 1}, // Note: ID 'alerts' used for label/tooltip keys
];
const MinNav = memo(function MinNav({locationPath}) {
    const t = useTranslations('MinNav'); //  Instantiate hook
    const router = useRouter();
    // Create translated navigation item components using useMemo
    const translatedNavItems = useMemo(() => {
        return baseNavItemsData.map((item) => {
            //   Translate label and tooltip using item.id
            const label = t(`item_${item.id}_label`);
            const tooltip = t(`item_${item.id}_tooltip`);
            // Handle back navigation as a special case
            if (item.type === "back") {
                return (
                    <MinNavItem
                        key={item.id}
                        icon={item.icon}
                        label={label} // Pass translated label
                        tooltip={tooltip} // Pass translated tooltip
                        onClick={() => router.back()}
                    />
                );
            }
            const isActive = locationPath === item.to;
            const badgeCount = item.badgeCount || 0;
            return (
                <MinNavItem
                    key={item.id}
                    to={item.to}
                    icon={item.icon}
                    label={label} // Pass translated label
                    tooltip={tooltip} // Pass translated tooltip
                    isActive={isActive}
                    badgeCount={badgeCount}
                    // Assume MinNavItem uses useTranslations internally for badge aria-label
                />
            );
        });
    }, [locationPath, router, t]); // Add t back since router.back could change
    return (
        //   Use translation for main nav aria-label
        <nav className={styles.nav} aria-label={t('navAriaLabel')}>
            <div className={styles.navContainer}>
                {/* Render translated nav items */}
                <ul className={styles.navList} role="menubar">
                    {translatedNavItems}
                </ul>
            </div>
        </nav>
    );
});
MinNav.propTypes = {
    locationPath: PropTypes.string.isRequired,
};
export default MinNav;