import { FileImage, X } from "lucide-react";

function formatBytes(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

export { formatBytes };

export default function ImagePreview({ file, previewUrl, onRemove }) {
  if (!file || !previewUrl) return null;
  return (
    <div className="selected-file-card">
      <div className="selected-preview"><img src={previewUrl} alt="Selected image preview" /></div>
      <div className="selected-file-info"><span className="file-type-icon"><FileImage size={18} /></span><div className="selected-file-text"><strong>{file.name}</strong><span>{formatBytes(file.size)} · {file.type || "Image file"}</span></div></div>
      <button className="remove-file-button" type="button" onClick={onRemove} aria-label="Remove selected image"><X size={17} /></button>
    </div>
  );
}

