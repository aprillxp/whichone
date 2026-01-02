import React from 'react';
import { Link, useNavigate, useNavigation } from 'react-router-dom';

export default function Navbar() {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();

  const logout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };
  return (
    <div className="navbar bg-base-100 border-b border-base-300 px-6">
      <div className="flex-1">
        <Link to="/" className="text-lg font-bold ">
          Whichone.
        </Link>
      </div>

      <div className="flex gap-2">
        {token ? (
          <>
            <Link to="/dashboard" className="btn btn-ghost btn-sm">
              Dashboard
            </Link>
            <button className="btn btn-outline btn-sm" onClick={logout}>
              Logout
            </button>
          </>
        ) : (
          <>
            <Link to="/login" className="btn btn-ghost btn-sm">
              Login
            </Link>
            <Link to="/register" className="btn btn-indigo btn-sm">
              Register
            </Link>
          </>
        )}
      </div>
    </div>
  );
}
