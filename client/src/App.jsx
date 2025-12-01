import { Link, Outlet } from 'react-router-dom';

import React from 'react';

export default function App() {
  return (
    <div className="max-w-3xl mx-auto p-4">
      <nav className="flex gap-4 mb-6">
        <Link to="/login" className="link">
          Login
        </Link>
        <Link to="/register" className="link">
          Register
        </Link>
        <Link to="/dashboard" className="link">
          Dashboard
        </Link>
      </nav>
      <Outlet />
    </div>
  );
}
