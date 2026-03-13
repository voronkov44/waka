import { useParams, useNavigate } from 'react-router';
import { ChevronLeft, Calendar } from 'lucide-react';
import { format } from 'date-fns';
import { ContentBlockComponent } from '../components/content-block';
import { useFAQArticleDetail } from '../hooks/useFAQArticleDetail';

export function FAQArticleDetail() {
  const { topicId, articleId } = useParams<{ topicId: string; articleId: string }>();
  const parsedTopicID = Number(topicId);
  const parsedArticleID = Number(articleId);
  const navigate = useNavigate();
  const { article, isLoading, error, notFound } = useFAQArticleDetail(parsedArticleID, parsedTopicID);

  const handleGoBack = () => {
    if (window.history.length > 1) {
      navigate(-1);
      return;
    }
    navigate('/faq');
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-muted-foreground">
        Loading article...
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6">
        <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
          {error}
        </div>
      </div>
    );
  }

  if (notFound || !article) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-xl mb-4">Article not found</p>
          <button type="button" onClick={handleGoBack} className="text-foreground underline">
            Go back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-24">
      <div className="sticky top-0 z-40 bg-background/80 backdrop-blur-xl border-b border-border pt-safe">
        <div className="flex items-center gap-4 px-4 py-4">
          <button
            type="button"
            onClick={handleGoBack}
            className="w-10 h-10 rounded-full bg-card border border-border flex items-center justify-center flex-shrink-0"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <div className="flex-1 min-w-0">
            <p className="text-sm text-muted-foreground truncate">Help Article</p>
          </div>
        </div>
      </div>

      <div className="px-6 pt-6">
        <div className="mb-8">
          <h1 className="text-2xl font-bold mb-3">{article.title}</h1>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Calendar className="w-4 h-4" />
            Last updated {format(new Date(article.updatedAt), 'MMMM d, yyyy')}
          </div>
        </div>

        <div className="space-y-5 max-w-2xl">
          {article.contentBlocks.map((block, index) => (
            <ContentBlockComponent key={index} block={block} />
          ))}
        </div>

        <div className="mt-12 pt-8 border-t border-border">
          <div className="bg-card/60 backdrop-blur-sm border border-border rounded-2xl p-6 text-center">
            <h3 className="font-semibold mb-2">Was this article helpful?</h3>
            <p className="text-sm text-muted-foreground mb-4">Let us know if you found the information you needed</p>
            <div className="flex gap-3 justify-center">
              <button
                type="button"
                className="px-6 py-2.5 bg-foreground text-background rounded-full font-semibold hover:opacity-90 transition-all"
              >
                Yes, helpful
              </button>
              <button
                type="button"
                className="px-6 py-2.5 bg-secondary text-foreground rounded-full font-semibold hover:bg-accent transition-all"
              >
                No, not helpful
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
