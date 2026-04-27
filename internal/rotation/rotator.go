package rotation

import (
	"fmt"
	"log"

	"github.com/yourusername/vaultrot/internal/backend"
	"github.com/yourusername/vaultrot/internal/config"
)

// Result holds the outcome of a single secret rotation attempt.
type Result struct {
	SecretName string
	Backend    string
	Success    bool
	DryRun     bool
	Err        error
}

// Rotator orchestrates secret rotation across configured backends.
type Rotator struct {
	cfg    *config.Config
	dryRun bool
}

// New creates a new Rotator from the provided config.
func New(cfg *config.Config, dryRun bool) *Rotator {
	return &Rotator{cfg: cfg, dryRun: dryRun}
}

// Run iterates over all secrets in the config and rotates each one.
func (r *Rotator) Run() []Result {
	results := make([]Result, 0, len(r.cfg.Secrets))

	for _, s := range r.cfg.Secrets {
		res := r.rotate(s)
		results = append(results, res)
		if res.Err != nil {
			log.Printf("[ERROR] secret=%s backend=%s err=%v", res.SecretName, res.Backend, res.Err)
		} else if res.DryRun {
			log.Printf("[DRY-RUN] would rotate secret=%s backend=%s", res.SecretName, res.Backend)
		} else {
			log.Printf("[OK] rotated secret=%s backend=%s", res.SecretName, res.Backend)
		}
	}

	return results
}

func (r *Rotator) rotate(s config.Secret) Result {
	res := Result{
		SecretName: s.Name,
		Backend:    s.Backend,
		DryRun:     r.dryRun,
	}

	if r.dryRun {
		res.Success = true
		return res
	}

	b, err := backend.New(s)
	if err != nil {
		res.Err = fmt.Errorf("initialising backend: %w", err)
		return res
	}

	newValue, err := b.Generate(s.Name)
	if err != nil {
		res.Err = fmt.Errorf("generating secret: %w", err)
		return res
	}

	if err := b.Set(s.Name, newValue); err != nil {
		res.Err = fmt.Errorf("setting secret: %w", err)
		return res
	}

	res.Success = true
	return res
}
