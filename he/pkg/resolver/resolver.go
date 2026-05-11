package resolver

import (
"errors"
"fmt"
"strings"
)

// Platform represents a target platform.
type Platform string

const (
PlatformWeb    Platform = "web"
PlatformMobile Platform = "mobile"
PlatformDesktop Platform = "desktop"
)

// Resolver handles platform-aware module/function resolution.
type Resolver struct {
platform   Platform
modules    map[string]*Module
precedence []Platform // resolution precedence order
}

// Module represents a resolved module.
type Module struct {
Name       string
Functions  map[string]bool
Submodules map[string]*Module
}

// New creates a new resolver.
func New(platform Platform) *Resolver {
r := &Resolver{
platform: platform,
modules:  make(map[string]*Module),
}

// Set platform precedence based on current platform
switch platform {
case PlatformMobile:
r.precedence = []Platform{PlatformMobile, PlatformWeb, PlatformDesktop}
case PlatformWeb:
r.precedence = []Platform{PlatformWeb, PlatformMobile, PlatformDesktop}
case PlatformDesktop:
r.precedence = []Platform{PlatformDesktop, PlatformWeb, PlatformMobile}
default:
r.precedence = []Platform{PlatformWeb, PlatformMobile, PlatformDesktop}
}

return r
}

// RegisterModule registers a module with the resolver.
func (r *Resolver) RegisterModule(name string, module *Module) {
r.modules[name] = module
}

// ResolveFunction resolves a function call with platform-aware routing.
func (r *Resolver) ResolveFunction(path string) (string, error) {
// Parse the path: ui.navbar or ui.mobile.navbar
parts := strings.Split(path, ".")
if len(parts) < 2 {
return "", fmt.Errorf("invalid path: %s", path)
}

funcName := parts[len(parts)-1]
modulePath := parts[:len(parts)-1]

// Try exact match first
exactPath := strings.Join(modulePath, ".")
if module, ok := r.modules[exactPath]; ok {
if _, hasFunc := module.Functions[funcName]; hasFunc {
return exactPath + "." + funcName, nil
}
}

// Platform-aware resolution for single-level modules
// e.g., ui.navbar -> try ui.mobile.navbar, ui.web.navbar, ui.desktop.navbar
if len(modulePath) == 1 {
baseModule := modulePath[0]

// Try each platform in precedence order
for _, platform := range r.precedence {
platformPath := baseModule + "." + string(platform)
if module, ok := r.modules[platformPath]; ok {
if _, hasFunc := module.Functions[funcName]; hasFunc {
return platformPath + "." + funcName, nil
}
}
}

// Try the base module with platform submodule
if baseModuleObj, ok := r.modules[baseModule]; ok {
for _, platform := range r.precedence {
if submodule, ok := baseModuleObj.Submodules[string(platform)]; ok {
if _, hasFunc := submodule.Functions[funcName]; hasFunc {
return baseModule + "." + string(platform) + "." + funcName, nil
}
}
}
}
}

// Try nested module resolution
// e.g., ui.components.navbar
for i := 1; i < len(modulePath); i++ {
parentPath := strings.Join(modulePath[:i], ".")
childName := modulePath[i]

if parent, ok := r.modules[parentPath]; ok {
if child, ok := parent.Submodules[childName]; ok {
if _, hasFunc := child.Functions[funcName]; hasFunc {
return parentPath + "." + childName + "." + funcName, nil
}
}
}
}

return "", fmt.Errorf("function not found: %s", path)
}

// ResolveModule resolves a module path to its canonical form.
func (r *Resolver) ResolveModule(path string) (string, error) {
parts := strings.Split(path, ".")

// Try exact match
if _, ok := r.modules[path]; ok {
return path, nil
}

// Platform-aware resolution for single module
if len(parts) == 1 {
baseModule := parts[0]

// Check if this module has platform-specific versions
for _, platform := range r.precedence {
platformPath := baseModule + "." + string(platform)
if _, ok := r.modules[platformPath]; ok {
return platformPath, nil
}
}
}

return "", fmt.Errorf("module not found: %s", path)
}
