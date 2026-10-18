import AuthPage from "./AuthPage.jsx";

export default function Signup({ onSwitchToLogin }) {
  return <AuthPage mode="signup" onLogin={() => {}} onSwitchToSignup={() => {}} onSwitchToLogin={onSwitchToLogin} />;
}
