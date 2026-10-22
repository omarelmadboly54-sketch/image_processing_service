import { useCallback, useEffect, useState } from "react";
import Login from "./pages/Login.jsx";
import Signup from "./pages/Signup.jsx";
import ImageUpload from "./pages/ImageUpload.jsx";
import ImageLibrary from "./pages/ImageLibrary.jsx";
import Navbar from "./components/Navbar.jsx";
import HomePage from "./pages/HomePage.jsx";
import WorkspacePage from "./pages/WorkspacePage.jsx";
import ResultPage from "./pages/ResultPage.jsx";
export default function App() {
  const [page, setPage] = useState(() => sessionStorage.getItem("access_token") ? "home" : "login");
  const [name,setName] = useState(()=>sessionStorage.getItem("omarix_display_name")||"");
  const [navigationImage,setNavigationImage] = useState(null);
  const navigate = useCallback((target, image = null) => { setNavigationImage(image); setPage(target); }, []);
  const [theme,setTheme] = useState(()=>localStorage.getItem("omarix_theme")||"light");
  useEffect(()=>{document.documentElement.dataset.theme=theme;localStorage.setItem("omarix_theme",theme)},[theme]);
  const handleLogin = useCallback((token, displayName) => { sessionStorage.setItem("access_token", token); sessionStorage.setItem("omarix_display_name",displayName||"there");setName(displayName||"there");setPage("home"); }, []);
  const handleLogout = useCallback(() => { sessionStorage.removeItem("access_token");sessionStorage.removeItem("omarix_display_name"); setPage("login"); }, []);
  if (page === "signup") return <Signup onSwitchToLogin={() => setPage("login")}/>;
  if (page === "login") return <Login onLogin={handleLogin} onSwitchToSignup={() => setPage("signup")}/>;
  return <div className="omarix-shell"><Navbar onLogout={handleLogout} onNavigate={navigate} activePage={page} name={name} theme={theme} onToggleTheme={()=>setTheme(t=>t==="dark"?"light":"dark")}/><div className="omarix-shell-content">
    {page === "home" ? <HomePage onNavigate={navigate} name={name}/> : page === "upload" ? <ImageUpload onLogout={handleLogout} onNavigate={navigate} activePage="upload" hideNavbar/> : page === "library" ? <ImageLibrary onLogout={handleLogout} onNavigate={navigate} initialImage={navigationImage} hideNavbar/> : page === "result" ? <ResultPage image={navigationImage} onNavigate={navigate}/> : <WorkspacePage onNavigate={navigate} onLogout={handleLogout} initialImage={navigationImage}/>}
  </div></div>;
}
