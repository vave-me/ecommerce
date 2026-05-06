// File: src/comments/Comments.jsx
"use client";
import React, { useState, useEffect, useRef, useCallback, useMemo, memo } from "react";
import PropTypes from "prop-types";
import { debounce } from '../../utils/lightweightLibraries';
import {useQuery, useQueryClient} from "@tanstack/react-query";
import {Loader2} from "@/icons";
import {useAuth} from "../../context/AuthContext";
import {useNATS} from "../../context/NATSContext";
import {jetstream as message_api} from "../../generated_proto/message_api_pb";
import {commentspb as comments_api} from "../../generated_proto/comments_api_pb";
import {message_type as message_types} from "../../generated_proto/message_types_pb";
import {getCommentsByProduct} from "../../api/commentsApi";
import {useCommentsActions, addCommentToCache} from "../../hooks/useCommentsActions";
import CommentsHeader from "./CommentsHeader";
import CommentsList from "./CommentsList";
import CommentsInput from "./CommentsInput";
import styles from "./Comments.module.css";
/* Build nested comment tree - memoized to prevent unnecessary recalculations */
const buildCommentTree = (comments) => {
    const cloned = comments.map((c) => ({...c, replies: []}));
    const commentMap = {};
    // index comments by id
    cloned.forEach((c) => {
        commentMap[c.id] = c;
    });
    const rootComments = [];
    cloned.forEach((c) => {
        if (c.parentId && commentMap[c.parentId]) {
            // it's a reply -> push into parent's replies
            commentMap[c.parentId].replies.push(c);
        } else if (!c.parentId || c.parentId === "") {
            // top-level
            rootComments.push(c);
        }
    });
    return rootComments;
};
const Comments = memo(function Comments({
                      itemId,
                      itemType,
                      categoryId,
                      toggleCommentsList,
                      productName,
                      productThumbnail,
                  }) {
    const {user} = useAuth() || {};
    const userId = user?.userId;
    const {subscribe, isConnected} = useNATS();
    const queryClient = useQueryClient();
    const {createComment} = useCommentsActions(itemId, userId, categoryId);
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
    // ---------------------------------------------------
    // 1) React Query fetch for existing comments
    // ---------------------------------------------------
    const {
        data: commentsData = [],
        isLoading,
        isFetching,
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
    const sortedComments = useMemo(() => {
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
        
        let unsubscribeFn = null;
        
        // Subscribe is async, so we need to handle the promise
        const setupSubscription = async () => {
            try {
                unsubscribeFn = await subscribe(natsName, handleIncomingMessage);
            } catch (error) {
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', error);
        }
    }
        };
        
        setupSubscription();
        
        // Cleanup function
        return () => {
            if (unsubscribeFn && typeof unsubscribeFn === 'function') {
                unsubscribeFn();
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
    return (
        <div className={styles.commentsContainer}>
            <div className={styles.filterBar}>
                <button
                    type="button"
                    className={`${styles.filterButton} ${
                        sortOption === "Newest" ? styles.filterButtonSelected : ""
                    }`}
                    onClick={() => handleSortChange("Newest")}
                >
                    Newest
                </button>
                <button
                    type="button"
                    className={`${styles.filterButton} ${
                        sortOption === "Top Reactions" ? styles.filterButtonSelected : ""
                    }`}
                    onClick={() => handleSortChange("Top Reactions")}
                >
                    Top Reactions
                </button>
                <button
                    type="button"
                    className={`${styles.filterButton} ${
                        sortOption === "Oldest" ? styles.filterButtonSelected : ""
                    }`}
                    onClick={() => handleSortChange("Oldest")}
                >
                    Oldest
                </button>
            </div>
            <div className={styles.commentsContent}>
                {isLoading || isFetching ? (
                    <div className={styles.loadingState}>
                        <Loader2 size={24} className={styles.loadingIcon}/>
                        <p className={styles.loadingMessage}>Loading comments...</p>
                    </div>
                ) : isError ? (
                    <div className={styles.errorState}>
                        <p className={styles.errorMessage}>Error loading comments. Please try again.</p>
                    </div>
                ) : (
                    <CommentsList comments={sortedComments}/>
                )}
            </div>
            {/* Authentication-aware comment input */}
            {userId ? (
                <CommentsInput
                    commentText={commentText}
                    setCommentText={setCommentText}
                    onSubmitComment={handleSubmitComment}
                    isConnected={isConnected}
                    placeholder="Add a comment..."
                    disabled={isLoading}
                />
            ) : (
                <div className={styles.loginPrompt}>
                    <p className={styles.loginMessage}>
                        Please <a href="/login" className={styles.loginLink}>log in</a> to add comments.
                    </p>
                </div>
            )}
        </div>
    );
});
// Add display name for better debugging
Comments.displayName = 'Comments';
Comments.propTypes = {
    toggleCommentsList: PropTypes.func.isRequired,
    itemId: PropTypes.string.isRequired,
    itemType: PropTypes.string.isRequired,
    categoryId: PropTypes.string.isRequired,
    productName: PropTypes.string,
    productThumbnail: PropTypes.string,
};
Comments.defaultProps = {
    productName: "",
    productThumbnail: "",
};
export default Comments;