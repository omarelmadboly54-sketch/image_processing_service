import { initializeApp, getApps } from 'firebase/app';
import { getAuth, GoogleAuthProvider, signInWithPopup } from 'firebase/auth';
const config = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID,
};
export const firebaseConfigured = Boolean(config.apiKey && config.authDomain && config.projectId && config.appId);
export async function getGoogleIdToken() {
  if (!firebaseConfigured) throw new Error('Google sign-in is not configured. Add the Firebase web app settings to your .env file.');
  const app = getApps().length ? getApps()[0] : initializeApp(config);
  const auth = getAuth(app);
  const provider = new GoogleAuthProvider();
  provider.setCustomParameters({ prompt: 'select_account' });
  const credential = await signInWithPopup(auth, provider);
  return { idToken: await credential.user.getIdToken(), name: credential.user.displayName || credential.user.email?.split('@')[0] || 'there' };
}
