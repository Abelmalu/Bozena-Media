import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createPost } from '../lib/api';
import { useAuth } from '../context/AuthContext';
import { FEED_REFRESH_EVENT } from '../lib/events';
import { PageFrame } from '../components/PageFrame';
import { StatCard } from '../components/StatCard';

export function CreatePostPage() {
  const navigate = useNavigate();
  const { sessionUser, username } = useAuth();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [posting, setPosting] = useState(false);
  const [error, setError] = useState('');

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPosting(true);
    setError('');

    try {
      await createPost({ title, content });
      window.dispatchEvent(new Event(FEED_REFRESH_EVENT));
      navigate('/app/feed', { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not create post');
    } finally {
      setPosting(false);
    }
  }

  return (
    <PageFrame
      eyebrow="Create Post"
      title="Write something new"
      subtitle="This is a dedicated page now, separate from the feed."
      aside={
        <div className="stack">
          <StatCard label="Author" value={username ?? 'unknown'} hint={sessionUser.userId ? `ID #${sessionUser.userId}` : sessionUser.role ?? 'role unavailable'} />
          <StatCard label="Destination" value="Feed" hint="New posts refresh the feed after submission" />
        </div>
      }
    >
      <section className="panel">
        <div className="panel-header">
          <div>
            <h2>New post</h2>
            <p>Publish a post to the backend and then return to the feed.</p>
          </div>
        </div>

        <form className="compose-form post-page-form" onSubmit={onSubmit}>
          <label>
            <span>Title</span>
            <input value={title} onChange={(event) => setTitle(event.target.value)} minLength={3} maxLength={30} required />
          </label>
          <label>
            <span>Content</span>
            <textarea value={content} onChange={(event) => setContent(event.target.value)} minLength={5} rows={10} required />
          </label>
          {error ? <div className="form-error">{error}</div> : null}
          <div className="form-actions">
            <button type="button" className="button button-soft" onClick={() => navigate('/app/feed')}>
              Cancel
            </button>
            <button type="submit" className="button" disabled={posting}>
              {posting ? 'Publishing...' : 'Create post'}
            </button>
          </div>
        </form>
      </section>
    </PageFrame>
  );
}
