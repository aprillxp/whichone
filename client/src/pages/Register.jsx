import React, { useState } from 'react';
import api from '../api/axios';

export default function Register() {
  const [form, setForm] = useState({ email: '', password: '' });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await api.post('/register', form);
      alert('Register is success');
      window.location.href = '/login';
    } catch (error) {
      alert(e?.response?.data?.error || 'Register is failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={submit} className="max-w-sm mx-auto">
      <h1 className="text-xl font-bold mb-4">Register</h1>
      <input
        className="input input-bordered w-full mb-2"
        placeholder="Email"
        type="email"
        value={form.email}
        onChange={(e) => setForm({ ...form, email: e.target.value })}
      />
      <input
        className="input input-bordered w-full mb-4"
        placeholder="Password"
        type="password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
      />
      <button className="btn btn-primary w-full" disabled={loading}>
        {loading ? '...' : 'Register'}
      </button>
    </form>
  );
}
