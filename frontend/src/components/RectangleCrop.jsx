import { useRef, useState } from 'react';
export default function RectangleCrop({src,selection,onSelection}) {
 const imageRef=useRef(null), startRef=useRef(null),[drawing,setDrawing]=useState(false),[draft,setDraft]=useState(null);
 const point=e=>{const img=imageRef.current;const rect=img.getBoundingClientRect();return {x:Math.max(0,Math.min(1,(e.clientX-rect.left)/rect.width)),y:Math.max(0,Math.min(1,(e.clientY-rect.top)/rect.height))};};
 const make=(a,b)=>({x:Math.min(a.x,b.x),y:Math.min(a.y,b.y),width:Math.abs(b.x-a.x),height:Math.abs(b.y-a.y)});
 const down=e=>{if(e.button!==0)return;e.preventDefault();e.currentTarget.setPointerCapture(e.pointerId);startRef.current=point(e);setDraft(null);setDrawing(true);};
 const move=e=>{if(drawing&&startRef.current)setDraft(make(startRef.current,point(e)));};
 const up=e=>{if(!drawing)return;const r=make(startRef.current,point(e));setDrawing(false);startRef.current=null;setDraft(null);if(r.width>0.005&&r.height>0.005){const img=imageRef.current;onSelection({x:Math.round(r.x*img.naturalWidth),y:Math.round(r.y*img.naturalHeight),width:Math.round(r.width*img.naturalWidth),height:Math.round(r.height*img.naturalHeight)});} };
 const displayed=draft||(selection&&imageRef.current?.naturalWidth?{x:selection.x/imageRef.current.naturalWidth,y:selection.y/imageRef.current.naturalHeight,width:selection.width/imageRef.current.naturalWidth,height:selection.height/imageRef.current.naturalHeight}:null);
 return <div className="omarix-crop-viewport"><div className="omarix-crop-surface" onPointerDown={down} onPointerMove={move} onPointerUp={up} onPointerCancel={()=>setDrawing(false)}><img ref={imageRef} src={src} draggable="false" alt="Select crop region" />{displayed&&<div className="omarix-crop-selection" style={{left:`${displayed.x*100}%`,top:`${displayed.y*100}%`,width:`${displayed.width*100}%`,height:`${displayed.height*100}%`}}><i/><i/><i/><i/></div>}</div><p>Drag across the photo to select the area you want to keep.</p></div>;
}
