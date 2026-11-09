package assetregistry

import "context"

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service             { return &Service{repository: repository} }
func (s *Service) List(ctx context.Context) ([]Asset, error) { return s.repository.List(ctx) }
func (s *Service) Find(ctx context.Context, id uint64) (Asset, error) {
	return s.repository.Find(ctx, id)
}

func Summarize(assets []Asset) Summary {
	x := Summary{Total: len(assets)}
	for _, asset := range assets {
		switch asset.StatusCode {
		case "MAINTENANCE":
			x.Maintenance++
		case "RETIRED", "DISPOSED", "INACTIVE":
			x.Inactive++
		default:
			x.Active++
		}
	}
	return x
}
