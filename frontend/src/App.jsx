import { useCallback, useState } from "react";
import Login from "./pages/Login.jsx";
import Signup from "./pages/Signup.jsx";
import ImageUpload from "./pages/ImageUpload.jsx";
import ImageLibrary from "./pages/ImageLibrary.jsx";
import Navbar from "./components/Navbar.jsx";
import HomePage from "./pages/HomePage.jsx";
import WorkspacePage from "./pages/WorkspacePage.jsx";
export default function App() {
  const [page, setPage] = useState(() => sessionStorage.getItem("access_token") ? "home" : "login");
  const handleLogin = useCallback(token => { sessionStorage.setItem("access_token", token); setPage("home"); }, []);
  const handleLogout = useCallback(() => { sessionStorage.removeItem("access_token"); setPage("login"); }, []);
  if (page === "signup") return <Signup onSwitchToLogin={() => setPage("login")}/>;
  if (page === "login") return <Login onLogin={handleLogin} onSwitchToSignup={() => setPage("signup")}/>;
  return <div className="omarix-shell"><Navbar onLogout={handleLogout} onNavigate={setPage} activePage={page}/><div className="omarix-shell-content">
    {page === "home" ? <HomePage onNavigate={setPage}/> : page === "upload" ? <ImageUpload onLogout={handleLogout} onNavigate={setPage} activePage="upload" hideNavbar/> : page === "library" ? <ImageLibrary onLogout={handleLogout} onNavigate={setPage} hideNavbar/> : <WorkspacePage onNavigate={setPage}/>}
  </div></div>;
}
