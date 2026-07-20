import { readAccessToken } from './session';

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') || 'http://localhost:8082';

/**
 * Connect to the SSE notification stream using fetch() so we can send the
 * Authorization header (EventSource does not support custom headers).
 *
 * Returns an AbortController – call controller.abort() to close the connection.
 */
export function connectNotificationStream(
  onEvent: (data: string) => void,
  onError?: (err: unknown) => void,
): AbortController {
  const controller = new AbortController();

  (async () => {
    try {
      const token = readAccessToken();
      const headers: Record<string, string> = {
        Accept: 'text/event-stream',
        'X-Client-Type': 'web',
      };
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(`${API_BASE_URL}/api/notification/stream`, {
        headers,
        credentials: 'include',
        signal: controller.signal,
      });

      if (!response.ok || !response.body) {
        throw new Error(`SSE connect failed: ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });

        // Parse SSE frames from the buffer
        const lines = buffer.split('\n');
        // Keep the last (potentially incomplete) line in the buffer
        buffer = lines.pop() ?? '';

        let currentData = '';
        for (const line of lines) {
          if (line.startsWith('data:')) {
            currentData = line.slice(5).trim();
          } else if (line === '' && currentData) {
            // Empty line = end of SSE event → dispatch
            onEvent(currentData);
            currentData = '';
          }
        }
      }
    } catch (err) {
      // AbortError is expected when we call controller.abort()
      if (err instanceof DOMException && err.name === 'AbortError') return;
      onError?.(err);
    }
  })();

  return controller;
}
