- [Conversion](#conversion)
  - [Introducing](#introducing)
  - [Usage](#usage)
  - [References](#references)

# Conversion

## Introducing

提供一个用来保存所有类型间转换函数的容器，用于简化在多个类型间转换时的包重复引用问题，以便对类型转换函数进行复用。

## Usage

Converter 定义如下：

```go
// Converter knows how to convert one type to another.
type Converter interface {
	// Convert will attempt to convert in into out. Both must be pointers.
	Convert(in interface{}, out interface{}, context interface{}) error

	// RegisterConversionFunc registers a function that converts between a and b by passing objects of those
	// types to the provided function. The function *must* accept objects of a and b - this machinery will not enforce
	// any other guarantee.
	RegisterConversionFunc(in interface{}, out interface{}, fn ConversionFunc) error

	// RegisterGenericConversionFunc registers a function that convert objects.
	// It will been called if there is no typed conversion func exists.
	// NOTE: All generic conversion funcs will been called one by one, if the function returns ErrUnspport
	//    the next one will been called. Otherwise Convert will return the result.
	RegisterGenericConversionFunc(fn ConversionFunc)
}
```

使用方式：

1. 生成一个 Converter 对象，同时该包提供一个 DefaultConverter 对象以供使用
2. 注册各种 ConversionFunc 函数，用于实现类型间的转换
3. 如果有必要，可以注册 GenericConversionFunc 函数，这些函数会在没有可用的 ConversionFunc 函数时被调用
4. 使用 Converter 对象的 Convert 函数，完成类型转换

一个简单的示例如下:
```go
type testConvert1 struct {
	Total1 int
}
type testConvert2 struct {
	Total2 int
}
var fn = func(in *testConvert1, out *testConvert2, _ Scope) error { //nolint
	out.Total2 = in.Total1
	return nil
}
    
c := NewConverter()
c.RegisterConversionFunc((*testConvert1)(nil), (*testConvert2)(nil), func(in interface{}, out interface{}, scope Scope) error {
	return fn(in.(*testConvert1), out.(*testConvert2), scope)
})

in := &testConvert1{
	Total1: 1,
}
out := testConvert2{
	Total2: 2,
}
err := c.Convert(in, out, nil)
```

## References

