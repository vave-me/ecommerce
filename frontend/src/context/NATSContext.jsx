// src/context/NATSContext.jsx
"use client"
import {connect, consumerOpts} from 'nats.ws';
import React, {createContext, useCallback, useContext, useRef, useState} from 'react';
import {v4 as uuidv4} from 'uuid';
import {useAuth} from './AuthContext';
import { ErrorBoundary } from '../components/ErrorBoundary';

const NATSContext = createContext();

export const useNATS = () => {
    const context = useContext(NATSContext);
    if (!context) {
        if (process.env.NODE_ENV === 'development') {
      // Error: 'useNATS must be used within a NATSProvider.'...
    }
        // Return dummy defaults so we don't crash
        return {
            isConnected: false,
            connectIfNeeded: async () => {
            },
            publish: async () => {
            },
            subscribe: async () => {
            },
        };
    }
    return context;
};

export function NATSProvider({children}) {
    const {user, setUserOnlineStatus} = useAuth();

    // NATS connection & state
    const ncRef = useRef(null);
    const [js, setJs] = useState(null);
    const [isConnected, setIsConnected] = useState(false);

    // Timer for marking user offline (if needed)
    const offlineTimerRef = useRef(null);
    const natsServersEnv = process.env.NEXT_NATS_SERVERS || 'http://192.168.178.84:9222';

    /**
     * connectIfNeeded()
     * - Connect to NATS (if not already connected)
     * - Initialize JetStream, etc.
     */
    const connectIfNeeded = useCallback(async () => {
        if (ncRef.current) {
            // Already connected
            return;
        }
        try {
            //const servers = 'http://192.168.178.84:9222'; // or from .env
            const servers = natsServersEnv;
            const connection = await connect({
                servers: [servers],
                reconnectTimeWait: 5000, // ms to wait before trying reconnect
                maxReconnectAttempts: -1,
            });
            ncRef.current = connection;
            setIsConnected(true);

            // If the user is authenticated => mark them as online
            if (user && user.userId) {
                // Clear any pending offline timer if we have one
                if (offlineTimerRef.current) {
                    clearTimeout(offlineTimerRef.current);
                }
                setUserOnlineStatus(true);
            }

            const jetstream = connection.jetstream();
            setJs(jetstream);

            // Handle closure
            connection.closed().then((err) => {
                setIsConnected(false);
                setJs(null);
                ncRef.current = null;

                // If user is logged in => schedule offline if no reconnect
                if (user && user.userId) {
                    offlineTimerRef.current = setTimeout(() => {
                        if (!ncRef.current && !isConnected) {
                            // Mark user offline
                            setUserOnlineStatus(false);
                        }
                    }, 3 * 60 * 1000); // 3 min
                }

                if (err) {
                    if (process.env.NODE_ENV === 'development') {
              // Error: 'NATS connection closed with error:', err...
            }
                }
            });
        } catch (error) {
            if (process.env.NODE_ENV === 'development') {
              // Error: 'Failed to connect to NATS:', error...
            }
        }
    }, [user, isConnected, setUserOnlineStatus]);

    /**
     * publish(subject, message)
     * - Makes sure we’re connected, then publish to JetStream.
     */
    const publish = useCallback(
        async (subject, message) => {
            // Ensure connection
            await connectIfNeeded();
            if (!js) {
                // Error: 'JetStream not initialized'...
                return;
            }
            try {
                await js.publish(subject, message);
            } catch (err) {
                // Initialization error - log but continue
                if (process.env.NODE_ENV === 'development') {
                    console.error('Initialization error:', err);
                }
                // Continue with default behavior
            }
        },
        [connectIfNeeded, js]
    );

    /**
     * subscribe(subject, onMessage)
     * - Connect if needed, create ephemeral subscription, pass raw bytes.
     * - Return an unsubscribe function.
     */
    const subscribe = useCallback(
        async (subject, onMessage) => {
            // Ensure connection
            await connectIfNeeded();
            if (!js) {
                // Error: 'JetStream not initialized. Cannot subscribe:', su...
                return;
            }
            try {
                const opts = consumerOpts()
                    .deliverTo(`_deliver_${uuidv4()}`) // ephemeral deliver subject
                    .ackExplicit()
                    .deliverNew(); // only new messages

                const sub = await js.subscribe(subject, opts);

                (async () => {
                    for await (const m of sub) {
                        onMessage(m.data);
                        m.ack();
                    }
                })().catch((err) => {
                    // Error details logged for debugging
                });

                // Return unsubscribe
                return () => {
                    sub.unsubscribe();
                };
            } catch (error) {
                // Error logged for debugging
                if (process.env.NODE_ENV === 'development') {
                    console.error('Error:', error);
                }
            }
        },
        [connectIfNeeded, js]
    );

    return (
        <ErrorBoundary 
            name="NATSProvider" 
            fallback={
                <NATSContext.Provider value={{
                    isConnected: false,
                    connectIfNeeded: async () => {},
                    publish: async () => {},
                    subscribe: async () => {}
                }}>
                    {children}
                </NATSContext.Provider>
            }
        >
            <NATSContext.Provider
                value={{
                    isConnected,
                    connectIfNeeded,
                    publish,
                    subscribe
                }}
            >
                {children}
            </NATSContext.Provider>
        </ErrorBoundary>
    );
}
