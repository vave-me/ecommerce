import React, { useState, useRef, useEffect } from 'react';
import { DndProvider, useDrag, useDrop } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Underline from '@tiptap/extension-underline';
import Link from '@tiptap/extension-link';
import TextAlign from '@tiptap/extension-text-align';
import { 
  Type, 
  Image, 
  Square, 
  Minus, 
  Link2, 
  Columns, 
  Code,
  Trash2,
  Copy,
  Eye,
  EyeOff,
  Undo,
  Redo,
  Save,
  Download,
  Upload,
  Smartphone,
  Monitor,
  Settings,
  Palette,
  ChevronUp,
  ChevronDown,
  X,
  Plus,
  Grip,
  Bold,
  Italic,
  Underline as UnderlineIcon,
  AlignLeft,
  AlignCenter,
  AlignRight,
  List,
  ListOrdered,
  MessageSquare,
  Move,
  FileText,
  Calendar,
  MapPin,
  Heart,
  Star,
  Package,
  ShoppingBag,
  Gift,
  Mail,
  Globe,
  Facebook,
  Twitter,
  Linkedin,
  Instagram,
  Youtube,
  Layers
} from 'lucide-react';
import { toast } from 'react-toastify';
import styles from './TemplateEditor.module.css';

// Block types
const BLOCK_TYPES = {
  TEXT: 'text',
  HEADING: 'heading',
  IMAGE: 'image',
  BUTTON: 'button',
  DIVIDER: 'divider',
  SPACER: 'spacer',
  COLUMNS: 'columns',
  SOCIAL: 'social',
  HTML: 'html'
};

// Rich Text Editor Component
const RichTextEditor = ({ content, onChange, style }) => {
  const editor = useEditor({
    extensions: [
      StarterKit,
      Underline,
      Link.configure({
        openOnClick: false,
      }),
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      })
    ],
    content: content || '<p>Add your text here...</p>',
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML());
    },
  });

  return (
    <div className={styles.richTextEditor}>
      <div className={styles.textToolbar}>
        <button
          onClick={() => editor.chain().focus().toggleBold().run()}
          className={editor?.isActive('bold') ? styles.isActive : ''}
          title="Bold"
        >
          <Bold size={16} />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleItalic().run()}
          className={editor?.isActive('italic') ? styles.isActive : ''}
          title="Italic"
        >
          <Italic size={16} />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleUnderline().run()}
          className={editor?.isActive('underline') ? styles.isActive : ''}
          title="Underline"
        >
          <UnderlineIcon size={16} />
        </button>
        <div className={styles.separator} />
        <button
          onClick={() => editor.chain().focus().setTextAlign('left').run()}
          className={editor?.isActive({ textAlign: 'left' }) ? styles.isActive : ''}
          title="Align Left"
        >
          <AlignLeft size={16} />
        </button>
        <button
          onClick={() => editor.chain().focus().setTextAlign('center').run()}
          className={editor?.isActive({ textAlign: 'center' }) ? styles.isActive : ''}
          title="Align Center"
        >
          <AlignCenter size={16} />
        </button>
        <button
          onClick={() => editor.chain().focus().setTextAlign('right').run()}
          className={editor?.isActive({ textAlign: 'right' }) ? styles.isActive : ''}
          title="Align Right"
        >
          <AlignRight size={16} />
        </button>
        <div className={styles.separator} />
        <button
          onClick={() => editor.chain().focus().toggleBulletList().run()}
          className={editor?.isActive('bulletList') ? styles.isActive : ''}
          title="Bullet List"
        >
          <List size={16} />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleOrderedList().run()}
          className={editor?.isActive('orderedList') ? styles.isActive : ''}
          title="Numbered List"
        >
          <ListOrdered size={16} />
        </button>
      </div>
      <EditorContent editor={editor} className={styles.textContent} style={style} />
    </div>
  );
};

// Draggable Block Component
const DraggableBlock = ({ type, icon: Icon, label }) => {
  const [{ isDragging }, drag] = useDrag({
    type: 'NEW_BLOCK',
    item: { blockType: type },
    collect: (monitor) => ({
      isDragging: monitor.isDragging()
    })
  });

  return (
    <div
      ref={drag}
      className={`${styles.blockItem} ${isDragging ? styles.dragging : ''}`}
      role="button"
      tabIndex={0}
      aria-label={`Drag ${label} block`}
    >
      <Icon size={20} />
      <span>{label}</span>
    </div>
  );
};

// Template Block Component
const TemplateBlock = ({ block, index, onUpdate, onDelete, onDuplicate, onMove }) => {
  const ref = useRef(null);
  const [isEditing, setIsEditing] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const [showStylePanel, setShowStylePanel] = useState(false);

  const [{ isDragging }, drag] = useDrag({
    type: 'TEMPLATE_BLOCK',
    item: { index, block },
    collect: (monitor) => ({
      isDragging: monitor.isDragging()
    })
  });

  const [{ isOver }, drop] = useDrop({
    accept: ['TEMPLATE_BLOCK', 'NEW_BLOCK'],
    hover: (item, monitor) => {
      if (!ref.current) return;
      
      if (item.index !== undefined && item.index !== index) {
        onMove(item.index, index);
        item.index = index;
      }
    },
    collect: (monitor) => ({
      isOver: monitor.isOver()
    })
  });

  drag(drop(ref));

  const updateBlockStyle = (styleKey, value) => {
    const updatedBlock = {
      ...block,
      style: {
        ...block.style,
        [styleKey]: value
      }
    };
    onUpdate(index, updatedBlock);
  };

  const renderBlockContent = () => {
    switch (block.type) {
      case BLOCK_TYPES.TEXT:
      case BLOCK_TYPES.HEADING:
        return (
          <RichTextEditor
            content={block.content}
            onChange={(html) => onUpdate(index, { ...block, content: html })}
            style={block.style}
          />
        );
      
      case BLOCK_TYPES.IMAGE:
        return (
          <div style={block.style} className={styles.imageBlock}>
            {block.src ? (
              <img src={block.src} alt={block.alt} style={{ maxWidth: '100%' }} />
            ) : (
              <div className={styles.imagePlaceholder} onClick={() => setIsEditing(true)}>
                <Image size={40} />
                <p>Click to add image</p>
              </div>
            )}
            {isEditing && (
              <div className={styles.imageEditPanel}>
                <input
                  type="url"
                  placeholder="Image URL"
                  value={block.src || ''}
                  onChange={(e) => onUpdate(index, { ...block, src: e.target.value })}
                />
                <input
                  type="text"
                  placeholder="Alt text"
                  value={block.alt || ''}
                  onChange={(e) => onUpdate(index, { ...block, alt: e.target.value })}
                />
                <button onClick={() => setIsEditing(false)}>Done</button>
              </div>
            )}
          </div>
        );
      
      case BLOCK_TYPES.BUTTON:
        return (
          <div style={{ textAlign: block.style.textAlign || 'center' }}>
            <a 
              href={block.url} 
              style={block.style}
              className={styles.buttonBlock}
              onClick={(e) => e.preventDefault()}
            >
              {block.text}
            </a>
            {isEditing && (
              <div className={styles.buttonEditPanel}>
                <input
                  type="text"
                  placeholder="Button text"
                  value={block.text || ''}
                  onChange={(e) => onUpdate(index, { ...block, text: e.target.value })}
                />
                <input
                  type="url"
                  placeholder="Button URL"
                  value={block.url || ''}
                  onChange={(e) => onUpdate(index, { ...block, url: e.target.value })}
                />
                <button onClick={() => setIsEditing(false)}>Done</button>
              </div>
            )}
          </div>
        );
      
      case BLOCK_TYPES.DIVIDER:
        return <hr style={block.style} />;
      
      case BLOCK_TYPES.SPACER:
        return (
          <div style={{ ...block.style, minHeight: block.height || 40 }} className={styles.spacerBlock}>
            {isEditing && (
              <input
                type="range"
                min="10"
                max="200"
                value={block.height || 40}
                onChange={(e) => onUpdate(index, { ...block, height: parseInt(e.target.value) })}
              />
            )}
          </div>
        );
      
      case BLOCK_TYPES.COLUMNS:
        return (
          <div className={styles.columnsBlock} style={block.style}>
            {block.columns.map((column, colIndex) => (
              <div 
                key={colIndex} 
                style={{ width: column.width }}
                className={styles.column}
              >
                <div className={styles.columnHeader}>
                  <input
                    type="text"
                    value={column.width}
                    onChange={(e) => {
                      const newColumns = [...block.columns];
                      newColumns[colIndex].width = e.target.value;
                      onUpdate(index, { ...block, columns: newColumns });
                    }}
                    className={styles.columnWidthInput}
                  />
                </div>
                {column.blocks.length === 0 && (
                  <div className={styles.emptyColumn}>
                    <Plus size={20} />
                    <span>Drop blocks here</span>
                  </div>
                )}
                {column.blocks.map((colBlock, blockIndex) => (
                  <TemplateBlock
                    key={blockIndex}
                    block={colBlock}
                    index={blockIndex}
                    onUpdate={(_, updatedBlock) => {
                      const newColumns = [...block.columns];
                      newColumns[colIndex].blocks[blockIndex] = updatedBlock;
                      onUpdate(index, { ...block, columns: newColumns });
                    }}
                    onDelete={() => {
                      const newColumns = [...block.columns];
                      newColumns[colIndex].blocks.splice(blockIndex, 1);
                      onUpdate(index, { ...block, columns: newColumns });
                    }}
                    onDuplicate={() => {
                      const newColumns = [...block.columns];
                      newColumns[colIndex].blocks.splice(blockIndex + 1, 0, { ...colBlock });
                      onUpdate(index, { ...block, columns: newColumns });
                    }}
                    onMove={() => {}}
                  />
                ))}
              </div>
            ))}
          </div>
        );
      
      case BLOCK_TYPES.SOCIAL:
        return (
          <div style={block.style} className={styles.socialBlock}>
            {block.platforms.map((platform) => {
              const Icon = {
                facebook: Facebook,
                twitter: Twitter,
                linkedin: Linkedin,
                instagram: Instagram,
                youtube: Youtube
              }[platform] || Globe;
              
              return (
                <a 
                  key={platform} 
                  href={block.links?.[platform] || '#'} 
                  className={styles.socialIcon}
                  aria-label={platform}
                >
                  <Icon size={20} />
                </a>
              );
            })}
            {isEditing && (
              <div className={styles.socialEditPanel}>
                {block.platforms.map((platform) => (
                  <div key={platform} className={styles.socialLinkInput}>
                    <label>{platform}</label>
                    <input
                      type="url"
                      placeholder={`${platform} URL`}
                      value={block.links?.[platform] || ''}
                      onChange={(e) => {
                        const links = { ...block.links, [platform]: e.target.value };
                        onUpdate(index, { ...block, links });
                      }}
                    />
                  </div>
                ))}
                <button onClick={() => setIsEditing(false)}>Done</button>
              </div>
            )}
          </div>
        );
      
      case BLOCK_TYPES.HTML:
        return (
          <div className={styles.htmlBlock}>
            {isEditing ? (
              <textarea
                value={block.content}
                onChange={(e) => onUpdate(index, { ...block, content: e.target.value })}
                onBlur={() => setIsEditing(false)}
                className={styles.htmlEditor}
                placeholder="Enter custom HTML..."
              />
            ) : (
              <div onClick={() => setIsEditing(true)}>
                <Code size={16} />
                <span>Custom HTML Block</span>
              </div>
            )}
          </div>
        );
      
      default:
        return null;
    }
  };

  return (
    <div
      ref={ref}
      className={`${styles.templateBlock} ${isDragging ? styles.dragging : ''} ${isOver ? styles.dragOver : ''}`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {isHovered && (
        <div className={styles.blockControls}>
          <button 
            onClick={() => onMove(index, Math.max(0, index - 1))}
            className={styles.controlButton}
            aria-label="Move up"
          >
            <ChevronUp size={16} />
          </button>
          <button 
            onClick={() => onMove(index, index + 1)}
            className={styles.controlButton}
            aria-label="Move down"
          >
            <ChevronDown size={16} />
          </button>
          <button 
            onClick={() => setIsEditing(!isEditing)}
            className={styles.controlButton}
            aria-label="Edit block"
          >
            <Settings size={16} />
          </button>
          <button 
            onClick={() => setShowStylePanel(!showStylePanel)}
            className={styles.controlButton}
            aria-label="Style block"
          >
            <Palette size={16} />
          </button>
          <button 
            onClick={() => onDuplicate(index)}
            className={styles.controlButton}
            aria-label="Duplicate block"
          >
            <Copy size={16} />
          </button>
          <button 
            onClick={() => onDelete(index)}
            className={`${styles.controlButton} ${styles.danger}`}
            aria-label="Delete block"
          >
            <Trash2 size={16} />
          </button>
          <div className={styles.dragHandle} aria-label="Drag to reorder">
            <Grip size={16} />
          </div>
        </div>
      )}
      
      {showStylePanel && (
        <div className={styles.stylePanel}>
          <h4>Block Styles</h4>
          <div className={styles.styleGroup}>
            <label>Background Color</label>
            <input
              type="color"
              value={block.style?.backgroundColor || '#ffffff'}
              onChange={(e) => updateBlockStyle('backgroundColor', e.target.value)}
            />
          </div>
          <div className={styles.styleGroup}>
            <label>Padding</label>
            <input
              type="range"
              min="0"
              max="60"
              value={parseInt(block.style?.padding) || 20}
              onChange={(e) => updateBlockStyle('padding', `${e.target.value}px`)}
            />
          </div>
          <div className={styles.styleGroup}>
            <label>Margin</label>
            <input
              type="range"
              min="0"
              max="60"
              value={parseInt(block.style?.margin) || 0}
              onChange={(e) => updateBlockStyle('margin', `${e.target.value}px`)}
            />
          </div>
          <button 
            onClick={() => setShowStylePanel(false)}
            className={styles.closePanelButton}
          >
            Done
          </button>
        </div>
      )}
      
      {renderBlockContent()}
    </div>
  );
};

// Canvas Component that uses drop zone
const CanvasDropZone = ({ blocks, addBlock, updateBlock, deleteBlock, duplicateBlock, moveBlock, previewMode, globalStyles }) => {
  const [{ isOver }, drop] = useDrop({
    accept: ['NEW_BLOCK', 'TEMPLATE_BLOCK'],
    drop: (item, monitor) => {
      if (item.blockType) {
        addBlock(item.blockType);
      }
    },
    collect: (monitor) => ({
      isOver: monitor.isOver()
    })
  });

  return (
    <div 
      ref={drop}
      className={`${styles.canvas} ${previewMode === 'mobile' ? styles.mobileView : ''} ${isOver ? styles.dragOver : ''}`}
      style={{ backgroundColor: globalStyles.backgroundColor }}
    >
      {blocks.length === 0 ? (
        <div className={styles.emptyCanvas}>
          <Plus size={40} />
          <h3>Start Building Your Template</h3>
          <p>Drag blocks from the sidebar or click the + button</p>
          <button 
            onClick={() => addBlock(BLOCK_TYPES.HEADING)}
            className={styles.addFirstBlockButton}
          >
            <Plus size={18} />
            Add First Block
          </button>
        </div>
      ) : (
        <div className={styles.blocksContainer}>
          {blocks.map((block, index) => (
            <TemplateBlock
              key={block.id || index}
              block={block}
              index={index}
              onUpdate={updateBlock}
              onDelete={deleteBlock}
              onDuplicate={duplicateBlock}
              onMove={moveBlock}
            />
          ))}
        </div>
      )}
    </div>
  );
};

// Main Template Editor Component
const TemplateEditor = ({ 
  initialTemplate = null, 
  onSave, 
  onCancel,
  isLoading = false 
}) => {
  const [blocks, setBlocks] = useState([]);
  const [history, setHistory] = useState([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [previewMode, setPreviewMode] = useState('desktop');
  const [showPreview, setShowPreview] = useState(false);
  const [templateName, setTemplateName] = useState('');
  const [templateDescription, setTemplateDescription] = useState('');
  const [isPublic, setIsPublic] = useState(false);
  const [globalStyles, setGlobalStyles] = useState({
    fontFamily: 'Arial, sans-serif',
    primaryColor: '#007bff',
    backgroundColor: '#ffffff',
    textColor: '#333333'
  });

  // Initialize template
  useEffect(() => {
    if (initialTemplate) {
      setTemplateName(initialTemplate.name || '');
      setTemplateDescription(initialTemplate.description || '');
      setIsPublic(initialTemplate.isPublic || false);
      
      // Parse blocks from template if available
      const parsedBlocks = initialTemplate.blocks || parseHtmlToBlocks(initialTemplate.htmlTemplate);
      setBlocks(parsedBlocks);
      addToHistory(parsedBlocks);
    }
  }, [initialTemplate]);

  // History management
  const addToHistory = (newBlocks) => {
    const newHistory = history.slice(0, historyIndex + 1);
    newHistory.push(JSON.parse(JSON.stringify(newBlocks)));
    setHistory(newHistory);
    setHistoryIndex(newHistory.length - 1);
  };

  const undo = () => {
    if (historyIndex > 0) {
      setHistoryIndex(historyIndex - 1);
      setBlocks(JSON.parse(JSON.stringify(history[historyIndex - 1])));
    }
  };

  const redo = () => {
    if (historyIndex < history.length - 1) {
      setHistoryIndex(historyIndex + 1);
      setBlocks(JSON.parse(JSON.stringify(history[historyIndex + 1])));
    }
  };

  // Block management
  const addBlock = (blockType, index = blocks.length) => {
    const newBlock = createDefaultBlock(blockType);
    const newBlocks = [...blocks];
    newBlocks.splice(index, 0, newBlock);
    setBlocks(newBlocks);
    addToHistory(newBlocks);
  };

  const createDefaultBlock = (blockType) => {
    const defaults = {
      [BLOCK_TYPES.TEXT]: {
        type: BLOCK_TYPES.TEXT,
        content: '<p>Add your text here...</p>',
        style: {
          fontSize: '16px',
          color: globalStyles.textColor,
          textAlign: 'left',
          padding: '20px',
          fontFamily: globalStyles.fontFamily
        }
      },
      [BLOCK_TYPES.HEADING]: {
        type: BLOCK_TYPES.HEADING,
        content: '<h2>Your Heading Here</h2>',
        style: {
          fontSize: '24px',
          color: globalStyles.textColor,
          textAlign: 'center',
          padding: '20px',
          fontFamily: globalStyles.fontFamily,
          fontWeight: 'bold'
        }
      },
      [BLOCK_TYPES.IMAGE]: {
        type: BLOCK_TYPES.IMAGE,
        src: '',
        alt: 'Image',
        style: {
          width: '100%',
          padding: '20px',
          textAlign: 'center'
        }
      },
      [BLOCK_TYPES.BUTTON]: {
        type: BLOCK_TYPES.BUTTON,
        text: 'Click Here',
        url: '#',
        style: {
          backgroundColor: globalStyles.primaryColor,
          color: '#ffffff',
          padding: '12px 24px',
          borderRadius: '4px',
          textAlign: 'center',
          fontSize: '16px',
          fontWeight: 'bold',
          textDecoration: 'none',
          display: 'inline-block',
          margin: '20px'
        }
      },
      [BLOCK_TYPES.DIVIDER]: {
        type: BLOCK_TYPES.DIVIDER,
        style: {
          borderTop: '1px solid #e0e0e0',
          margin: '20px 0',
          width: '100%'
        }
      },
      [BLOCK_TYPES.SPACER]: {
        type: BLOCK_TYPES.SPACER,
        height: 40,
        style: {
          height: '40px'
        }
      },
      [BLOCK_TYPES.COLUMNS]: {
        type: BLOCK_TYPES.COLUMNS,
        columns: [
          { blocks: [], width: '50%' },
          { blocks: [], width: '50%' }
        ],
        style: {
          padding: '20px'
        }
      },
      [BLOCK_TYPES.SOCIAL]: {
        type: BLOCK_TYPES.SOCIAL,
        platforms: ['facebook', 'twitter', 'linkedin', 'instagram'],
        links: {},
        style: {
          textAlign: 'center',
          padding: '20px'
        }
      },
      [BLOCK_TYPES.HTML]: {
        type: BLOCK_TYPES.HTML,
        content: '<!-- Custom HTML -->',
        style: {}
      }
    };

    return {
      ...defaults[blockType],
      id: `${blockType}_${Date.now()}`
    };
  };

  const updateBlock = (index, updatedBlock) => {
    const newBlocks = [...blocks];
    newBlocks[index] = updatedBlock;
    setBlocks(newBlocks);
    addToHistory(newBlocks);
  };

  const deleteBlock = (index) => {
    const newBlocks = blocks.filter((_, i) => i !== index);
    setBlocks(newBlocks);
    addToHistory(newBlocks);
  };

  const duplicateBlock = (index) => {
    const newBlocks = [...blocks];
    const duplicatedBlock = { 
      ...JSON.parse(JSON.stringify(blocks[index])), 
      id: `${blocks[index].type}_${Date.now()}` 
    };
    newBlocks.splice(index + 1, 0, duplicatedBlock);
    setBlocks(newBlocks);
    addToHistory(newBlocks);
  };

  const moveBlock = (fromIndex, toIndex) => {
    if (toIndex < 0 || toIndex >= blocks.length) return;
    
    const newBlocks = [...blocks];
    const [movedBlock] = newBlocks.splice(fromIndex, 1);
    newBlocks.splice(toIndex, 0, movedBlock);
    setBlocks(newBlocks);
    addToHistory(newBlocks);
  };

  // Convert blocks to HTML
  const blocksToHtml = () => {
    const html = blocks.map(block => blockToHtml(block)).join('\n');
    return wrapInEmailTemplate(html);
  };

  const blockToHtml = (block) => {
    const styleString = styleToString(block.style);
    
    switch (block.type) {
      case BLOCK_TYPES.TEXT:
      case BLOCK_TYPES.HEADING:
        return `<div style="${styleString}">${block.content}</div>`;
      
      case BLOCK_TYPES.IMAGE:
        return block.src ? 
          `<div style="${styleString}"><img src="${block.src}" alt="${block.alt}" style="max-width: 100%; height: auto;" /></div>` : '';
      
      case BLOCK_TYPES.BUTTON:
        return `<div style="text-align: ${block.style.textAlign || 'center'};">
          <a href="${block.url}" style="${styleString}">${block.text}</a>
        </div>`;
      
      case BLOCK_TYPES.DIVIDER:
        return `<hr style="${styleString}" />`;
      
      case BLOCK_TYPES.SPACER:
        return `<div style="${styleString}"></div>`;
      
      case BLOCK_TYPES.COLUMNS:
        const columnsHtml = block.columns.map(column => 
          `<td style="width: ${column.width}; vertical-align: top; padding: 10px;">
            ${column.blocks.map(b => blockToHtml(b)).join('\n')}
          </td>`
        ).join('\n');
        return `<table width="100%" cellpadding="0" cellspacing="0" style="${styleString}">
          <tr>${columnsHtml}</tr>
        </table>`;
      
      case BLOCK_TYPES.SOCIAL:
        const socialHtml = block.platforms.map(platform => {
          const url = block.links?.[platform] || '#';
          return `<a href="${url}" style="display: inline-block; margin: 0 10px;">
            <img src="https://via.placeholder.com/32x32/333333/ffffff?text=${platform[0].toUpperCase()}" 
                 alt="${platform}" width="32" height="32" style="border-radius: 50%;" />
          </a>`;
        }).join('');
        return `<div style="${styleString}">${socialHtml}</div>`;
      
      case BLOCK_TYPES.HTML:
        return block.content;
      
      default:
        return '';
    }
  };

  const styleToString = (style) => {
    return Object.entries(style || {})
      .map(([key, value]) => `${key.replace(/([A-Z])/g, '-$1').toLowerCase()}: ${value}`)
      .join('; ');
  };

  const wrapInEmailTemplate = (content) => {
    return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${templateName || 'Newsletter'}</title>
  <style>
    body { margin: 0; padding: 0; font-family: ${globalStyles.fontFamily}; background-color: #f4f4f4; }
    .email-container { max-width: 600px; margin: 0 auto; background-color: ${globalStyles.backgroundColor}; }
    @media only screen and (max-width: 600px) {
      .email-container { width: 100% !important; }
      table { width: 100% !important; }
      td { display: block !important; width: 100% !important; }
    }
  </style>
</head>
<body>
  <div class="email-container">
    ${content}
  </div>
</body>
</html>`;
  };

  const parseHtmlToBlocks = (html) => {
    // Simple parser - in production, use a proper HTML parser
    return [];
  };

  // Handle save
  const handleSave = async () => {
    if (!templateName) {
      toast.error('Please enter a template name');
      return;
    }

    const templateData = {
      name: templateName,
      description: templateDescription,
      htmlTemplate: blocksToHtml(),
      textTemplate: blocks.map(b => {
        if (b.type === BLOCK_TYPES.TEXT || b.type === BLOCK_TYPES.HEADING) {
          return b.content.replace(/<[^>]*>/g, '');
        }
        if (b.type === BLOCK_TYPES.BUTTON) return `${b.text}: ${b.url}`;
        return '';
      }).filter(Boolean).join('\n\n'),
      isPublic: isPublic,
      variables: {},
      previewData: {},
      blocks: blocks // Save block structure for future editing
    };

    onSave(templateData);
  };

  // Export template
  const exportTemplate = () => {
    const templateData = {
      name: templateName,
      description: templateDescription,
      blocks: blocks,
      globalStyles: globalStyles,
      version: '1.0'
    };
    
    const dataStr = JSON.stringify(templateData, null, 2);
    const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr);
    
    const exportFileDefaultName = `${templateName.replace(/\s+/g, '-').toLowerCase()}-template.json`;
    
    const linkElement = document.createElement('a');
    linkElement.setAttribute('href', dataUri);
    linkElement.setAttribute('download', exportFileDefaultName);
    linkElement.click();
  };

  // Import template
  const importTemplate = (event) => {
    const file = event.target.files[0];
    if (!file) return;
    
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const templateData = JSON.parse(e.target.result);
        setBlocks(templateData.blocks || []);
        setTemplateName(templateData.name || 'Imported Template');
        setTemplateDescription(templateData.description || '');
        setGlobalStyles(templateData.globalStyles || globalStyles);
        addToHistory(templateData.blocks || []);
        toast.success('Template imported successfully');
      } catch (error) {
        toast.error('Failed to import template');
      }
    };
    reader.readAsText(file);
  };

  return (
    <DndProvider backend={HTML5Backend}>
      <div className={styles.editorContainer}>
        {/* Header */}
        <div className={styles.editorHeader}>
          <div className={styles.headerLeft}>
            <button onClick={onCancel} className={styles.cancelButton}>
              <X size={18} />
              Cancel
            </button>
            <input
              type="text"
              placeholder="Template Name"
              value={templateName}
              onChange={(e) => setTemplateName(e.target.value)}
              className={styles.templateNameInput}
              aria-label="Template name"
            />
          </div>
          
          <div className={styles.headerCenter}>
            <button 
              onClick={undo}
              disabled={historyIndex <= 0}
              className={styles.toolButton}
              aria-label="Undo"
            >
              <Undo size={18} />
            </button>
            <button 
              onClick={redo}
              disabled={historyIndex >= history.length - 1}
              className={styles.toolButton}
              aria-label="Redo"
            >
              <Redo size={18} />
            </button>
            
            <div className={styles.viewToggle}>
              <button 
                onClick={() => setPreviewMode('desktop')}
                className={`${styles.viewButton} ${previewMode === 'desktop' ? styles.active : ''}`}
                aria-label="Desktop view"
              >
                <Monitor size={18} />
              </button>
              <button 
                onClick={() => setPreviewMode('mobile')}
                className={`${styles.viewButton} ${previewMode === 'mobile' ? styles.active : ''}`}
                aria-label="Mobile view"
              >
                <Smartphone size={18} />
              </button>
            </div>
            
            <button 
              onClick={() => setShowPreview(!showPreview)}
              className={styles.toolButton}
              aria-label={showPreview ? 'Hide preview' : 'Show preview'}
            >
              {showPreview ? <EyeOff size={18} /> : <Eye size={18} />}
              Preview
            </button>
            
            <button 
              onClick={exportTemplate}
              className={styles.toolButton}
              aria-label="Export template"
            >
              <Download size={18} />
            </button>
            
            <label className={styles.toolButton} aria-label="Import template">
              <Upload size={18} />
              <input 
                type="file" 
                accept=".json"
                onChange={importTemplate}
                style={{ display: 'none' }}
              />
            </label>
          </div>
          
          <div className={styles.headerRight}>
            <button 
              onClick={handleSave}
              disabled={isLoading || !templateName}
              className={styles.saveButton}
            >
              <Save size={18} />
              {isLoading ? 'Saving...' : 'Save Template'}
            </button>
          </div>
        </div>

        <div className={styles.editorBody}>
          {/* Sidebar */}
          <div className={styles.sidebar}>
            <h3 className={styles.sidebarTitle}>Blocks</h3>
            <div className={styles.blocksList}>
              <DraggableBlock type={BLOCK_TYPES.TEXT} icon={Type} label="Text" />
              <DraggableBlock type={BLOCK_TYPES.HEADING} icon={FileText} label="Heading" />
              <DraggableBlock type={BLOCK_TYPES.IMAGE} icon={Image} label="Image" />
              <DraggableBlock type={BLOCK_TYPES.BUTTON} icon={Square} label="Button" />
              <DraggableBlock type={BLOCK_TYPES.DIVIDER} icon={Minus} label="Divider" />
              <DraggableBlock type={BLOCK_TYPES.SPACER} icon={Move} label="Spacer" />
              <DraggableBlock type={BLOCK_TYPES.COLUMNS} icon={Columns} label="Columns" />
              <DraggableBlock type={BLOCK_TYPES.SOCIAL} icon={Link2} label="Social" />
              <DraggableBlock type={BLOCK_TYPES.HTML} icon={Code} label="HTML" />
            </div>

            <div className={styles.templateSettings}>
              <h3 className={styles.sidebarTitle}>Global Styles</h3>
              <div className={styles.settingsGroup}>
                <label>Font Family</label>
                <select
                  value={globalStyles.fontFamily}
                  onChange={(e) => setGlobalStyles({ ...globalStyles, fontFamily: e.target.value })}
                >
                  <option value="Arial, sans-serif">Arial</option>
                  <option value="Georgia, serif">Georgia</option>
                  <option value="'Times New Roman', serif">Times New Roman</option>
                  <option value="'Courier New', monospace">Courier New</option>
                  <option value="Verdana, sans-serif">Verdana</option>
                </select>
              </div>
              <div className={styles.settingsGroup}>
                <label>Primary Color</label>
                <input
                  type="color"
                  value={globalStyles.primaryColor}
                  onChange={(e) => setGlobalStyles({ ...globalStyles, primaryColor: e.target.value })}
                />
              </div>
              <div className={styles.settingsGroup}>
                <label>Background Color</label>
                <input
                  type="color"
                  value={globalStyles.backgroundColor}
                  onChange={(e) => setGlobalStyles({ ...globalStyles, backgroundColor: e.target.value })}
                />
              </div>
              <div className={styles.settingsGroup}>
                <label>Text Color</label>
                <input
                  type="color"
                  value={globalStyles.textColor}
                  onChange={(e) => setGlobalStyles({ ...globalStyles, textColor: e.target.value })}
                />
              </div>
              
              <h3 className={styles.sidebarTitle}>Template Settings</h3>
              <div className={styles.settingsGroup}>
                <label>Description</label>
                <textarea
                  value={templateDescription}
                  onChange={(e) => setTemplateDescription(e.target.value)}
                  placeholder="Template description..."
                  rows={3}
                />
              </div>
              <div className={styles.settingsGroup}>
                <label className={styles.checkboxLabel}>
                  <input
                    type="checkbox"
                    checked={isPublic}
                    onChange={(e) => setIsPublic(e.target.checked)}
                  />
                  Make template public
                </label>
              </div>
            </div>
          </div>

          {/* Canvas */}
          <div className={styles.canvasContainer}>
            <CanvasDropZone
              blocks={blocks}
              addBlock={addBlock}
              updateBlock={updateBlock}
              deleteBlock={deleteBlock}
              duplicateBlock={duplicateBlock}
              moveBlock={moveBlock}
              previewMode={previewMode}
              globalStyles={globalStyles}
            />
          </div>

          {/* Preview Panel */}
          {showPreview && (
            <div className={styles.previewPanel}>
              <div className={styles.previewHeader}>
                <h3>Preview</h3>
                <button 
                  onClick={() => setShowPreview(false)}
                  className={styles.closeButton}
                >
                  <X size={18} />
                </button>
              </div>
              <iframe
                srcDoc={blocksToHtml()}
                className={`${styles.previewFrame} ${previewMode === 'mobile' ? styles.mobilePreview : ''}`}
                title="Template preview"
              />
            </div>
          )}
        </div>
      </div>
    </DndProvider>
  );
};

export default TemplateEditor;