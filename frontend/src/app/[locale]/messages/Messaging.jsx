"use client";
import React, {useCallback, useEffect, useState, useMemo} from "react";
import {useRouter} from "next/navigation";
import {useQuery} from "@tanstack/react-query";
import {useNATS} from "../../../context/NATSContext";
import {useAuth} from "../../../context/AuthContext";
import {getConversations, getMessagesByConversation, startConversation,} from "../../../api/client/messagingApi";
import {getBaseUserById} from "../../../api/client/userApi";
import MessageList from "../../../components/Messaging/MessageList";
import ChatWindow from "../../../components/Messaging/ChatWindow";
import styles from "./Messaging.module.css";
import MinNav from "../../../components/Header/MinNav";
function Messaging() {
    const router = useRouter();
    const {user} = useAuth();
    const {isConnected} = useNATS();
    const userId = user?.userId; // might be undefined
    const [isMessagingOpen, setIsMessagingOpen] = useState(true);
    const [selectedConversationId, setSelectedConversationId] = useState(null);
    const [messages, setMessages] = useState([]);
    const [recipient, setRecipient] = useState(null);
    const [conversations, setConversations] = useState([]);
    const [conversationsLoading, setConversationsLoading] = useState(false);
    const [conversationsError, setConversationsError] = useState(null);
    // Enhanced error handling utility
    const handleApiError = useCallback((error, context) => {
        const errorMessage = error?.response?.data?.message || error?.message || `Failed to ${context.toLowerCase()}`;
        return errorMessage;
    }, []);
    // Fetch conversations with proper error handling
    useEffect(() => {
        if (!userId) return;
        let isMounted = true;
        setConversationsLoading(true);
        setConversationsError(null);
        const fetchConversations = async () => {
            try {
                const data = await getConversations(userId);
                if (isMounted) {
                    setConversations(data || []);
                    setConversationsError(null);
                }
            } catch (err) {
                if (isMounted) {
                    const errorMessage = handleApiError(err, "Load conversations");
                    setConversationsError(errorMessage);
                }
            } finally {
                if (isMounted) {
                    setConversationsLoading(false);
                }
            }
        };
        fetchConversations();
        return () => {
            isMounted = false;
        };
    }, [userId, handleApiError]);
    // Fixed useQuery - removed deprecated onSuccess callback
    const {
        data: fetchedMessages,
        isLoading: messagesLoading,
        isError: messagesError,
        error: messagesErrorDetails
    } = useQuery({
        queryKey: ["messages", selectedConversationId],
        queryFn: () => getMessagesByConversation(selectedConversationId),
        enabled: !!selectedConversationId && !!userId,
        retry: 2,
        retryDelay: 1000,
    });
    // Handle messages transformation when data changes
    useEffect(() => {
        if (fetchedMessages && userId) {
            const newTransformed = fetchedMessages.map((m) => ({
                id: m.id,
                text: m.body,
                recipientId: m.recipientId,
                senderId: m.senderId,
                isUserMessage: m.senderId === userId,
                time: new Date().toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"}),
                createdAt: m.createdAt,
                itemId: m.itemId,
            }));
            setMessages(newTransformed);
        }
    }, [fetchedMessages, userId]);
    // Enhanced conversation starter
    const handleStartConversation = useCallback(
        async (recipientId, itemId) => {
            if (!userId) {
                return;
            }
            try {
                const result = await startConversation(userId, recipientId, itemId);
                if (result?.id) {
                    setSelectedConversationId(result.id);
                    setConversations((prev) => {
                        // Check if conversation already exists
                        const existingIndex = prev.findIndex(c => c.id === result.id);
                        if (existingIndex >= 0) {
                            return prev; // Don't duplicate
                        }
                        return [
                            {
                                id: result.id,
                                senderId: result.senderId,
                                recipientId: result.recipientId,
                                itemId: result.itemId,
                                createdAt: result.createdAt || new Date().toISOString(),
                                updatedAt: result.updatedAt || new Date().toISOString(),
                            },
                            ...prev
                        ];
                    });
                }
            } catch (err) {
                const errorMessage = handleApiError(err, "Start conversation");
            }
        },
        [userId, handleApiError]
    );
    // Enhanced conversation selector with better error handling
    const handleSelectConversation = useCallback(
        async (conversationId) => {
            if (!conversationId || !userId) return;
            setSelectedConversationId(conversationId);
            setRecipient(null); // Clear previous recipient
            try {
                // Find conversation in local state first
                const convo = conversations.find((c) => c.id === conversationId);
                let recipientId = null;
                if (convo && convo.senderId && convo.recipientId) {
                    recipientId = convo.senderId === userId ? convo.recipientId : convo.senderId;
                } else if (fetchedMessages && fetchedMessages.length > 0) {
                    // Fallback: Use message data to determine recipient ID
                    const firstMessage = fetchedMessages[0];
                    recipientId = firstMessage.senderId === userId ? firstMessage.recipientId : firstMessage.senderId;
                }
                if (recipientId) {
                    try {
                        const userResponse = await getBaseUserById(recipientId);
                        setRecipient({
                            id: recipientId,
                            name: userResponse?.user?.userName || userResponse?.userName || "Unknown user",
                            avatar: userResponse?.user?.avatar || userResponse?.avatar || "/images/user-user.webp",
                            online: !!isConnected,
                        });
                    } catch (error) {
                        // Fallback to basic recipient info
                        setRecipient({
                            id: recipientId,
                            name: "Unknown user",
                            avatar: "/images/user-user.webp",
                            online: !!isConnected,
                        });
                    }
                }
            } catch (err) {
                const errorMessage = handleApiError(err, "Select conversation");
            }
        },
        [conversations, userId, isConnected, fetchedMessages, handleApiError]
    );
    // Completed mapConversationsToUI function with proper async user data fetching
    const mapConversationsToUI = useMemo(() => {
        return async (convos, currentUserId) => {
            if (!convos || !currentUserId) return [];
            const conversationsWithUsers = await Promise.allSettled(
                convos.map(async (c) => {
                    const recipientId = c.senderId === currentUserId ? c.recipientId : c.senderId;
                    try {
                        const userResponse = await getBaseUserById(recipientId);
                        return {
                            id: c.id,
                            user: {
                                id: recipientId,
                                name: userResponse?.user?.userName || userResponse?.userName || "Unknown User",
                                avatar: userResponse?.user?.avatar || userResponse?.avatar || "/images/user-user.webp",
                                online: false, // Could be enhanced with real-time status
                            },
                            lastMessage: c.lastMessage || "No recent message",
                            time: c.updatedAt || c.createdAt || "N/A",
                            unreadCount: c.unreadCount || 0,
                        };
                    } catch (error) {
                        return {
                            id: c.id,
                            user: {
                                id: recipientId,
                                name: "Unknown User",
                                avatar: "/images/user-user.webp",
                                online: false,
                            },
                            lastMessage: c.lastMessage || "No recent message",
                            time: c.updatedAt || c.createdAt || "N/A",
                            unreadCount: c.unreadCount || 0,
                        };
                    }
                })
            );
            // Filter out rejected promises and return only fulfilled ones
            return conversationsWithUsers
                .filter(result => result.status === 'fulfilled')
                .map(result => result.value);
        };
    }, []);
    // Enhanced UI conversations with loading state
    const [uiConversations, setUiConversations] = useState([]);
    const [uiConversationsLoading, setUiConversationsLoading] = useState(false);
    useEffect(() => {
        if (!conversations.length || !userId) {
            setUiConversations([]);
            return;
        }
        let isMounted = true;
        setUiConversationsLoading(true);
        const loadUiConversations = async () => {
            try {
                const mappedConversations = await mapConversationsToUI(conversations, userId);
                if (isMounted) {
                    setUiConversations(mappedConversations);
                }
            } catch (error) {
                if (isMounted) {
                    setUiConversations([]);
                }
            } finally {
                if (isMounted) {
                    setUiConversationsLoading(false);
                }
            }
        };
        loadUiConversations();
        return () => {
            isMounted = false;
        };
    }, [conversations, userId, mapConversationsToUI]);
    // Early return if messaging is closed
    if (!isMessagingOpen) {
        return null;
    }
    // Early return if user isn't logged in
    if (!userId) {
        return (
            <div className={styles.notLoggedIn}>
                <h2>You must be logged in to view or send messages.</h2>
                <p>Please log in to access your conversations and send messages.</p>
            </div>
        );
    }
    // Determine loading state for conversations
    const isConversationsLoading = conversationsLoading || uiConversationsLoading;
    const hasConversationsError = conversationsError;
    return (
        <div className={styles.messagingContainer}>
            <div className={styles.mobileHeader}>
                <div className={styles.mobileHeaderContainer}>
                    <MinNav locationPath="/messages"/>
                </div>
            </div>
            <div className={styles.messagingContent}>
                {/* LEFT SIDEBAR */}
                <div
                    className={`${styles.sidebar} ${
                        !selectedConversationId ? "" : styles.sidebarHiddenOnMobile
                    }`}
                >
                    <MessageList
                        conversations={uiConversations}
                        onSelectConversation={handleSelectConversation}
                        selectedConversationId={selectedConversationId}
                        onClose={() => setIsMessagingOpen(false)}
                        loading={isConversationsLoading}
                        error={hasConversationsError}
                    />
                </div>
                {/* RIGHT CHAT AREA */}
                <div
                    className={`${styles.chatSection} ${
                        selectedConversationId ? "" : styles.chatSectionHiddenOnMobile
                    }`}
                >
                    {selectedConversationId ? (
                        messagesLoading ? (
                            <div className={styles.loadingContainer}>
                                <div className={styles.loadingSpinner}></div>
                                <span>Loading messages...</span>
                            </div>
                        ) : messagesError ? (
                            <div className={styles.errorContainer}>
                                <span className={styles.errorIcon}>⚠️</span>
                                <h3>Error loading messages</h3>
                                <p>{messagesErrorDetails?.message || "Failed to load messages. Please try again."}</p>
                                <button 
                                    className={styles.retryButton}
                                    onClick={() => window.location.reload()}
                                >
                                    Retry
                                </button>
                            </div>
                        ) : (
                            <ChatWindow
                                messages={messages}
                                recipient={recipient}
                                conversationId={selectedConversationId}
                                senderId={userId}
                                isConnected={isConnected}
                                onCloseChat={() => setSelectedConversationId(null)}
                                onClose={() => setIsMessagingOpen(false)}
                                onPrevious={() => setSelectedConversationId(null)}
                                onStartConversation={handleStartConversation}
                                itemId={selectedConversationId}
                                recipientId={recipient?.id}
                            />
                        )
                    ) : (
                        <div className={styles.noConversationSelected}>
                            <h2>Select a conversation</h2>
                            <p>Choose a conversation from the list on the left to start messaging.</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
export default Messaging;
