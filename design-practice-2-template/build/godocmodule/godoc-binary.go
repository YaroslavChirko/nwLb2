package godocmodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/YaroslavChirko/design-practice-2-template/build/godocmodule")

	// Ninja rule to execute go build.
	goDoc= pctx.StaticRule("godoc", blueprint.RuleParams{
		Command:     "cd $workDir&& mkdir -p ./out/docs && godoc -url=./>'./${outputPath}.html' -goroot='./${pkg}'",
		Description: "command to create html file",
	}, "workDir", "outputPath", "name", "pkg")

)


type goDocModule struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		// Example of how to specify dependencies.
		Deps []string
	}
}

func (gb *goDocModule) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return gb.properties.Deps
}

func (gb *goDocModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "docs", name)

	var inputs []string
	inputErors := false
	for _, src := range gb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gb.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}
	
	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Godoc"),
		Rule:        goDoc,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gb.properties.Pkg,
			"name":	name,
		},
	})
	
	
	

}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &goDocModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
