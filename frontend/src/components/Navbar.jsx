import { Image as ImageIcon, LogOut } from "lucide-react";

export default function Navbar({ onLogout }) {
  return (
    <header className="app-navbar">
      <a className="app-brand" href="#top" onClick={(event) => event.preventDefault()}>
        <span className="brand-icon"><ImageIcon size={20} strokeWidth={2.4} /></span>
        <span>Omarix</span>
      </a>
      <div className="nav-right"><span className="nav-status"><i /> Workspace</span><button className="logout-button" type="button" onClick={onLogout}><LogOut size={16} /> Sign out</button></div>
    </header>
  );
}
