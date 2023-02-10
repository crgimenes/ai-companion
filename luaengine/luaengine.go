package luaengine

import (
	"bufio"
	"os"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type LuaExtender struct {
	luaState *lua.LState
	proto    *lua.FunctionProto
}

func New() *LuaExtender {
	return &LuaExtender{
		luaState: lua.NewState(),
	}
}

// CompileLua reads the passed lua file from disk and compiles it.
func (le *LuaExtender) Compile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return err
	}

	le.proto, err = lua.Compile(chunk, filePath)
	return err
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func (le *LuaExtender) DoCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}

// InitState starts the lua interpreter with a script.
func (le *LuaExtender) InitState() error {
	return le.DoCompiledFile(le.luaState, le.proto)
}

func (le *LuaExtender) Close() error {
	le.luaState.Close()
	return nil
}
