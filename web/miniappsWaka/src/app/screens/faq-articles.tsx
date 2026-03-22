import { useMemo } from 'react';
import { useParams, Link, useNavigate } from 'react-router';
import { ChevronLeft, ChevronRight, Calendar } from 'lucide-react';
import { useFAQArticles } from '../hooks/useFAQArticles';
import { resolveI18nText, useI18n } from '../../shared/i18n';

export function FAQArticles() {
  const { t, tp, localeCode } = useI18n();
  const { topicId } = useParams<{ topicId: string }>();
  const parsedTopicID = Number(topicId);
  const navigate = useNavigate();
  const { topic, articles, isLoading, error, notFound } = useFAQArticles(parsedTopicID);
  const localizedError = resolveI18nText(error, t);

  const handleGoBack = () => {
    if (window.history.length > 1) {
      navigate(-1);
      return;
    }
    navigate('/faq');
  };

  const sortedArticles = useMemo(
    () =>
      [...articles].sort((a, b) => {
        return new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime();
      }),
    [articles],
  );
  const articleCountLabel = tp('nouns.article', sortedArticles.length);
  const articleDateFormatter = useMemo(
    () =>
      new Intl.DateTimeFormat(localeCode, {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      }),
    [localeCode],
  );

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-muted-foreground">
        {t('faqArticles.loading')}
      </div>
    );
  }

  if (localizedError) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6">
        <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
          {localizedError}
        </div>
      </div>
    );
  }

  if (notFound || !topic) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-xl mb-4">{t('faqArticles.topicNotFound')}</p>
          <Link to="/faq" className="text-foreground underline">
            {t('actions.backToFaq')}
          </Link>
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
            className="w-10 h-10 rounded-full bg-card border border-border flex items-center justify-center"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <div className="flex-1">
            <h1 className="text-xl font-bold">{topic.title}</h1>
            <p className="text-sm text-muted-foreground">
              {t('faqArticles.headerCount', {
                count: sortedArticles.length.toLocaleString(localeCode),
                unit: articleCountLabel,
              })}
            </p>
          </div>
        </div>
      </div>

      <div className="px-6 pt-6">
        {sortedArticles.length > 0 ? (
          <div className="space-y-3">
            {sortedArticles.map((article) => (
              <Link
                key={article.id}
                to={`/faq/${topic.id}/${article.id}`}
                className="block bg-card/60 backdrop-blur-sm border border-border rounded-2xl p-4 hover:border-foreground/20 transition-all"
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="flex-1">
                    <h3 className="font-semibold mb-2">{article.title}</h3>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <Calendar className="w-3 h-3" />
                      {t('common.updatedAt', {
                        date: articleDateFormatter.format(new Date(article.updatedAt)),
                      })}
                    </div>
                  </div>
                  <ChevronRight className="w-5 h-5 text-muted-foreground flex-shrink-0" />
                </div>
              </Link>
            ))}
          </div>
        ) : (
          <div className="text-center py-16">
            <p className="text-muted-foreground">{t('faqArticles.noArticlesAvailable')}</p>
          </div>
        )}
      </div>
    </div>
  );
}
