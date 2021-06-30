// Package global provides structs, interfaces, paths and locks which the should be accessible fpr the hole project.
package global

type Module interface {
	New(map[string]interface{}) (Module, error)
	Execute(currentStep *Step) (string, error)
}
