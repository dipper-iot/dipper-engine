package core

import (
	"fmt"
	"plugin"
)

type PluginDipper interface {
	Rules() []Rule
}

func (d *DipperEngine) LoadPlugin() {
	fmt.Println("Load Plugin...")
	fmt.Println("-----------------------------------------------------------")
	for _, pluginName := range d.config.Plugins {
		raw, err := plugin.Open(pluginName)
		if err != nil {
			fmt.Printf("Load Plugin: %s have error %v \n", pluginName, err)
			continue
		}
		symPlugin, err := raw.Lookup("Plugin")
		if err != nil {
			fmt.Printf("Load Plugin: %s lookup have error %v \n", pluginName, err)
			continue
		}

		pluginDipper, ok := symPlugin.(PluginDipper)
		if !ok {
			fmt.Printf("Load Plugin: %s unexpected type from module symbol \n", pluginName)
			continue
		}

		rules := pluginDipper.Rules()

		for _, rule := range rules {
			d.mapRule[rule.Id()] = rule
		}

		fmt.Printf("Load Plugin: %s success\n", pluginName)
	}
	fmt.Println()
}
