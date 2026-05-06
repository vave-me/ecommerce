"use client";
import {useEffect, useState, useCallback, useRef} from "react";
import {v4 as uuidv4} from "uuid";
import {useNATS} from "../context/NATSContext";
import {getMessagesByConversation} from "../api/client/messagingApi";
import {mes as messages_api} from "../generated_proto/messages_api_events_pb";
import {message_type} from "../generated_proto/message_types_pb";
import {jetstream as message_api} from "../generated_proto/message_api_pb";
import {useAuth} from "../context/AuthContext";
/**
 * useChatHistory
 * - Loads historical messages from your REST endpoint (getMessagesByConversation).
 * - Subscribes via ephemeral push for new messages in real-time.
 * - Provides sendMessage() for publishing new messages to NATS.
 *
 * Now includes ID checks to prevent duplicates from multi-subscriptions or
 * from the "optimistic" local append.
 */
export default function useChatHistory(
    conversationId,
    {
        recipientId, // required for building new message
        itemId,      // optional
        metadata,    // optional
    } = {}
) {
    const {isConnected, publish, subscribe} = useNATS();
    // 1) Safely destructure user
    const {user} = useAuth();
    const userId = user?.userId; // may be undefined if not logged in
    const [messages, setMessages] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const natsName = process.env.NEXT_NATS_SM_NAME || "messenger.SendMessage";
    // Keep track of subscription cleanup
    const unsubscribeRef = useRef(null);
    // ------------------------------
    // 1) Load historical messages from REST
    // ------------------------------
    useEffect(() => {
        // If conversation or user is missing, skip
        if (!conversationId || !userId) return;
        let isMounted = true;
        const fetchHistory = async () => {
            setIsLoading(true);
            setError(null);
            try {
                const fetched = await getMessagesByConversation(conversationId);
                // Transform/rename fields if needed (body -> text, etc.)
                const mapped = fetched.map((m) => ({
                    id: m.id,
                    text: m.body,
                    senderId: m.senderId,
                    recipientId: m.recipientId,
                    conversationId: m.conversationId,
                    createdAt: m.createdAt || Date.now(),
                    // For UI: local user check
                    isUserMessage: m.senderId === userId,
                }));
                if (isMounted) {
                    setMessages(mapped);
                }
            } catch (err) {
                if (isMounted) {
                    setError(err);
                }
            } finally {
                if (isMounted) {
                    setIsLoading(false);
                }
            }
        };
        fetchHistory();
        return () => {
            isMounted = false;
        };
    }, [conversationId, userId]);
    // ------------------------------
    // 2) Subscribe for real-time updates
    // ------------------------------
    useEffect(() => {
        // If conversation or user is missing or not connected, skip
        if (!conversationId || !isConnected || !userId) return;
        let isMounted = true;
        // Unsubscribe old subscription if it exists
        if (unsubscribeRef.current) {
            unsubscribeRef.current();
            unsubscribeRef.current = null;
        }
        const subject = `${natsName}.${conversationId}`;
        (async () => {
            try {
                const unsub = await subscribe(subject, (rawBytes) => {
                    if (!isMounted) return;
                    try {
                        // 1) Decode the outer StreamMessage
                        const decodedStreamMessage = message_api.StreamMessage.decode(rawBytes);
                        // 2) Decode the WebsocketMessageData
                        const wsData = message_type.WebsocketMessageData.decode(
                            decodedStreamMessage.data
                        );
                        // 3) Decode the final SendMessage
                        const msg = messages_api.SendMessage.decode(wsData.payload);
                        // Build a local shape
                        const newIncoming = {
                            id: msg.id,
                            text: msg.body,
                            senderId: msg.senderId,
                            recipientId: msg.recipientId,
                            itemId: msg.itemId,
                            conversationId: msg.conversationId,
                            createdAt: Date.now(),
                            isUserMessage: msg.senderId === userId,
                        };
                        // >>> ID CHECK to prevent duplicates <<<
                        setMessages((prev) => {
                            // If we already have this exact ID, skip
                            if (prev.some((m) => m.id === newIncoming.id)) {
                                return prev;
                            }
                            // Otherwise add it
                            return [...prev, newIncoming];
                        });
                    } catch (decodingErr) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', decodingErr);
        }
    }
                });
                unsubscribeRef.current = unsub;
            } catch (err) {
                setError(err);
            }
        })();
        return () => {
            isMounted = false;
            if (unsubscribeRef.current) {
                unsubscribeRef.current();
                unsubscribeRef.current = null;
            }
        };
    }, [conversationId, isConnected, subscribe, userId]);
    // ------------------------------
    // 3) sendMessage: publish a new message
    // ------------------------------
    const sendMessage = useCallback(
        async (text) => {
            // Check for required data
            if (!isConnected) {
                setError(new Error("Not connected to NATS"));
                return;
            }
            if (!conversationId || !userId || !recipientId) {
                setError(new Error("Missing conversationId / userId / recipientId"));
                return;
            }
            const trimmed = (text || "").trim();
            if (!trimmed) return;
            try {
                // Build SendMessage proto
                const msgId = uuidv4();
                const messagePayload = messages_api.SendMessage.create({
                    id: msgId,
                    conversationId: conversationId,
                    senderId: userId,
                    recipientId: recipientId,
                    itemId: itemId || "",
                    body: trimmed,
                    isRead: false,
                });
                const msgBytes = messages_api.SendMessage.encode(messagePayload).finish();
                // Wrap with WebsocketMessageData
                const wsCommand = message_type.WebsocketMessageData.create({
                    payload: msgBytes,
                    occurred_at: {seconds: Math.floor(Date.now() / 1000), nanos: 0},
                });
                const encCommand = message_type.WebsocketMessageData.encode(wsCommand).finish();
                // Build the final StreamMessage
                const streamMessage = message_api.StreamMessage.create({
                    id: uuidv4(),
                    name: natsName,
                    data: encCommand,
                    metadata: metadata || {user: userId, role: "sender"},
                    sent_at: {
                        seconds: Math.floor(Date.now() / 1000),
                        nanos: 0,
                    },
                });
                const encodedStreamMessage =
                    message_api.StreamMessage.encode(streamMessage).finish();
                // Publish
                const subject = `${natsName}.${conversationId}`;
                await publish(subject, encodedStreamMessage);
                // (Optional) Optimistic append
                const newLocalMsg = {
                    id: msgId,
                    text: trimmed,
                    senderId: userId,
                    recipientId,
                    conversationId,
                    itemId: itemId,
                    createdAt: Date.now(),
                    isUserMessage: true,
                };
                // >>> ID CHECK to prevent duplicates if subscription arrives quickly <<<
                setMessages((prev) => {
                    if (prev.some((m) => m.id === newLocalMsg.id)) {
                        return prev;
                    }
                    return [...prev, newLocalMsg];
                });
            } catch (err) {
                setError(err);
            }
        },
        [conversationId, userId, recipientId, itemId, publish, isConnected, metadata]
    );
    // Return the hook's data & methods
    return {
        messages,
        isLoading,
        error,
        sendMessage,
    };
}
