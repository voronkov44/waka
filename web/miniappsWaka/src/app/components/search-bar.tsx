import { Search } from 'lucide-react';

interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

export function SearchBar({ value, onChange, placeholder = 'Search...' }: SearchBarProps) {
  return (
    <div className="relative group">
      <div className="absolute inset-0 bg-foreground/5 rounded-[20px] blur-xl opacity-0 group-focus-within:opacity-100 transition-opacity duration-500" />
      <div className="relative">
        <Search className="absolute left-5 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground group-focus-within:text-foreground transition-colors duration-300" />
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className="w-full bg-card border border-border/50 rounded-[20px] pl-14 pr-5 py-4 text-[13px] font-bold tracking-wide text-foreground placeholder:text-muted-foreground focus:outline-none focus:border-foreground/50 focus:ring-4 focus:ring-foreground/10 transition-all shadow-sm"
        />
      </div>
    </div>
  );
}
