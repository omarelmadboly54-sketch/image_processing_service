
import React, { useState } from "react";
import { createRoot } from "react-dom/client";
import {
  ArrowRight,
  Eye,
  EyeOff,
  Image as ImageIcon,
  Sparkles,
} from "lucide-react";
import "./styles.css";

const API_BASE_URL = (
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"
).replace(/\/+$/, "");

function AuthApp() {
  const [mode, setMode] = useState("login");
  const [showPassword, setShowPassword] = useState(false);
  const [busy, setBusy] = useState(false);
  const [notice, setNotice] = useState({ type: "", text: "" });
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
  });

  const isLogin = mode === "login";

  const update = (event) => {
    setForm((previous) => ({
      ...previous,
      [event.target.name]: event.target.value,
    }));
    setNotice({ type: "", text: "" });
  };

  async function handleSubmit(event) {
    event.preventDefault();
    setNotice({ type: "", text: "" });
    setBusy(true);

    
    const endpoint = isLogin ? "/api/login" : "/api/users";

    const payload = isLogin
      ? {
          email: form.email.trim(),
          password: form.password,
        }
      : {
          name: form.name.trim(),
          email: form.email.trim(),
          password: form.password,
        };

    try {
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      const raw = await response.text();
      let data = {};

      try {
        data = raw ? JSON.parse(raw) : {};
      } catch {
        data = { message: raw };
      }

      if (!response.ok) {
        throw new Error(
          data.error ||
            data.message ||
            `Request failed with status ${response.status}`
        );
      }

      if (isLogin) {
        if (data.access_token) {
          sessionStorage.setItem("access_token", data.access_token);
        }

        setNotice({
          type: "success",
          text: "Signed in successfully.",
        });
      } else {
        setNotice({
          type: "success",
          text: "Account created successfully. You can now sign in.",
        });

        setMode("login");
        setForm((previous) => ({
          ...previous,
          password: "",
        }));
      }
    } catch (error) {
      const isNetworkError = error instanceof TypeError;

      setNotice({
        type: "error",
        text: isNetworkError
          ? `Could not complete the request to ${API_BASE_URL}${endpoint}. Check that the Go backend is running and that CORS allows http://localhost:5173.`
          : error.message,
      });
    } finally {
      setBusy(false);
    }
  }

  return (
    <main className="auth-page">
      <section className="visual-side" aria-label="Omarix introduction">
        <a
          className="brand"
          href="#"
          onClick={(event) => event.preventDefault()}
        >
          <span className="brand-icon">
            <ImageIcon size={21} strokeWidth={2.4} />
          </span>
          <span>Omarix</span>
        </a>

        <div className="visual-copy">
          <span className="kicker">
            <Sparkles size={14} /> YOUR CREATIVE SPACE
          </span>
          <h1>
            Every image has
            <br />a <span>new perspective.</span>
          </h1>
          <p>
            A simple workspace for your image processing journey. Sign in and
            pick up where your creativity begins.
          </p>
        </div>

        <div className="artwork" aria-hidden="true">
          <div className="art-glow" />
          <div className="art-orbit orbit-one" />
          <div className="art-orbit orbit-two" />
          <div className="art-card card-back">
            <div className="abstract-shape shape-one" />
            <div className="abstract-shape shape-two" />
          </div>
          <div className="art-card card-front">
            <div className="image-sun" />
            <div className="image-mountain mountain-one" />
            <div className="image-mountain mountain-two" />
            <div className="image-spark">✦</div>
          </div>
          <span className="floating-star star-one">✳</span>
          <span className="floating-star star-two">✦</span>
        </div>

        <div className="visual-footer">
          <span>DESIGNED FOR YOUR NEXT IDEA</span>
          <span className="footer-dots">
            <i />
            <i />
            <i />
          </span>
        </div>
      </section>

      <section className="form-side">
        <div className="mobile-brand">
          <span className="brand-icon">
            <ImageIcon size={20} />
          </span>
          <span>Omarix</span>
        </div>

        <div className="auth-card">
          <div className="form-heading">
            <span className="form-kicker">
              {isLogin ? "WELCOME BACK" : "GET STARTED"}
            </span>
            <h2>
              {isLogin ? "Sign in to your account" : "Create your account"}
            </h2>
            <p>
              {isLogin
                ? "Enter your details below to continue."
                : "A few details and you're ready to go."}
            </p>
          </div>

          <div className="mode-switch" role="tablist" aria-label="Account action">
            <button
              type="button"
              role="tab"
              aria-selected={isLogin}
              className={isLogin ? "selected" : ""}
              onClick={() => {
                setMode("login");
                setNotice({ type: "", text: "" });
              }}
            >
              Sign In
            </button>

            <button
              type="button"
              role="tab"
              aria-selected={!isLogin}
              className={!isLogin ? "selected" : ""}
              onClick={() => {
                setMode("signup");
                setNotice({ type: "", text: "" });
              }}
            >
              Sign Up
            </button>
          </div>

          <form onSubmit={handleSubmit}>
            {!isLogin && (
              <label className="input-label">
                Full name
                <input
                  name="name"
                  type="text"
                  autoComplete="name"
                  placeholder="Your name"
                  value={form.name}
                  onChange={update}
                  required
                />
              </label>
            )}

            <label className="input-label">
              Email address
              <input
                name="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                value={form.email}
                onChange={update}
                required
              />
            </label>

            <label className="input-label">
              Password
              <span className="password-wrap">
                <input
                  name="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete={isLogin ? "current-password" : "new-password"}
                  placeholder="Enter your password"
                  value={form.password}
                  onChange={update}
                  required
                  minLength={6}
                />
                <button
                  className="eye-button"
                  type="button"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                  onClick={() => setShowPassword((value) => !value)}
                >
                  {showPassword ? <EyeOff size={17} /> : <Eye size={17} />}
                </button>
              </span>
            </label>

            {isLogin && (
              <div className="form-options">
                <span className="secure-note">
                  <span className="tiny-lock">●</span> Secure sign in
                </span>
                <span className="password-hint">Use your registered email</span>
              </div>
            )}

            <button className="submit-button" type="submit" disabled={busy}>
              {busy
                ? "Please wait..."
                : isLogin
                  ? "Sign In"
                  : "Create account"}
              <ArrowRight size={17} />
            </button>

            {notice.text && (
              <div className={`notice ${notice.type}`} role="status">
                {notice.text}
              </div>
            )}
          </form>

          <p className="switch-prompt">
            {isLogin ? "Don't have an account?" : "Already have an account?"}{" "}
            <button
              type="button"
              onClick={() => {
                setMode(isLogin ? "signup" : "login");
                setNotice({ type: "", text: "" });
              }}
            >
              {isLogin ? "Sign up" : "Sign in"}
            </button>
          </p>

          <div className="privacy-note">
            <span className="privacy-dot" /> Your information is sent to your
            configured backend.
          </div>
        </div>

        <div className="copyright">
          © {new Date().getFullYear()} Omarix <span>·</span> IMAGE PROCESSING
          SERVICE
        </div>
      </section>
    </main>
  );
}

createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <AuthApp />
  </React.StrictMode>
);
