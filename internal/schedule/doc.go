// Package schedule provides cron-based scheduling for periodic secret rotation.
//
// It wraps the robfig/cron library with a simplified API that associates
// human-readable names with cron expressions, making it straightforward to
// register and inspect rotation jobs.
//
// Example usage:
//
//	s := schedule.New()
//	err := s.Add("rotate-db-password", "0 3 * * *", func() {
//		// trigger rotation logic here
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	s.Start()
//	defer s.Stop()
package schedule
