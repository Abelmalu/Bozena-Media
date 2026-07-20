import { useEffect, useState } from 'react';
import { getNotifications } from '../lib/api';
import type { Notification } from '../types';
import { PageFrame } from '../components/PageFrame';
import { NOTIFICATION_RECEIVED_EVENT } from '../lib/events';

let _liveId = 0;
function nextLiveId() {
  _liveId -= 1;
  return _liveId;
}

export function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [cursor, setCursor] = useState('');
  const [hasNext, setHasNext] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [liveCount, setLiveCount] = useState(0);

  async function loadNotifications(nextCursor = '', append = false) {
    setLoading(true);
    setError('');
    try {
      const response = await getNotifications(nextCursor, 10);
      // Backend returns PascalCase field names (Go struct without json tags)
      const fetched = response.UserNotifications ?? [];
      setNotifications((current) => (append ? [...current, ...fetched] : fetched));
      setCursor(response.Cursor ?? '');
      setHasNext(Boolean(response.HasNext));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load notifications');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    // 1. Load historical notifications from the REST API
    void loadNotifications();

    // 2. Listen to the global notification event dispatched from AppShell
    const handleNotification = (e: Event) => {
      const data = (e as CustomEvent).detail as string;
      const liveNotif: Notification = {
        id: nextLiveId(),
        username: '',
        actor_id: 0,
        message: data,
        created_at: new Date().toISOString(),
        is_read: false,
      };
      setNotifications((current) => [liveNotif, ...current]);
      setLiveCount((c) => c + 1);
    };

    window.addEventListener(NOTIFICATION_RECEIVED_EVENT, handleNotification);
    return () => window.removeEventListener(NOTIFICATION_RECEIVED_EVENT, handleNotification);
  }, []);

  return (
    <PageFrame
      eyebrow="Updates"
      title="Notifications"
      subtitle="Historical notifications loaded on mount; new ones arrive in real-time via SSE."
      aside={
        liveCount > 0 ? (
          <div
            style={{
              background: '#0ea5e9',
              color: '#fff',
              borderRadius: '8px',
              padding: '12px 16px',
              fontSize: '0.875rem',
              fontWeight: 600,
            }}
          >
            🔔 {liveCount} new since you opened this page
          </div>
        ) : undefined
      }
    >
      <section className="panel">
        <div className="panel-header">
          <div>
            <h2>Recent Activity</h2>
            <p>New follower notifications appear instantly. Older ones are loaded from the server.</p>
          </div>
          <button
            type="button"
            className="button button-soft"
            onClick={() => {
              setLiveCount(0);
              void loadNotifications();
            }}
          >
            Refresh
          </button>
        </div>

        {error ? <div className="form-error">{error}</div> : null}

        {loading && notifications.length === 0 ? (
          <div className="empty-state">Loading notifications...</div>
        ) : notifications.length === 0 ? (
          <div className="empty-state">No notifications yet — follow someone and wait for them to follow back!</div>
        ) : (
          <div className="feed-list stack">
            {notifications.map((notif, i) => (
              <div
                key={notif.id ?? i}
                className="feed-card"
                style={{
                  padding: '16px',
                  borderLeft: notif.is_read ? 'none' : '4px solid #0ea5e9',
                  transition: 'border-color 0.2s',
                }}
              >
                <p style={{ margin: 0 }}>{notif.message}</p>
                {notif.username && (
                  <small style={{ color: 'var(--text-secondary)', display: 'block', marginTop: '4px' }}>
                    from @{notif.username}
                  </small>
                )}
                {notif.created_at && (
                  <small style={{ color: 'var(--text-secondary)', display: 'block', marginTop: '4px' }}>
                    {new Date(notif.created_at).toLocaleString()}
                  </small>
                )}
              </div>
            ))}
          </div>
        )}

        {hasNext ? (
          <div className="load-more-row">
            <button type="button" className="button button-soft" onClick={() => void loadNotifications(cursor, true)}>
              Load more
            </button>
          </div>
        ) : null}
      </section>
    </PageFrame>
  );
}
