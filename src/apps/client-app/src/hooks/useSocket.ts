'use client';

import { useEffect, useRef, useState } from 'react';
import { realtimeService } from '@/services/realtime.service';
import { ENV } from '@/constants/env';

export interface SocketMessage {
  type: string;
  data: any;
  flow_id?: string;
}

export function useSocket(flowId?: string, isPublic: boolean = false) {
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);
  const [lastMessage, setLastMessage] = useState<SocketMessage | null>(null);

  useEffect(() => {
    let mounted = true;
    let ticket: string | null = null;

    async function connect() {
      try {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = ENV.GATEWAY_URL.replace(/^https?:\/\//, '');
        let url: string;

        if (isPublic) {
          url = `${protocol}//${host}/v1/realtime/public/topology/ws`;
        } else {
          ticket = await realtimeService.getTicket();
          if (!mounted) return;
          url = `${protocol}//${host}/v1/realtime/stream?ticket=${ticket}`;
        }

        if (flowId) {
          url += (url.includes('?') ? '&' : '?') + `flow_id=${flowId}`;
        }

        const socket = new WebSocket(url);
        socketRef.current = socket;

        socket.onopen = () => {
          if (mounted) setIsConnected(true);
        };

        socket.onclose = () => {
          if (mounted) setIsConnected(false);
          // Simple reconnect logic
          setTimeout(() => {
            if (mounted) connect();
          }, 3000);
        };

        socket.onmessage = (event) => {
          if (!mounted) return;
          try {
            const data = JSON.parse(event.data);
            setLastMessage(data);
          } catch (e) {
            console.error('Failed to parse socket message', e);
          }
        };
      } catch (e) {
        console.error('Failed to connect to socket', e);
        if (mounted) {
          setTimeout(() => {
            if (mounted) connect();
          }, 5000);
        }
      }
    }

    connect();

    return () => {
      mounted = false;
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [flowId]);

  const sendMessage = (msg: SocketMessage) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(msg));
    }
  };

  return { isConnected, lastMessage, sendMessage };
}
