import React, { useState } from 'react';
import api from '../api/axios';

export default function Login() {
  const [form, setForm] = useState({ email: '', password: '' });

  const submit = async (e) => {
    e.preventDefault();
    try {
      const { data } = await api.post('/login', form);
      const token = data?.token || data?.access_token;
      if (!token) throw new Error('token is empty');
      localStorage.setItem('token: ', token);
      window.location.href = '/dashboard';
    } catch (error) {
      alert(e?.response.data?.error || 'Login failed');
    }
  };

  return (
    <form onSubmit={submit} className="max-w-sm mx-auto">
      <h1 className="text-md font-semibold mb-3">Login</h1>
      <input
        className="input input-bordered w-full mb-2"
        placeholder="Email"
        type="email"
        value={form.email}
        onChange={(e) => setForm({ ...form, email: e.target.value })}
      />
      <input
        className="input input-bordered w-full mb-2"
        placeholder="Password"
        type="password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
      />
      <button className="btn btn-primary w-full">Login</button>
    </form>
  );
}
