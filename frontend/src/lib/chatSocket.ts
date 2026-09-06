import { readAccessToken } from './session';
import type { ChatMessageEvent } from '../types';

const API_BASE_URL =
  (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') || 'http://localhost:8082';

// ws(s):// equivalent of the API base URL
const WS_BASE_URL = API_BASE_URL.replace(/^http/, 'ws');

export type ChatSocketHandlers = {
  onMessage: (event: ChatMessageEvent) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (err: unknown) => void;
};

export type ChatSocket = {
  send: (receiverId: number, message: string) => boolean;
  close: () => void;
};

/**
 * Open the chat websocket through the API Gateway (`/api/chat/ws`).
 *
 * The browser WebSocket API cannot send an Authorization header, so the
 * access token is passed as a `token` query param, which the gateway's
 * AuthMiddleware accepts as a fallback.
 */
export function connectChatSocket(handlers: ChatSocketHandlers): ChatSocket {
  const token = readAccessToken();
  const url = `${WS_BASE_URL}/api/chat/ws${token ? `?token=${encodeURIComponent(token)}` : ''}`;

  let closedByUser = false;
  const ws = new WebSocket(url);

  ws.onopen = () => handlers.onOpen?.();

  ws.onmessage = (event: MessageEvent) => {
    const raw = typeof event.data === 'string' ? event.data : '';
    if (!raw) return;

    try {
      const parsed = JSON.parse(raw) as Partial<ChatMessageEvent>;
      if (typeof parsed.sender_id === 'number' && typeof parsed.message === 'string') {
        handlers.onMessage({ sender_id: parsed.sender_id, message: parsed.message });
        return;
      }
    } catch {
      // fall through: plain-text frame from older backend versions
    }
    handlers.onMessage({ sender_id: 0, message: raw });
  };

  ws.onerror = (err) => {
    if (!closedByUser) handlers.onError?.(err);
  };

  ws.onclose = () => {
    if (!closedByUser) handlers.onClose?.();
  };

  return {
    send(receiverId: number, message: string) {
      if (ws.readyState !== WebSocket.OPEN) return false;
      ws.send(JSON.stringify({ message, receiver_id: receiverId }));
      return true;
    },
    close() {
      closedByUser = true;
      ws.close();
    },
  };
}
