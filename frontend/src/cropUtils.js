export async function cropImageToBlob(url, pixels) {
  const image = new Image();
  image.crossOrigin = 'anonymous';
  image.src = url;
  await new Promise((resolve, reject) => { image.onload = resolve; image.onerror = () => reject(new Error('Could not load image for crop. Check Cloudinary CORS settings.')); });
  const canvas = document.createElement('canvas');
  canvas.width = Math.max(1, Math.round(pixels.width));
  canvas.height = Math.max(1, Math.round(pixels.height));
  const context = canvas.getContext('2d');
  context.drawImage(image, pixels.x, pixels.y, pixels.width, pixels.height, 0, 0, canvas.width, canvas.height);
  return await new Promise((resolve, reject) => canvas.toBlob(blob => blob ? resolve(blob) : reject(new Error('Could not export cropped image.')), 'image/png'));
}
