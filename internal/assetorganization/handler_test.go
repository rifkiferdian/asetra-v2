package assetorganization

import (
	"strings"
	"testing"
)

func TestEnrichLocationHierarchy(t *testing.T) {
	locations := []AssetLocation{
		{ID: 3, ParentID: 2, BranchName: "Pusat", Name: "Ruang Rapat"},
		{ID: 1, BranchName: "Pusat", Name: "Gedung A", ChildCount: 1},
		{ID: 2, ParentID: 1, BranchName: "Pusat", Name: "Lantai 1", ChildCount: 1},
	}

	if roots := enrichLocationHierarchy(locations); roots != 1 {
		t.Fatalf("root count = %d, want 1", roots)
	}
	if locations[0].ID != 1 || locations[1].ID != 2 || locations[2].ID != 3 {
		t.Fatalf("locations not sorted hierarchically: %#v", locations)
	}
	if locations[0].DescendantCount != 2 {
		t.Fatalf("root descendant count = %d, want 2", locations[0].DescendantCount)
	}
	if locations[2].Depth != 2 || locations[2].HierarchyPath != "Gedung A / Lantai 1 / Ruang Rapat" {
		t.Fatalf("unexpected leaf hierarchy: %#v", locations[2])
	}
	if !strings.Contains(locations[2].AncestorIDs, " 1 ") || !strings.Contains(locations[2].AncestorIDs, " 2 ") {
		t.Fatalf("leaf ancestors = %q, want root and parent IDs", locations[2].AncestorIDs)
	}
}
