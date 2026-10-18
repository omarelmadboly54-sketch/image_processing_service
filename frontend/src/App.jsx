import { useState } from "react";
import Login from "./pages/Login.jsx";
import Signup from "./pages/Signup.jsx";
import ImageUpload from "./pages/ImageUpload.jsx";

export default function App() {
  const [page, setPage] = useState(() =>
    sessionStorage.getItem("access_token") ? "upload" : "login"
  );

  function handleLogin(token) {
    sessionStorage.setItem("access_token", token);
    setPage("upload");
  }

  function handleLogout() {
    sessionStorage.removeItem("access_token");
    setPage("login");
  }

  if (page === "upload") {
    return <ImageUpload onLogout={handleLogout} />;
  }

  if (page === "signup") {
    return <Signup onSwitchToLogin={() => setPage("login")} />;
  }

  return (
    <Login
      onLogin={handleLogin}
      onSwitchToSignup={() => setPage("signup")}
    />
  );
}
