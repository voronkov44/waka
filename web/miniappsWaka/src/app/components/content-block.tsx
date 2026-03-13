import { AlertCircle, Info, CheckCircle, ExternalLink } from 'lucide-react';
import { Link } from 'react-router';
import type { ContentBlock } from '../types/domain';

interface ContentBlockProps {
  block: ContentBlock;
}

export function ContentBlockComponent({ block }: ContentBlockProps) {
  switch (block.type) {
    case 'text':
      return <p className="text-foreground/80 leading-relaxed">{block.content}</p>;

    case 'image':
      return (
        <div className="rounded-2xl overflow-hidden bg-secondary">
          <img
            src={block.url}
            alt=""
            className="w-full h-auto"
          />
        </div>
      );

    case 'link':
      if (!block.url) {
        return null;
      }

      if (block.url.startsWith('/')) {
        return (
          <Link
            to={block.url}
            className="inline-flex items-center gap-2 font-medium text-foreground underline underline-offset-4 hover:opacity-70"
          >
            {block.content}
            <ExternalLink className="h-4 w-4" />
          </Link>
        );
      }

      return (
        <a
          href={block.url}
          target="_blank"
          rel="noreferrer"
          className="inline-flex items-center gap-2 text-foreground font-medium underline underline-offset-4 hover:opacity-70"
        >
          {block.content}
          <ExternalLink className="w-4 h-4" />
        </a>
      );

    case 'bullets':
      return (
        <ul className="space-y-2 pl-1">
          {block.items?.map((item, index) => (
            <li key={index} className="flex items-start gap-3">
              <div className="w-1.5 h-1.5 rounded-full bg-foreground mt-2.5 flex-shrink-0" />
              <span className="text-foreground/80">{item}</span>
            </li>
          ))}
        </ul>
      );

    case 'divider':
      return <hr className="border-border my-6" />;

    case 'callout':
      const variants = {
        info: {
          bg: 'bg-foreground/5',
          border: 'border-foreground/15',
          icon: Info,
          iconColor: 'text-foreground/60',
        },
        warning: {
          bg: 'bg-foreground/5',
          border: 'border-foreground/20',
          icon: AlertCircle,
          iconColor: 'text-foreground/70',
        },
        success: {
          bg: 'bg-foreground/5',
          border: 'border-foreground/15',
          icon: CheckCircle,
          iconColor: 'text-foreground/60',
        },
      };

      const variant = variants[block.variant || 'info'];
      const Icon = variant.icon;

      return (
        <div className={`${variant.bg} border ${variant.border} rounded-2xl p-4 flex gap-3`}>
          <Icon className={`w-5 h-5 ${variant.iconColor} flex-shrink-0 mt-0.5`} />
          <p className="text-foreground/80 leading-relaxed">{block.content}</p>
        </div>
      );

    default:
      return null;
  }
}
