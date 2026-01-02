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

  if (!user) return <div className="p-4">Loading...</div>;

  return (
    <div className="p-4">
      <h1 className="text-2xl font-semibold mb-2">Dashboard</h1>
      <p className="mb-4">Welcome, {user?.username}</p>
    </div>
  );
}
