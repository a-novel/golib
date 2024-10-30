package formatters

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"google.golang.org/grpc"
)

type LogDiscoverGRPC interface {
	LogContent
}

type logDiscoverGRPCImpl struct {
	services []grpc.ServiceDesc
	port     int
}

func (logDiscoverGRPC *logDiscoverGRPCImpl) RenderConsole() string {
	servicesList := NewList().
		SetEnumerator(NoEnumerator).
		SetIndenter(func(_ list.Items, _ int) string { return "" }).
		SetStyle(lipgloss.NewStyle().MarginLeft(4)).
		SetItemStyle(lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color("#FFCC66")))

	for _, service := range logDiscoverGRPC.services {
		servicesList.Append(service.ServiceName)

		methods := NewList().
			SetEnumerator(NoEnumerator).
			SetStyle(lipgloss.NewStyle().MarginLeft(4)).
			SetIndenter(func(_ list.Items, _ int) string { return "" }).
			SetItemStyle(lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#FFCC66")))

		for _, method := range service.Methods {
			methods.Append(method.MethodName)
		}
		for _, method := range service.Streams {
			methods.Append("[" + method.StreamName + "]")
		}

		servicesList.Nest(methods)
	}

	description := fmt.Sprintf(
		"%v services registered on port :%v",
		len(logDiscoverGRPC.services), logDiscoverGRPC.port,
	)

	return NewTitle("RPC services successfully registered.").
		SetDescription(description).
		SetChild(servicesList).
		RenderConsole()
}

func (logDiscoverGRPC *logDiscoverGRPCImpl) RenderJSON() interface{} {
	servicesList := map[string]interface{}{}

	for _, service := range logDiscoverGRPC.services {
		methods := make([]interface{}, 0)
		streams := make([]interface{}, 0)

		for _, method := range service.Methods {
			methods = append(methods, method.MethodName)
		}
		for _, stream := range service.Streams {
			streams = append(streams, stream.StreamName)
		}

		servicesList[service.ServiceName] = map[string]interface{}{
			"methods": methods,
			"streams": streams,
		}
	}

	return servicesList
}

func NewDiscoverGRPC(services []grpc.ServiceDesc, port int) LogDiscoverGRPC {
	return &logDiscoverGRPCImpl{
		services: services,
		port:     port,
	}
}
