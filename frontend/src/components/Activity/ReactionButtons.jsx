"use client";
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { Heart, ThumbsUp } from '@/icons';
import { FaLaugh, FaSadTear, FaAngry } from '../../utils/iconImports';
import styles from './ReactionButtons.module.css';
import { useTranslations } from 'next-intl'; // <-- import translations hook
/**
 * ReactionButtons - Component for rendering reaction buttons (like, love, etc.)
 *
 * @param {Function} onReact - Handler function when a reaction is clicked
 */
const ReactionButtons = memo(({ onReact }) => {
  const t = useTranslations('ReactionButtons'); // <-- translation hook
  const reactions = [
    { type: 'like', icon: <ThumbsUp />, label: t('like') },
    { type: 'love', icon: <Heart />, label: t('love') },
    { type: 'haha', icon: <FaLaugh />, label: t('haha') },
    { type: 'sad', icon: <FaSadTear />, label: t('sad') },
    { type: 'angry', icon: <FaAngry />, label: t('angry') },
  ];
  return (
      <div className={styles.buttonsContainer}>
        {reactions.map((reaction) => (
            <button
                key={reaction.type}
                className={styles.reactionButton}
                onClick={() => onReact(reaction.type)}
                aria-label={reaction.label}
            >
              {reaction.icon}
            </button>
        ))}
      </div>
  );
});
ReactionButtons.displayName = 'ReactionButtons';
ReactionButtons.propTypes = {
  onReact: PropTypes.func.isRequired,
};
export default ReactionButtons;
