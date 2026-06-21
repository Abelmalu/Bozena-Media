import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getFollowers, getFollowings, getUserPosts, toggleFollow } from '../lib/api';
import type { FollowUser, Post } from '../types';
import { useAuth } from '../context/AuthContext';
import { PageFrame } from '../components/PageFrame';
import { PostCard } from '../components/PostCard';
import { Modal } from '../components/Modal';

type TabName = 'posts' | 'followers' | 'followings';

export function ProfilePage() {
  const { id } = useParams();
  const userId = Number(id);
  const { sessionUser } = useAuth();
  const [tab, setTab] = useState<TabName>('posts');
  const [posts, setPosts] = useState<Post[]>([]);
  const [followers, setFollowers] = useState<FollowUser[]>([]);
  const [followings, setFollowings] = useState<FollowUser[]>([]);
  const [postCursor, setPostCursor] = useState('');
  const [followerCursor, setFollowerCursor] = useState('');
  const [followingCursor, setFollowingCursor] = useState('');
  const [postHasNext, setPostHasNext] = useState(false);
  const [followerHasNext, setFollowerHasNext] = useState(false);
  const [followingHasNext, setFollowingHasNext] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [likedMap, setLikedMap] = useState<Record<number, boolean>>({});
  const [isFollowing, setIsFollowing] = useState(false);
  const [actionBusy, setActionBusy] = useState(false);
  const [selectedUser, setSelectedUser] = useState<FollowUser | null>(null);

  useEffect(() => {
    async function loadProfile() {
      setLoading(true);
      setError('');
      try {
        const [postsResponse, followersResponse, followingsResponse] = await Promise.all([
          getUserPosts(userId, '', 10),
          getFollowers(userId, '', 10),
          getFollowings(userId, '', 10),
        ]);
        setPosts(postsResponse.posts ?? []);
        setFollowers(followersResponse.followers ?? []);
        setFollowings(followingsResponse.followings ?? []);
        setPostCursor(postsResponse.cursor ?? '');
        setFollowerCursor(followersResponse.cursor ?? '');
        setFollowingCursor(followingsResponse.cursor ?? '');
        setPostHasNext(Boolean(postsResponse.has_next));
        setFollowerHasNext(Boolean(followersResponse.has_next));
        setFollowingHasNext(Boolean(followingsResponse.has_next));
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Could not load profile');
      } finally {
        setLoading(false);
      }
    }

    if (!Number.isNaN(userId)) {
      void loadProfile();
    }
  }, [userId]);

  async function loadMorePosts() {
    try {
      const response = await getUserPosts(userId, postCursor, 10);
      setPosts((current) => [...current, ...(response.posts ?? [])]);
      setPostCursor(response.cursor ?? '');
      setPostHasNext(Boolean(response.has_next));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load more posts');
    }
  }

  async function loadMoreFollowers() {
    try {
      const response = await getFollowers(userId, followerCursor, 10);
      setFollowers((current) => [...current, ...(response.followers ?? [])]);
      setFollowerCursor(response.cursor ?? '');
      setFollowerHasNext(Boolean(response.has_next));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load more followers');
    }
  }

  async function loadMoreFollowings() {
    try {
      const response = await getFollowings(userId, followingCursor, 10);
      setFollowings((current) => [...current, ...(response.followings ?? [])]);
      setFollowingCursor(response.cursor ?? '');
      setFollowingHasNext(Boolean(response.has_next));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not load more followings');
    }
  }

  async function handleFollow() {
    setActionBusy(true);
    try {
      const nextState = !isFollowing;
      await toggleFollow(userId, nextState);
      setIsFollowing(nextState);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not toggle follow');
    } finally {
      setActionBusy(false);
    }
  }

  return (
    <PageFrame
      eyebrow="Profile"
      title={`User #${userId}`}
      subtitle="The backend currently exposes user relationships by numeric ID, so this page is designed around that contract."
      aside={
        <div className="stack">
          <button type="button" className="button" onClick={() => void handleFollow()} disabled={actionBusy || Number.isNaN(userId)}>
            {isFollowing ? 'Unfollow' : 'Follow'}
          </button>
          <div className="profile-meta">
            <span>Viewer</span>
            <strong>{sessionUser.userId ? `#${sessionUser.userId}` : 'anonymous'}</strong>
          </div>
        </div>
      }
    >
      <section className="panel">
        <div className="tabs">
          <button type="button" className={tab === 'posts' ? 'tab active' : 'tab'} onClick={() => setTab('posts')}>
            Posts
          </button>
          <button type="button" className={tab === 'followers' ? 'tab active' : 'tab'} onClick={() => setTab('followers')}>
            Followers
          </button>
          <button type="button" className={tab === 'followings' ? 'tab active' : 'tab'} onClick={() => setTab('followings')}>
            Following
          </button>
        </div>

        {error ? <div className="form-error">{error}</div> : null}

        {loading ? (
          <div className="empty-state">Loading profile...</div>
        ) : tab === 'posts' ? (
          <>
            <div className="post-list">
              {posts.length ? (
                posts.map((post) => (
                  <PostCard
                    key={post.id}
                    post={post}
                    isOwner={post.user_id === sessionUser.userId}
                    liked={Boolean(likedMap[post.id])}
                    onToggleLike={() => setLikedMap((current) => ({ ...current, [post.id]: !current[post.id] }))}
                  />
                ))
              ) : (
                <div className="empty-state">No posts found for this user.</div>
              )}
            </div>

            {postHasNext ? (
              <div className="load-more-row">
                <button type="button" className="button button-soft" onClick={() => void loadMorePosts()}>
                  Load more posts
                </button>
              </div>
            ) : null}
          </>
        ) : tab === 'followers' ? (
          followers.length ? (
            <>
              <div className="compact-grid">
                {followers.map((user) => (
                  <button key={`${user.username}-${user.name}`} type="button" className="user-card" onClick={() => setSelectedUser(user)}>
                    <div className="user-name">{user.name}</div>
                    <div className="user-handle">@{user.username}</div>
                  </button>
                ))}
              </div>

              {followerHasNext ? (
                <div className="load-more-row">
                  <button type="button" className="button button-soft" onClick={() => void loadMoreFollowers()}>
                    Load more followers
                  </button>
                </div>
              ) : null}
            </>
          ) : (
            <div className="empty-state">No followers returned.</div>
          )
        ) : followings.length ? (
          <>
            <div className="compact-grid">
              {followings.map((user) => (
                <button key={`${user.username}-${user.name}`} type="button" className="user-card" onClick={() => setSelectedUser(user)}>
                  <div className="user-name">{user.name}</div>
                  <div className="user-handle">@{user.username}</div>
                </button>
              ))}
            </div>

            {followingHasNext ? (
              <div className="load-more-row">
                <button type="button" className="button button-soft" onClick={() => void loadMoreFollowings()}>
                  Load more followings
                </button>
              </div>
            ) : null}
          </>
        ) : (
          <div className="empty-state">No followings returned.</div>
        )}
      </section>

      {selectedUser ? (
        <Modal title={selectedUser.name} onClose={() => setSelectedUser(null)}>
          <div className="modal-body">
            <div className="profile-meta">
              <span>Username</span>
              <strong>@{selectedUser.username}</strong>
            </div>
            <p className="hero-copy">
              The current gateway response for followers and followings does not include IDs, so these cards stay informational.
            </p>
          </div>
        </Modal>
      ) : null}
    </PageFrame>
  );
}
