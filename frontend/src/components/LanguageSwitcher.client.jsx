// src/components/LanguageSwitcher.client.jsx
"use client";
import React, {memo} from 'react';
import {useRouter, usePathname} from 'next/navigation';
import {useLocale} from 'next-intl';
const LanguageSwitcher = memo(function LanguageSwitcher() {
    const router = useRouter();
    const pathname = usePathname();
    const locale = useLocale();
    const switchLanguage = (newLocale) => {
        // Remove current locale from pathname and add new one
        const pathWithoutLocale = pathname.replace(/^\/[a-z]{2}/, '');
        const newPath = `/${newLocale}${pathWithoutLocale}`;
        router.push(newPath);
    };
    return (
        <select
            value={locale}
            onChange={(e) => switchLanguage(e.target.value)}
            className="language-switcher"
        >
            <option value="en">English</option>
            <option value="pl">Polski</option>
            <option value="de">Deutsch</option>
        </select>
    );
});
export default LanguageSwitcher;
