// File: src/comments/CommentsFull.jsx
"use client";
import React, { useState, useEffect, useRef, useCallback, useMemo, memo } from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";
import { debounce } from '../../utils/lightweightLibraries';
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useAuth } from "../../context/AuthContext";
import { useNATS } from "../../context/NATSContext";
import { jetstream as message_api } from "../../generated_proto/message_api_pb";
import { commentspb as comments_api } from "../../generated_proto/comments_api_pb";
import { message_type as message_types } from "../../generated_proto/message_types_pb";
import { getCommentsByProduct } from "../../api/commentsApi";
import { useCommentsActions, addCommentToCache } from "../../hooks/useCommentsActions";
import CommentsFullHeader from "./CommentsFullHeader";
import CommentsFullList from "./CommentsFullList";
import CommentsFullInput from "./CommentsFullInput";
// Import the new CSS module
import styles from "./CommentsFull.module.css";
/* Build nested comment tree */
function buildCommentTree(comments) {
    const cloned = comments.map((c) => ({ ...c, replies: [] }));
    const commentMap = {};
    // index comments by id
    cloned.forEach((c) => {
        commentMap[c.id] = c;
    });
    const rootCommentsFull = [];
    cloned.forEach((c) => {
        if (c.parentId && commentMap[c.parentId]) {
            // it's a reply -> push into parent's replies
            commentMap[c.parentId].replies.push(c);
        } else if (!c.parentId || c.parentId === "") {
            // top-level
            rootCommentsFull.push(c);
        }
    });
    return rootCommentsFull;
}
function CommentsFull({
                      itemId,
                      itemType,
                      categoryId,
                      toggleCommentsFullList,
                  }) {
    const { user } = useAuth() || {};
    const userId = user?.userId;
    const { subscribe, isConnected } = useNATS();
    const queryClient = useQueryClient();
    const { createComment } = useCommentsActions(itemId, userId, categoryId);
    const natsName = process.env.NEXT_NATS_CM_NAME || "comments.AddComment";
    const isMountedRef = useRef(true);
    // Local state
    const [commentText, setCommentText] = useState("");
    const [sortOption, setSortOption] = useState("Newest");
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
        };
    }, []);
    // Prevent background scrolling when CommentsFull are open
    useEffect(() => {
        const prevOverflow = document.body.style.overflow;
        document.body.style.overflow = "hidden";
        return () => {
            document.body.style.overflow = prevOverflow;
        };
    }, []);
    // ---------------------------------------------------
    // 1) React Query fetch for existing comments
    // ---------------------------------------------------
    const {
        data: commentsData = [],
        isLoading,
        isError,
    } = useQuery({
        queryKey: ["comments", itemId],
        queryFn: async () => {
            const result = await getCommentsByProduct(itemId);
            // result.comments is presumably an array of {id, parentId, etc.}
            return result.comments;
        },
        staleTime: 5 * 60 * 1000,
        select: buildCommentTree,
    });
    // ---------------------------------------------------
    // 2) Sorting logic
    // ---------------------------------------------------
    const sortedCommentsFull = useMemo(() => {
        const copy = [...commentsData];
        if (sortOption === "Newest") {
            copy.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
        } else if (sortOption === "Oldest") {
            copy.sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt));
        } else if (sortOption === "Top Reactions") {
            copy.sort((a, b) => (b.reactions || 0) - (a.reactions || 0));
        }
        return copy;
    }, [commentsData, sortOption]);
    // Memoize the debounced message handler to prevent recreating it on every render
    const handleIncomingMessage = useMemo(() => 
        debounce((rawBytes) => {
            if (!isMountedRef.current) return;
            try {
                // Decode raw message from NATS
                const streamMsg = message_api.StreamMessage.decode(rawBytes);
                const wsData = message_types.WebsocketMessageData.decode(streamMsg.data);
                const addCommentPb = comments_api.AddComment.decode(wsData.payload);
                const newComment = {
                    id: addCommentPb.id,
                    senderId: addCommentPb.senderId,
                    itemId: addCommentPb.itemId,
                    itemType: addCommentPb.itemType,
                    content: addCommentPb.content,
                    categoryId: addCommentPb.categoryId,
                    parentId: addCommentPb.parentId === "" ? null : addCommentPb.parentId,
                    createdAt: new Date().toISOString(),
                    replies: [],
                };
                // Insert into the React Query cache
                queryClient.setQueryData(["comments", itemId], (old = []) =>
                    addCommentToCache(old, newComment)
                );
            } catch (err) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', err);
        }
    }
        }, 300), 
    [queryClient, itemId]);
    // ---------------------------------------------------
    // 3) NATS subscription
    // ---------------------------------------------------
    useEffect(() => {
        if (!isConnected || !natsName) return;
        const unsubscribe = subscribe(natsName, handleIncomingMessage);
        return () => {
            if (unsubscribe) {
                unsubscribe();
            }
        };
    }, [isConnected, natsName, subscribe, handleIncomingMessage]);
    // ---------------------------------------------------
    // Handle comment submission
    // ---------------------------------------------------
    const handleSubmitComment = useCallback(
        (text) => {
            if (!userId) {
                // TODO: Replace with toast notification
        if (typeof window !== 'undefined' && window.toast) {
            window.toast.warn('Please log in to add comments.');
        }
                return;
            }
            createComment(text, null); // null means top-level comment
        },
        [createComment, userId]
    );
    const handleSortChange = useCallback((newSortOption) => {
        setSortOption(newSortOption);
    }, []);
    // ---------------------------------------------------
    // 4) Render the component
    // ---------------------------------------------------
    return ReactDOM.createPortal(
        <div className={styles.overlay} onClick={toggleCommentsFullList}>
            <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
                <CommentsFullHeader
                    toggleCommentsFullList={toggleCommentsFullList}
                    sortOption={sortOption}
                    onSortChange={handleSortChange}
                />
                <div className={styles.commentsContent}>
                    {isLoading ? (
                        <p className={styles.loadingMessage}>Loading comments...</p>
                    ) : isError ? (
                        <p className={styles.errorMessage}>Error loading comments.</p>
                    ) : (
                        <CommentsFullList comments={sortedCommentsFull} />
                    )}
                </div>
                {/* Authentication-aware comment input */}
                {userId ? (
                    <CommentsFullInput
                        commentText={commentText}
                        setCommentText={setCommentText}
                        onSubmitComment={handleSubmitComment}
                        isConnected={isConnected}
                        placeholder="Add a comment..."
                    />
                ) : (
                    <div className={styles.loginPrompt}>
                        <p className={styles.loginMessage}>
                            Please <a href="/login" className={styles.loginLink}>log in</a> to add comments.
                        </p>
                    </div>
                )}
            </div>
        </div>,
        document.body
    );
}
CommentsFull.propTypes = {
    toggleCommentsFullList: PropTypes.func.isRequired,
    itemId: PropTypes.string.isRequired,
    itemType: PropTypes.string.isRequired,
    categoryId: PropTypes.string.isRequired,
};
export default memo(CommentsFull);
