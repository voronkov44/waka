import { useEffect, useMemo, useRef, useState } from 'react';
import { Check, Search, X } from 'lucide-react';
import { type Locale, useI18n } from '../../shared/i18n';

export interface LanguageOption {
  code: Locale;
  label: string;
  labelEn: string;
  flag: string;
}

export const LANGUAGE_OPTIONS: LanguageOption[] = [
  { code: 'en', label: 'English', labelEn: 'English', flag: '🇬🇧' },
  { code: 'ru', label: 'Русский', labelEn: 'Russian', flag: '🇷🇺' },
];

interface LanguagePickerProps {
  open: boolean;
  current: Locale;
  onSelect: (lang: Locale) => void;
  onClose: () => void;
}

export function LanguagePicker({ open, current, onSelect, onClose }: LanguagePickerProps) {
  const { t } = useI18n();
  const [query, setQuery] = useState('');
  const [visible, setVisible] = useState(false);
  const [animating, setAnimating] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (open) {
      setVisible(true);
      window.requestAnimationFrame(() => {
        window.requestAnimationFrame(() => setAnimating(true));
      });
      window.setTimeout(() => inputRef.current?.focus(), 300);
      return;
    }

    setAnimating(false);
    const timer = window.setTimeout(() => {
      setVisible(false);
      setQuery('');
    }, 320);

    return () => window.clearTimeout(timer);
  }, [open]);

  const filteredOptions = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase();
    if (!normalizedQuery) {
      return LANGUAGE_OPTIONS;
    }

    return LANGUAGE_OPTIONS.filter((item) => {
      return (
        item.label.toLowerCase().includes(normalizedQuery) ||
        item.labelEn.toLowerCase().includes(normalizedQuery) ||
        item.code.toLowerCase().includes(normalizedQuery)
      );
    });
  }, [query]);

  const handleSelect = (language: Locale) => {
    onSelect(language);
    onClose();
  };

  if (!visible) {
    return null;
  }

  return (
    <>
      <div
        className="fixed inset-0 z-[60] transition-opacity duration-300"
        style={{ backgroundColor: `rgba(0,0,0,${animating ? 0.4 : 0})` }}
        onClick={onClose}
      />

      <div
        className="fixed inset-x-0 bottom-0 z-[70] transition-transform duration-300 ease-[cubic-bezier(0.2,0.8,0.2,1)]"
        style={{ transform: animating ? 'translateY(0)' : 'translateY(100%)' }}
      >
        <div className="max-h-[80vh] overflow-hidden rounded-t-[36px] border-t border-border/50 bg-card shadow-2xl flex flex-col">
          <div className="flex justify-center pt-3 pb-1 shrink-0">
            <div className="h-1 w-10 rounded-full bg-border/60" />
          </div>

          <div className="flex items-center gap-4 px-6 pt-4 pb-5 shrink-0">
            <div className="flex-1">
              <h2 className="text-xl font-extrabold tracking-tight text-foreground">{t('profile.languagePicker.title')}</h2>
              <p className="mt-0.5 text-[11px] font-bold tracking-[0.08em] text-muted-foreground">
                {t('profile.languagePicker.availableCount', { count: LANGUAGE_OPTIONS.length })}
              </p>
            </div>
            <button
              type="button"
              onClick={onClose}
              aria-label={t('profile.languagePicker.closeAria')}
              className="h-9 w-9 rounded-[12px] border border-border/50 bg-background flex items-center justify-center hover:bg-foreground/5 transition-colors"
            >
              <X className="h-4 w-4 text-foreground" />
            </button>
          </div>

          <div className="px-6 pb-4 shrink-0">
            <div className="flex items-center gap-3 rounded-[18px] border border-border/50 bg-background px-4 py-3 shadow-sm focus-within:border-foreground/30 transition-colors">
              <Search className="h-4 w-4 text-muted-foreground shrink-0" />
              <input
                ref={inputRef}
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder={t('profile.languagePicker.searchPlaceholder')}
                className="flex-1 bg-transparent text-sm font-medium text-foreground placeholder:text-muted-foreground/50 outline-none"
              />
              {query && (
                <button
                  type="button"
                  onClick={() => setQuery('')}
                  aria-label={t('actions.clear')}
                >
                  <X className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground transition-colors" />
                </button>
              )}
            </div>
          </div>

          <div className="overflow-y-auto overscroll-contain px-6 pb-10">
            {filteredOptions.length === 0 ? (
              <div className="py-12 text-center">
                <p className="text-sm font-bold text-muted-foreground">{t('profile.languagePicker.empty')}</p>
              </div>
            ) : (
              <div className="overflow-hidden rounded-[28px] border border-border/50 bg-background shadow-sm">
                {filteredOptions.map((language, index) => {
                  const isLast = index === filteredOptions.length - 1;
                  const isSelected = language.code === current;

                  return (
                    <button
                      key={language.code}
                      type="button"
                      onClick={() => handleSelect(language.code)}
                      className={`w-full flex items-center gap-4 px-5 py-4 text-left transition-all hover:bg-foreground/5 ${
                        !isLast ? 'border-b border-border/40' : ''
                      }`}
                    >
                      <span className="w-8 shrink-0 text-center text-2xl leading-none">{language.flag}</span>

                      <div className="min-w-0 flex-1">
                        <p className="font-bold leading-tight tracking-tight text-foreground">{language.label}</p>
                        <p className="mt-0.5 text-[11px] font-bold tracking-[0.06em] text-muted-foreground">
                          {language.labelEn}
                        </p>
                      </div>

                      <div
                        className={`h-6 w-6 shrink-0 rounded-full flex items-center justify-center transition-all duration-300 ${
                          isSelected ? 'bg-foreground' : 'border border-border/50'
                        }`}
                      >
                        {isSelected && <Check className="h-3 w-3 text-background" strokeWidth={3} />}
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
