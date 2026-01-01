import React, { useState } from 'react';
import api from '../api/axios';
import { useNavigate } from 'react-router-dom';

export default function Register() {
  const [form, setForm] = useState({ username: '', password: '' });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const submit = async (e) => {
    e.preventDefault();
    try {
      setLoading(true);
      await api.post('/register', form);
      alert('Register is success');

      navigate('/login');
    } catch (error) {
      console.log(error.response);

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
        type="username"
        value={form.username}
        onChange={(e) => setForm({ ...form, username: e.target.value })}
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
