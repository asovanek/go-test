const TOKEN_KEY = 'auth_token'

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export type AuthUser = {
  id: string
  email: string
}

export type AuthResponse = {
  token: string
  user: AuthUser
}

export type ApiError = {
  error?: string
  errors?: Record<string, string>
}

function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  if (!headers.has('Content-Type') && init.body) {
    headers.set('Content-Type', 'application/json')
  }
  const token = getToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const res = await fetch(`${baseURL}${path}`, { ...init, headers })
  const data = (await res.json().catch(() => ({}))) as T & ApiError

  if (!res.ok) {
    const message =
      data.error ||
      (data.errors ? Object.values(data.errors).join(', ') : 'request failed')
    throw new Error(message)
  }

  return data as T
}

export function signup(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>('/api/v1/auth/signup', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function signin(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>('/api/v1/auth/signin', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function getMe(): Promise<AuthUser> {
  return request<AuthUser>('/api/v1/me')
}
