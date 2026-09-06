import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  getProfileUploadUrl,
  getUserFollowers,
  getUserFollowings,
  getUserPosts,
  getUserProfile,
  toggleFollow,
  updatePost,
  deletePost,
  uploadFileToMinio,
  resolveImageUrl,
} from '../lib/api';
import type { FollowersResponse, FollowingsResponse, ProfileUser, UserPostsResponse, UserPost, UserProfile } from '../types';
import { useAuth } from '../context/AuthContext';
import { PageFrame } from '../components/PageFrame';
import { StatCard } from '../components/StatCard';
import { Modal } from '../components/Modal';
import { FEED_REFRESH_EVENT } from '../lib/events';

export function ProfilePage() {
  const params = useParams();
  const navigate = useNavigate();
  const { sessionUser, username, avatarUrl, updateAvatarUrl } = useAuth();
  const userId = params.id ? Number(params.id) : sessionUser.userId ?? NaN;
  const isOwnProfile = !params.id || (sessionUser.userId !== null && userId === sessionUser.userId);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [followers, setFollowers] = useState<FollowersResponse | null>(null);
  const [followings, setFollowings] = useState<FollowingsResponse | null>(null);
  const [posts, setPosts] = useState<UserPostsResponse | null>(null);
  const [modalOpen, setModalOpen] = useState<'none' | 'followers' | 'followings'>('none');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [followerCursor, setFollowerCursor] = useState('');
  const [followingCursor, setFollowingCursor] = useState('');
  const [postsCursor, setPostsCursor] = useState('');
  const [isFollowing, setIsFollowing] = useState(false);
  const [editingPost, setEditingPost] = useState<UserPost | null>(null);
  const [editTitle, setEditTitle] = useState('');
  const [editContent, setEditContent] = useState('');
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState('');
  const [avatarUploadStep, setAvatarUploadStep] = useState('');
  const [avatarUploading, setAvatarUploading] = useState(false);
  const [avatarError, setAvatarError] = useState('');
  useEffect(() => {
    async function loadProfile() {
      setLoading(true);
      setError('');
      try {
        const [followersResponse, followingsResponse, postsResponse] = await Promise.all([
          getUserFollowers(userId, '', 10),
          getUserFollowings(userId, '', 10),
          getUserPosts(userId, '', 10),
        ]);
        setFollowers(followersResponse);
        setFollowings(followingsResponse);
        setPosts(postsResponse);
        setFollowerCursor(followersResponse.cursor ?? '');
        setFollowingCursor(followingsResponse.cursor ?? '');
        setPostsCursor(postsResponse.cursor ?? '');
        setIsFollowing(followersResponse.followers?.some((f) => f.username === username) ?? false);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Could not load profile');
      } finally {
        setLoading(false);
      }
    }

    async function loadProfileUser() {
      try {
        setProfile(await getUserProfile(userId));
      } catch {
        setProfile(null);
      }
    }

    if (!Number.isNaN(userId)) {
      void loadProfile();
      if (!isOwnProfile) {
        void loadProfileUser();
      } else {
        setProfile(null);
      }
    }
  }, [userId, username, isOwnProfile]);

  async function handleToggleFollow() {
    try {
      const nextState = !isFollowing;
      await toggleFollow(userId, nextState);
      setIsFollowing(nextState);
      // Reload followers to reflect changes
      const updatedFollowers = await getUserFollowers(userId, '', 10);
      setFollowers(updatedFollowers);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not toggle follow');
    }
  }

  async function loadMoreFollowers() {
    try {
      const response = await getUserFollowers(userId, followerCursor, 10);
      setFollowers((current) => ({
        ...response,
        followers: [...(current?.followers ?? []), ...(response.followers ?? [])],
      }));
      setFollowerCursor(response.cursor ?? '');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load more followers');
    }
  }

  async function loadMoreFollowings() {
    try {
      const response = await getUserFollowings(userId, followingCursor, 10);
      setFollowings((current) => ({
        ...response,
        Followings: [...(current?.Followings ?? []), ...(response.Followings ?? [])],
      }));
      setFollowingCursor(response.cursor ?? '');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load more followings');
    }
  }

  async function loadMorePosts() {
    try {
      const response = await getUserPosts(userId, postsCursor, 10);
      setPosts((current) => ({
        ...response,
        posts: [...(current?.posts ?? []), ...(response.posts ?? [])],
      }));
      setPostsCursor(response.cursor ?? '');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load more posts');
    }
  }

  const followersCount = followers?.followers?.length ?? 0;
  const followingsCount = followings?.Followings?.length ?? 0;
  const postsCount = posts?.posts?.length ?? 0;

  function openEditModal(post: UserPost) {
    setEditingPost(post);
    setEditTitle(post.title);
    setEditContent(post.content);
  }

  async function handleUpdatePost(e: React.FormEvent) {
    e.preventDefault();
    if (!editingPost) return;
    try {
      await updatePost(editingPost.id, { title: editTitle, content: editContent });
      setPosts((current) => {
        if (!current) return current;
        return {
          ...current,
          posts: current.posts.map((p) => (p.id === editingPost.id ? { ...p, title: editTitle, content: editContent } : p)),
        };
      });
      setEditingPost(null);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Could not update post');
    }
  }

  async function handleDeletePost(postId: number) {
    if (!confirm('Are you sure you want to delete this post?')) return;
    try {
      await deletePost(postId);
      setPosts((current) => {
        if (!current) return current;
        return {
          ...current,
          posts: current.posts.filter((p) => p.id !== postId),
        };
      });
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Could not delete post');
    }
  }

  function onAvatarFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0] ?? null;
    setAvatarFile(file);
    setAvatarError('');
    if (file) {
      setAvatarPreview(URL.createObjectURL(file));
    } else {
      setAvatarPreview('');
    }
  }

  function clearAvatarSelection() {
    setAvatarFile(null);
    setAvatarPreview('');
  }

  async function handleAvatarUpload() {
    if (!avatarFile) return;
    setAvatarUploading(true);
    setAvatarUploadStep('');
    setAvatarError('');

    try {
      setAvatarUploadStep('Getting upload URL...');
      const presigned = await getProfileUploadUrl(avatarFile.name, avatarFile.type);
      setAvatarUploadStep('Uploading avatar...');
      await uploadFileToMinio(avatarFile, presigned);
      updateAvatarUrl(avatarPreview);
      setAvatarFile(null);
      setAvatarPreview('');
      window.dispatchEvent(new Event(FEED_REFRESH_EVENT));
      setAvatarUploadStep('Avatar updated');
    } catch (err) {
      setAvatarError(err instanceof Error ? err.message : 'Could not upload avatar');
    } finally {
      setAvatarUploading(false);
    }
  }

  return (
    <PageFrame
      eyebrow="Profile"
      title={
        isOwnProfile && username
          ? `${username}'s profile`
          : profile?.username
            ? `@${profile.username}`
            : userId
              ? `User #${userId}`
              : 'Profile'
      }
      subtitle="Counts below use exactly what the backend returned, including the current page limit."
      aside={
        <div className="stack">
          {avatarUrl ? (
            <div className="profile-avatar-card">
              <img className="profile-avatar" src={resolveImageUrl(avatarUrl)} alt={`${username ?? 'User'} avatar`} />
            </div>
          ) : isOwnProfile ? (
            <div className="profile-avatar-empty">
              <div className="profile-avatar-fallback profile-avatar-empty-mark">?</div>
              <div className="profile-avatar-empty-text">No profile picture yet</div>
            </div>
          ) : null}
          {!avatarUrl && isOwnProfile ? (
            <div className="profile-avatar-upload">
              {avatarPreview ? (
                <div className="profile-avatar-preview-wrap">
                  <img className="profile-avatar-preview" src={avatarPreview} alt="Avatar preview" />
                  <button type="button" className="profile-avatar-clear" onClick={clearAvatarSelection}>
                    Remove
                  </button>
                </div>
              ) : (
                <label htmlFor="profile-avatar-input" className="profile-avatar-upload-dropzone">
                  <span className="profile-avatar-upload-icon">+</span>
                  <span>Upload profile photo</span>
                </label>
              )}
              <input id="profile-avatar-input" type="file" accept="image/*" hidden onChange={onAvatarFileChange} />
              {avatarFile ? (
                <div className="profile-avatar-upload-actions">
                  <div className="profile-avatar-upload-name">{avatarFile.name}</div>
                  <button type="button" className="button button-soft" onClick={() => void handleAvatarUpload()} disabled={avatarUploading}>
                    {avatarUploading ? (avatarUploadStep || 'Uploading...') : 'Save photo'}
                  </button>
                </div>
              ) : null}
              {avatarError ? <div className="form-error">{avatarError}</div> : null}
              {avatarUploadStep ? <div className="profile-avatar-step">{avatarUploadStep}</div> : null}
            </div>
          ) : null}
          {!isOwnProfile && (
            <>
              <button
                type="button"
                className={`button ${isFollowing ? 'button-soft' : ''}`}
                onClick={() => void handleToggleFollow()}
                style={{ width: '100%' }}
              >
                {isFollowing ? 'Unfollow' : 'Follow'}
              </button>
              <button
                type="button"
                className="button button-soft profile-message-button"
                onClick={() =>
                  navigate(`/app/chat/${userId}`, {
                    state: { username: profile?.username },
                  })
                }
                style={{ width: '100%' }}
              >
                <svg
                  className="profile-message-icon"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  aria-hidden="true"
                >
                  <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
                </svg>
                Message
              </button>
            </>
          )}
          <StatCard label="Followers" value={String(followersCount)} hint="Click to view" onClick={() => setModalOpen('followers')} />
          <StatCard label="Following" value={String(followingsCount)} hint="Click to view" onClick={() => setModalOpen('followings')} />
          <StatCard label="Posts returned" value={String(postsCount)} hint={posts?.cursor ? 'More posts available' : 'No more posts'} />
          <StatCard label="Viewer" value={username ?? 'unknown'} hint={sessionUser.userId ? `ID #${sessionUser.userId}` : sessionUser.role ?? 'role unavailable'} />
        </div>
      }
    >
      {error ? <div className="form-error">{error}</div> : null}

      {modalOpen === 'followers' && (
        <Modal title="Followers" onClose={() => setModalOpen('none')}>
          <div className="modal-body" style={{ maxHeight: '400px', overflowY: 'auto', padding: '16px' }}>
            <div className="compact-grid">
              {followers?.followers?.length ? (
                followers.followers.map((user: ProfileUser) => (
                  <div key={`${user.username}-${user.name}`} className="user-card">
                    <div className="user-name">{user.name}</div>
                    <div className="user-handle">@{user.username}</div>
                  </div>
                ))
              ) : (
                <div className="empty-state">0 followers</div>
              )}
            </div>
            {followers?.cursor && (
              <div className="load-more-row" style={{ marginTop: '16px' }}>
                <button type="button" className="button button-soft" onClick={() => void loadMoreFollowers()}>
                  Load more
                </button>
              </div>
            )}
          </div>
        </Modal>
      )}

      {modalOpen === 'followings' && (
        <Modal title="Following" onClose={() => setModalOpen('none')}>
          <div className="modal-body" style={{ maxHeight: '400px', overflowY: 'auto', padding: '16px' }}>
            <div className="compact-grid">
              {followings?.Followings?.length ? (
                followings.Followings.map((user: ProfileUser) => (
                  <div key={`${user.username}-${user.name}`} className="user-card">
                    <div className="user-name">{user.name}</div>
                    <div className="user-handle">@{user.username}</div>
                  </div>
                ))
              ) : (
                <div className="empty-state">0 following</div>
              )}
            </div>
            {followings?.cursor && (
              <div className="load-more-row" style={{ marginTop: '16px' }}>
                <button type="button" className="button button-soft" onClick={() => void loadMoreFollowings()}>
                  Load more
                </button>
              </div>
            )}
          </div>
        </Modal>
      )}

      <section className="panel">
        <div className="panel-header">
          <div>
            <h2>{isOwnProfile ? 'Your posts' : 'User posts'}</h2>
            <p>Posts returned by the backend for this profile.</p>
          </div>
          <div className="panel-kicker">{posts?.has_next ? 'More available' : 'End of list'}</div>
        </div>

        {error ? <div className="form-error">{error}</div> : null}

        {loading ? (
          <div className="empty-state">Loading posts...</div>
        ) : posts?.posts?.length ? (
          <div className="feed-list">
            {posts.posts.map((post) => (
              <article key={post.id} className="feed-card profile-post-card">
                <div className="feed-card-head">
                  <div>
                    <div className="feed-meta">Post #{post.id}</div>
                    <h3>{post.title}</h3>
                  </div>
                  {isOwnProfile && (
                    <div style={{ display: 'flex', gap: '8px' }}>
                      <button type="button" className="button button-soft" style={{ padding: '4px 8px', fontSize: '0.875rem' }} onClick={() => openEditModal(post)}>Edit</button>
                      <button type="button" className="button" style={{ padding: '4px 8px', fontSize: '0.875rem', backgroundColor: '#ef4444', color: 'white', border: 'none' }} onClick={() => void handleDeletePost(post.id)}>Delete</button>
                    </div>
                  )}
                </div>
                <p className="feed-content">{post.content}</p>
                {post.post_image_url || post.image ? (
                  <div className="feed-image-wrap">
                    <img
                      className="feed-image"
                      src={resolveImageUrl(post.post_image_url ?? post.image)}
                      alt={post.title ? `Post image for ${post.title}` : 'Post image'}
                      loading="lazy"
                    />
                  </div>
                ) : null}
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '12px' }}>
                  <div className="feed-byline" style={{ margin: 0 }}>
                    {isOwnProfile && username ? `By @${username}` : `Author ID #${post.user_id}`}
                  </div>
                  {((post.like_count ?? post.likeCount) ?? 0) > 0 ? (
                    <div style={{ fontSize: '0.875rem', color: '#6b7280', display: 'flex', alignItems: 'center', gap: '4px' }}>
                      ❤️ {post.like_count ?? post.likeCount} {((post.like_count ?? post.likeCount) === 1) ? 'like' : 'likes'}
                    </div>
                  ) : null}
                </div>
              </article>
            ))}
          </div>
        ) : (
          <div className="empty-state">No posts returned on this page.</div>
        )}

        {posts?.cursor ? (
          <div className="load-more-row">
            <button type="button" className="button button-soft" onClick={() => void loadMorePosts()}>
              Load more posts
            </button>
          </div>
        ) : null}
      </section>

      {editingPost && (
        <Modal title="Edit Post" onClose={() => setEditingPost(null)}>
          <form className="modal-body stack" onSubmit={(e) => void handleUpdatePost(e)} style={{ padding: '16px' }}>
            <div className="form-group">
              <label htmlFor="edit-title">Title</label>
              <input id="edit-title" type="text" className="input" value={editTitle} onChange={(e) => setEditTitle(e.target.value)} required />
            </div>
            <div className="form-group">
              <label htmlFor="edit-content">Content</label>
              <textarea id="edit-content" className="input" value={editContent} onChange={(e) => setEditContent(e.target.value)} rows={4} required />
            </div>
            <div className="form-actions" style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px', marginTop: '16px' }}>
              <button type="button" className="button button-soft" onClick={() => setEditingPost(null)}>Cancel</button>
              <button type="submit" className="button">Save Changes</button>
            </div>
          </form>
        </Modal>
      )}
    </PageFrame>
  );
}
