package domain

type Deployment struct {
	Platform string
	Image string
	ImageDigest string
}

func NewDeployment(platform, image, imageDigest string) (*Deployment, error) {
	return &Deployment{
		Platform: platform,
		Image: image,
		ImageDigest: imageDigest,
	}, nil
}