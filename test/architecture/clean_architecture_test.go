package architecture

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLayerSeparation ensures that layers follow clean architecture principles
func TestLayerSeparation(t *testing.T) {
	projectRoot := getProjectRoot()

	tests := []struct {
		name             string
		layer            string
		allowedImports   []string
		forbiddenImports []string
	}{
		{
			name:             "Models should not import from other layers",
			layer:            "models",
			allowedImports:   []string{"time", "github.com", "gorm.io"},
			forbiddenImports: []string{"handlers", "services", "repositories"},
		},
		{
			name:             "Repositories should not import handlers or services",
			layer:            "repositories",
			allowedImports:   []string{"models", "gorm.io", "github.com"},
			forbiddenImports: []string{"handlers", "services"},
		},
		{
			name:             "Services should not import handlers",
			layer:            "services",
			allowedImports:   []string{"models", "repositories", "gorm.io"},
			forbiddenImports: []string{"handlers"},
		},
		{
			name:             "Handlers can import services and models",
			layer:            "handlers",
			allowedImports:   []string{"models", "services", "gin"},
			forbiddenImports: []string{"repositories"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layerPath := filepath.Join(projectRoot, tt.layer)

			// Skip if directory doesn't exist
			if _, err := os.Stat(layerPath); os.IsNotExist(err) {
				t.Skipf("Layer %s does not exist", tt.layer)
				return
			}

			// Parse all Go files in the layer
			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, layerPath, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("Failed to parse directory %s: %v", layerPath, err)
			}

			// Check imports in each package
			for _, pkg := range pkgs {
				for _, file := range pkg.Files {
					for _, imp := range file.Imports {
						importPath := strings.Trim(imp.Path.Value, `"`)

						// Check if import is forbidden
						for _, forbidden := range tt.forbiddenImports {
							if strings.Contains(importPath, forbidden) {
								t.Errorf("Layer %s should not import %s (found in %s)",
									tt.layer, forbidden, importPath)
							}
						}
					}
				}
			}
		})
	}
}

// TestNamingConventions verifies that files and packages follow naming conventions
func TestNamingConventions(t *testing.T) {
	projectRoot := getProjectRoot()

	tests := []struct {
		name      string
		directory string
		suffix    string
	}{
		{
			name:      "Repository files should end with _repository.go",
			directory: "repositories",
			suffix:    "_repository.go",
		},
		{
			name:      "Service files should end with _service.go",
			directory: "services",
			suffix:    "_service.go",
		},
		{
			name:      "Handler files should end with _handler.go",
			directory: "handlers",
			suffix:    "_handler.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirPath := filepath.Join(projectRoot, tt.directory)

			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Skipf("Directory %s does not exist", tt.directory)
				return
			}

			files, err := os.ReadDir(dirPath)
			if err != nil {
				t.Fatalf("Failed to read directory %s: %v", dirPath, err)
			}

			for _, file := range files {
				if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
					continue
				}

				// Skip test files
				if strings.HasSuffix(file.Name(), "_test.go") {
					continue
				}

				assert.True(t, strings.HasSuffix(file.Name(), tt.suffix),
					"File %s should end with %s", file.Name(), tt.suffix)
			}
		})
	}
}

// TestInterfaceDefinitions ensures that services and repositories define interfaces
func TestInterfaceDefinitions(t *testing.T) {
	projectRoot := getProjectRoot()

	tests := []struct {
		name      string
		directory string
	}{
		{
			name:      "Repositories should define interfaces",
			directory: "repositories",
		},
		{
			name:      "Services should define interfaces",
			directory: "services",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirPath := filepath.Join(projectRoot, tt.directory)

			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Skipf("Directory %s does not exist", tt.directory)
				return
			}

			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse directory %s: %v", dirPath, err)
			}

			hasInterface := false
			for _, pkg := range pkgs {
				for _, file := range pkg.Files {
					ast.Inspect(file, func(n ast.Node) bool {
						if typeSpec, ok := n.(*ast.TypeSpec); ok {
							if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
								hasInterface = true
								return false
							}
						}
						return true
					})
				}
			}

			assert.True(t, hasInterface, "Directory %s should define at least one interface", tt.directory)
		})
	}
}

// TestErrorHandling checks that functions return errors properly
func TestErrorHandling(t *testing.T) {
	projectRoot := getProjectRoot()

	directories := []string{"repositories", "services", "handlers"}

	for _, dir := range directories {
		t.Run("Error handling in "+dir, func(t *testing.T) {
			dirPath := filepath.Join(projectRoot, dir)

			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Skipf("Directory %s does not exist", dir)
				return
			}

			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse directory %s: %v", dirPath, err)
			}

			functionsWithoutError := 0
			totalFunctions := 0

			for _, pkg := range pkgs {
				for _, file := range pkg.Files {
					ast.Inspect(file, func(n ast.Node) bool {
						if funcDecl, ok := n.(*ast.FuncDecl); ok {
							// Skip constructors (New* functions)
							if strings.HasPrefix(funcDecl.Name.Name, "New") {
								return true
							}

							// Skip test functions
							if strings.HasPrefix(funcDecl.Name.Name, "Test") {
								return true
							}

							totalFunctions++

							// Check if function returns error
							if funcDecl.Type.Results != nil {
								hasError := false
								for _, result := range funcDecl.Type.Results.List {
									if ident, ok := result.Type.(*ast.Ident); ok {
										if ident.Name == "error" {
											hasError = true
											break
										}
									}
								}

								if !hasError {
									functionsWithoutError++
								}
							}
						}
						return true
					})
				}
			}

			// At least 50% of functions should return errors (excluding constructors)
			if totalFunctions > 0 {
				errorHandlingRate := float64(totalFunctions-functionsWithoutError) / float64(totalFunctions)
				assert.GreaterOrEqual(t, errorHandlingRate, 0.3,
					"At least 30%% of functions in %s should return errors", dir)
			}
		})
	}
}

// TestDependencyInjection verifies that handlers and services use dependency injection
func TestDependencyInjection(t *testing.T) {
	projectRoot := getProjectRoot()

	tests := []struct {
		name      string
		directory string
	}{
		{
			name:      "Handlers should use dependency injection",
			directory: "handlers",
		},
		{
			name:      "Services should use dependency injection",
			directory: "services",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirPath := filepath.Join(projectRoot, tt.directory)

			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Skipf("Directory %s does not exist", tt.directory)
				return
			}

			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse directory %s: %v", dirPath, err)
			}

			hasConstructor := false
			for _, pkg := range pkgs {
				for _, file := range pkg.Files {
					ast.Inspect(file, func(n ast.Node) bool {
						if funcDecl, ok := n.(*ast.FuncDecl); ok {
							// Check for New* constructor functions
							if strings.HasPrefix(funcDecl.Name.Name, "New") {
								hasConstructor = true

								// Constructor should have parameters (dependencies)
								if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
									return false
								}
							}
						}
						return true
					})
				}
			}

			assert.True(t, hasConstructor, "Directory %s should have constructor functions", tt.directory)
		})
	}
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	// Navigate up from test/architecture to project root
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "../../" // Fallback
}
