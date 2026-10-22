import { useState } from 'react';
import { ArrowLeft, CheckCircle2, Download, Images, LoaderCircle } from 'lucide-react';
import { downloadImage } from '../cloudinaryPreview.js';
export default function ResultPage({ image, onNavigate }) {
  const [downloading, setDownloading] = useState(false);
  const [error, setError] = useState('');
  if (!image) return <main className="omarix-result"><h1>No image selected</h1><button onClick={() => onNavigate('library')}>Open gallery</button></main>;
  const handleDownload = async () => { setError(''); try { await downloadImage(image, setDownloading); } catch (e) { setError('Could not download automatically. Please try again.'); } };
  return <main className="omarix-result">
    <div className="omarix-result-heading"><span><CheckCircle2 size={18}/> Your creation is ready</span><h1>Looking good.</h1><p>Here's your finished image, ready to keep.</p></div>
    <div className="omarix-result-frame"><img src={image.file_url} alt="Your finished edit" /></div>
    {error && <p role="alert" className="notice error">{error}</p>}
    <button className="omarix-result-download" disabled={downloading} onClick={handleDownload}>{downloading?<LoaderCircle className="spin" size={20}/>:<Download size={20}/>} {downloading?'Preparing download…':'Download image'}</button>
    <div className="omarix-result-links"><button onClick={() => onNavigate('library')}><Images size={16}/> Gallery</button><button onClick={() => onNavigate('workspace',image)}><ArrowLeft size={16}/> Edit again</button></div>
  </main>;
}
