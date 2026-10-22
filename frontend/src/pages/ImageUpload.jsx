import { useEffect, useRef, useState } from "react";
import { ArrowUpFromLine, CheckCircle2, Image as ImageIcon, LoaderCircle, UploadCloud } from "lucide-react";
import Navbar from "../components/Navbar.jsx";
import ImagePreview, { formatBytes } from "../components/ImagePreview.jsx";

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || "http://localhost:8080").replace(/\/+$/, "");
const MAX_FILE_SIZE = 5 * 1024 * 1024;
const ALLOWED_EXTENSIONS = /\.(jpe?g|png|webq)$/i;

export default function ImageUpload({ onLogout, onNavigate, activePage = "upload" }) {
  const inputRef = useRef(null);
  const [file, setFile] = useState(null);
  const [previewUrl, setPreviewUrl] = useState("");
  const [dragging, setDragging] = useState(false);
  const [busy, setBusy] = useState(false);
  const [notice, setNotice] = useState({ type: "", text: "" });
  const [uploadedImage, setUploadedImage] = useState(null);

  useEffect(() => () => { if (previewUrl) URL.revokeObjectURL(previewUrl); }, [previewUrl]);

  function chooseFile(nextFile) {
    setNotice({ type: "", text: "" });
    setUploadedImage(null);
    if (!nextFile) return;
    if (!ALLOWED_EXTENSIONS.test(nextFile.name)) {
      setNotice({ type: "error", text: "Unsupported file type. Choose a JPG, JPEG, or PNG image." });
      return;
    }
    if (nextFile.size > MAX_FILE_SIZE) {
      setNotice({ type: "error", text: "This image is too large. The maximum file size is 5 MB." });
      return;
    }
    setFile(nextFile);
    setPreviewUrl(URL.createObjectURL(nextFile));
  }

  function clearSelection() {
    setFile(null);
    setPreviewUrl("");
    if (inputRef.current) inputRef.current.value = "";
    setNotice({ type: "", text: "" });
  }

  async function uploadImage(event) {
    event.preventDefault();
    if (!file) {
      setNotice({ type: "error", text: "Choose an image before uploading." });
      return;
    }
    const accessToken = sessionStorage.getItem("access_token");
    if (!accessToken) {
      onLogout();
      return;
    }

    setBusy(true);
    setNotice({ type: "", text: "" });
    const formData = new FormData();
    formData.append("image", file);

    try {
      const response = await fetch(`${API_BASE_URL}/api/images`, {
        method: "POST",
        headers: { Authorization: `Bearer ${accessToken}` },
        body: formData,
      });
      const raw = await response.text();
      let data = {};
      try { data = raw ? JSON.parse(raw) : {}; } catch { data = { message: raw }; }
      if (response.status === 401) {
        sessionStorage.removeItem("access_token");
        onLogout();
        throw new Error("Your session expired. Please sign in again.");
      }
      if (!response.ok) throw new Error(data.error || data.message || `Upload failed (${response.status})`);
      setUploadedImage(data.data || null);
      setFile(null);
      setPreviewUrl("");
      if (inputRef.current) inputRef.current.value = "";
      setNotice({ type: "success", text: data.message || "Image uploaded successfully." });
    } catch (error) {
      setNotice({ type: "error", text: error instanceof TypeError ? `Could not connect to ${API_BASE_URL}. Check that your Go backend is running.` : error.message });
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="workspace-page" id="top">
      <main className="workspace-main">
        <div className="workspace-heading"><span className="workspace-eyebrow"><ImageIcon size={14} /> OMARIX WORKSPACE</span><h1>Bring your images <span>to life.</span></h1><p>Upload an image to securely store it in your Omarix library.</p></div>

        {notice.text && <div className={`notice ${notice.type} upload-feedback`} role="status">{notice.text}</div>}
        {uploadedImage && <section className="upload-success-card"><img src={uploadedImage.file_url} alt={uploadedImage.original_name || "Recently uploaded"}/><div><span className="workspace-eyebrow"><CheckCircle2 size={14}/> UPLOAD COMPLETE</span><h3>{uploadedImage.original_name || "Your image is ready"}</h3><p>Your image has been added to your gallery.</p><button type="button" className="library-primary-button" onClick={() => onNavigate("library")}>View gallery →</button></div></section>}

        <section className="upload-panel">
          <div className="panel-title"><div className="panel-title-icon"><UploadCloud size={21} /></div><div><h2>Upload an image</h2><p>Select a file from your device to get started.</p></div></div>
          <form onSubmit={uploadImage}>
            <div className={`dropzone ${dragging ? "is-dragging" : ""} ${file ? "has-file" : ""}`} onDragOver={(event) => { event.preventDefault(); setDragging(true); }} onDragLeave={(event) => { event.preventDefault(); setDragging(false); }} onDrop={(event) => { event.preventDefault(); setDragging(false); chooseFile(event.dataTransfer.files?.[0]); }}>
              {!file ? <>
                <div className="dropzone-icon"><UploadCloud size={29} strokeWidth={1.8} /></div>
                <h3>Drag and drop your image here</h3><p>or choose a file from your computer</p>
                <button className="choose-file-button" type="button" onClick={() => inputRef.current?.click()}>Browse images</button>
                <span className="file-requirements">JPG, JPEG, PNG · Maximum size 5 MB</span>
              </> : <ImagePreview file={file} previewUrl={previewUrl} onRemove={clearSelection} />}
              <input ref={inputRef} className="visually-hidden" type="file" accept=".jpg,.jpeg,.png,.webq,image/jpeg,image/png" onChange={(event) => chooseFile(event.target.files?.[0])} />
            </div>
            <div className="upload-actions"><span className="upload-size-note">{file ? `${formatBytes(file.size)} selected` : "Your image will be uploaded securely."}</span><button className="upload-button" type="submit" disabled={busy || !file}>{busy ? <><LoaderCircle className="spin" size={17} /> Uploading...</> : <><ArrowUpFromLine size={17} /> Upload image</>}</button></div>

          </form>
        </section>


      </main>
    </div>
  );
}
