package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if file.Generate {
				generateFile(plugin, file)
			}
		}
		return nil
	})
}

func generateFile(plugin *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + ".qk.go"
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-qkrpc. DO NOT EDIT.")
	g.P("package ", file.GoPackageName)
	g.P()

	g.P("import (")
	g.P("  \"context\"")
	g.P("  \"github.com/ajaypanthagani/qkrpc\"")
	g.P("  \"github.com/quic-go/quic-go\"")
	g.P(")")
	g.P()

	for _, service := range file.Services {
		genService(g, file, service)
	}
}

func genService(g *protogen.GeneratedFile, file *protogen.File, service *protogen.Service) {
	serviceName := service.GoName
	fullServiceName := fmt.Sprintf("%s.%s", file.Desc.Package(), service.Desc.Name())

	// Service Interface
	g.P("type ", serviceName, " interface {")
	for _, method := range service.Methods {
		g.P(method.GoName, "(ctx context.Context, req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ", error)")
	}
	g.P("}")
	g.P()

	// Server registration
	regFuncName := "Register" + serviceName
	g.P("func ", regFuncName, "(s qkrpc.QkServer, impl ", serviceName, ") {")
	for _, method := range service.Methods {
		methodName := fmt.Sprintf("%s.%s", fullServiceName, method.Desc.Name())
		g.P("  s.RegisterHandler(\"", methodName, "\", func(ctx context.Context, stream *quic.Stream) error {")
		g.P("    var req ", method.Input.GoIdent)
		g.P("    if err := qkrpc.ReadProtobuf(stream, &req); err != nil {")
		g.P("      return err")
		g.P("    }")
		g.P("    resp, err := impl.", method.GoName, "(ctx, &req)")
		g.P("    if err != nil {")
		g.P("      return err")
		g.P("    }")
		g.P("    return qkrpc.WriteProtobuf(stream, resp)")
		g.P("  })")
	}
	g.P("}")
	g.P()

	// Client stub
	clientName := serviceName + "Client"
	g.P("type ", clientName, " struct {")
	g.P("  Conn *qkrpc.ClientConn")
	g.P("}")
	g.P()

	for _, method := range service.Methods {
		methodName := fmt.Sprintf("%s.%s", fullServiceName, method.Desc.Name())
		g.P("func (c *", clientName, ") ", method.GoName, "(ctx context.Context, req *", method.Input.GoIdent, ") (*", method.Output.GoIdent, ", error) {")
		g.P("  stream, err := c.Conn.Call(ctx, \"", methodName, "\")")
		g.P("  if err != nil {")
		g.P("    return nil, err")
		g.P("  }")
		g.P("  if err := qkrpc.WriteProtobuf(stream, req); err != nil {")
		g.P("    return nil, err")
		g.P("  }")
		g.P("  var resp ", method.Output.GoIdent)
		g.P("  if err := qkrpc.ReadProtobuf(stream, &resp); err != nil {")
		g.P("    return nil, err")
		g.P("  }")
		g.P("  return &resp, nil")
		g.P("}")
		g.P()
	}
}
