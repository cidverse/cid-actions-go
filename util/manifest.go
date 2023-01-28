package util

import (
	"github.com/containers/image/v5/manifest"
)

type ManifestLayer struct {
	Digest           string
	PlatformOS       string
	PlatformArch     string
	PlatformVariant  string
	PlatformFilename string
}

func GetManifestLayersV2S2(mf *manifest.Schema2List) []ManifestLayer {
	var result []ManifestLayer

	for _, m := range mf.Manifests {
		fileName := m.Platform.OS + "_" + m.Platform.Architecture
		if len(m.Platform.Variant) > 0 {
			fileName = fileName + "_" + m.Platform.Variant
		}
		result = append(result, ManifestLayer{
			Digest:           m.Digest.String(),
			PlatformOS:       m.Platform.OS,
			PlatformArch:     m.Platform.Architecture,
			PlatformVariant:  m.Platform.Variant,
			PlatformFilename: fileName,
		})
	}

	return result
}
