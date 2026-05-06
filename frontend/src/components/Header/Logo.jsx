"use client";
import React, { memo } from "react";
import Link from "next/link";
import Image from "next/image";
import styles from "./Logo.module.css";
const Logo = memo(function Logo({ size = "default", aiMode = false }) {
    const containerClasses = [
        styles.logoContainer,
        styles[size],
        aiMode ? styles.aiMode : ''
    ].filter(Boolean).join(' ');
    // Direct path to sfx.png - no complex config
    const logoSrc = "/images/sfx.png";
    // Simple size mapping based on the actual sfx.png dimensions (59x25)
    const sizeMap = {
        small: { width: 50, height: 21 },
        default: { width: 90, height: 38 }, // 1.5x scale
        large: { width: 110, height: 46 }, // ~2x scale
        aiMode: { width: 80, height: 34 } // Smaller AI mode scale
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
                        style={{
                            width: currentSize.width + 'px',
                            height: currentSize.height + 'px'
                        }}
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
export default Logo;