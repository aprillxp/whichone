import React, { useState } from "react";
import api from "../api/axios";

export const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await api.post("/login", { email, password });
      alert("Login OK");
      window.location.href = "/dashboard";
    } catch (error) {
      console.log(error, "< login");
      alert("Login NOT OK");
    }
  };
  return (
    <div className="max-w-sm mx-auto mt-10 p-4 shadow rounded bg-base-200">
      <h1 className="text-xl font-bold mb-4">Login</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          placeholder="Email"
          className="input input-bordered w-full mb-2"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          className="input input-bordered w-full mb-4"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button className="btn btn-primary w-full">Login</button>
      </form>
    </div>
  );
};
