'use client';

import { useEffect, useRef, type CSSProperties } from 'react';

interface SaturnOrbProps {
  size?: number;
  className?: string;
  style?: CSSProperties;
}

/**
 * SaturnOrb — animated ASCII Saturn for the login screen.
 *
 * A shaded planet with TILTED rings: the near (front) arc draws over the planet,
 * the far (back) arc passes behind it. A comet-style glint orbits the rings,
 * the planet's bands drift, stars twinkle, and the occasional shooting star
 * streaks past. Geometry is procedural; x-radii are widened because monospace
 * cells are ~0.6× as wide as tall.
 */

const COLS = 90;
const ROWS = 42;
const CX = 45;
const CY = 21;

const P_RX = 17, P_RY = 9.5;      // planet radii (cells)
const R_RXO = 42, R_RYO = 7.5;    // ring outer radii (wide + flat)
const R_RIN = 0.55;               // ring inner edge as a fraction of outer
const TILT = 0.34;                // ring tilt in radians (~19°)
const COS_T = Math.cos(TILT), SIN_T = Math.sin(TILT);

// Light direction for sphere shading (normalized), from the upper-left.
const LX = -0.42, LY = -0.5, LZ = 0.76;

const RAMP = ' .:-=+*oa#%@';

type Kind = 'void' | 'planet' | 'ringFront' | 'ringBack' | 'star';

interface Built {
  kind: Kind;
  ch: string;
  lum: number;
  ang: number;   // ring angle in ring-local space (for the glint sweep)
  pny: number;   // planet normalized y (for band drift)
  baseLum: number;
}

function rampChar(l: number): string {
  const i = Math.max(0, Math.min(RAMP.length - 1, Math.round(l * (RAMP.length - 1))));
  return RAMP[i];
}

function planetBucket(l: number): string {
  if (l < 0.18) return 'sat-p0';
  if (l < 0.34) return 'sat-p1';
  if (l < 0.5) return 'sat-p2';
  if (l < 0.66) return 'sat-p3';
  if (l < 0.8) return 'sat-p4';
  if (l < 0.92) return 'sat-p5';
  return 'sat-p6';
}

function hash(x: number, y: number): number {
  let h = x * 374761393 + y * 668265263;
  h = (h ^ (h >>> 13)) * 1274126177;
  h = (h ^ (h >>> 16)) >>> 0;
  return (h % 1000) / 1000;
}

function buildModel(): Built[][] {
  const grid: Built[][] = [];
  for (let y = 0; y < ROWS; y++) {
    const row: Built[] = [];
    for (let x = 0; x < COLS; x++) {
      let cell: Built = { kind: 'void', ch: ' ', lum: 0, ang: 0, pny: 0, baseLum: 0 };

      const pnx = (x - CX) / P_RX;
      const pny = (y - CY) / P_RY;
      const pr2 = pnx * pnx + pny * pny;
      const onPlanet = pr2 <= 1;

      // Tilted ring: rotate the screen offset into ring-local space.
      const dx = x - CX, dy = y - CY;
      const lx = dx * COS_T + dy * SIN_T;
      const ly = -dx * SIN_T + dy * COS_T;
      const rnx = lx / R_RXO, rny = ly / R_RYO;
      const rd = Math.sqrt(rnx * rnx + rny * rny);
      const onRing = rd > R_RIN && rd <= 1;
      const front = ly > 0; // near (front) half of the tilted ring

      const ringCell = (): Built => {
        const bandT = (rd - R_RIN) / (1 - R_RIN); // 0 inner .. 1 outer
        if (bandT > 0.44 && bandT < 0.58) return { kind: 'void', ch: ' ', lum: 0, ang: 0, pny: 0, baseLum: 0 }; // Cassini gap
        const edge = bandT < 0.1 || bandT > 0.92;
        const ch = edge ? '-' : bandT < 0.44 ? '=' : '≡';
        return {
          kind: front ? 'ringFront' : 'ringBack',
          ch,
          lum: front ? 0.85 : 0.4,
          ang: Math.atan2(rny, rnx),
          pny: 0,
          baseLum: 0,
        };
      };

      if (onRing && front) {
        cell = ringCell();
      } else if (onPlanet) {
        const nz = Math.sqrt(Math.max(0, 1 - pr2));
        const lambert = Math.max(0, pnx * LX + pny * LY + nz * LZ);
        const baseLum = Math.max(0, Math.min(1, 0.12 + lambert * 0.92));
        cell = { kind: 'planet', ch: rampChar(baseLum), lum: baseLum, ang: 0, pny, baseLum };
      } else if (onRing) {
        cell = ringCell(); // back arc (planet didn't cover it)
      } else {
        const sd = Math.hypot((x - CX) / 44, (y - CY) / 20);
        if (sd <= 1 && hash(x, y) > 0.955) {
          const h2 = hash(y, x);
          cell = { kind: 'star', ch: h2 > 0.55 ? '+' : '.', lum: 0.5, ang: 0, pny: 0, baseLum: 0 };
        }
      }

      row.push(cell);
    }
    grid.push(row);
  }
  return grid;
}

export default function SaturnOrb({ size = 320, className, style }: SaturnOrbProps) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const orbRef = useRef<HTMLDivElement | null>(null);
  const preRef = useRef<HTMLPreElement | null>(null);

  useEffect(() => {
    const pre = preRef.current;
    const orb = orbRef.current;
    const container = containerRef.current;
    if (!pre || !orb || !container) return;

    const prefersReduced =
      typeof window !== 'undefined' &&
      typeof window.matchMedia === 'function' &&
      window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    const model = buildModel();

    pre.innerHTML = '';
    const frag = document.createDocumentFragment();

    interface PlanetCell { el: HTMLElement; pny: number; baseLum: number; cls: string; ch: string; }
    interface RingCell { el: HTMLElement; ang: number; baseCls: string; cur: string; }
    interface StarCell { el: HTMLElement; phase: number; on: boolean; }

    const planetCells: PlanetCell[] = [];
    const ringCells: RingCell[] = [];
    const starCells: StarCell[] = [];

    for (let y = 0; y < ROWS; y++) {
      for (let x = 0; x < COLS; x++) {
        const m = model[y][x];
        const e = document.createElement('i');
        e.textContent = m.ch;
        switch (m.kind) {
          case 'planet': {
            const cls = planetBucket(m.baseLum);
            e.className = cls;
            planetCells.push({ el: e, pny: m.pny, baseLum: m.baseLum, cls, ch: m.ch });
            break;
          }
          case 'ringFront':
            e.className = 'sat-rf';
            ringCells.push({ el: e, ang: m.ang, baseCls: 'sat-rf', cur: 'sat-rf' });
            break;
          case 'ringBack':
            e.className = 'sat-rb';
            ringCells.push({ el: e, ang: m.ang, baseCls: 'sat-rb', cur: 'sat-rb' });
            break;
          case 'star':
            e.className = 'sat-star';
            starCells.push({ el: e, phase: hash(x, y) * Math.PI * 2, on: false });
            break;
          default:
            e.className = 'sat-void';
        }
        frag.appendChild(e);
      }
      frag.appendChild(document.createTextNode('\n'));
    }
    pre.appendChild(frag);

    const fontPx = 11;
    const gridW = COLS * fontPx * 0.6;
    orb.style.transform = `translate(-50%, -50%) scale(${size / gridW})`;

    if (prefersReduced) return;

    let raf = 0;
    let last = 0;
    let nextShoot = 2200;
    const shooters = new Set<HTMLElement>();

    function spawnShootingStar() {
      const rect = container!.getBoundingClientRect();
      if (rect.width === 0) return;
      const fromTop = Math.random() < 0.5;
      const sx = rect.width * (0.1 + Math.random() * 0.3);
      const sy = fromTop ? -6 : rect.height * (0.1 + Math.random() * 0.25);
      const ex = sx + rect.width * (0.45 + Math.random() * 0.3);
      const ey = sy + rect.height * (0.4 + Math.random() * 0.35);
      const d = document.createElement('div');
      d.className = 'sat-shoot';
      d.textContent = '✧';
      container!.appendChild(d);
      shooters.add(d);
      const a = d.animate(
        [
          { transform: `translate(${sx}px, ${sy}px) scale(.6)`, opacity: 0 },
          { opacity: 1, offset: 0.2 },
          { opacity: 1, offset: 0.7 },
          { transform: `translate(${ex}px, ${ey}px) scale(1)`, opacity: 0 },
        ],
        { duration: 850 + Math.random() * 400, easing: 'cubic-bezier(.4,.1,.7,.9)' },
      );
      a.onfinish = () => { d.remove(); shooters.delete(d); };
    }

    function tick(now: number) {
      if (now - last > 70) {
        last = now;
        const t = now / 1000;

        // Planet: drifting bands + subtle terminator shimmer.
        for (let i = 0; i < planetCells.length; i++) {
          const c = planetCells[i];
          const band = 0.07 * Math.sin(c.pny * Math.PI * 3.4 + t * 0.9);
          const flick = 0.03 * Math.sin(t * 1.7 + c.pny * 6.0 + i);
          const lum = Math.max(0, Math.min(1, c.baseLum * (1 + band) + flick));
          const cls = planetBucket(lum);
          if (cls !== c.cls) { c.el.className = cls; c.cls = cls; }
          const ch = rampChar(lum);
          if (ch !== c.ch) { c.el.textContent = ch; c.ch = ch; }
        }

        // Comet glint orbiting the rings (head + trailing tail).
        const glint = ((t * 1.0) % (Math.PI * 2)) - Math.PI;
        for (let i = 0; i < ringCells.length; i++) {
          const r = ringCells[i];
          let d = r.ang - glint;
          while (d > Math.PI) d -= Math.PI * 2;
          while (d < -Math.PI) d += Math.PI * 2;
          let cls = r.baseCls;
          if (Math.abs(d) < 0.16) cls = 'sat-rg';
          else if (d < 0 && d > -0.75) cls = 'sat-rg2';
          if (cls !== r.cur) { r.el.className = cls; r.cur = cls; }
        }

        // Twinkle stars.
        for (let i = 0; i < starCells.length; i++) {
          const s = starCells[i];
          const on = Math.sin(t * 2.4 + s.phase) > 0.35;
          if (on !== s.on) {
            s.el.className = on ? 'sat-star sat-star-on' : 'sat-star';
            s.on = on;
          }
        }

        // Occasional shooting star.
        if (now > nextShoot) {
          spawnShootingStar();
          nextShoot = now + 3200 + Math.random() * 4200;
        }
      }
      raf = requestAnimationFrame(tick);
    }
    raf = requestAnimationFrame(tick);

    return () => {
      cancelAnimationFrame(raf);
      shooters.forEach((d) => d.remove());
      shooters.clear();
    };
  }, [size]);

  return (
    <div
      ref={containerRef}
      className={className}
      style={{ width: size, height: size, position: 'relative', ...style }}
    >
      <div className="sat-halo sat-halo-outer" aria-hidden />
      <div className="sat-halo sat-halo-inner" aria-hidden />

      <div style={{ position: 'absolute', inset: 0, overflow: 'hidden' }}>
        <div
          ref={orbRef}
          style={{
            position: 'absolute',
            left: '50%',
            top: '50%',
            transformOrigin: 'center center',
            willChange: 'transform',
            pointerEvents: 'none',
          }}
        >
          <pre
            ref={preRef}
            className="sat-surface"
            style={{
              fontFamily: '"SF Mono", "JetBrains Mono", Menlo, Consolas, monospace',
              fontSize: 11,
              lineHeight: '11px',
              letterSpacing: 0,
              whiteSpace: 'pre',
              margin: 0,
              fontWeight: 700,
            }}
          />
        </div>
      </div>

      <style jsx global>{`
        @keyframes sat-halo-pulse {
          0%, 100% { transform: translate(-50%, -50%) scale(1) rotate(-17deg);    opacity: .28; }
          50%      { transform: translate(-50%, -50%) scale(1.05) rotate(-17deg); opacity: .46; }
        }
        .sat-halo {
          position: absolute; left: 50%; top: 50%; border-radius: 50%;
          transform: translate(-50%, -50%) rotate(-17deg); pointer-events: none;
          mix-blend-mode: screen; will-change: transform, opacity;
        }
        .sat-halo-inner {
          width: 60%; height: 60%;
          background: radial-gradient(circle at center,
            rgba(255,210,120,.32) 0%, rgba(255,190,90,.13) 42%, transparent 72%);
          filter: blur(18px);
          animation: sat-halo-pulse 5.5s ease-in-out infinite;
        }
        .sat-halo-outer {
          width: 160%; height: 95%;
          background: radial-gradient(ellipse at center,
            rgba(255,200,110,.11) 0%, rgba(255,200,110,.045) 42%, transparent 72%);
          filter: blur(34px);
          animation: sat-halo-pulse 9s ease-in-out infinite;
        }
        .sat-surface i { font-style: normal; display: inline; }
        .sat-surface .sat-void { color: transparent; }
        .sat-surface .sat-p0 { color: #3a2a12; }
        .sat-surface .sat-p1 { color: #6a4a1c; }
        .sat-surface .sat-p2 { color: #9a6a26; }
        .sat-surface .sat-p3 { color: #c89038; }
        .sat-surface .sat-p4 { color: #e8b85a; }
        .sat-surface .sat-p5 { color: #ffd884; text-shadow: 0 0 3px rgba(255,216,132,.4); }
        .sat-surface .sat-p6 { color: #fff0c4; text-shadow: 0 0 5px rgba(255,240,196,.6); }
        .sat-surface .sat-rb { color: #7a663e; }
        .sat-surface .sat-rf { color: #ffdf9e; text-shadow: 0 0 4px rgba(255,223,158,.45); }
        .sat-surface .sat-rg2 { color: #fff0c0; text-shadow: 0 0 6px rgba(255,235,170,.7); }
        .sat-surface .sat-rg {
          color: #fffbe6;
          text-shadow: 0 0 8px rgba(255,245,200,.95), 0 0 16px rgba(255,220,140,.65);
        }
        .sat-surface .sat-star { color: #5a6a82; }
        .sat-surface .sat-star-on {
          color: #d6e8ff; text-shadow: 0 0 5px rgba(200,228,255,.7);
        }
        .sat-shoot {
          position: absolute; left: 0; top: 0;
          font-family: "SF Mono", "JetBrains Mono", Menlo, Consolas, monospace;
          font-size: 13px; color: #eaf4ff;
          text-shadow: 0 0 6px rgba(210,232,255,.9), 0 0 14px rgba(150,200,255,.5);
          pointer-events: none; will-change: transform, opacity;
        }
        @media (prefers-reduced-motion: reduce) {
          .sat-halo-inner, .sat-halo-outer { animation: none; }
        }
      `}</style>
    </div>
  );
}
