// frontend/src/api/client.js
// Centralized API client for building query strings & handling JSON responses.

// Normalize base URL (remove any trailing slashes)
const rawBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const apiBase = rawBase.replace(/\/+$/, '');

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
  if (body) headers.set('Content-Type', 'application/json');

  const res = await fetch(url, { method, headers, body: body ? JSON.stringify(body) : undefined, signal });
  if (!res.ok) {
    let detail = '';
    const text = await res.text().catch(() => '');
    detail = text || '';
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
