import AuthPage from "./AuthPage.jsx";

export default function Login({ onLogin, onSwitchToSignup }) {
  return <AuthPage mode="login" onLogin={onLogin} onSwitchToSignup={onSwitchToSignup} onSwitchToLogin={() => {}} />;
}