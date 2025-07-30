package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	binaryName = "techircd"
	version    = "1.0.0"
)

func main() {
	var (
		buildFlag     = flag.Bool("build", false, "Build the binary")
		runFlag       = flag.Bool("run", false, "Build and run the server")
		testFlag      = flag.Bool("test", false, "Run all tests")
		cleanFlag     = flag.Bool("clean", false, "Clean build artifacts")
		fmtFlag       = flag.Bool("fmt", false, "Format Go code")
		lintFlag      = flag.Bool("lint", false, "Run linters")
		buildAllFlag  = flag.Bool("build-all", false, "Build for multiple platforms")
		releaseFlag   = flag.Bool("release", false, "Create optimized release build")
		helpFlag      = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *helpFlag || flag.NFlag() == 0 {
		showHelp()
		return
	}

	switch {
	case *buildFlag:
		build()
	case *runFlag:
		build()
		run()
	case *testFlag:
		test()
	case *cleanFlag:
		clean()
	case *fmtFlag:
		format()
	case *lintFlag:
		lint()
	case *buildAllFlag:
		buildAll()
	case *releaseFlag:
		release()
	}
}

func showHelp() {
	fmt.Println("TechIRCd Build Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run build.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -build       Build the binary")
	fmt.Println("  -run         Build and run the server")
	fmt.Println("  -test        Run all tests")
	fmt.Println("  -clean       Clean build artifacts")
	fmt.Println("  -fmt         Format Go code")
	fmt.Println("  -lint        Run linters")
	fmt.Println("  -build-all   Build for multiple platforms")
	fmt.Println("  -release     Create optimized release build")
	fmt.Println("  -help        Show this help message")
}

func build() {
	fmt.Println("Building TechIRCd...")
	
	gitVersion, err := exec.Command("git", "describe", "--tags", "--always", "--dirty").Output()
	var versionStr string
	if err != nil {
		versionStr = version
	} else {
		versionStr = strings.TrimSpace(string(gitVersion))
	}
	
	ldflags := fmt.Sprintf("-ldflags=-X main.version=%s", versionStr)
	
	cmd := exec.Command("go", "build", ldflags, "-o", binaryName, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Build completed successfully!")
}

func run() {
	fmt.Println("Starting TechIRCd...")
	
	cmd := exec.Command("./" + binaryName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Run failed: %v\n", err)
		os.Exit(1)
	}
}

func test() {
	fmt.Println("Running tests...")
	
	cmd := exec.Command("go", "test", "-v", "-race", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("All tests passed!")
}

func clean() {
	fmt.Println("Cleaning build artifacts...")
	
	// Remove binary files
	patterns := []string{
		binaryName + "*",
		"coverage.out",
		"coverage.html",
	}
	
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		
		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				fmt.Printf("Failed to remove %s: %v\n", match, err)
			} else {
				fmt.Printf("Removed %s\n", match)
			}
		}
	}
	
	fmt.Println("Clean completed!")
}

func format() {
	fmt.Println("Formatting Go code...")
	
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Format failed: %v\n", err)
		os.Exit(1)
	}
	
	// Try to run goimports if available
	if _, err := exec.LookPath("goimports"); err == nil {
		fmt.Println("Running goimports...")
		cmd := exec.Command("goimports", "-w", "-local", "github.com/ComputerTech312/TechIRCd", ".")
		cmd.Run() // Don't fail if this doesn't work
	}
	
	fmt.Println("Format completed!")
}

func lint() {
	fmt.Println("Running linters...")
	
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		fmt.Println("golangci-lint not found, skipping...")
		return
	}
	
	cmd := exec.Command("golangci-lint", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Linting found issues: %v\n", err)
		// Don't exit on lint errors, just report them
	} else {
		fmt.Println("No linting issues found!")
	}
}

func buildAll() {
	fmt.Println("Building for multiple platforms...")
	
	platforms := []struct {
		goos   string
		goarch string
		ext    string
	}{
		{"linux", "amd64", ""},
		{"windows", "amd64", ".exe"},
		{"darwin", "amd64", ""},
		{"darwin", "arm64", ""},
	}
	
	gitVersion, err := exec.Command("git", "describe", "--tags", "--always", "--dirty").Output()
	var versionStr string
	if err != nil {
		versionStr = version
	} else {
		versionStr = strings.TrimSpace(string(gitVersion))
	}
	
	for _, platform := range platforms {
		outputName := fmt.Sprintf("%s-%s-%s%s", binaryName, platform.goos, platform.goarch, platform.ext)
		fmt.Printf("Building %s...\n", outputName)
		
		ldflags := fmt.Sprintf("-ldflags=-X main.version=%s", versionStr)
		
		cmd := exec.Command("go", "build", ldflags, "-o", outputName, ".")
		cmd.Env = append(os.Environ(),
			"GOOS="+platform.goos,
			"GOARCH="+platform.goarch,
		)
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to build %s: %v\n", outputName, err)
		} else {
			fmt.Printf("Built %s successfully!\n", outputName)
		}
	}
	
	fmt.Println("Cross-platform build completed!")
}

func release() {
	fmt.Println("Creating optimized release build...")
	
	gitVersion, err := exec.Command("git", "describe", "--tags", "--always", "--dirty").Output()
	var versionStr string
	if err != nil {
		versionStr = version
	} else {
		versionStr = strings.TrimSpace(string(gitVersion))
	}
	
	ldflags := fmt.Sprintf("-ldflags=-X main.version=%s", versionStr)
	
	cmd := exec.Command("go", "build", ldflags, "-a", "-installsuffix", "cgo", "-o", binaryName, ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("Release build failed: %v\n", err)
		os.Exit(1)
	}
	
	// Get file info to show size
	if info, err := os.Stat(binaryName); err == nil {
		fmt.Printf("Release build completed! Binary size: %.2f MB\n", float64(info.Size())/1024/1024)
	} else {
		fmt.Println("Release build completed!")
	}
}
