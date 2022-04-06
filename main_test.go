package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-search-data-finder/features/steps"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	component, err := steps.NewSearchDataFinderComponent()
	if err != nil {
		fmt.Printf("failed to create search data finder component - error: %v", err)
		os.Exit(1)
	}

	apiFeature := component.InitAPIFeature()

	beforeScenarioHook := func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		apiFeature.Reset()
		component.Reset()
		return ctx, nil
	}

	ctx.Before(beforeScenarioHook)

	apiFeature.RegisterSteps(ctx)
	component.RegisterSteps(ctx)

	afterScenarioHook := func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// component.Close()
		return ctx, nil
	}

	ctx.After(afterScenarioHook)
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		status := 0

		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{}

		status = godog.TestSuite{
			Name:                "feature_tests",
			ScenarioInitializer: f.InitializeScenario,
			Options:             &opts,
		}.Run()

		fmt.Println("=================================")
		fmt.Printf("Component test coverage: %.2f%%\n", testing.Coverage()*100)
		fmt.Println("=================================")

		if status > 0 {
			t.FailNow()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
