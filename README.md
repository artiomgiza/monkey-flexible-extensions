## Make monkey patch easier

Instead of:
```go

    monkey.PatchInstanceMethod(SomeType, "SomeFunc", func(a int, b string, c somePackage.SomeStruct, d []string, e interface{}) error {
        return nil
    })
```

You can use:
```go

    monkeyExtensions.PatchInstanceMethodFlexible(SomeType, "SomeFunc", func() error {
                return nil
    })
```


## Make monkey patch possible

Extension to `PatchInstanceMethod` functionality of `https://github.com/bouk/monkey`

Allows to patch an instance method with providing part (or not providing at all) input parameters.
- must when we **don't have access to the private receiver type**
- might be useful when we interested only in return value

Added `PatchInstanceMethodFlexible`, similar to `PatchInstanceMethod`, but doesn't require input parameters be the same type as patched method input parameters.

Example:
```go

    type MyStruct struct {...}
    func (s *MyStruct) SomeFunc(a int, b string) error {...}
    ...

    my := MyStruct{}

    monkey.PatchInstanceMethodFlexible(reflect.TypeOf(my), "SomeFunc", func( /*NO INPUT PARAMETERS*/) error {
        ...
    })

    // or if we are interested in some input parameter
    monkey.PatchInstanceMethodFlexible(reflect.TypeOf(my), "SomeFunc", func(_, _ interface{}, b string) error {
        // can use "a"
    })
```
In both examples above, the first input param in the replacement func is the "receiver".
If "receiver" is private - widely used when exposing by interface and initiated internally with "private" implementation - regular monkey patch can't be done :/
`PatchInstanceMethodFlexible` allows us to define replacement func without specify **real** inputs type (that sometimes is not accessible), so we can patch also functions with "private" receiver



## Known issues:

1. UnpatchInstanceMethod instance method caused by bug in base [project](https://github.com/bouk/monkey/pull/16).
 - You can use UnpatchAll method instead of UnpatchInstanceMethod
 - Upd: 23.05.2018, unmerged [fix](https://github.com/bouk/monkey/pull/16)).


Tests:

> \> go get -v github.com/onsi/ginkgo/ginkgo<br>
> \> go get -v github.com/onsi/gomega <br>
> \> ginkgo