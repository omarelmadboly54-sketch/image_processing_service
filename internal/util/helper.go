package util

import (
	"errors"
	"fmt"
	"image_processing_service/internal/presenters"
	"strconv"
	"strings"
)

func BuildCloudinaryTransformationString(req *presenters.TransformImageRequest) string {
	var parts []string

	if req.Resize != nil {
		if req.Resize.Width > 0 {
			parts = append(parts, fmt.Sprintf("w_%d", req.Resize.Width))
		}
		if req.Resize.Height > 0 {
			parts = append(parts, fmt.Sprintf("h_%d", req.Resize.Height))
		}
		if req.Resize.Crop != "" {
			parts = append(parts, fmt.Sprintf("c_%s", req.Resize.Crop))
			if req.Resize.Crop == "crop" && req.Resize.Width > 0 && req.Resize.Height > 0 {
				parts = append(parts, fmt.Sprintf("x_%d", req.Resize.X), fmt.Sprintf("y_%d", req.Resize.Y))
			}
		}
		if req.Resize.Gravity != "" {
			parts = append(parts, fmt.Sprintf("g_%s", req.Resize.Gravity))
		}
		if req.Resize.AspectRatio != "" {
			parts = append(parts, fmt.Sprintf("ar_%s", req.Resize.AspectRatio))
		}
	}

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
		if req.Effects.Pixelate > 0 {
			parts = append(parts, fmt.Sprintf("e_pixelate:%d", req.Effects.Pixelate))
		}

		if req.Effects.Brightness != 0 {
			parts = append(parts, fmt.Sprintf("e_brightness:%d", req.Effects.Brightness))
		}
		if req.Effects.Contrast != 0 {
			parts = append(parts, fmt.Sprintf("e_contrast:%d", req.Effects.Contrast))
		}
		if req.Effects.Saturation != 0 {
			parts = append(parts, fmt.Sprintf("e_saturation:%d", req.Effects.Saturation))
		}
	}

	if req.AI != nil {
		if req.AI.BackgroundRemoval {
			parts = append(parts, "e_background_removal")
		}
		if req.AI.GenerativeRestore {
			parts = append(parts, "e_gen_restore")
		}
		if req.AI.AIUpscale {
			parts = append(parts, "e_upscale")
		}
	}

	if req.Geometry != nil {
		if req.Geometry.Angle != 0 {
			parts = append(parts, fmt.Sprintf("a_%d", req.Geometry.Angle))
		}
		if req.Geometry.Flip != "" {

			if req.Geometry.Flip == "horizontal" {
				parts = append(parts, "a_hflip")
			} else if req.Geometry.Flip == "vertical" {
				parts = append(parts, "a_vflip")
			}
		}
	}

	if req.Optimization != nil {
		if req.Optimization.Format != "" {
			parts = append(parts, fmt.Sprintf("f_%s", req.Optimization.Format))
		}
		if req.Optimization.Quality != "" {
			parts = append(parts, fmt.Sprintf("q_%s", req.Optimization.Quality))
		}
	}

	var stages []string
	var group []string
	for _, part := range parts {
		if strings.HasPrefix(part, "e_") {
			if len(group) > 0 {
				stages = append(stages, strings.Join(group, ","))
				group = nil
			}
			stages = append(stages, part)
		} else {
			group = append(group, part)
		}
	}
	if len(group) > 0 {
		stages = append(stages, strings.Join(group, ","))
	}
	return strings.Join(stages, "/")
}

func ValidateTransformRequest(req *presenters.TransformImageRequest) error {
	if req == nil {
		return errors.New("transformation request is required")
	}

	if r := req.Resize; r != nil {
		if r.X < 0 || r.Y < 0 || r.X > 4096 || r.Y > 4096 {
			return errors.New("crop coordinates must be between 0 and 4096")
		}
		if r.Crop == "crop" && (r.Width == 0 || r.Height == 0) {
			return errors.New("manual crop requires width and height")
		}
		if r.Width < 0 || r.Width > 4096 {
			return errors.New("width must be between 0 and 4096")
		}
		if r.Height < 0 || r.Height > 4096 {
			return errors.New("height must be between 0 and 4096")
		}

		if !allowed(r.Crop, "", "fill", "fit", "scale", "crop", "pad") {
			return errors.New("unsupported crop mode")
		}

		if !allowed(r.Gravity, "", "auto", "center", "face", "faces",
			"north", "south", "east", "west",
			"north_east", "north_west", "south_east", "south_west") {
			return errors.New("unsupported gravity")
		}

		if r.AspectRatio != "" {
			if err := validateAspectRatio(r.AspectRatio); err != nil {
				return err
			}
		}
	}

	if e := req.Effects; e != nil {
		if e.Blur < 0 || e.Blur > 2000 {
			return errors.New("blur must be between 0 and 2000")
		}
		if e.Pixelate < 0 || e.Pixelate > 200 {
			return errors.New("pixelate must be between 0 and 200")
		}
		if e.Brightness < -100 || e.Brightness > 100 {
			return errors.New("brightness must be between -100 and 100")
		}
		if e.Contrast < -100 || e.Contrast > 100 {
			return errors.New("contrast must be between -100 and 100")
		}
		if e.Saturation < -100 || e.Saturation > 100 {
			return errors.New("saturation must be between -100 and 100")
		}
	}

	if g := req.Geometry; g != nil {
		if g.Angle < -360 || g.Angle > 360 {
			return errors.New("angle must be between -360 and 360")
		}
		if !allowed(g.Flip, "", "horizontal", "vertical") {
			return errors.New("unsupported flip direction")
		}
	}

	if o := req.Optimization; o != nil {
		if !allowed(o.Format, "", "jpg", "jpeg", "png", "webp", "avif", "auto") {
			return errors.New("unsupported image format")
		}

		if o.Quality != "" && o.Quality != "auto" {
			quality, err := strconv.Atoi(o.Quality)
			if err != nil || quality < 1 || quality > 100 {
				return errors.New("quality must be auto or a number between 1 and 100")
			}
		}
	}

	return nil
}

func allowed(value string, values ...string) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}
	return false
}

func validateAspectRatio(value string) error {
	parts := strings.Split(value, ":")

	if len(parts) != 2 {
		return errors.New("aspect_ratio must use width:height format, e.g. 16:9")
	}

	width, errW := strconv.Atoi(parts[0])
	height, errH := strconv.Atoi(parts[1])

	if errW != nil || errH != nil ||
		width <= 0 || height <= 0 ||
		width > 10000 || height > 10000 {
		return fmt.Errorf("invalid aspect_ratio: %s", value)
	}

	return nil
}
