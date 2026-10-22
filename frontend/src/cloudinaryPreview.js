// Keep the transformation syntax identical to the Go BuildCloudinaryTransformationString helper.
export function buildTransformation(values) {
  const p = [];
  const r = values.resize || {};
  if (Number(r.width) > 0) p.push(`w_${Number(r.width)}`);
  if (Number(r.height) > 0) p.push(`h_${Number(r.height)}`);
  if (r.crop) p.push(`c_${r.crop}`);
  if (r.crop === 'crop' && Number(r.width)>0 && Number(r.height)>0) {p.push(`x_${Math.round(r.x||0)}`);p.push(`y_${Math.round(r.y||0)}`);}
  if (r.gravity) p.push(`g_${r.gravity}`);
  if (r.aspect_ratio) p.push(`ar_${r.aspect_ratio}`);
  const e = values.effects || {};
  if (e.blur > 0) p.push(`e_blur:${e.blur}`);
  if (e.grayscale) p.push('e_grayscale');
  if (e.sepia) p.push('e_sepia');
  if (e.cartoonify) p.push('e_cartoonify');
  if (e.pixelate > 0) p.push(`e_pixelate:${e.pixelate}`);
  if (e.brightness !== 0 && e.brightness != null) p.push(`e_brightness:${e.brightness}`);
  if (e.contrast !== 0 && e.contrast != null) p.push(`e_contrast:${e.contrast}`);
  if (e.saturation !== 0 && e.saturation != null) p.push(`e_saturation:${e.saturation}`);
  const a = values.ai || {};
  if (a.background_removal) p.push('e_background_removal');
  if (a.generative_restore) p.push('e_gen_restore');
  if (a.ai_upscale) p.push('e_upscale');
  const g = values.geometry || {};
  if (g.angle) p.push(`a_${g.angle}`);
  if (g.flip === 'horizontal') p.push('a_hflip');
  if (g.flip === 'vertical') p.push('a_vflip');
  const o = values.optimization || {};
  if (o.format) p.push(`f_${o.format}`);
  if (o.quality) p.push(`q_${o.quality}`);
  const stages = [];
  let group = [];
  for (const part of p) {
    if (part.startsWith('e_')) {
      if (group.length) { stages.push(group.join(',')); group = []; }
      stages.push(part);
    } else { group.push(part); }
  }
  if (group.length) stages.push(group.join(','));
  return stages.join('/');
}

// Cloudinary applies the same transformation on delivery without creating an upload or DB row.
// This requires a public Cloudinary image URL with delivery transformations enabled.
export function cloudinaryPreviewUrl(sourceUrl, transformation) {
  if (!transformation) return sourceUrl;
  const url = new URL(sourceUrl);
  if (url.hostname !== 'res.cloudinary.com') throw new Error('Live preview requires a Cloudinary image URL.');
  const marker = '/image/upload/';
  const index = url.pathname.indexOf(marker);
  if (index === -1) throw new Error('Unsupported image URL for live preview.');
  url.pathname = url.pathname.slice(0, index + marker.length) + transformation + '/' + url.pathname.slice(index + marker.length);
  return url.toString();
}

export async function downloadImage(image, onProgress) {
  if (!image?.file_url) throw new Error('Image is not available.');
  onProgress?.(true);
  const ext = (image.format || image.file_url.match(/\.(png|jpe?g|webp|avif)(?:\?|$)/i)?.[1] || 'jpg').replace('jpeg', 'jpg');
  const filename = (image.original_name || 'omarix-image').replace(/\.[^.]+$/, '').replace(/[^a-z0-9_. -]/gi, '_') + '.' + ext;
  try {
    // Fetch the bytes, then use a same-origin blob URL so Chrome downloads instead of opening a tab.
    const response = await fetch(image.file_url, { mode: 'cors' });
    if (!response.ok) throw new Error(`Download failed (${response.status})`);
    const blob = await response.blob();
    const objectUrl = URL.createObjectURL(blob);
    const anchor = document.createElement('a');
    anchor.href = objectUrl;
    anchor.download = filename;
    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
    setTimeout(() => URL.revokeObjectURL(objectUrl), 30000);
  } catch (err) {
    // If Cloudinary's CORS policy prevents fetching bytes, use its delivery attachment flag.
    // This requests an attachment directly without opening a new browser tab.
    let fallback;
    try { fallback = cloudinaryPreviewUrl(image.file_url, 'fl_attachment'); }
    catch { throw err; }
    const a = document.createElement('a');
    a.href = fallback;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
  } finally {
    onProgress?.(false);
  }
}
