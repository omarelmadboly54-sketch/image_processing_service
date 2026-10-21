import { Home, Images, LogOut, UploadCloud, WandSparkles } from "lucide-react";
const links = [
  { id: "home", label: "Home", icon: Home },
  { id: "upload", label: "Upload image", icon: UploadCloud },
  { id: "library", label: "My gallery", icon: Images },
  { id: "workspace", label: "Workspace", icon: WandSparkles },
];
export default function Navbar({ onLogout, onNavigate, activePage = "home" }) {
  return <aside className="omarix-sidebar" aria-label="Main navigation">
    <button className="sidebar-brand" onClick={() => onNavigate("home")} type="button"><img src="/omarix-logo.jfif" alt="Omarix" /><span>omarix<small>IMAGE PROCESSING</small></span></button>
    <div className="sidebar-caption">YOUR SPACE</div>
    <nav className="sidebar-links">{links.map(({ id, label, icon: Icon }) => <button key={id} type="button" className={`sidebar-link ${activePage === id ? "active" : ""}`} aria-current={activePage === id ? "page" : undefined} onClick={() => onNavigate(id)}><Icon size={19}/><span>{label}</span>{id === "workspace" && <small>Soon</small>}</button>)}</nav>
    <div className="sidebar-bottom"><button type="button" className="sidebar-signout" onClick={onLogout}><LogOut size={18}/> Sign out</button></div>
  </aside>;
}

