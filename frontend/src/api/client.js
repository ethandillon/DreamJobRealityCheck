// frontend/src/api/client.js
// Centralized API client that automatically attaches the API key header.
// NOTE: Exposing an API key in client-side code only provides light obfuscation.
// For true secrecy, move privileged operations server-side or use a backend proxy.

// Normalize base URL (remove any trailing slashes)
const rawBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const apiBase = rawBase.replace(/\/+$/, '');
const apiKey = import.meta.env.VITE_API_KEY; // define in your frontend .env (VITE_*)

function buildQuery(params) {
  const search = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') search.append(k, v);
  });
  return search.toString();
}

async function request(path, { method = 'GET', query, body, signal } = {}) {
  const url = query ? `${apiBase}${path}?${buildQuery(query)}` : `${apiBase}${path}`;
  const headers = new Headers();
  if (apiKey) headers.set('X-API-Key', apiKey);
  if (body) headers.set('Content-Type', 'application/json');

  const res = await fetch(url, { method, headers, body: body ? JSON.stringify(body) : undefined, signal });
  if (!res.ok) {
    let detail = '';
    try { detail = JSON.stringify(await res.json()); } catch (_) {}
    throw new Error(`API ${res.status} ${res.statusText} ${detail}`.trim());
  }
  return res.json();
}

export async function calculate(filters) {
  return request('/api/calculate', { query: filters });
}

export async function health() { return request('/api/health'); }

export async function getOccupations() { return request('/api/occupations'); }
export async function getStates() { return request('/api/states'); }
export async function getAreasByState(state, { signal } = {}) { return request('/api/areas-by-state', { query: { state }, signal }); }

export { apiBase };
