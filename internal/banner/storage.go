package banner

type Storage interface {
	GetBanner(tagID, featureID int, banner *Banner) error
	GetAllBanners(tagID, featureID, limit, offset int) ([]Banner, error)
	CreateBanner(requestBody *CreateBannerDTO) (int, error)
	UpdateBanner(id int, requestBody *CreateBannerDTO) (int, error)
	DeleteBanner(id int) (int, error)
}