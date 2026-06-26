import { useEffect, useMemo, useState } from 'react';
import { getFeed, getPostLikes, toggleLike } from '../lib/api';
import type { FeedItem, LikesResponse } from '../types';
import { PageFrame } from '../components/PageFrame';
import { StatCard } from '../components/StatCard';
import { FeedCard } from '../components/Card';
import { Modal } from '../components/Modal';
import { useAuth } from '../context/AuthContext';
import { FEED_REFRESH_EVENT } from '../lib/events';

export function FeedPage() {
  const { sessionUser, username } = useAuth();
  const [items, setItems] = useState<FeedItem[]>([]);
  const [cursor, setCursor] = useState('');
  const [limit, setLimit] = useState<number | null>(null);
  const [hasNext, setHasNext] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [likedMap, setLikedMap] = useState<Record<number, boolean>>({});
  const [likesModal, setLikesModal] = useState<{ postId: number; data: LikesResponse | null }>({ postId: 0, data: null });

  async function loadFeed(nextCursor = '', append = false) {
    setLoading(true);
    setError('');
    try {
      const response = await getFeed(nextCursor, 10);
      setItems((current) => (append ? [...current, ...(response.userfeeds ?? [])] : response.userfeeds ?? []));
      setCursor(response.cursor ?? '');
      setLimit(typeof response.limit === 'number' ? response.limit : null);
      setHasNext(Boolean(response.cursor));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load feed');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadFeed();
    const handler = () => {
      void loadFeed();
    };

    window.addEventListener(FEED_REFRESH_EVENT, handler);
    return () => window.removeEventListener(FEED_REFRESH_EVENT, handler);
  }, []);

  const totalReturned = items.length;

  function resolvePostId(item: FeedItem) {
    return item.PostID ?? item.post_id ?? item.postId ?? item.id ?? 0;
  }

  function resolveOwnerId(item: FeedItem) {
    return item.PostOwnerID ?? item.post_owner_id ?? item.postOwnerId ?? item.user_id ?? item.userId ?? 0;
  }

  const ownerPosts = useMemo(
    () =>
      items.filter((item) => {
        const ownerId = resolveOwnerId(item) || null;
        return ownerId !== null && ownerId === sessionUser.userId;
      }).length,
    [items, sessionUser.userId],
  );

  async function openLikes(postId: number) {
    setLikesModal({ postId, data: null });
    try {
      const response = await getPostLikes(postId, '', 10);
      setLikesModal({ postId, data: response });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load likes');
    }
  }

  async function handleLike(postId: number) {
    const nextState = !likedMap[postId];
    setLikedMap((current) => ({ ...current, [postId]: nextState }));
    try {
      await toggleLike(postId, nextState);
    } catch (err) {
      setLikedMap((current) => ({ ...current, [postId]: !nextState }));
      setError(err instanceof Error ? err.message : 'Could not toggle like');
    }
  }

  return (
    <PageFrame
      eyebrow="Feed"
      title="Your timeline"
      subtitle="This loads from `/api/feed/` after login and stays in sync with the backend cursor."
      aside={
        <div className="stack">
          <StatCard label="Returned items" value={String(totalReturned)} hint={`Limit ${limit ?? 10}`} />
          <StatCard label="Your posts in feed" value={String(ownerPosts)} hint="Matched from access token user ID" />
          <StatCard label="Logged in user" value={username ?? 'unknown'} hint={sessionUser.userId ? `ID #${sessionUser.userId}` : sessionUser.role ?? 'role unavailable'} />
        </div>
      }
    >
      <section className="panel">
        <div className="panel-header">
          <div>
            <h2>Feed items</h2>
            <p>Each card is rendered from the feed service response.</p>
          </div>
          <button type="button" className="button button-soft" onClick={() => void loadFeed()}>
            Refresh
          </button>
        </div>

        {error ? <div className="form-error">{error}</div> : null}

        {loading ? (
          <div className="empty-state">Loading feed...</div>
        ) : items.length === 0 ? (
          <div className="empty-state">No feed items returned yet.</div>
        ) : (
          <div className="feed-list">
            {items.map((item) => (
              <FeedCard
                key={resolvePostId(item) || `${item.UserName}-${item.PostTitle}`}
                item={item}
                postId={resolvePostId(item)}
                ownerId={resolveOwnerId(item)}
                liked={Boolean(likedMap[resolvePostId(item)])}
                onLike={() => {
                  const postId = resolvePostId(item);
                  if (!postId) {
                    setError('Feed item is missing a post ID');
                    return;
                  }
                  void handleLike(postId);
                }}
                onOpenLikes={() => {
                  const postId = resolvePostId(item);
                  if (!postId) {
                    setError('Feed item is missing a post ID');
                    return;
                  }
                  void openLikes(postId);
                }}
              />
            ))}
          </div>
        )}

        {hasNext ? (
          <div className="load-more-row">
            <button type="button" className="button button-soft" onClick={() => void loadFeed(cursor, true)}>
              Load more
            </button>
          </div>
        ) : null}
      </section>

      {likesModal.postId ? (
        <Modal title={`Likes for post #${likesModal.postId}`} onClose={() => setLikesModal({ postId: 0, data: null })}>
          {likesModal.data ? (
            <div className="modal-body">
              {likesModal.data.users?.length ? (
                <div className="compact-grid">
                  {likesModal.data.users.map((user) => (
                    <div key={`${user.username}-${user.name}`} className="user-card">
                      <div className="user-name">{user.name}</div>
                      <div className="user-handle">@{user.username}</div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="empty-state">No users returned.</div>
              )}
            </div>
          ) : (
            <div className="empty-state">Loading...</div>
          )}
        </Modal>
      ) : null}
    </PageFrame>
  );
}
