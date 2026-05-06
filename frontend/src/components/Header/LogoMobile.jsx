"use client";
import React, { memo } from "react";
import Link from "next/link";
import Image from "next/image";
import styles from "./Logo.module.css";
const LogoMobile = memo(function LogoMobile({ size = "default", aiMode = false }) {
    const containerClasses = [
        styles.logoContainer,
        styles.mobile,
        styles[size],
        aiMode ? styles.aiMode : ''
    ].filter(Boolean).join(' ');
    // Direct path to sfx.png - no complex config
    const logoSrc = "/images/sfx.png";
    // Mobile-optimized size mapping based on the actual sfx.png dimensions (59x25)
    const sizeMap = {
        small: { width: 40, height: 17 }, // 0.7x scale
        default: { width: 50, height: 21 }, // 0.85x scale
        large: { width: 65, height: 27 }, // 1.1x scale
        aiMode: { width: 60, height: 25 } // 1x scale
    };
    const currentSize = aiMode ? sizeMap.aiMode : (sizeMap[size] || sizeMap.default);
    return (
        <div className={containerClasses}>
            <Link href="/" className={styles.logoLink} aria-label="SFX - go to home page">
                <div className={styles.logoImageWrapper}>
                    <Image
                        src={logoSrc}
                        alt="SFX Logo"
                        width={currentSize.width}
                        height={currentSize.height}
                        className={styles.logoImage}
                        priority
                        onError={(e) => {
                        }}
                        onLoad={() => {
                        }}
                    />
                </div>
            </Link>
        </div>
    );
});
export default LogoMobile;