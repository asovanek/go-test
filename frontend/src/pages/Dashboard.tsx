import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { clearToken, getMe } from '../api/client'

export default function Dashboard() {
  const navigate = useNavigate()
  const [email, setEmail] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    getMe()
      .then((user) => {
        if (!cancelled) setEmail(user.email)
      })
      .catch(() => {
        if (!cancelled) navigate('/signin', { replace: true })
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [navigate])

  function signOut() {
    clearToken()
    navigate('/signin', { replace: true })
  }

  if (loading) {
    return (
      <div className="card">
        <p>Loading…</p>
      </div>
    )
  }

  return (
    <div className="card">
      <h1>Dashboard</h1>
      <p>Signed in as <strong>{email}</strong></p>
      <button type="button" onClick={signOut}>
        Sign out
      </button>
    </div>
  )
}
