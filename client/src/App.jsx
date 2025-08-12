import { BrowserRouter, Router, Route } from "react-router-dom";
import { Login } from "./page/Login";
import { Register } from "./page/Register";
import { Dashboard } from "./page/Dashboard";

import "./App.css";

export default function App() {
  return (
    <BrowserRouter>
      <Router>
        <Route path="/" element={<Login />}></Route>
        <Route path="/register" element={<Register />}></Route>
        <Route path="/dashboard" element={<Dashboard />}></Route>
      </Router>
    </BrowserRouter>
  );
}
