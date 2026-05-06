// src/comments/CommentsSetup.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import Comments from './Comments';
import LazyNATSProvider from '../../components/Utils/LazyNATSProvider';
/**
 * CommentsSetup
 * Wraps the Comments component in a LazyNATSProvider
 * so that NATS is only loaded when needed.
 */
const CommentsSetup = memo(function CommentsSetup({itemId, itemType, userId, categoryId, toggleCommentsList, dealName, dealThumbnail}) {
    return (
        <LazyNATSProvider>
            <Comments
                itemId={itemId}
                userId={userId}
                categoryId={categoryId}
                toggleCommentsList={toggleCommentsList}
                itemType={itemType}
                dealName={dealName}
                dealThumbnail={dealThumbnail}
            />
        </LazyNATSProvider>
    );
});
// Add display name for debugging
CommentsSetup.displayName = 'CommentsSetup';
CommentsSetup.propTypes = {
    itemId: PropTypes.string.isRequired,
    itemType: PropTypes.string.isRequired,
    categoryId: PropTypes.string.isRequired,
    toggleCommentsList: PropTypes.func.isRequired,
    userId: PropTypes.string,
    dealName: PropTypes.string,
    dealThumbnail: PropTypes.string,
};
export default CommentsSetup;
