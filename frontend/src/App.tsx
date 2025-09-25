// frontend/src/App.tsx
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Home from "./pages/Home";
import CreateApp from "./pages/apps/CreateApp";
import AppDetail from "./pages/apps/AppDetail";
import Translator from "./pages/apps/Translator";

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/login" />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<Home />} />
        <Route path="/apps/create" element={<CreateApp />} />
        <Route path="/apps/translator" element={<Translator />} />
        <Route path="/apps/:id" element={<AppDetail />} />
      </Routes>
    </Router>
  );
}
