import React, { useState, useEffect } from 'react';
import { 
  NatsProvider, 
  useSubscription, 
  usePublish, 
  ConnectionIndicator,
  createMessageType,
  MessageEncoder,
  MessageDeduplicator
} from 'react-nats-messaging';

// Define comment types
interface Comment {
  id: string;
  itemId: string;
  userId: string;
  parentId?: string;
  text: string;
  createdAt: string;
  children?: Comment[];
}

// Create simple JSON-based encoder
const encoder = new MessageEncoder({
  useStreamWrapper: false,
  useWebsocketWrapper: false
});

// Create message type with JSON encoding
const commentMessageType = createMessageType<Comment>(
  'comments.AddComment',
  {
    encode: (msg: Comment) => new TextEncoder().encode(JSON.stringify(msg)),
    decode: (data: Uint8Array) => JSON.parse(new TextDecoder().decode(data))
  }
);

// Comment component
function CommentItem({ comment, onReply }: { comment: Comment; onReply: (parentId: string) => void }) {
  return (
    <div style={{ marginLeft: comment.parentId ? '20px' : '0', marginBottom: '12px' }}>
      <div style={{ 
        padding: '12px',
        backgroundColor: '#f9fafb',
        borderRadius: '8px',
        border: '1px solid #e5e7eb'
      }}>
        <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>
          {comment.userId}
        </div>
        <div style={{ marginBottom: '8px' }}>{comment.text}</div>
        <div style={{ 
          display: 'flex', 
          gap: '16px',
          fontSize: '14px',
          color: '#6b7280'
        }}>
          <span>{new Date(comment.createdAt).toLocaleString()}</span>
          <button
            onClick={() => onReply(comment.id)}
            style={{ 
              color: '#3b82f6',
              cursor: 'pointer',
              background: 'none',
              border: 'none'
            }}
          >
            Reply
          </button>
        </div>
      </div>
      {comment.children && comment.children.map(child => (
        <CommentItem key={child.id} comment={child} onReply={onReply} />
      ))}
    </div>
  );
}

function CommentsSection({ itemId }: { itemId: string }) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [inputText, setInputText] = useState('');
  const [replyToId, setReplyToId] = useState<string | null>(null);
  const deduplicator = new MessageDeduplicator<Comment>();

  // Subscribe to comments for this item
  const { messages: incomingComments, error: subError } = useSubscription(
    {
      ...commentMessageType,
      subject: `${commentMessageType.subject}.${itemId}`
    },
    {
      encoder,
      deduplicate: true,
      deduplicationKey: (msg) => msg.id
    }
  );

  // Publish hook
  const { publish, isPublishing, error: pubError } = usePublish(
    {
      ...commentMessageType,
      subject: `${commentMessageType.subject}.${itemId}`
    },
    { encoder }
  );

  // Build comment tree
  const buildCommentTree = (comments: Comment[]): Comment[] => {
    const map = new Map<string, Comment>();
    const roots: Comment[] = [];

    // Create a map of all comments
    comments.forEach(comment => {
      map.set(comment.id, { ...comment, children: [] });
    });

    // Build the tree
    comments.forEach(comment => {
      const node = map.get(comment.id)!;
      if (comment.parentId) {
        const parent = map.get(comment.parentId);
        if (parent) {
          parent.children = parent.children || [];
          parent.children.push(node);
        } else {
          roots.push(node);
        }
      } else {
        roots.push(node);
      }
    });

    return roots;
  };

  // Update comments when new ones arrive
  useEffect(() => {
    if (incomingComments.length > 0) {
      const newComments = incomingComments
        .filter(msg => !deduplicator.isDuplicate(msg));
      
      setComments(prev => [...prev, ...newComments]);
    }
  }, [incomingComments, deduplicator]);

  const handleSubmit = async () => {
    if (!inputText.trim()) return;

    const comment: Comment = {
      id: `comment-${Date.now()}`,
      itemId,
      userId: 'current-user',
      parentId: replyToId || undefined,
      text: inputText,
      createdAt: new Date().toISOString()
    };

    try {
      await publish(comment);
      setInputText('');
      setReplyToId(null);
    } catch (err) {
      console.error('Failed to add comment:', err);
    }
  };

  const handleReply = (parentId: string) => {
    setReplyToId(parentId);
    // Focus input (in real app)
  };

  const commentTree = buildCommentTree(comments);

  return (
    <div>
      <h3 style={{ marginBottom: '16px' }}>
        Comments ({comments.length})
      </h3>
      
      <div style={{ marginBottom: '24px' }}>
        {replyToId && (
          <div style={{ 
            marginBottom: '8px',
            padding: '8px',
            backgroundColor: '#dbeafe',
            borderRadius: '4px',
            fontSize: '14px'
          }}>
            Replying to comment...
            <button
              onClick={() => setReplyToId(null)}
              style={{ 
                marginLeft: '8px',
                color: '#3b82f6',
                cursor: 'pointer',
                background: 'none',
                border: 'none'
              }}
            >
              Cancel
            </button>
          </div>
        )}
        
        <textarea
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          placeholder="Add a comment..."
          style={{ 
            width: '100%',
            padding: '12px',
            borderRadius: '8px',
            border: '1px solid #d1d5db',
            minHeight: '80px',
            resize: 'vertical'
          }}
        />
        
        <button
          onClick={handleSubmit}
          disabled={isPublishing || !inputText.trim()}
          style={{
            marginTop: '8px',
            padding: '8px 16px',
            borderRadius: '4px',
            backgroundColor: '#3b82f6',
            color: 'white',
            border: 'none',
            cursor: isPublishing ? 'not-allowed' : 'pointer',
            opacity: isPublishing ? 0.6 : 1
          }}
        >
          {isPublishing ? 'Posting...' : 'Post Comment'}
        </button>
        
        {(subError || pubError) && (
          <div style={{ color: 'red', marginTop: '8px' }}>
            Error: {(subError || pubError)?.message}
          </div>
        )}
      </div>
      
      <div>
        {commentTree.map(comment => (
          <CommentItem 
            key={comment.id} 
            comment={comment} 
            onReply={handleReply}
          />
        ))}
      </div>
    </div>
  );
}

export function CommentsExample() {
  const [itemId] = useState('item-456');

  return (
    <NatsProvider
      config={{
        servers: 'wss://nats-ws.example.com',
        options: {
          reconnect: true,
          maxReconnectAttempts: 10
        },
        jetstream: {
          enabled: true
        }
      }}
      autoConnect={true}
    >
      <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          marginBottom: '20px'
        }}>
          <h1>Comments Example</h1>
          <ConnectionIndicator />
        </div>
        
        <CommentsSection itemId={itemId} />
      </div>
    </NatsProvider>
  );
}