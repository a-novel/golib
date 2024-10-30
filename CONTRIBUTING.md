# Contributing

> This document is a WIP for the project. We are still working a lot of things for now (also I was bored).

## Reporting bugs

There is never too much information. Provide everything you can, too much is preferable to not enough. Also
include reproducible steps if possible. Investigation on your own is a bonus.

## Requesting features

Go to discussions for this. Look for someone who requested a similar feature, before opening anything.

## Submitting code

NEVER submit anything before your idea was even approved, be it discussions, issues, or private communication with the
maintainers. We will not accept any pull request that was not discussed beforehand.

### Readme first

Everything you need to install for the project to run is already included there. If something is lacking, then
it should be noticed to the maintainers and added.

### Write tests

We don't care about TDD. Write code and tests in the order you want. The only important thing for your contribution
is to be covered past the required coverage (80% of new code). Sometimes it is not possible, it will be up to
maintainers to decide if your work is still mergeable. But there is never too much tests.

### Read the code

The best way to follow coding practices is to read the previous code and understand it. Try to blend in with what is 
already there, not to stand out.

### Magic make

```bash
make format
make lint
make test
```

As long as those commands succeed, and your code is properly covered, then you are guaranteed to have your code
merged. Those commands are your new Ctrl+S. Learn to love them.

### Comment, comment, comment

Your audience (maintainers) have the maturity of 8yo. Add documentation to every interface or function, with 
explaination about what it does and how to use it. Morever, if you write some complex logic somewhere, explain. No one 
should ever need to ask you about a piece of code, ever.

### Write. Actual. Names.

2-letters variables or acronyms are a headache. Be as specific as possible when naming variables or functions.
Again, it might seem obvious in your context, but someone new might lack information. Especially:

- In for loops, write actual names
  ```go
  ❌ for k, v := range items
  
  ✅ for itemName, itemData := range items
  ```
  Even when looping arrays, take time to give proper names
  ```go
  ❌ for i, v := range items
  
  ✅ for index, item := range items
  ```
- Same goes for receivers name on a struct method
  ```go
  ❌ func (h *myHandler) DoSomething(...)
  
  ✅ func (handler *myHandler) DoSomething(...)
  ```

### Interfaces all the way

Duck-typing is very powerful in Go, especially for writing tests. Whenever you need to expose more than a variable
or function, make it an interface. Then don't export the implementation. That way, other project can easily generate
their own mocks for it, or even reuse yours!

```go
// ❌ Exported struct

type MyHandler struct {
  ...
}

func (h *MyHandler) DoSomething(...) {
  ...
}

// ✅ Hide behind an interface

type MyHandler interface {
  DoSomething(...)
}

type myHandlerImpl struct {
  ...
}

func (h *myHandlerImpl) DoSomething(...) {
  ...
}

func NewMyHandler() MyHandler {
  return &myHandlerImpl{}
}
```
