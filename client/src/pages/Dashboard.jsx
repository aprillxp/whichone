import React, { useEffect, useState } from 'react';
import api from '../api/axios';

export default function Dashboard() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    api
      .get('/players/me')
      .then((res) => setUser(res.data))
      .catch(() => {});
  }, []);

  const logout = async () => {
    try {
      await api.post('/logout');
    } catch (_) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
  };
  if (!user) return <div className="p-4">Loading...</div>;

  return (
    <div className="p-4">
      <h1 className="text-2xl font-semibold mb-2">Dashboard</h1>
      <p className="mb-4">Welcome, {user?.username}</p>
      <div className="flex gap-2">
        <button className="btn btn-outline" onClick={logout}>
          Logout
        </button>
      </div>
    </div>
  );
}
