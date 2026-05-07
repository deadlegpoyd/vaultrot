// Package observe implements a fan-out event observer for vaultrot rotation
// lifecycle events.
//
// Callers register one or more Handler functions and then call Emit with an
// Event value. Handlers are invoked synchronously in registration order so
// they should be fast; offload heavy work (e.g. HTTP calls) to a goroutine
// inside the handler if necessary.
//
// Example:
//
//	obs := observe.New()
//	obs.Register(func(e observe.Event) {
//		fmt.Println(e)
//	})
//	obs.Emit(observe.Event{
//		Kind:    observe.EventSecretRotated,
//		Secret:  "db/password",
//		Backend: "vault",
//	})
package observe
