import { Link } from 'react-router-dom';
import type { FeedItem, ProfileUser } from '../types';

export function FeedCard({
  item,
  postId,
  ownerId,
  liked,
  onLike,
  onOpenLikes,
}: {
  item: FeedItem;
  postId: number;
  ownerId: number;
  liked: boolean;
  onLike: () => void;
  onOpenLikes: () => void;
}) {
  const canOpenProfile = ownerId > 0;

  return (
    <article className="feed-card">
      <div className="feed-card-head">
        <div>
          <div className="feed-meta">@{item.UserName}</div>
          <h3>{item.PostTitle}</h3>
        </div>
        {canOpenProfile ? (
          <Link className="pill-link" to={`/app/profile/${ownerId}`}>
            View profile
          </Link>
        ) : (
          <span className="pill-link pill-link-disabled">View profile</span>
        )}
      </div>

      <p className="feed-content">{item.PostContent}</p>

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
