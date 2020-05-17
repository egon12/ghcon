package coverage

import (
	"path"

	epath "github.com/egon12/ghr/path"

	"golang.org/x/tools/cover"
)

type Range struct {
	From int
	To   int
}

type GoCoverageInGit interface {
	Percentage() float32
	PercentagePackage(packageName string) float32
	PercentageFile(filename string) float32
	NotInCoverageLines(filename string) []Range
}

func FromProfile(profileFileName string) (GoCoverageInGit, error) {
	profiles, err := cover.ParseProfiles(profileFileName)
	return &goCoverageInGit{
		profiles: profiles,
	}, err
}

type goCoverageInGit struct {
	profiles []*cover.Profile
}

func (g *goCoverageInGit) Percentage() float32 {
	var total, covered int
	for _, p := range g.profiles {
		c, t := countCovered(p)
		total += t
		covered += c
	}
	if total == 0 {
		return 0
	}
	return float32(covered) / float32(total) * 100
}

func (g *goCoverageInGit) PercentagePackage(packageName string) float32 {
	var total, covered int
	for _, p := range g.profiles {
		if path.Dir(p.FileName) == packageName {
			c, t := countCovered(p)
			total += t
			covered += c
		}
	}
	if total == 0 {
		return 0
	}
	return float32(covered) / float32(total) * 100
}

func (g *goCoverageInGit) PercentageFile(fileName string) float32 {
	for _, p := range g.profiles {
		if p.FileName == fileName {
			return float32(percentCovered(p))
		}
	}
	return 0
}

func (g *goCoverageInGit) NotInCoverageLines(fileName string) []Range {
	var result []Range

	f, err := epath.GetFileWithPackagePath(fileName)
	if err != nil {
		return result
	}

	for _, p := range g.profiles {
		if p.FileName == f {
			for _, b := range p.Blocks {
				if b.Count != 0 {
					continue
				}
				r := Range{
					From: b.StartLine,
					To:   b.EndLine,
				}
				result = append(result, r)
			}
		}
	}
	return result
}

// Copied from /src/cmd/cover/html.go
// percentCovered returns, as a percentage, the fraction of the statements in
// the profile covered by the test run.
// In effect, it reports the coverage of a given source file.
func percentCovered(p *cover.Profile) float64 {
	var total, covered int
	for _, b := range p.Blocks {
		total += b.NumStmt
		if b.Count > 0 {
			covered += b.NumStmt
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}

func countCovered(p *cover.Profile) (covered, total int) {
	for _, b := range p.Blocks {
		total += b.NumStmt
		if b.Count > 0 {
			covered += b.NumStmt
		}
	}
	return
}
