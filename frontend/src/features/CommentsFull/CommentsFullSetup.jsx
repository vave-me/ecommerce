// src/comments/CommentsFullSetup.jsx
import React from 'react';
import PropTypes from 'prop-types';
import CommentsFull from './CommentsFull';
import LazyNATSProvider from '../../components/Utils/LazyNATSProvider';
/**
 * CommentsFullSetup
 * Wraps the CommentsFull component in a LazyNATSProvider
 * so that NATS is only loaded when needed.
 */
function CommentsFullSetup({itemId, itemType, userId, categoryId, toggleCommentsFullList}) {
    return (
        <LazyNATSProvider>
            <CommentsFull
                itemId={itemId}
                userId={userId}
                categoryId={categoryId}
                toggleCommentsFullList={toggleCommentsFullList}
                itemType={itemType}
            />
        </LazyNATSProvider>
    );
}
CommentsFullSetup.propTypes = {
    itemId: PropTypes.string.isRequired,
    categoryId: PropTypes.string.isRequired,
    toggleCommentsFullList: PropTypes.func.isRequired,
    userId: PropTypes.string,
};
export default CommentsFullSetup;
