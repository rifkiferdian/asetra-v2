package assetlocationreport

import (
	"context"
	"database/sql"
)

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func (s *Service) Report(ctx context.Context) (Report, error) {
	assets, err := s.repository.List(ctx)
	if err != nil {
		return Report{}, err
	}
	return BuildReport(assets), nil
}

func (s *Service) Area(ctx context.Context, id uint64) (Location, Area, error) {
	report, err := s.Report(ctx)
	if err != nil {
		return Location{}, Area{}, err
	}
	for _, location := range report.Locations {
		for _, area := range location.Areas {
			if area.ID == id {
				return location, area, nil
			}
		}
	}
	return Location{}, Area{}, sql.ErrNoRows
}

func BuildReport(assets []Asset) Report {
	report := Report{}
	for _, asset := range assets {
		if len(report.Locations) == 0 || report.Locations[len(report.Locations)-1].BranchName != asset.BranchName || report.Locations[len(report.Locations)-1].Name != asset.Location {
			report.Locations = append(report.Locations, Location{BranchName: asset.BranchName, Name: asset.Location})
		}
		location := &report.Locations[len(report.Locations)-1]
		if len(location.Areas) == 0 || location.Areas[len(location.Areas)-1].Name != asset.Area {
			location.Areas = append(location.Areas, Area{ID: asset.LocationID, Name: asset.Area})
		}
		area := &location.Areas[len(location.Areas)-1]
		if asset.ID != 0 || asset.Name != "" {
			area.Assets = append(area.Assets, asset)
			location.Total++
			report.TotalAssets++
		}
	}
	report.TotalLocations = len(report.Locations)
	return report
}
