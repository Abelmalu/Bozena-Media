import { Link } from 'react-router-dom';
import type { Post } from '../types';

type Props = {
  post: Post;
  isOwner: boolean;
  liked: boolean;
  onToggleLike: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  onViewLikes?: () => void;
};

export function PostCard({ post, isOwner, liked, onToggleLike, onEdit, onDelete, onViewLikes }: Props) {
  return (
    <article className="post-card">
      <div className="post-card-top">
        <div>
          <div className="post-author">User #{post.user_id}</div>
          <h3>{post.title}</h3>
        </div>
        <Link className="pill-link" to={`/app/users/${post.user_id}`}>
          Open profile
        </Link>
      </div>

      <p className="post-content">{post.content}</p>

      <div className="post-actions">
        <button type="button" className={`button button-soft ${liked ? 'button-soft-active' : ''}`} onClick={onToggleLike}>
          {liked ? 'Unlike' : 'Like'}
        </button>
        {onViewLikes ? (
          <button type="button" className="button button-soft" onClick={onViewLikes}>
            Likes
          </button>
        ) : null}
        {isOwner ? (
          <>
            <button type="button" className="button button-soft" onClick={onEdit}>
              Edit
            </button>
            <button type="button" className="button button-danger" onClick={onDelete}>
              Delete
            </button>
          </>
        ) : null}
      </div>
    </article>
  );
}
