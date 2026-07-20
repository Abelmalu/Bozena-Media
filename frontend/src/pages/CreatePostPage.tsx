import { FormEvent, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createPost, getPostUploadUrl, uploadFileToMinio } from '../lib/api';
import { useAuth } from '../context/AuthContext';
import { FEED_REFRESH_EVENT } from '../lib/events';
import { PageFrame } from '../components/PageFrame';
import { StatCard } from '../components/StatCard';

export function CreatePostPage() {
  const navigate = useNavigate();
  const { sessionUser, username } = useAuth();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState('');
  const [posting, setPosting] = useState(false);
  const [uploadStep, setUploadStep] = useState('');
  const [error, setError] = useState('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  function onFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0] ?? null;
    setImageFile(file);
    if (file) {
      setImagePreview(URL.createObjectURL(file));
    } else {
      setImagePreview('');
    }
  }

  function removeImage() {
    setImageFile(null);
    setImagePreview('');
    if (fileInputRef.current) fileInputRef.current.value = '';
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPosting(true);
    setError('');

    try {
      let objectName = '';

      if (imageFile) {
        // Step 1: Get presigned upload URL from Go backend
        setUploadStep('Getting upload URL...');
        const presigned = await getPostUploadUrl(imageFile.name, imageFile.type);

        // Step 2: Upload image bytes directly to MinIO
        setUploadStep('Uploading image...');
        objectName = await uploadFileToMinio(imageFile, presigned);
      }

      // Step 3: Create the post (with optional object_name)
      setUploadStep('Publishing post...');
      await createPost({ title, content, object_name: objectName });
      window.dispatchEvent(new Event(FEED_REFRESH_EVENT));
      navigate('/app/feed', { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not create post');
    } finally {
      setPosting(false);
      setUploadStep('');
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
          {imageFile && <StatCard label="Image" value={imageFile.name} hint={`${(imageFile.size / 1024).toFixed(1)} KB`} />}
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

        <form className="compose-form post-page-form" onSubmit={(e) => void onSubmit(e)}>
          <label>
            <span>Title</span>
            <input value={title} onChange={(event) => setTitle(event.target.value)} minLength={3} maxLength={30} required />
          </label>
          <label>
            <span>Content</span>
            <textarea value={content} onChange={(event) => setContent(event.target.value)} minLength={5} rows={10} required />
          </label>

          {/* Optional image upload */}
          <div className="form-group">
            <span style={{ fontWeight: 600, fontSize: '0.875rem', display: 'block', marginBottom: '8px' }}>
              Photo (optional)
            </span>
            {imagePreview ? (
              <div style={{ position: 'relative', display: 'inline-block' }}>
                <img
                  src={imagePreview}
                  alt="Preview"
                  style={{ maxWidth: '100%', maxHeight: '240px', borderRadius: '8px', display: 'block' }}
                />
                <button
                  type="button"
                  onClick={removeImage}
                  style={{ position: 'absolute', top: 6, right: 6, background: '#ef4444', color: '#fff', border: 'none', borderRadius: '50%', width: 28, height: 28, cursor: 'pointer', fontWeight: 700 }}
                >
                  ✕
                </button>
              </div>
            ) : (
              <label
                htmlFor="post-image-input"
                style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer', padding: '12px 16px', border: '2px dashed var(--border)', borderRadius: '8px' }}
              >
                <span style={{ fontSize: '1.5rem' }}>📷</span>
                <span style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>Click to attach an image</span>
              </label>
            )}
            <input
              id="post-image-input"
              ref={fileInputRef}
              type="file"
              accept="image/*"
              style={{ display: 'none' }}
              onChange={onFileChange}
            />
          </div>

          {error ? <div className="form-error">{error}</div> : null}
          {uploadStep ? <div style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>{uploadStep}</div> : null}

          <div className="form-actions">
            <button type="button" className="button button-soft" onClick={() => navigate('/app/feed')}>
              Cancel
            </button>
            <button type="submit" className="button" disabled={posting}>
              {posting ? (uploadStep || 'Publishing...') : 'Create post'}
            </button>
          </div>
        </form>
      </section>
    </PageFrame>
  );
}
