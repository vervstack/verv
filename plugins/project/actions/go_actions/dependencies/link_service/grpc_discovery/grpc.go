package grpc_discovery

import (
	"bytes"
	stderrs "errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"go.redsock.ru/rerrors"

	rscliconfig "go.vervstack.ru/verv/internal/config"
)

type GrpcDiscovery struct {
	Cfg *rscliconfig.VervConfig
}

var defaultDiscoverer GrpcDiscovery

func DiscoverPackage(packageName string) (*GrpcPackage, error) {
	if defaultDiscoverer.Cfg == nil {
		defaultDiscoverer.Cfg = rscliconfig.GetConfig()
	}

	return defaultDiscoverer.DiscoverPackage(packageName)
}

func (g GrpcDiscovery) DiscoverPackage(packageName string) (*GrpcPackage, error) {
	pathToPackageInMod, err := GetPathToGlobalModule(packageName)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting path to package in global module")
	}

	grpcPackage, err := g.getGrpcPackageFromMod(pathToPackageInMod)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting grpc package from mod")
	}

	return grpcPackage, nil
}

type GrpcPackage struct {
	ImportPath  string
	Constructor string
	ClientName  string
}

func (g GrpcDiscovery) getGrpcPackageFromMod(packagePath string) (*GrpcPackage, error) {
	var errs error

	for _, compiledClientPath := range g.Cfg.Env.PathsToCompiledClients {
		apiPath := path.Join(packagePath, compiledClientPath)
		files, err := os.ReadDir(apiPath)
		if err != nil {
			errs = stderrs.Join(errs, err)
		}

		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			var pkg *GrpcPackage
			pkg, err = readGrpcPackageFromPackageClientPath(packagePath, path.Join(compiledClientPath, file.Name()))
			if err != nil {
				return nil, rerrors.Wrap(err, "error reading package from path")
			}

			if pkg != nil {
				packagePath = packagePath[len(modFolderPath):]
				atIdx := strings.Index(packagePath, "@")
				if atIdx != -1 {
					packagePath = packagePath[:atIdx]
				}

				pkg.ImportPath = path.Join(packagePath, pkg.ImportPath)

				return pkg, nil
			}
		}
	}

	return nil, nil
}

func readGrpcPackageFromPackageClientPath(projectPath, apiContractPath string) (*GrpcPackage, error) {
	packagePath := path.Join(projectPath, apiContractPath)

	clientContractPath, err := findGrpcClientContractFile(packagePath)
	if err != nil {
		return nil, err
	}
	if clientContractPath == "" {
		return nil, nil
	}

	clientContractB, err := os.ReadFile(clientContractPath)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading client contract files")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path.Base(clientContractPath), clientContractB, 0)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing go contract file")
	}

	out := findClientNameAndConstructor(f, clientContractB)
	if out.Constructor == "" {
		return nil, nil
	}

	packageName, ok := extractPackageName(f, clientContractB)
	if !ok {
		return nil, nil
	}

	out.ImportPath = path.Join(path.Dir(apiContractPath), packageName)

	return out, nil
}

func findGrpcClientContractFile(packagePath string) (string, error) {
	files, err := os.ReadDir(packagePath)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading contracts dir")
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		fileName := f.Name()
		if strings.HasSuffix(fileName, "_grpc.pb.go") {
			return path.Join(packagePath, fileName), nil
		}
	}

	return "", nil
}

// findClientNameAndConstructor walks the parsed grpc contract file looking for the
// Client interface type and its constructor func, stopping as soon as both are found.
func findClientNameAndConstructor(f *ast.File, src []byte) *GrpcPackage {
	out := &GrpcPackage{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.GenDecl:
			inspectClientTypeDecl(fn, out)
		case *ast.FuncDecl:
			inspectClientConstructorDecl(fn, src, out)
		}

		return out.ClientName == "" || out.Constructor == ""
	})

	return out
}

func inspectClientTypeDecl(fn *ast.GenDecl, out *GrpcPackage) {
	if out.ClientName != "" || fn.Tok != token.TYPE || len(fn.Specs) == 0 {
		return
	}

	spec, ok := fn.Specs[0].(*ast.TypeSpec)
	if !ok {
		return
	}

	if strings.HasSuffix(spec.Name.Name, "Client") {
		out.ClientName = spec.Name.Name
	}
}

func inspectClientConstructorDecl(fn *ast.FuncDecl, src []byte, out *GrpcPackage) {
	if out.Constructor != "" {
		return
	}

	if !strings.HasPrefix(fn.Name.Name, "New") || !strings.HasSuffix(fn.Name.Name, "Client") {
		return
	}

	if fn.Type == nil || fn.Type.Params == nil || len(fn.Type.Params.List) != 1 || fn.Type.Params.List[0] == nil {
		return
	}

	startIdx, endIdx := int(fn.Type.Params.List[0].Pos()), int(fn.Type.Params.List[0].End())
	if bytes.Contains(src[startIdx:endIdx], []byte("grpc.ClientConnInterface")) {
		out.Constructor = fn.Name.Name
	}
}

func extractPackageName(f *ast.File, src []byte) (string, bool) {
	if !f.Package.IsValid() {
		return "", false
	}

	packageSubSet := src[f.Package:]
	packageSubSet = packageSubSet[:bytes.IndexByte(packageSubSet, '\n')]
	packageSubSet = packageSubSet[bytes.IndexByte(packageSubSet, ' ')+1:]

	return string(packageSubSet), true
}
