import { FormEvent, useState } from 'react';
import { searchUsers } from '../lib/api';
import type { SearchUser } from '../types';
import { PageFrame } from '../components/PageFrame';

export function SearchPage() {
  const [query, setQuery] = useState('');
  const [activeQuery, setActiveQuery] = useState('');
  const [results, setResults] = useState<SearchUser[]>([]);
  const [cursor, setCursor] = useState('');
  const [hasNext, setHasNext] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function runSearch(nextQuery: string, nextCursor = '', append = false) {
    if (!nextQuery.trim()) {
      setResults([]);
      setActiveQuery('');
      setCursor('');
      setHasNext(false);
      return;
    }

    setLoading(true);
    setError('');
    try {
      const response = await searchUsers(nextQuery.trim(), nextCursor, 10);
      setActiveQuery(nextQuery.trim());
      setResults((current) => (append ? [...current, ...(response.users ?? [])] : response.users ?? []));
      setCursor(response.cursor ?? '');
      setHasNext(Boolean(response.has_next));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not search users');
    } finally {
      setLoading(false);
    }
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    await runSearch(query);
  }

  return (
    <PageFrame
      eyebrow="Search"
      title="Find users"
      subtitle="The backend currently returns username and name only, so results are informational for now."
    >
      <section className="panel">
        <form className="search-form" onSubmit={onSubmit}>
          <input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Search username" />
          <button type="submit" className="button" disabled={loading}>
            {loading ? 'Searching...' : 'Search'}
          </button>
        </form>

        {error ? <div className="form-error">{error}</div> : null}

        {results.length === 0 ? (
          <div className="empty-state">No search results yet.</div>
        ) : (
          <>
            <div className="compact-grid">
              {results.map((user) => (
                <div key={`${user.username}-${user.name}`} className="user-card">
                  <div className="user-name">{user.name}</div>
                  <div className="user-handle">@{user.username}</div>
                </div>
              ))}
            </div>

            {hasNext ? (
              <div className="load-more-row">
                <button type="button" className="button button-soft" onClick={() => void runSearch(activeQuery, cursor, true)} disabled={loading}>
                  Load more
                </button>
              </div>
            ) : null}
          </>
        )}
      </section>
    </PageFrame>
  );
}
