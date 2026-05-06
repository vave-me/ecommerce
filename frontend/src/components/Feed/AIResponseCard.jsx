import React, { memo } from 'react';
import { Bot, ThumbsUp, ThumbsDown, Share2, Bookmark } from '@/icons';
import styles from './AIResponseCard.module.css';
import { useAuth } from '../../context/AuthContext';

/**
 * Card component for displaying AI responses in the feed
 * Styled to match other feed cards but with AI-specific indicators
 */
const AIResponseCard = ({ 
  response, 
  query, 
  timestamp, 
  metadata = {},
  onFeedback = null,
  onShare = null,
  onSave = null 
}) => {
  const { user } = useAuth();
  
  const handleFeedback = (type) => {
    if (onFeedback) {
      onFeedback(response.id, type);
    }
  };

  const formatTimestamp = (ts) => {
    const date = new Date(ts);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) return 'Just now';
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`;
    return date.toLocaleDateString();
  };

  return (
    <div className={styles.card}>
      <div className={styles.header}>
        <div className={styles.aiInfo}>
          <div className={styles.aiAvatar}>
            <Bot size={24} />
          </div>
          <div className={styles.aiMeta}>
            <span className={styles.aiName}>AI Assistant</span>
            <span className={styles.timestamp}>{formatTimestamp(timestamp)}</span>
          </div>
        </div>
        <div className={styles.aiIndicator}>
          <span>AI Response</span>
        </div>
      </div>

      {query && (
        <div className={styles.querySection}>
          <span className={styles.queryLabel}>You asked:</span>
          <p className={styles.queryText}>{query}</p>
        </div>
      )}

      <div className={styles.responseSection}>
        <div className={styles.responseContent}>
          {response}
        </div>
        
        {metadata.sources && metadata.sources.length > 0 && (
          <div className={styles.sources}>
            <span className={styles.sourcesLabel}>Sources:</span>
            <ul className={styles.sourcesList}>
              {metadata.sources.map((source, idx) => (
                <li key={idx}>
                  <a href={source.url} target="_blank" rel="noopener noreferrer">
                    {source.title}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        )}

        {metadata.relatedProducts && metadata.relatedProducts.length > 0 && (
          <div className={styles.relatedProducts}>
            <h4>Related Products</h4>
            <div className={styles.productGrid}>
              {metadata.relatedProducts.slice(0, 4).map(product => (
                <a 
                  key={product.id} 
                  href={`/product/${product.id}`} 
                  className={styles.relatedProduct}
                >
                  <img src={product.thumbnail} alt={product.name} />
                  <span className={styles.productName}>{product.name}</span>
                  <span className={styles.productPrice}>${product.price}</span>
                </a>
              ))}
            </div>
          </div>
        )}

        {metadata.suggestions && metadata.suggestions.length > 0 && (
          <div className={styles.suggestions}>
            <h4>You might also ask:</h4>
            <ul>
              {metadata.suggestions.map((suggestion, idx) => (
                <li key={idx}>
                  <button className={styles.suggestionButton}>
                    {suggestion}
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      <div className={styles.actions}>
        <button 
          className={styles.actionButton}
          onClick={() => handleFeedback('helpful')}
          title="This was helpful"
        >
          <ThumbsUp size={18} />
          <span>Helpful</span>
        </button>
        
        <button 
          className={styles.actionButton}
          onClick={() => handleFeedback('not-helpful')}
          title="This wasn't helpful"
        >
          <ThumbsDown size={18} />
          <span>Not helpful</span>
        </button>
        
        <button 
          className={styles.actionButton}
          onClick={() => onShare && onShare(response)}
          title="Share response"
        >
          <Share2 size={18} />
          <span>Share</span>
        </button>
        
        <button 
          className={styles.actionButton}
          onClick={() => onSave && onSave(response)}
          title="Save response"
        >
          <Bookmark size={18} />
          <span>Save</span>
        </button>
      </div>
    </div>
  );
};

export default memo(AIResponseCard);