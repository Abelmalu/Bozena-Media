import { FormEvent, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { getChatMessages, getUserChats, getUserProfile } from '../lib/api';
import { connectChatSocket, ChatSocket } from '../lib/chatSocket';
import { useAuth } from '../context/AuthContext';
import { PageFrame } from '../components/PageFrame';
import type { ChatMessage, ChatParticipant, Conversation } from '../types';

type LocationState = { username?: string } | null;

function otherParticipant(conversation: Conversation, myId: number | null): ChatParticipant | null {
  const others = conversation.participants?.filter((p) => p.userId !== myId) ?? [];
  return others[0] ?? conversation.participants?.[0] ?? null;
}

function formatTime(value?: string) {
  if (!value) return '';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleString();
}

export function ChatPage() {
  const { userId: peerIdParam } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { sessionUser } = useAuth();
  const myId = sessionUser.userId;

  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [chatsCursor, setChatsCursor] = useState('');
  const [chatsHasNext, setChatsHasNext] = useState(false);
  const [chatsLoading, setChatsLoading] = useState(true);
  const [chatsError, setChatsError] = useState('');

  const [activePeer, setActivePeer] = useState<ChatParticipant | null>(null);
  const [activeChatId, setActiveChatId] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [messagesCursor, setMessagesCursor] = useState('');
  const [messagesHasNext, setMessagesHasNext] = useState(false);
  const [messagesLoading, setMessagesLoading] = useState(false);
  const [messagesError, setMessagesError] = useState('');

  const [draft, setDraft] = useState('');
  const [connected, setConnected] = useState(false);
  const [unreadPeers, setUnreadPeers] = useState<Set<number>>(new Set());

  const socketRef = useRef<ChatSocket | null>(null);
  const activePeerRef = useRef<ChatParticipant | null>(null);
  const scrollRef = useRef<HTMLDivElement | null>(null);
  const localIdRef = useRef(0);

  activePeerRef.current = activePeer;

  const loadChats = useCallback(async (cursor = '', append = false) => {
    setChatsLoading(true);
    setChatsError('');
    try {
      const response = await getUserChats(cursor, 20);
      const chats = response.chats ?? [];
      setConversations((current) => (append ? [...current, ...chats] : chats));
      setChatsCursor(response.cursor ?? '');
      setChatsHasNext(Boolean(response.has_next));
    } catch (err) {
      setChatsError(err instanceof Error ? err.message : 'Could not load chats');
    } finally {
      setChatsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadChats();
  }, [loadChats]);

  // Open the websocket through the API gateway once per page mount.
  useEffect(() => {
    if (myId === null) return;

    const socket = connectChatSocket({
      onOpen: () => setConnected(true),
      onClose: () => setConnected(false),
      onError: () => setConnected(false),
      onMessage: (event) => {
        const peer = activePeerRef.current;
        const senderId = event.sender_id || peer?.userId || 0;

        if (peer && senderId === peer.userId) {
          localIdRef.current += 1;
          setMessages((current) => [
            ...current,
            {
              id: `incoming-${localIdRef.current}`,
              chatID: '',
              senderId,
              content: event.message,
              createdAt: new Date().toISOString(),
            },
          ]);
        } else if (senderId) {
          setUnreadPeers((current) => new Set(current).add(senderId));
        }

        // Keep the conversation list preview fresh.
        setConversations((current) =>
          current.map((conversation) => {
            const other = otherParticipant(conversation, myId);
            if (other?.userId !== senderId) return conversation;
            return {
              ...conversation,
              lastMessage: { text: event.message, senderId, createdAt: new Date().toISOString() },
            };
          }),
        );
      },
    });

    socketRef.current = socket;
    return () => {
      socketRef.current = null;
      socket.close();
    };
  }, [myId]);

  const openConversation = useCallback(
    async (conversation: Conversation) => {
      const peer = otherParticipant(conversation, myId);
      if (!peer) return;

      setActivePeer(peer);
      setActiveChatId(conversation.id);
      setMessages([]);
      setMessagesCursor('');
      setMessagesHasNext(false);
      setMessagesError('');
      setUnreadPeers((current) => {
        const next = new Set(current);
        next.delete(peer.userId);
        return next;
      });
      navigate(`/app/chat/${peer.userId}`, { replace: true, state: { username: peer.username } });

      setMessagesLoading(true);
      try {
        const response = await getChatMessages(conversation.id, '', 20);
        const fetched = (response.Messages ?? []).slice().sort((a, b) => a.id.localeCompare(b.id));
        setMessages(fetched);
        setMessagesCursor(response.cursor ?? '');
        setMessagesHasNext(Boolean(response.has_next));
      } catch (err) {
        setMessagesError(err instanceof Error ? err.message : 'Could not load messages');
      } finally {
        setMessagesLoading(false);
      }
    },
    [myId, navigate],
  );

  async function loadOlderMessages() {
    if (!activeChatId || !messagesCursor) return;
    setMessagesLoading(true);
    try {
      const response = await getChatMessages(activeChatId, messagesCursor, 20);
      const older = (response.Messages ?? []).slice().sort((a, b) => a.id.localeCompare(b.id));
      setMessages((current) => [...older, ...current]);
      setMessagesCursor(response.cursor ?? '');
      setMessagesHasNext(Boolean(response.has_next));
    } catch (err) {
      setMessagesError(err instanceof Error ? err.message : 'Could not load older messages');
    } finally {
      setMessagesLoading(false);
    }
  }

  // Deep link: /app/chat/:userId opens (or starts) a chat with that user.
  useEffect(() => {
    if (!peerIdParam || chatsLoading) return;
    const peerId = Number(peerIdParam);
    if (!Number.isFinite(peerId) || peerId === activePeer?.userId) return;

    const existing = conversations.find((c) => otherParticipant(c, myId)?.userId === peerId);
    if (existing) {
      void openConversation(existing);
      return;
    }

    const stateUsername = (location.state as LocationState)?.username;
    setActivePeer({ userId: peerId, username: stateUsername ?? `User #${peerId}` });
    setActiveChatId(null);
    setMessages([]);
    setMessagesCursor('');
    setMessagesHasNext(false);
    setMessagesError('');

    // No username passed via navigation state (e.g. direct URL) — resolve it
    // from the profile endpoint so we never show a bare user ID.
    if (!stateUsername) {
      void getUserProfile(peerId)
        .then((profile) => {
          if (profile?.username) {
            setActivePeer((current) =>
              current && current.userId === peerId ? { ...current, username: profile.username! } : current,
            );
          }
        })
        .catch(() => undefined);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [peerIdParam, chatsLoading]);

  // Keep the thread scrolled to the latest message.
  useEffect(() => {
    const el = scrollRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [messages]);

  function onSend(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const text = draft.trim();
    if (!text || !activePeer || myId === null) return;

    const sent = socketRef.current?.send(activePeer.userId, text);
    if (!sent) {
      setMessagesError('Chat connection is not open. Try reloading the page.');
      return;
    }

    localIdRef.current += 1;
    setMessages((current) => [
      ...current,
      {
        id: `outgoing-${localIdRef.current}`,
        chatID: activeChatId ?? '',
        senderId: myId,
        content: text,
        createdAt: new Date().toISOString(),
      },
    ]);
    setConversations((current) =>
      current.map((conversation) => {
        if (conversation.id !== activeChatId) return conversation;
        return {
          ...conversation,
          lastMessage: { text, senderId: myId, createdAt: new Date().toISOString() },
        };
      }),
    );
    setDraft('');
  }

  const sortedConversations = useMemo(
    () =>
      conversations
        .slice()
        .sort((a, b) => (b.updatedAt ?? '').localeCompare(a.updatedAt ?? '')),
    [conversations],
  );

  return (
    <PageFrame
      eyebrow="Chat"
      title="Messages"
      subtitle="Direct messages delivered in real time over a websocket through the API gateway."
    >
      <section className="panel chat-layout">
        <aside className="chat-sidebar">
          <div className="chat-sidebar-header">
            <span>Conversations</span>
            <span className={connected ? 'chat-status chat-status-online' : 'chat-status'}>
              {connected ? 'live' : 'offline'}
            </span>
          </div>

          {chatsError ? <div className="form-error">{chatsError}</div> : null}

          {chatsLoading && conversations.length === 0 ? (
            <div className="empty-state">Loading chats...</div>
          ) : sortedConversations.length === 0 ? (
            <div className="empty-state">
              No conversations yet. Find someone via Search and hit Message.
            </div>
          ) : (
            <div className="chat-conversation-list">
              {sortedConversations.map((conversation) => {
                const peer = otherParticipant(conversation, myId);
                if (!peer) return null;
                const isActive = peer.userId === activePeer?.userId;
                const unread = unreadPeers.has(peer.userId);
                return (
                  <button
                    key={conversation.id}
                    type="button"
                    className={`chat-conversation ${isActive ? 'active' : ''}`}
                    onClick={() => void openConversation(conversation)}
                  >
                    {peer.avatar ? (
                      <img className="chat-avatar" src={peer.avatar} alt={`${peer.username} avatar`} />
                    ) : (
                      <div className="chat-avatar chat-avatar-fallback">
                        {peer.username?.[0]?.toUpperCase() ?? '?'}
                      </div>
                    )}
                    <div className="chat-conversation-body">
                      <div className="chat-conversation-name">
                        @{peer.username}
                        {unread ? <span className="chat-unread-dot" /> : null}
                      </div>
                      <div className="chat-conversation-preview">
                        {conversation.lastMessage?.senderId === myId ? 'You: ' : ''}
                        {conversation.lastMessage?.text ?? 'No messages yet'}
                      </div>
                    </div>
                  </button>
                );
              })}
            </div>
          )}

          {chatsHasNext ? (
            <div className="load-more-row">
              <button
                type="button"
                className="button button-soft"
                onClick={() => void loadChats(chatsCursor, true)}
                disabled={chatsLoading}
              >
                Load more
              </button>
            </div>
          ) : null}
        </aside>

        <div className="chat-thread">
          {!activePeer ? (
            <div className="empty-state chat-thread-empty">Select a conversation to start chatting.</div>
          ) : (
            <>
              <div className="chat-thread-header">
                <div className="chat-thread-title">@{activePeer.username}</div>
                <div className="chat-thread-subtitle">Direct message</div>
              </div>

              <div className="chat-messages" ref={scrollRef}>
                {messagesHasNext ? (
                  <div className="load-more-row">
                    <button
                      type="button"
                      className="button button-soft"
                      onClick={() => void loadOlderMessages()}
                      disabled={messagesLoading}
                    >
                      Load older
                    </button>
                  </div>
                ) : null}

                {messagesLoading && messages.length === 0 ? (
                  <div className="empty-state">Loading messages...</div>
                ) : messages.length === 0 ? (
                  <div className="empty-state">No messages yet. Say hi!</div>
                ) : (
                  messages.map((message) => (
                    <div
                      key={message.id}
                      className={`chat-bubble ${message.senderId === myId ? 'chat-bubble-own' : ''}`}
                    >
                      <div className="chat-bubble-text">{message.content}</div>
                      {message.createdAt ? (
                        <div className="chat-bubble-time">{formatTime(message.createdAt)}</div>
                      ) : null}
                    </div>
                  ))
                )}
              </div>

              {messagesError ? <div className="form-error">{messagesError}</div> : null}

              <form className="chat-composer" onSubmit={onSend}>
                <input
                  value={draft}
                  onChange={(event) => setDraft(event.target.value)}
                  placeholder={`Message @${activePeer.username}`}
                  maxLength={2000}
                />
                <button type="submit" className="button" disabled={!draft.trim() || !connected}>
                  Send
                </button>
              </form>
            </>
          )}
        </div>
      </section>
    </PageFrame>
  );
}
