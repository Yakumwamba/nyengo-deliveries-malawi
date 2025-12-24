import { useEffect, useRef, useCallback } from 'react';
import { WS_URL } from '../constants/config';
import { useAuthStore } from '../stores/authStore';

interface WebSocketMessage {
    type: string;
    data: any;
}

export function useWebSocket(onMessage: (message: WebSocketMessage) => void) {
    const wsRef = useRef<WebSocket | null>(null);
    const { token, isAuthenticated } = useAuthStore();

    const connect = useCallback(() => {
        if (!token || !isAuthenticated) return;

        const ws = new WebSocket(`${WS_URL}?token=${token}`);

        ws.onopen = () => {
            console.log('WebSocket connected');
        };

        ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                onMessage(message);
            } catch (error) {
                console.error('Failed to parse WebSocket message:', error);
            }
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
            // Reconnect after 5 seconds
            setTimeout(connect, 5000);
        };

        wsRef.current = ws;
    }, [token, isAuthenticated, onMessage]);

    const disconnect = useCallback(() => {
        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }
    }, []);

    const sendMessage = useCallback((message: any) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify(message));
        }
    }, []);

    useEffect(() => {
        connect();
        return disconnect;
    }, [connect, disconnect]);

    return { sendMessage, disconnect, reconnect: connect };
}
