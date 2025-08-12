import React, { useEffect, useState } from "react";
import api from "../api/axios";

export const Dashboard = () => {
  const [user, setUser] = useState("");

  useEffect(() => {
    api
      .get("/players/me")
      .then((res) => setUser(res.data))
      .catch(() => {
        alert("NOT login yet");
        window.location.href = "/";
      });
  }, []);

  if (!user) {
    return <p className="p-4">Loading...</p>;
  }
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Dashboard</h1>
      <p>Welcome, {user.email}</p>
    </div>
  );
};
