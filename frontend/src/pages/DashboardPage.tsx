import { FormEvent, useEffect, useMemo, useState } from 'react';
import { getPostLikes, listPosts, createPost, deletePost, toggleLike, updatePost } from '../lib/api';
import type { LikesResponse, Post, PostDraft } from '../types';
import { useAuth } from '../context/AuthContext';
import { PageFrame } from '../components/PageFrame';
import { StatCard } from '../components/StatCard';
import { PostCard } from '../components/PostCard';
import { Modal } from '../components/Modal';

type DraftState = PostDraft & { id?: number };

export function DashboardPage() {
  const { sessionUser } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [likedMap, setLikedMap] = useState<Record<number, boolean>>({});
  const [draft, setDraft] = useState<DraftState>({ title: '', content: '' });
  const [likesModal, setLikesModal] = useState<{ postId: number; data: LikesResponse | null }>({ postId: 0, data: null });
  const [busyPostId, setBusyPostId] = useState<number | null>(null);

  async function loadPosts() {
    setLoading(true);
    setError('');
    try {
      const response = await listPosts();
      setPosts(response.posts ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load posts');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadPosts();
  }, []);

  const totalPosts = posts.length;
  const ownerPosts = useMemo(() => posts.filter((post) => post.user_id === sessionUser.userId).length, [posts, sessionUser.userId]);

  async function submitDraft(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusyPostId(draft.id ?? -1);
    try {
      if (draft.id) {
        await updatePost(draft.id, { title: draft.title, content: draft.content });
      } else {
        await createPost({ title: draft.title, content: draft.content });
      }
      setDraft({ title: '', content: '' });
      await loadPosts();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not save post');
    } finally {
      setBusyPostId(null);
    }
  }

  function startEdit(post: Post) {
    setDraft({ id: post.id, title: post.title, content: post.content });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  async function removePost(postId: number) {
    if (!window.confirm('Delete this post?')) {
      return;
    }

    setBusyPostId(postId);
    try {
      await deletePost(postId);
      await loadPosts();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not delete post');
    } finally {
      setBusyPostId(null);
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

  async function openLikes(postId: number) {
    setLikesModal({ postId, data: null });
    try {
      const data = await getPostLikes(postId, '', 10);
      setLikesModal({ postId, data });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load likes');
    }
  }

  return (
    <PageFrame
      eyebrow="Dashboard"
      title="Posts from the gateway"
      subtitle="This feed is wired to the current list-posts endpoint and the existing social actions."
      aside={
        <div className="stack">
          <StatCard label="Session user" value={sessionUser.userId ? `#${sessionUser.userId}` : 'anonymous'} hint={sessionUser.role ?? 'role unavailable'} />
          <StatCard label="Posts loaded" value={String(totalPosts)} hint="Current gateway response" />
          <StatCard label="Your posts" value={String(ownerPosts)} hint="Matched by access token user_id" />
        </div>
      }
    >
      <div className="two-column-layout">
        <section className="panel">
          <div className="panel-header">
            <div>
              <h2>{draft.id ? `Edit post #${draft.id}` : 'Create post'}</h2>
              <p>Draft locally, then push to the post service through the gateway.</p>
            </div>
          </div>

          <form className="compose-form" onSubmit={submitDraft}>
            <label>
              <span>Title</span>
              <input
                value={draft.title}
                onChange={(event) => setDraft((current) => ({ ...current, title: event.target.value }))}
                minLength={3}
                maxLength={30}
                required
              />
            </label>
            <label>
              <span>Content</span>
              <textarea
                value={draft.content}
                onChange={(event) => setDraft((current) => ({ ...current, content: event.target.value }))}
                minLength={5}
                rows={5}
                required
              />
            </label>
            <div className="form-actions">
              <button type="submit" className="button" disabled={busyPostId !== null}>
                {busyPostId !== null ? 'Saving...' : draft.id ? 'Update post' : 'Publish post'}
              </button>
              {draft.id ? (
                <button type="button" className="button button-soft" onClick={() => setDraft({ title: '', content: '' })}>
                  Cancel edit
                </button>
              ) : null}
            </div>
          </form>
        </section>

        <section className="panel">
          <div className="panel-header">
            <div>
              <h2>Feed</h2>
              <p>Posts are rendered from the current API response without assuming a future feed service.</p>
            </div>
            <button type="button" className="button button-soft" onClick={() => void loadPosts()}>
              Refresh
            </button>
          </div>

          {error ? <div className="form-error">{error}</div> : null}

          {loading ? (
            <div className="empty-state">Loading posts...</div>
          ) : posts.length === 0 ? (
            <div className="empty-state">No posts yet. Create the first one.</div>
          ) : (
            <div className="post-list">
              {posts.map((post) => (
                <PostCard
                  key={post.id}
                  post={post}
                  isOwner={post.user_id === sessionUser.userId}
                  liked={Boolean(likedMap[post.id])}
                  onToggleLike={() => void handleLike(post.id)}
                  onEdit={post.user_id === sessionUser.userId ? () => startEdit(post) : undefined}
                  onDelete={post.user_id === sessionUser.userId ? () => void removePost(post.id) : undefined}
                  onViewLikes={() => void openLikes(post.id)}
                />
              ))}
            </div>
          )}
        </section>
      </div>

      {likesModal.postId ? (
        <Modal title={`Likes for post #${likesModal.postId}`} onClose={() => setLikesModal({ postId: 0, data: null })}>
          {likesModal.data ? (
            <div className="modal-body">
              {likesModal.data.users?.length ? (
                <ul className="compact-list">
                  {likesModal.data.users.map((user, index) => (
                    <li key={`${user.username}-${index}`}>
                      <strong>{user.name}</strong>
                      <span>@{user.username}</span>
                    </li>
                  ))}
                </ul>
              ) : (
                <div className="empty-state">No likes returned for this post.</div>
              )}
            </div>
          ) : (
            <div className="empty-state">Loading likes...</div>
          )}
        </Modal>
      ) : null}
    </PageFrame>
  );
}
