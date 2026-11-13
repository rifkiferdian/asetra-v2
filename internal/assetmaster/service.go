package assetmaster

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalid   = errors.New("invalid asset master data")
	ErrDuplicate = errors.New("duplicate asset master data")
	ErrInUse     = errors.New("asset master data in use")
)
var codePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{0,49}$`)

type Service struct{ repository *Repository }

func NewService(r *Repository) *Service { return &Service{repository: r} }
func (s *Service) Items(ctx context.Context, kind string) ([]Item, error) {
	return s.repository.Items(ctx, kind)
}
func (s *Service) Groups(ctx context.Context) ([]Option, error) { return s.repository.Groups(ctx) }
func (s *Service) Parents(ctx context.Context, kind string) ([]Option, error) {
	switch kind {
	case "types":
		return s.repository.Categories(ctx)
	case "models":
		return s.repository.Types(ctx)
	case "variants":
		return s.repository.Models(ctx)
	}
	return nil, nil
}
func (s *Service) Brands(ctx context.Context) ([]Option, error) { return s.repository.Brands(ctx) }
func (s *Service) TypeBrands(ctx context.Context) ([]Option, error) {
	return s.repository.TypeBrands(ctx)
}
func (s *Service) ReplaceTypeBrands(ctx context.Context, typeID uint64, brandIDs []uint64) error {
	if typeID == 0 {
		return ErrInvalid
	}
	seen := make(map[uint64]bool, len(brandIDs))
	unique := make([]uint64, 0, len(brandIDs))
	for _, id := range brandIDs {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		unique = append(unique, id)
	}
	return mapDB(s.repository.ReplaceTypeBrands(ctx, typeID, unique))
}
func (s *Service) Create(ctx context.Context, kind string, x Input) error {
	x = normalize(x)
	if !valid(kind, x) {
		return ErrInvalid
	}
	ok, err := s.repository.RelationsValid(ctx, kind, x)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalid
	}
	return mapDB(s.repository.Create(ctx, kind, x))
}
func (s *Service) Update(ctx context.Context, kind string, id uint64, x Input) error {
	x = normalize(x)
	if id == 0 || !valid(kind, x) {
		return ErrInvalid
	}
	ok, err := s.repository.RelationsValid(ctx, kind, x)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalid
	}
	return mapDB(s.repository.Update(ctx, kind, id, x))
}
func (s *Service) Delete(ctx context.Context, kind string, id uint64) error {
	if id == 0 {
		return ErrInvalid
	}
	return mapDB(s.repository.Delete(ctx, kind, id))
}
func normalize(x Input) Input {
	x.Code = strings.ToUpper(strings.TrimSpace(x.Code))
	x.Name = strings.TrimSpace(x.Name)
	x.Description = strings.TrimSpace(x.Description)
	return x
}
func valid(kind string, x Input) bool {
	if !codePattern.MatchString(x.Code) || x.Name == "" || utf8.RuneCountInString(x.Name) > 150 || utf8.RuneCountInString(x.Description) > 65535 {
		return false
	}
	if kind == "categories" || kind == "brands" {
		return x.GroupID > 0
	}
	if kind == "models" {
		return x.ParentID > 0 && x.BrandID > 0
	}
	return (kind == "types" || kind == "models" || kind == "variants") && x.ParentID > 0
}
func mapDB(err error) error {
	var e *mysql.MySQLError
	if errors.As(err, &e) {
		if e.Number == 1062 {
			return ErrDuplicate
		}
		if e.Number == 1451 || e.Number == 1452 {
			return ErrInUse
		}
	}
	return err
}
