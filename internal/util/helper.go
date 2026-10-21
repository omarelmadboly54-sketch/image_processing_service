package util

import (
	"fmt"
	"image_processing_service/internal/presenters"
	"strings"
)

func BuildCloudinaryTransformationString(req *presenters.TransformImageRequest) string {

	var parts []string

	// Resize options
	if req.Resize != nil {
		if req.Resize.Width > 0 {
			parts = append(parts, fmt.Sprintf("w_%d", req.Resize.Width))
		}
		if req.Resize.Height > 0 {
			parts = append(parts, fmt.Sprintf("h_%d", req.Resize.Height))
		}
		if req.Resize.Crop != "" {
			parts = append(parts, fmt.Sprintf("c_%s", req.Resize.Crop))
		}
		if req.Resize.Gravity != "" {
			parts = append(parts, fmt.Sprintf("g_%s", req.Resize.Gravity))
		}
	}

	// Effect options
	if req.Effects != nil {
		if req.Effects.Blur > 0 {
			parts = append(parts, fmt.Sprintf("e_blur:%d", req.Effects.Blur))
		}
		if req.Effects.Grayscale {
			parts = append(parts, "e_grayscale")
		}
		if req.Effects.Sepia {
			parts = append(parts, "e_sepia")
		}
		if req.Effects.Cartoonify {
			parts = append(parts, "e_cartoonify")
		}
	}

	// AI options
	if req.AI != nil {
		if req.AI.BackgroundRemoval {
			parts = append(parts, "e_background_removal")
		}
		if req.AI.AIUpscale {
			parts = append(parts, "e_upscale")
		}
	}

	// Optimization / Format
	if req.Optimization != nil {
		if req.Optimization.Format != "" {
			parts = append(parts, fmt.Sprintf("f_%s", req.Optimization.Format))
		}
		if req.Optimization.Quality != "" {
			parts = append(parts, fmt.Sprintf("q_%s", req.Optimization.Quality))
		}
	}

	
	return strings.Join(parts,",")
}

