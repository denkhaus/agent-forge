// Package main provides a test runner for TUI components.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/denkhaus/agentforge/internal/config"
	"github.com/denkhaus/agentforge/internal/container"
	"github.com/denkhaus/agentforge/internal/tui"
	"go.uber.org/zap"
)

func main() {
	// Command line flags
	var (
		mode     = flag.String("mode", "interactive", "Test mode: interactive, automated, headless, demo")
		prompt   = flag.String("prompt", "test-prompt", "Prompt name to test")
		duration = flag.Duration("duration", 30*time.Second, "Test duration/timeout")
		verbose  = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	// Setup logger
	var logger *zap.Logger
	var err error
	if *verbose {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer logger.Sync()

	// Setup configuration and container
	cfg := &config.Config{
		LogLevel: "info",
	}
	
	diContainer := container.Setup(cfg)
	defer diContainer.Shutdown()

	// Parse test mode
	var testMode tui.TestMode
	switch *mode {
	case "interactive":
		testMode = tui.InteractiveMode
	case "automated":
		testMode = tui.AutomatedMode
	case "headless":
		testMode = tui.HeadlessMode
	case "demo":
		testMode = tui.DemoMode
	default:
		log.Fatal("Invalid mode. Use: interactive, automated, headless, demo")
	}

	// Create test config
	testConfig := tui.TestConfig{
		Mode:     testMode,
		Duration: *duration,
		AutoActions: []string{
			"tab", "wait:1s",
			"tab", "wait:1s", 
			"tab", "wait:1s",
			"tab", "wait:1s",
			"quit",
		},
		OutputCapture: true,
	}

	// Create test manager
	testManager := tui.NewTestManager(logger, testConfig)

	// Run test
	fmt.Printf("🧪 Running TUI test: %s mode for prompt '%s'\n", *mode, *prompt)
	fmt.Printf("⏱️  Duration: %v\n", *duration)
	fmt.Println()

	result, err := testManager.RunWorkbenchTest(*prompt)
	if err != nil {
		logger.Error("Test failed", zap.Error(err))
		fmt.Printf("❌ Test failed: %v\n", err)
		return
	}

	// Print results
	fmt.Println("📊 Test Results:")
	fmt.Println(result.String())
	fmt.Printf("   Duration: %v\n", result.Duration)
	if result.Error != "" {
		fmt.Printf("   Error: %s\n", result.Error)
	}
	if len(result.Actions) > 0 {
		fmt.Printf("   Actions: %v\n", result.Actions)
	}

	if result.Success {
		fmt.Println("\n🎉 All tests passed!")
	} else {
		fmt.Println("\n💥 Tests failed!")
	}
}