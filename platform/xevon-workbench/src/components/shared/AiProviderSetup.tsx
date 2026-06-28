'use client';

import { useEffect, useMemo, useState } from 'react';
import { useConfig, useUpdateConfig } from '@/api/hooks';
import { apiPost, ApiError } from '@/api/client';
import { useToast } from '@/contexts/ToastContext';
import type { ConfigEntry } from '@/api/types';

/**
 * AiProviderSetup — friendly dropdown configurator for the agentic-scan AI
 * provider (olium). Lets a user pick a provider + model instead of hand-typing
 * raw config keys. Writes the matching `agent.olium.*` keys on save. Styled with
 * the shared CSS variables so it renders correctly in both dark and light.
 */

interface Preset {
  name: string;
  base_url: string;
  model: string;
  needsKey: boolean;
  keyPlaceholder?: string;
}

interface ProviderDef {
  value: string;
  label: string;
  kind: 'apikey' | 'compatible';
  models?: string[];
  keyPlaceholder?: string;
  keyHelp?: string;
  presets?: Preset[];
}

const PROVIDERS: ProviderDef[] = [
  {
    value: 'anthropic-api-key',
    label: 'Anthropic — Claude (API key)',
    kind: 'apikey',
    models: ['claude-opus-4-8', 'claude-sonnet-4-6', 'claude-haiku-4-5-20251001'],
    keyPlaceholder: 'sk-ant-api03-...',
    keyHelp: 'Get a key at console.anthropic.com',
  },
  {
    value: 'openai-api-key',
    label: 'OpenAI — GPT (API key)',
    kind: 'apikey',
    models: ['gpt-5.5', 'gpt-4o', 'gpt-4o-mini'],
    keyPlaceholder: 'sk-...',
    keyHelp: 'Get a key at platform.openai.com',
  },
  {
    value: 'openai-compatible',
    label: 'Gemini / Ollama / Custom (OpenAI-compatible)',
    kind: 'compatible',
    presets: [
      {
        name: 'Google Gemini',
        base_url: 'https://generativelanguage.googleapis.com/v1beta/openai/chat/completions',
        model: 'gemini-2.5-flash',
        needsKey: true,
        keyPlaceholder: 'AIza...',
      },
      {
        name: 'Local Ollama (free)',
        base_url: 'http://localhost:11434/v1/chat/completions',
        model: 'qwen2.5:14b',
        needsKey: false,
      },
      {
        name: 'OpenRouter',
        base_url: 'https://openrouter.ai/api/v1/chat/completions',
        model: 'anthropic/claude-3.5-sonnet',
        needsKey: true,
        keyPlaceholder: 'sk-or-...',
      },
      { name: 'Custom…', base_url: '', model: '', needsKey: true },
    ],
  },
];

const inputStyle: React.CSSProperties = {
  backgroundColor: 'var(--v-bg)',
  borderColor: 'var(--v-border)',
  color: 'var(--v-text)',
};

export default function AiProviderSetup() {
  const { data } = useConfig('olium');
  const update = useUpdateConfig();
  const { toast } = useToast();

  const cfg = useMemo(() => {
    const m: Record<string, string> = {};
    (data?.entries ?? []).forEach((e) => { m[e.key] = e.value; });
    return m;
  }, [data]);

  const [provider, setProvider] = useState('anthropic-api-key');
  const [model, setModel] = useState('claude-opus-4-8');
  const [apiKey, setApiKey] = useState('');
  const [presetName, setPresetName] = useState('Google Gemini');
  const [baseUrl, setBaseUrl] = useState('');
  const [modelId, setModelId] = useState('');
  const [hydrated, setHydrated] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testMsg, setTestMsg] = useState<{ ok: boolean; text: string } | null>(null);

  // Pre-fill provider/model from existing config once it loads (never the key).
  useEffect(() => {
    if (hydrated || !data) return;
    const p = cfg['agent.olium.provider'] || 'anthropic-api-key';
    const def = PROVIDERS.find((x) => x.value === p);
    if (def) {
      setProvider(p);
      if (def.kind === 'apikey') {
        setModel(cfg['agent.olium.model'] || def.models?.[0] || '');
      } else {
        const bu = cfg['agent.olium.custom_provider.base_url'] || '';
        const mid = cfg['agent.olium.custom_provider.model_id'] || '';
        setBaseUrl(bu);
        setModelId(mid);
        const match = def.presets?.find((pr) => pr.base_url === bu);
        if (match) setPresetName(match.name);
        else if (bu) setPresetName('Custom…');
      }
    }
    setHydrated(true);
  }, [data, cfg, hydrated]);

  const def = PROVIDERS.find((x) => x.value === provider) ?? PROVIDERS[0];
  const preset = def.presets?.find((p) => p.name === presetName);
  const compatNeedsKey = def.kind === 'compatible' && (preset?.needsKey ?? true);

  function onProviderChange(v: string) {
    setProvider(v);
    const d = PROVIDERS.find((x) => x.value === v)!;
    if (d.kind === 'apikey') setModel(d.models?.[0] || '');
  }

  function onPresetChange(name: string) {
    setPresetName(name);
    const pr = def.presets?.find((p) => p.name === name);
    if (pr) { setBaseUrl(pr.base_url); setModelId(pr.model); }
  }

  function buildEntries(): ConfigEntry[] {
    const entries: ConfigEntry[] = [{ key: 'agent.olium.provider', value: provider }];
    if (def.kind === 'apikey') {
      entries.push({ key: 'agent.olium.model', value: model });
      if (apiKey.trim()) entries.push({ key: 'agent.olium.llm_api_key', value: apiKey.trim() });
    } else {
      entries.push({ key: 'agent.olium.model', value: '' }); // let model_id win
      entries.push({ key: 'agent.olium.custom_provider.base_url', value: baseUrl.trim() });
      entries.push({ key: 'agent.olium.custom_provider.model_id', value: modelId.trim() });
      if (apiKey.trim()) entries.push({ key: 'agent.olium.custom_provider.api_key', value: apiKey.trim() });
    }
    return entries;
  }

  async function save() {
    try {
      await update.mutateAsync(buildEntries());
      toast('AI provider saved — run an Agentic Scan to use it', 'success');
      setApiKey('');
    } catch (e) {
      toast('Failed to save: ' + ((e as Error)?.message ?? 'error'), 'error');
    }
  }

  // Save the chosen settings (applied live), then ping the provider with a tiny
  // prompt via the OpenAI-compatible endpoint to confirm the key/model work.
  async function testConnection() {
    setTesting(true);
    setTestMsg(null);
    try {
      await update.mutateAsync(buildEntries());
      await apiPost('/api/agent/chat/completions', {
        messages: [{ role: 'user', content: 'Reply with the single word: OK' }],
        max_tokens: 8,
      });
      setTestMsg({ ok: true, text: '✓ Connected — the provider responded successfully.' });
    } catch (e) {
      const err = e as ApiError;
      const detail = err?.details ? ' — ' + err.details : '';
      setTestMsg({ ok: false, text: '✗ ' + (err?.message || 'connection failed') + detail });
    } finally {
      setTesting(false);
    }
  }

  const labelCls = 'text-[10px] uppercase block mb-0.5';
  const labelStyle: React.CSSProperties = { color: 'var(--v-text-muted)' };

  return (
    <div className="border" style={{ borderColor: 'var(--v-border)', backgroundColor: 'var(--v-bg)' }}>
      <div className="px-3 py-1.5 border-b flex items-center justify-between" style={{ borderColor: 'var(--v-border)' }}>
        <span className="text-xs font-bold" style={{ color: 'var(--v-accent)' }}>
          AI PROVIDER · for Agentic Scan (Autopilot / Swarm)
        </span>
      </div>

      <div className="p-3 space-y-3">
        {/* Provider */}
        <div>
          <label className={labelCls} style={labelStyle}>Provider</label>
          <select
            value={provider}
            onChange={(e) => onProviderChange(e.target.value)}
            className="w-full border text-xs px-2 py-1 focus:outline-none"
            style={inputStyle}
          >
            {PROVIDERS.map((p) => (
              <option key={p.value} value={p.value} style={{ backgroundColor: 'var(--v-surface)', color: 'var(--v-text)' }}>
                {p.label}
              </option>
            ))}
          </select>
        </div>

        {def.kind === 'apikey' ? (
          <>
            {/* Model */}
            <div>
              <label className={labelCls} style={labelStyle}>Model</label>
              <select
                value={model}
                onChange={(e) => setModel(e.target.value)}
                className="w-full border text-xs px-2 py-1 focus:outline-none"
                style={inputStyle}
              >
                {def.models?.map((m) => (
                  <option key={m} value={m} style={{ backgroundColor: 'var(--v-surface)', color: 'var(--v-text)' }}>{m}</option>
                ))}
              </select>
            </div>
            {/* API key */}
            <div>
              <label className={labelCls} style={labelStyle}>API key</label>
              <input
                type="password"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                placeholder={def.keyPlaceholder}
                autoComplete="off"
                className="w-full border text-xs px-2 py-1 focus:outline-none font-mono"
                style={inputStyle}
              />
              <span className="text-[10px]" style={labelStyle}>
                {def.keyHelp} · leave blank to keep the existing key
              </span>
            </div>
          </>
        ) : (
          <>
            {/* Compatible preset */}
            <div>
              <label className={labelCls} style={labelStyle}>Service</label>
              <select
                value={presetName}
                onChange={(e) => onPresetChange(e.target.value)}
                className="w-full border text-xs px-2 py-1 focus:outline-none"
                style={inputStyle}
              >
                {def.presets?.map((p) => (
                  <option key={p.name} value={p.name} style={{ backgroundColor: 'var(--v-surface)', color: 'var(--v-text)' }}>{p.name}</option>
                ))}
              </select>
            </div>
            <div className="flex gap-2">
              <div className="flex-1">
                <label className={labelCls} style={labelStyle}>Endpoint URL</label>
                <input
                  type="text"
                  value={baseUrl}
                  onChange={(e) => setBaseUrl(e.target.value)}
                  placeholder="https://.../chat/completions"
                  className="w-full border text-xs px-2 py-1 focus:outline-none font-mono"
                  style={inputStyle}
                />
              </div>
              <div className="w-44">
                <label className={labelCls} style={labelStyle}>Model</label>
                <input
                  type="text"
                  value={modelId}
                  onChange={(e) => setModelId(e.target.value)}
                  placeholder="gemini-2.5-flash"
                  className="w-full border text-xs px-2 py-1 focus:outline-none font-mono"
                  style={inputStyle}
                />
              </div>
            </div>
            {compatNeedsKey && (
              <div>
                <label className={labelCls} style={labelStyle}>API key {preset && !preset.needsKey ? '(optional)' : ''}</label>
                <input
                  type="password"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder={preset?.keyPlaceholder || 'API key (optional for local)'}
                  autoComplete="off"
                  className="w-full border text-xs px-2 py-1 focus:outline-none font-mono"
                  style={inputStyle}
                />
                <span className="text-[10px]" style={labelStyle}>leave blank to keep the existing key</span>
              </div>
            )}
          </>
        )}

        <div className="flex items-center gap-2 pt-1 flex-wrap">
          <button
            onClick={save}
            disabled={update.isPending || testing}
            className="text-xs px-4 py-1 border font-semibold disabled:opacity-50"
            style={{ borderColor: 'var(--v-accent)', color: 'var(--v-accent)', backgroundColor: 'color-mix(in srgb, var(--v-accent) 10%, transparent)' }}
          >
            {update.isPending && !testing ? 'saving...' : 'Save AI provider'}
          </button>
          <button
            onClick={testConnection}
            disabled={testing || update.isPending}
            className="text-xs px-4 py-1 border font-semibold disabled:opacity-50"
            style={{ borderColor: 'var(--v-border)', color: 'var(--v-text)', backgroundColor: 'transparent' }}
          >
            {testing ? 'testing...' : 'Test connection'}
          </button>
          <span className="text-[10px]" style={labelStyle}>
            Test saves the settings, then pings the model to verify the key works.
          </span>
        </div>

        {testMsg && (
          <div
            className="text-[11px] px-2 py-1 border break-words"
            style={{
              color: testMsg.ok ? 'var(--v-success, #7fd962)' : 'var(--v-error, #ef2f27)',
              borderColor: testMsg.ok ? 'color-mix(in srgb, var(--v-success, #7fd962) 40%, transparent)' : 'color-mix(in srgb, var(--v-error, #ef2f27) 40%, transparent)',
              backgroundColor: testMsg.ok ? 'color-mix(in srgb, var(--v-success, #7fd962) 8%, transparent)' : 'color-mix(in srgb, var(--v-error, #ef2f27) 8%, transparent)',
            }}
          >
            {testMsg.text}
          </div>
        )}
      </div>
    </div>
  );
}
