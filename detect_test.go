package main_test

import (
	"errors"
	"testing"

	main "github.com/paketo-buildpacks/go-compiler"
	"github.com/paketo-buildpacks/go-compiler/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		buildpackYAMLParser *fakes.VersionParser

		detect packit.DetectFunc
	)

	it.Before(func() {
		buildpackYAMLParser = &fakes.VersionParser{}

		detect = main.Detect(buildpackYAMLParser)
	})

	it("returns a plan that provides go", func() {
		result, err := detect(packit.DetectContext{
			WorkingDir: "/working-dir",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "go"},
				},
			},
		}))

		Expect(buildpackYAMLParser.ParseVersionCall.Receives.Path).To(Equal("/working-dir/buildpack.yml"))
	})

	context("when the buildpack.yml contains a version", func() {
		it.Before(func() {
			buildpackYAMLParser.ParseVersionCall.Returns.Version = "some-version"
		})

		it("returns a plan that requires that version of go", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: "/working-dir",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: "go"},
					},
					Requires: []packit.BuildPlanRequirement{
						{
							Name:    "go",
							Version: "some-version",
							Metadata: main.BuildPlanMetadata{
								VersionSource: "buildpack.yml",
							},
						},
					},
				},
			}))

			Expect(buildpackYAMLParser.ParseVersionCall.Receives.Path).To(Equal("/working-dir/buildpack.yml"))
		})
	})

	context("failure cases", func() {
		context("when the buildpack YAML parser errors", func() {
			it.Before(func() {
				buildpackYAMLParser.ParseVersionCall.Returns.Err = errors.New("failed to parse buildpack.yml")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: "/working-dir",
				})
				Expect(err).To(MatchError("failed to parse buildpack.yml"))
			})
		})
	})
}
