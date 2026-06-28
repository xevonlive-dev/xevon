// Ambient declaration for lucide-react v1.14.0
// The installed package is missing its dist/lucide-react.d.ts file.
// This declaration silences TS2307/TS7016 and types all named icon imports.
declare module 'lucide-react' {
  import type { ComponentType, SVGProps } from 'react';

  /** Props accepted by every Lucide icon component. */
  export interface LucideProps extends SVGProps<SVGSVGElement> {
    size?: number | string;
    strokeWidth?: number | string;
    absoluteStrokeWidth?: boolean;
    color?: string;
  }

  /** A single Lucide icon React component. */
  export type LucideIcon = ComponentType<LucideProps>;

  // ── All icons used in this project ─────────────────────────────────────────
  export const Activity: LucideIcon;
  export const ArrowDown: LucideIcon;
  export const ArrowUp: LucideIcon;
  export const ArrowUpDown: LucideIcon;
  export const Blocks: LucideIcon;
  export const Bot: LucideIcon;
  export const Box: LucideIcon;
  export const Bug: LucideIcon;
  export const Check: LucideIcon;
  export const CheckCircle: LucideIcon;
  export const ChevronDown: LucideIcon;
  export const ChevronLeft: LucideIcon;
  export const ChevronRight: LucideIcon;
  export const Clock: LucideIcon;
  export const Cloud: LucideIcon;
  export const Code: LucideIcon;
  export const Coins: LucideIcon;
  export const Columns: LucideIcon;
  export const Copy: LucideIcon;
  export const CreditCard: LucideIcon;
  export const Crosshair: LucideIcon;
  export const Database: LucideIcon;
  export const Eye: LucideIcon;
  export const Filter: LucideIcon;
  export const FolderKanban: LucideIcon;
  export const Globe: LucideIcon;
  export const Import: LucideIcon;
  export const Info: LucideIcon;
  export const KeyRound: LucideIcon;
  export const LayoutDashboard: LucideIcon;
  export const Layers: LucideIcon;
  export const Link: LucideIcon;
  export const List: LucideIcon;
  export const Loader2: LucideIcon;
  export const LogIn: LucideIcon;
  export const Mail: LucideIcon;
  export const Menu: LucideIcon;
  export const MessageSquare: LucideIcon;
  export const Monitor: LucideIcon;
  export const Moon: LucideIcon;
  export const Network: LucideIcon;
  export const PanelLeftClose: LucideIcon;
  export const PanelLeftOpen: LucideIcon;
  export const Palette: LucideIcon;
  export const Pencil: LucideIcon;
  export const Play: LucideIcon;
  export const Plus: LucideIcon;
  export const Puzzle: LucideIcon;
  export const Radar: LucideIcon;
  export const Radio: LucideIcon;
  export const RefreshCw: LucideIcon;
  export const Scale: LucideIcon;
  export const ScrollText: LucideIcon;
  export const Search: LucideIcon;
  export const Send: LucideIcon;
  export const Server: LucideIcon;
  export const Settings: LucideIcon;
  export const Settings2: LucideIcon;
  export const Shield: LucideIcon;
  export const ShieldAlert: LucideIcon;
  export const ShieldCheck: LucideIcon;
  export const SlidersHorizontal: LucideIcon;
  export const Square: LucideIcon;
  export const Sun: LucideIcon;
  export const Terminal: LucideIcon;
  export const Trash2: LucideIcon;
  export const Upload: LucideIcon;
  export const User: LucideIcon;
  export const Users: LucideIcon;
  export const X: LucideIcon;
  export const XCircle: LucideIcon;
  export const Zap: LucideIcon;
}
