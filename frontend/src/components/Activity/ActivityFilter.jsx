"use client";
import React, { memo } from 'react';
import PropTypes from "prop-types";
import { useTranslations } from "next-intl";
import styles from "./ActivityFilter.module.css";
const ActivityFilter = memo(function ActivityFilter({ currentFilter, onFilterChange }) {
    // 🔑 namespace "ActivityFilter" lives in /messages/…/activityFilter.json
    const t = useTranslations("ActivityFilter");
    const filters = [
        { label: t("all"),      value: "all" },
        { label: t("likes"),    value: "like" },
        { label: t("disliked"), value: "dislike" },
        { label: t("visited"),  value: "visit" },
        { label: t("commented"),value: "comment" },
        { label: t("shared"),   value: "share" },
        { label: t("followed"), value: "follow" }
    ];
    return (
        <div className={styles.filterContainer}>
            {filters.map(({ label, value }) => {
                const isActive = currentFilter === value;
                return (
                    <button
                        key={value}
                        className={`${styles.filterOption} ${isActive ? styles.active : ""}`}
                        onClick={() => onFilterChange(value)}
                        aria-label={t("aria", { label })}
                    >
                        {label}
                    </button>
                );
            })}
        </div>
    );
});
ActivityFilter.propTypes = {
    currentFilter:  PropTypes.string.isRequired,
    onFilterChange: PropTypes.func.isRequired
};
export default ActivityFilter;
