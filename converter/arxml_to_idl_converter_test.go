package converter

import (
	"testing"

	"arxml-converter/ast"
	"github.com/stretchr/testify/require"
)

func TestArXMLToIDLConverter(t *testing.T) {
	// 创建测试用的 parser
	parser, err := ast.NewParser("../ast/example.xml")
	require.NoError(t, err)
	require.NoError(t, parser.Parse())

	// 创建转换器
	converter := NewArXMLToIDLConverter(parser)

	// 执行转换
	module, err := converter.ConvertToIDLModule()
	require.NoError(t, err)
	require.NotNil(t, module)

	// 验证模块基本信息
	require.Equal(t, "ArXMLDataTypes", module.Name)
	require.Equal(t, "Module", module.Type)
	require.NotEmpty(t, module.Content)

	// 打印转换结果用于调试
	t.Logf("Converted module: %+v", module)
	for i, content := range module.Content {
		t.Logf("Content[%d]: %+v", i, content)
	}
}

func TestArXMLConverter_ToIDLModule(t *testing.T) {
	// 创建 ArXMLConverter
	converter, err := NewConverter("../ast/example.xml", false)
	require.NoError(t, err)
	require.NotNil(t, converter)

	// 执行转换
	module, err := converter.ToIDLModule()
	require.NoError(t, err)
	require.NotNil(t, module)

	// 验证模块基本信息
	require.Equal(t, "ArXMLDataTypes", module.Name)
	require.Equal(t, "Module", module.Type)
	require.NotEmpty(t, module.Content)

	// 打印转换结果用于调试
	t.Logf("Converted module: %+v", module)
	for i, content := range module.Content {
		t.Logf("Content[%d]: %+v", i, content)
	}
}