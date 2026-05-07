package ownership_test

import (
	"testing"

	"github.com/yourusername/vaultrot/internal/ownership"
)

func baseOwner() ownership.Owner {
	return ownership.Owner{
		Name:    "platform-team",
		Contact: "platform@example.com",
		Tags:    map[string]string{"env": "prod"},
	}
}

func TestRegister_And_Lookup(t *testing.T) {
	r := ownership.New()
	if err := r.Register("db/password", baseOwner()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	owner, ok := r.Lookup("db/password")
	if !ok {
		t.Fatal("expected owner to be found")
	}
	if owner.Name != "platform-team" {
		t.Errorf("got name %q, want %q", owner.Name, "platform-team")
	}
}

func TestRegister_EmptyKey_ReturnsError(t *testing.T) {
	r := ownership.New()
	if err := r.Register("", baseOwner()); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestRegister_EmptyOwnerName_ReturnsError(t *testing.T) {
	r := ownership.New()
	owner := baseOwner()
	owner.Name = ""
	if err := r.Register("db/password", owner); err == nil {
		t.Fatal("expected error for empty owner name")
	}
}

func TestLookup_MissingKey_ReturnsFalse(t *testing.T) {
	r := ownership.New()
	_, ok := r.Lookup("nonexistent")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestMustLookup_Panics_WhenAbsent(t *testing.T) {
	r := ownership.New()
	defer func() {
		if rec := recover(); rec == nil {
			t.Fatal("expected panic for missing key")
		}
	}()
	r.MustLookup("ghost/secret")
}

func TestLen_ReturnsCorrectCount(t *testing.T) {
	r := ownership.New()
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
	_ = r.Register("a", baseOwner())
	_ = r.Register("b", baseOwner())
	if r.Len() != 2 {
		t.Fatalf("expected 2, got %d", r.Len())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := ownership.New()
	_ = r.Register("svc/token", baseOwner())
	all := r.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	// Mutating the returned map must not affect the registry.
	delete(all, "svc/token")
	if r.Len() != 1 {
		t.Fatal("registry was mutated via All() return value")
	}
}

func TestRegister_OverwritesPreviousOwner(t *testing.T) {
	r := ownership.New()
	_ = r.Register("db/password", baseOwner())
	newOwner := ownership.Owner{Name: "security-team", Contact: "sec@example.com"}
	_ = r.Register("db/password", newOwner)
	owner, _ := r.Lookup("db/password")
	if owner.Name != "security-team" {
		t.Errorf("got %q, want %q", owner.Name, "security-team")
	}
}
