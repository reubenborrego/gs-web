package web

import "path/filepath"

type Resolver interface {
	directory() string
	pathPrefix(path string) Resolver
	pathPostfix(path string) Resolver
	resolve([]string) string
}

type FileResolver struct {
	path string
}

func NewFileResolver(path string) *FileResolver {
	return &FileResolver{path: path}
}

func (resolver *FileResolver) pathPrefix(path string) Resolver {
	resolver.path = filepath.Join(path, resolver.path)
	return resolver
}

func (resolver *FileResolver) pathPostfix(path string) Resolver {
	resolver.path = filepath.Join(resolver.path, path)
	return resolver
}

func (resolver FileResolver) directory() string {
	return resolver.path
}

func (resolver FileResolver) resolve(path []string) string {
	fileName := path[len(path)-1]
	return filepath.Join(resolver.path, fileName)
}

type ConstantResolver struct {
	FileResolver
	literal string
}

func NewConstantResolver(path string, literal string) *ConstantResolver {
	return &ConstantResolver{FileResolver: FileResolver{path: path}, literal: literal}
}

func (resolver ConstantResolver) resolve(path []string) string {
	return filepath.Join(resolver.path, resolver.literal)
}
