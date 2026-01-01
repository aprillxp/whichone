import React, { useEffect, useState } from 'react';
import api from '../api/axios';

export default function Dashboard() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const data = async () => {
      try {
        const res = await api.get('/players/me');
        setUser(res.data);
      } catch (err) {
        if (err.response?.status === 401) {
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
      }
    };
    data();
  }, []);

  useEffect(() => {
    console.log('DASHBOARD MOUNTED');
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
