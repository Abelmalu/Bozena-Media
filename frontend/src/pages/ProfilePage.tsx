import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getUserFollowers, getUserFollowings, getUserPosts } from '../lib/api';
import type { FollowersResponse, FollowingsResponse, ProfileUser, UserPostsResponse } from '../types';
import { useAuth } from '../context/AuthContext';
import { PageFrame } from '../components/PageFrame';
import { StatCard } from '../components/StatCard';

type Section = 'followers' | 'followings';

export function ProfilePage() {
  const params = useParams();
  const { sessionUser, username } = useAuth();
  const userId = params.id ? Number(params.id) : sessionUser.userId ?? NaN;
  const isOwnProfile = !params.id || (sessionUser.userId !== null && userId === sessionUser.userId);
  const [followers, setFollowers] = useState<FollowersResponse | null>(null);
  const [followings, setFollowings] = useState<FollowingsResponse | null>(null);
  const [posts, setPosts] = useState<UserPostsResponse | null>(null);
  const [section, setSection] = useState<Section>('followers');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [followerCursor, setFollowerCursor] = useState('');
  const [followingCursor, setFollowingCursor] = useState('');
  const [postsCursor, setPostsCursor] = useState('');

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

  return (
    <PageFrame
      eyebrow="Profile"
      title={isOwnProfile && username ? `${username}'s profile` : userId ? `User #${userId}` : 'Profile'}
      subtitle="Counts below use exactly what the backend returned, including the current page limit."
      aside={
        <div className="stack">
          <StatCard label="Followers returned" value={String(followersCount)} hint={`Limit ${followers?.limit ?? 10}`} />
          <StatCard label="Following returned" value={String(followingsCount)} hint={`Limit ${followings?.limit ?? 10}`} />
          <StatCard label="Posts returned" value={String(postsCount)} hint={posts?.cursor ? 'More posts available' : 'No more posts'} />
          <StatCard label="Viewer" value={username ?? 'unknown'} hint={sessionUser.userId ? `ID #${sessionUser.userId}` : sessionUser.role ?? 'role unavailable'} />
        </div>
      }
    >
      <section className="panel">
        <div className="tabs">
          <button type="button" className={section === 'followers' ? 'tab active' : 'tab'} onClick={() => setSection('followers')}>
            Followers
          </button>
          <button type="button" className={section === 'followings' ? 'tab active' : 'tab'} onClick={() => setSection('followings')}>
            Following
          </button>
        </div>

        {error ? <div className="form-error">{error}</div> : null}

        {loading ? (
          <div className="empty-state">Loading profile...</div>
        ) : section === 'followers' ? (
          <>
            <div className="compact-grid">
              {followers?.followers?.length ? (
                followers.followers.map((user: ProfileUser) => (
                  <div key={`${user.username}-${user.name}`} className="user-card">
                    <div className="user-name">{user.name}</div>
                    <div className="user-handle">@{user.username}</div>
                  </div>
                ))
              ) : (
                <div className="empty-state">No followers returned on this page.</div>
              )}
            </div>

            {followers?.cursor ? (
              <div className="load-more-row">
                <button type="button" className="button button-soft" onClick={() => void loadMoreFollowers()}>
                  Load more followers
                </button>
              </div>
            ) : null}
          </>
        ) : (
          <>
            <div className="compact-grid">
              {followings?.Followings?.length ? (
                followings.Followings.map((user: ProfileUser) => (
                  <div key={`${user.username}-${user.name}`} className="user-card">
                    <div className="user-name">{user.name}</div>
                    <div className="user-handle">@{user.username}</div>
                  </div>
                ))
              ) : (
                <div className="empty-state">No followings returned on this page.</div>
              )}
            </div>

            {followings?.cursor ? (
              <div className="load-more-row">
                <button type="button" className="button button-soft" onClick={() => void loadMoreFollowings()}>
                  Load more followings
                </button>
              </div>
            ) : null}
          </>
        )}
      </section>

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
                </div>
                <p className="feed-content">{post.content}</p>
                <div className="feed-byline">{isOwnProfile && username ? `By @${username}` : `Author ID #${post.user_id}`}</div>
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
    </PageFrame>
  );
}
