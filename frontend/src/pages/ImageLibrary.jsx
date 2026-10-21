
import { useCallback, useEffect, useState } from "react";
import {
  ArrowLeft,
  ChevronLeft,
  ChevronRight,
  Image as ImageIcon,
  LoaderCircle,
  RefreshCw,
  Search,
  Trash2,
  X,
} from "lucide-react";
import Navbar from "../components/Navbar.jsx";

const API_BASE_URL = (
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"
).replace(/\/+$/, "");

const PAGE_SIZE = 5;


function formatBytes(bytes) {
  if (typeof bytes !== "number") return "—";
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

async function readResponse(response) {
  const raw = await response.text();
  try {
    return raw ? JSON.parse(raw) : {};
  } catch {
    return { message: raw };
  }
}

export default function MyImages({ onLogout, onNavigate }) {
  const [totalImages, setTotalImages] = useState(0);
  const [images, setImages] = useState([]);
  const [page, setPage] = useState(1);
  const [selectedImage, setSelectedImage] = useState(null);
  const [loading, setLoading] = useState(true);
  const [detailLoading, setDetailLoading] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [error, setError] = useState("");
  const [detailError, setDetailError] = useState("");

  const getToken = () => sessionStorage.getItem("access_token");

  const handleUnauthorized = () => {
    sessionStorage.removeItem("access_token");
    onLogout();
  };

  const fetchImages = useCallback(async (pageNumber = 1) => {
    const token = getToken();

    if (!token) {
      handleUnauthorized();
      return;
    }

    setLoading(true);
    setError("");

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/images?page=${pageNumber}&limit=${PAGE_SIZE}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const data = await readResponse(response);

      if (response.status === 401 || response.status === 403) {
        handleUnauthorized();
        return;
      }

      if (!response.ok) {
        throw new Error(data.error || data.message || `Request failed (${response.status})`);
      }

    const result = Array.isArray(data.images) ? data.images : [];

    setImages(result);
    setPage(pageNumber);
    setTotalImages(Number(data.total) || 0);
     
} catch (err) {
      setError(
        err instanceof TypeError
          ? `Couldn't connect to ${API_BASE_URL}. Check that your Go backend is running.`
          : err.message
      );
      setImages([]);
    } finally {
      setLoading(false);
    }
  }, [onLogout]);

  useEffect(() => {
    fetchImages(1);
  }, [fetchImages]);

  async function openImage(imageId) {
    if (imageId == null) {
      setDetailError("This image has no ID in the API response. Make sure the model uses json:\"id\".");
      return;
    }

    const token = getToken();
    if (!token) {
      handleUnauthorized();
      return;
    }

    setSelectedImage(null);
    setDetailError("");
    setDetailLoading(true);

    try {
      const response = await fetch(
        `${API_BASE_URL}/api/images/${encodeURIComponent(imageId)}`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const data = await readResponse(response);

      if (response.status === 401 || response.status === 403) {
        handleUnauthorized();
        return;
      }

      if (!response.ok) {
        throw new Error(data.error || data.message || `Request failed (${response.status})`);
      }

      setSelectedImage(data.image || null);

      if (!data.image) {
        throw new Error("The server response did not include an image.");
      }
    } catch (err) {
      setDetailError(
        err instanceof TypeError
          ? "Couldn't connect to the backend."
          : err.message
      );
    } finally {
      setDetailLoading(false);
    }
  }

  async function deleteSelectedImage() {
    if (!selectedImage?.id || deleting) return;
    if (!window.confirm("Delete this image permanently? This cannot be undone.")) return;
    setDeleting(true);
    setDeleteError("");
    setSuccessMessage("");
    try {
      const response = await fetch(`${API_BASE_URL}/api/images/${encodeURIComponent(selectedImage.id)}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      const data = await readResponse(response);
      if (response.status === 401 || response.status === 403) {
        handleUnauthorized();
        return;
      }
      if (!response.ok) throw new Error(data.error || data.message || `Delete failed (${response.status})`);
      setSelectedImage(null);
      setSuccessMessage("Image deleted successfully.");
      const targetPage = images.length === 1 && page > 1 ? page - 1 : page;
      await fetchImages(targetPage);
    } catch (err) {
      setDeleteError(err.message || "Could not delete image.");
    } finally {
      setDeleting(false);
    }
  }

  async function goToPage(nextPage) {
    if (nextPage < 1 || loading) return;

    await fetchImages(nextPage);

    // If a page has no results, go back to the previous page.
    // This can happen when the last page contained exactly five images.
    // We check the updated state on the next render instead of guessing
    // the number of total images.
  }

  return (
    <div className="workspace-page" id="top">

      <main className="library-main">
        <div className="library-heading">
          <span className="workspace-eyebrow">
            <ImageIcon size={14} /> YOUR OMARIX LIBRARY
          </span>
          <h1>My <span>Images.</span></h1>
          <p>Your uploaded images, all in one place.</p>
        </div>

        <section className="library-panel">
          <div className="library-toolbar">
            <div>
              <h2>Image library</h2>
                <p>
                 Page {page} · {images.length} shown · {totalImages} total images
                </p>
            </div>
            <button
              className="library-refresh"
              type="button"
              onClick={() => fetchImages(page)}
              disabled={loading}
            >
              <RefreshCw size={15} className={loading ? "spin" : ""} />
              Refresh
            </button>
          </div>

          {successMessage && <div className="notice success" role="status">{successMessage}</div>}
          {error && (
            <div className="notice error" role="alert">
              {error}
              <button type="button" className="inline-retry" onClick={() => fetchImages(page)}>
                Try again
              </button>
            </div>
          )}

          {loading ? (
            <div className="library-state">
              <LoaderCircle className="spin" size={28} />
              <p>Loading your images...</p>
            </div>
          ) : !error && images.length === 0 ? (
            <div className="library-state">
              <div className="library-empty-icon"><ImageIcon size={27} /></div>
              <h3>{page === 1 ? "Your library is empty" : "No more images"}</h3>
              <p>
                {page === 1
                  ? "Upload your first image and it will appear here."
                  : "There are no images on this page."}
              </p>
              {page > 1 ? (
                <button className="library-primary-button" onClick={() => goToPage(page - 1)}>
                  <ArrowLeft size={16} /> Previous page
                </button>
              ) : (
                <button className="library-primary-button" onClick={() => onNavigate("upload")}>
                  Upload an image
                </button>
              )}
            </div>
          ) : (
            <>
              <div className="image-library-grid">
                {images.map((image) => (
                  <article className="library-image-card" key={image.id}>
                    <button
                      type="button"
                      className="library-image-button"
                      onClick={() => openImage(image.id)}
                      aria-label={`View details for ${image.original_name || "image"}`}
                    >
                      <img
                        src={image.file_url}
                        alt={image.original_name || "Uploaded image"}
                        loading="lazy"
                      />
                      <span className="library-image-overlay">
                        <Search size={17} /> View details
                      </span>
                    </button>
                    <div className="library-card-info">
                      <strong title={image.original_name}>
                        {image.original_name || "Untitled image"}
                      </strong>
                      <span>
                        {image.format?.toUpperCase() || "IMAGE"} · {formatBytes(image.file_size)}
                      </span>
                    </div>
                  </article>
                ))}
              </div>

              <div className="library-pagination">
                <button
                  type="button"
                  onClick={() => goToPage(page - 1)}
                  disabled={page === 1 || loading}
                >
                  <ChevronLeft size={17} /> Previous
                </button>
                <span>
                    Page {page} of {Math.max(1, Math.ceil(totalImages / PAGE_SIZE))}
                </span>
                <button
                  type="button"
                  onClick={() => goToPage(page + 1)}
                  disabled={loading || page * PAGE_SIZE >= totalImages}
                >
                  Next <ChevronRight size={17} />
                </button>
              </div>
            </>
          )}
        </section>
      </main>

      {(detailLoading || selectedImage || detailError) && (
        <div
          className="image-modal-backdrop"
          role="presentation"
          onMouseDown={(event) => {
            if (event.target === event.currentTarget) {
              setSelectedImage(null);
              setDetailError("");
            }
          }}
        >
          <section className="image-modal" role="dialog" aria-modal="true" aria-labelledby="image-modal-title">
            <div className="image-modal-header">
              <div>
                <span className="workspace-eyebrow">IMAGE DETAILS</span>
                <h2 id="image-modal-title">Your image</h2>
              </div>
              <button
                className="image-modal-close"
                type="button"
                aria-label="Close image details"
                onClick={() => {
                  setSelectedImage(null);
                  setDetailError("");
                }}
              >
                <X size={19} />
              </button>
            </div>

            {deleteError && <div className="notice error" role="alert">{deleteError}</div>}
            {detailLoading ? (
              <div className="library-state">
                <LoaderCircle className="spin" size={28} />
                <p>Fetching image details...</p>
              </div>
            ) : detailError ? (
              <div className="notice error" role="alert">{detailError}</div>
            ) : selectedImage ? (
              <div className="image-detail-layout">
                <div className="image-detail-preview">
                  <img
                    src={selectedImage.file_url}
                    alt={selectedImage.original_name || "Selected image"}
                  />
                </div>
                <div className="image-detail-metadata">
                  <div><span>File name</span><strong>{selectedImage.original_name || "—"}</strong></div>
                  <div><span>Format</span><strong>{selectedImage.format?.toUpperCase() || "—"}</strong></div>
                  <div><span>File size</span><strong>{formatBytes(selectedImage.file_size)}</strong></div>
                  <div><span>Dimensions</span><strong>{selectedImage.width && selectedImage.height ? `${selectedImage.width} × ${selectedImage.height} px` : "—"}</strong></div>
                  
                  <button type="button" className="delete-image-button" disabled={deleting} onClick={deleteSelectedImage}><Trash2 size={16} /> {deleting ? "Deleting..." : "Delete image"}</button>
                  {selectedImage.file_url && (
                    <a
                      className="image-url-link"
                      href={selectedImage.file_url}
                      target="_blank"
                      rel="noreferrer"
                    >
                      Open original image ↗
                    </a>
                  )}
                </div>
              </div>
            ) : null}
          </section>
        </div>
      )}
    </div>
  );
}
