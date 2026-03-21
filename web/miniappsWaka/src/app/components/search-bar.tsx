import { Search, X } from 'lucide-react';
import { useRef, type FormEvent } from 'react';

interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

export function SearchBar({ value, onChange, placeholder = 'Search...' }: SearchBarProps) {
  const inputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    window.requestAnimationFrame(() => {
      inputRef.current?.blur();
    });
  };

  return (
    <form className="relative group" role="search" onSubmit={handleSubmit}>
      <div className="absolute inset-0 bg-foreground/5 rounded-[20px] blur-xl opacity-0 group-focus-within:opacity-100 transition-opacity duration-500" />
      <div className="relative">
        <button
          type="submit"
          aria-label="Submit search"
          className="absolute left-5 top-1/2 -translate-y-1/2 text-muted-foreground group-focus-within:text-foreground transition-colors duration-300 p-0 m-0 border-0 bg-transparent"
        >
          <Search className="w-5 h-5" />
        </button>
        <input
          ref={inputRef}
          type="text"
          inputMode="search"
          enterKeyHint="search"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className="w-full bg-card border border-border/50 rounded-[20px] pl-14 pr-14 py-4 text-[13px] font-bold tracking-wide text-foreground placeholder:text-muted-foreground focus:outline-none focus:border-foreground/50 focus:ring-4 focus:ring-foreground/10 transition-all shadow-sm"
        />
        {value.length > 0 && (
          <button
            type="button"
            aria-label="Clear search"
            onClick={() => {
              onChange('');
              window.requestAnimationFrame(() => {
                inputRef.current?.focus();
              });
            }}
            className="absolute right-2 top-1/2 -translate-y-1/2 inline-flex h-9 w-9 items-center justify-center rounded-full text-muted-foreground hover:text-foreground active:scale-95 transition-all"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
    </form>
  );
}
