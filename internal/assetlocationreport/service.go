package assetlocationreport

import "context"

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func (s *Service) Report(ctx context.Context) (Report, error) {
	assets, err := s.repository.List(ctx)
	if err != nil {
		return Report{}, err
	}
	return BuildReport(assets), nil
}

func BuildReport(assets []Asset) Report {
	report := Report{TotalAssets: len(assets)}
	for _, asset := range assets {
		if len(report.Locations) == 0 || report.Locations[len(report.Locations)-1].BranchName != asset.BranchName || report.Locations[len(report.Locations)-1].Name != asset.Location {
			report.Locations = append(report.Locations, Location{BranchName: asset.BranchName, Name: asset.Location})
		}
		location := &report.Locations[len(report.Locations)-1]
		location.Total++
		if len(location.Areas) == 0 || location.Areas[len(location.Areas)-1].Name != asset.Area {
			location.Areas = append(location.Areas, Area{Name: asset.Area})
		}
		area := &location.Areas[len(location.Areas)-1]
		area.Assets = append(area.Assets, asset)
	}
	report.TotalLocations = len(report.Locations)
	return report
}
