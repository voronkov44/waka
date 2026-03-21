import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router';
import { ChevronRight, Rocket, Info, Sparkles, Wrench, AlertCircle, Shield } from 'lucide-react';
import { SearchBar } from '../components/search-bar';
import { apiClient } from '../api/client';
import { useFAQTopics } from '../hooks/useFAQTopics';
import type { FAQTopicIcon } from '../types/domain';

const iconMap: Record<FAQTopicIcon, typeof Rocket> = {
  Rocket,
  Info,
  Sparkles,
  Wrench,
  AlertCircle,
  Shield,
};

export function FAQTopics() {
  const [searchQuery, setSearchQuery] = useState('');
  const [apiMatchedTopicIDs, setApiMatchedTopicIDs] = useState<Set<number> | null>(null);
  const { topics, isLoading, error } = useFAQTopics();

  useEffect(() => {
    const query = searchQuery.trim();
    if (!query) {
      setApiMatchedTopicIDs(null);
      return;
    }

    setApiMatchedTopicIDs(null);
    let isCancelled = false;
    const timeoutID = window.setTimeout(() => {
      void apiClient
        .searchFAQArticles(query)
        .then((articles) => {
          if (isCancelled) {
            return;
          }
          setApiMatchedTopicIDs(new Set(articles.map((article) => article.topic_id)));
        })
        .catch(() => {
          if (isCancelled) {
            return;
          }
          setApiMatchedTopicIDs(new Set());
        });
    }, 300);

    return () => {
      isCancelled = true;
      window.clearTimeout(timeoutID);
    };
  }, [searchQuery]);

  const filteredTopics = useMemo(() => {
    const query = searchQuery.trim().toLowerCase();
    if (!query) {
      return topics;
    }

    return topics.filter((topic) => {
      const localMatch = topic.title.toLowerCase().includes(query) || topic.description.toLowerCase().includes(query);
      if (localMatch) {
        return true;
      }

      return apiMatchedTopicIDs?.has(topic.id) ?? false;
    });
  }, [topics, searchQuery, apiMatchedTopicIDs]);

  return (
    <div className="min-h-screen pb-32">
      <div className="px-6 pt-14 pb-8">
        <h1 className="text-4xl font-extrabold tracking-tighter leading-none mb-3">Help Center</h1>
        <p className="text-[11px] font-bold tracking-[0.1em] uppercase text-muted-foreground">Find answers & get support</p>
      </div>

      <div className="px-6 mb-8">
        <SearchBar value={searchQuery} onChange={setSearchQuery} placeholder="Search help topics..." />
      </div>

      <div className="px-6">
        {isLoading && <p className="text-sm text-muted-foreground py-8 text-center">Loading FAQ topics...</p>}

        {error && (
          <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
            {error}
          </div>
        )}

        {!isLoading && !error && (
          <div className="space-y-4">
            {filteredTopics.map((topic) => {
              const Icon = iconMap[topic.icon];
              return (
                <Link
                  key={topic.id}
                  to={`/faq/${topic.id}`}
                  className="group block bg-card border border-border/50 shadow-sm hover:shadow-lg dark:shadow-none rounded-[32px] p-6 transition-all duration-500 overflow-hidden relative"
                >
                  <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-foreground/5 to-transparent rounded-full blur-2xl -translate-y-1/2 translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                  <div className="flex items-start gap-5 relative z-10">
                    <div className="w-14 h-14 rounded-[18px] bg-background border border-border/50 flex items-center justify-center flex-shrink-0 shadow-sm group-hover:scale-105 group-hover:border-foreground/30 transition-all duration-500">
                      <Icon className="w-6 h-6 text-foreground" />
                    </div>
                    <div className="flex-1 min-w-0 pt-1">
                      <div className="flex items-center justify-between mb-1.5">
                        <h3 className="font-bold text-lg tracking-tight text-foreground">{topic.title}</h3>
                        <ChevronRight className="w-5 h-5 text-muted-foreground/50 group-hover:text-foreground transition-colors" />
                      </div>
                      <p className="text-sm font-medium text-muted-foreground mb-3 leading-relaxed">{topic.description}</p>
                      <p className="text-[9px] font-bold tracking-[0.2em] uppercase text-foreground/70">
                        {topic.articleCount} articles
                      </p>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        )}

        {!isLoading && !error && filteredTopics.length === 0 && (
          <div className="text-center py-16">
            <p className="text-[11px] font-bold tracking-[0.1em] uppercase text-muted-foreground">No topics found</p>
          </div>
        )}
      </div>

      <div className="px-6 mt-10">
        <div className="bg-card border border-border/50 rounded-[32px] p-8 text-center shadow-sm">
          <h3 className="text-xl font-bold tracking-tight mb-2 text-foreground">Still need help?</h3>
          <p className="text-[11px] font-bold tracking-[0.1em] uppercase text-muted-foreground mb-8">
            Our premium support team is here to assist you
          </p>
          <Link
            to="/profile"
            className="inline-flex items-center gap-3 px-8 py-4 bg-foreground text-background rounded-full text-[11px] font-bold tracking-[0.2em] uppercase hover:scale-105 hover:bg-foreground/90 transition-all duration-500 shadow-md hover:shadow-lg"
          >
            Contact Support
            <ChevronRight className="w-4 h-4" />
          </Link>
        </div>
      </div>
    </div>
  );
}
