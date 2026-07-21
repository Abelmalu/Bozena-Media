import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import type { FeedItem, ProfileUser } from '../types';
import { getPostLikes, resolveImageUrl } from '../lib/api';
import { useAuth } from '../context/AuthContext';

export function FeedCard({
  item,
  ownerId,
  liked,
  onLike,
  onOpenLikes,
  onLikeStatusFetched,
}: {
  item: FeedItem;
  ownerId: number;
  liked: boolean;
  onLike: () => void;
  onOpenLikes: () => void;
  onLikeStatusFetched?: (postId: number, liked: boolean) => void;
}) {
  const canOpenProfile = ownerId > 0;
  const [likesCount, setLikesCount] = useState<number | null>(null);
  const postId = item.PostID ?? item.post_id ?? item.postId ?? item.id ?? 0;
  const imageUrl = resolveImageUrl(item.Image ?? item.image);
  const { username: currentUsername } = useAuth();

  useEffect(() => {
    let active = true;
    if (postId > 0) {
      getPostLikes(postId, '', 100)
        .then((res) => {
          if (active) {
            setLikesCount(res.users?.length ?? 0);
            if (currentUsername && onLikeStatusFetched) {
              const isLikedByMe = res.users?.some((u) => u.username === currentUsername) ?? false;
              onLikeStatusFetched(postId, isLikedByMe);
            }
          }
        })
        .catch(() => {});
    }
    return () => {
      active = false;
    };
  }, [postId, liked, currentUsername, onLikeStatusFetched]);

  return (
    <article className="feed-card">
      <div className="feed-card-head">
        <div>
          <div className="feed-meta">@{item.UserName}</div>
          <h3>{item.PostTitle}</h3>
        </div>
        <div className="feed-card-head-right">
          <div className="feed-like-pill" aria-label={likesCount !== null ? `${likesCount} likes` : 'Likes'}>
            <span className="feed-like-heart">♥</span>
            <span>{likesCount !== null ? likesCount : 0}</span>
          </div>
          {canOpenProfile ? (
            <Link className="pill-link" to={`/app/profile/${ownerId}`}>
              View profile
            </Link>
          ) : (
            <span className="pill-link pill-link-disabled">View profile</span>
          )}
        </div>
      </div>

      <p className="feed-content">{item.PostContent}</p>

      {imageUrl ? (
        <div className="feed-image-wrap">
          <img className="feed-image" src={imageUrl} alt={item.PostTitle ? `Post image for ${item.PostTitle}` : 'Post image'} loading="lazy" />
        </div>
      ) : null}

      <div className="feed-byline">{item.Name}</div>

      <div className="feed-actions">
        <button type="button" className={`button button-soft ${liked ? 'button-soft-active' : ''}`} onClick={onLike}>
          {liked ? 'Unlike' : 'Like'}
        </button>
        <button type="button" className="button button-soft" onClick={onOpenLikes}>
          Likes
        </button>
      </div>
    </article>
  );
}

export function UserChip({ user }: { user: ProfileUser }) {
  return (
    <div className="user-chip">
      <div className="user-chip-name">{user.name}</div>
      <div className="user-chip-handle">@{user.username}</div>
    </div>
  );
}
