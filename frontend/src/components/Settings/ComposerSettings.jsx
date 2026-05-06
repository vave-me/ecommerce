import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { 
  toggleUnifiedComposer,
  toggleComposerOnPage,
  selectShowUnifiedComposer,
  selectShowComposerOnHomepage,
  selectShowComposerOnMarketplace,
  selectShowComposerOnProducts,
  toggleAutoEnableAIMode,
  selectAutoEnableAIMode
} from '../../redux/slices/uiPreferencesSlice';
import styles from './ComposerSettings.module.css';

const ComposerSettings = () => {
  const dispatch = useDispatch();
  
  // Selectors
  const showUnifiedComposer = useSelector(selectShowUnifiedComposer);
  const showComposerOnHomepage = useSelector(selectShowComposerOnHomepage);
  const showComposerOnMarketplace = useSelector(selectShowComposerOnMarketplace);
  const showComposerOnProducts = useSelector(selectShowComposerOnProducts);
  const autoEnableAIMode = useSelector(selectAutoEnableAIMode);

  return (
    <div className={styles.container}>
      <h3 className={styles.title}>Composer Settings</h3>
      
      <div className={styles.section}>
        <h4 className={styles.sectionTitle}>General</h4>
        
        <div className={styles.setting}>
          <label className={styles.label}>
            <input
              type="checkbox"
              checked={showUnifiedComposer}
              onChange={() => dispatch(toggleUnifiedComposer())}
              className={styles.checkbox}
            />
            <span>Show Unified Composer</span>
          </label>
          <p className={styles.description}>
            Display the composer bar for creating posts and searching
          </p>
        </div>

        <div className={styles.setting}>
          <label className={styles.label}>
            <input
              type="checkbox"
              checked={autoEnableAIMode}
              onChange={() => dispatch(toggleAutoEnableAIMode())}
              className={styles.checkbox}
            />
            <span>Auto-enable AI Mode</span>
          </label>
          <p className={styles.description}>
            Automatically switch to AI mode when typing questions
          </p>
        </div>
      </div>

      <div className={styles.section}>
        <h4 className={styles.sectionTitle}>Page-Specific Settings</h4>
        
        <div className={styles.setting}>
          <label className={styles.label}>
            <input
              type="checkbox"
              checked={showComposerOnHomepage}
              onChange={() => dispatch(toggleComposerOnPage({ page: 'homepage' }))}
              className={styles.checkbox}
              disabled={!showUnifiedComposer}
            />
            <span>Show on Homepage</span>
          </label>
        </div>

        <div className={styles.setting}>
          <label className={styles.label}>
            <input
              type="checkbox"
              checked={showComposerOnMarketplace}
              onChange={() => dispatch(toggleComposerOnPage({ page: 'marketplace' }))}
              className={styles.checkbox}
              disabled={!showUnifiedComposer}
            />
            <span>Show on Marketplace</span>
          </label>
        </div>

        <div className={styles.setting}>
          <label className={styles.label}>
            <input
              type="checkbox"
              checked={showComposerOnProducts}
              onChange={() => dispatch(toggleComposerOnPage({ page: 'products' }))}
              className={styles.checkbox}
              disabled={!showUnifiedComposer}
            />
            <span>Show on Products Page</span>
          </label>
        </div>
      </div>

      {!showUnifiedComposer && (
        <div className={styles.notice}>
          <p>The composer is currently hidden. Enable it above to show it again.</p>
        </div>
      )}
    </div>
  );
};

export default ComposerSettings;