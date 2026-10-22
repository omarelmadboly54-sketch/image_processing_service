package presenters

type TransformImageRequest struct {
	Resize       *ResizeOptions       `json:"resize"`
	Effects      *EffectOptions       `json:"effects"`
	AI           *AIOptions           `json:"ai"`
	Geometry     *GeometryOptions     `json:"geometry"`
	Optimization *OptimizationOptions `json:"optimization"`
}

type ResizeOptions struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Crop        string `json:"crop"`
	Gravity     string `json:"gravity"`
	AspectRatio string `json:"aspect_ratio"`
}

type EffectOptions struct {
	Blur       int  `json:"blur"`
	Grayscale  bool `json:"grayscale"`
	Sepia      bool `json:"sepia"`
	Cartoonify bool `json:"cartoonify"`
	Pixelate   int  `json:"pixelate"`
	Brightness int  `json:"brightness"`
	Contrast   int  `json:"contrast"`
	Saturation int  `json:"saturation"`
}

type AIOptions struct {
	BackgroundRemoval bool `json:"background_removal"`
	GenerativeRestore bool `json:"generative_restore"`
	AIUpscale         bool `json:"ai_upscale"`
}

type GeometryOptions struct {
	Angle int    `json:"angle"`
	Flip  string `json:"flip"`
}

type OptimizationOptions struct {
	Format  string `json:"format"`
	Quality string `json:"quality"`
}


