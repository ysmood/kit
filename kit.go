package kit

import (
	"github.com/ysmood/gokit/pkg/http"
	"github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/run"
	"github.com/ysmood/gokit/pkg/utils"
)

// All imported
var All = utils.All

// E imported
var E = utils.E

// E1 imported
var E1 = utils.E1

// ErrArg imported
var ErrArg = utils.ErrArg

// JSON imported
var JSON = utils.JSON

// Nil imported
type Nil = utils.Nil

// Noop imported
var Noop = utils.Noop

// RandBytes imported
var RandBytes = utils.RandBytes

// RandString imported
var RandString = utils.RandString

// S imported
var S = utils.S

// Try imported
var Try = utils.Try

// Version imported
var Version = utils.Version

// GinContext imported
type GinContext = http.GinContext

// MustServer imported
var MustServer = http.MustServer

// Req imported
var Req = http.Req

// ReqContext imported
type ReqContext = http.ReqContext

// Server imported
var Server = http.Server

// ServerContext imported
type ServerContext = http.ServerContext

// C imported
var C = os.C

// Chmod imported
var Chmod = os.Chmod

// ClearScreen imported
var ClearScreen = os.ClearScreen

// Copy imported
var Copy = os.Copy

// DirExists imported
var DirExists = os.DirExists

// Dump imported
var Dump = os.Dump

// Err imported
var Err = os.Err

// ExecutableExt imported
var ExecutableExt = os.ExecutableExt

// Exists imported
var Exists = os.Exists

// FileExists imported
var FileExists = os.FileExists

// GoPath imported
var GoPath = os.GoPath

// HomeDir imported
var HomeDir = os.HomeDir

// Log imported
var Log = os.Log

// Matcher imported
type Matcher = os.Matcher

// Mkdir imported
var Mkdir = os.Mkdir

// MkdirOptions imported
type MkdirOptions = os.MkdirOptions

// Move imported
var Move = os.Move

// NewMatcher imported
var NewMatcher = os.NewMatcher

// OutputFile imported
var OutputFile = os.OutputFile

// OutputFileOptions imported
type OutputFileOptions = os.OutputFileOptions

// ReadFile imported
var ReadFile = os.ReadFile

// ReadJSON imported
var ReadJSON = os.ReadJSON

// ReadString imported
var ReadString = os.ReadString

// Remove imported
var Remove = os.Remove

// Retry imported
var Retry = os.Retry

// Sdump imported
var Sdump = os.Sdump

// SendSigInt imported
var SendSigInt = os.SendSigInt

// Stderr imported
var Stderr = os.Stderr

// Stdout imported
var Stdout = os.Stdout

// ThisDirPath imported
var ThisDirPath = os.ThisDirPath

// ThisFilePath imported
var ThisFilePath = os.ThisFilePath

// WaitSignal imported
var WaitSignal = os.WaitSignal

// Walk imported
var Walk = os.Walk

// WalkContext imported
type WalkContext = os.WalkContext

// WalkDirent imported
type WalkDirent = os.WalkDirent

// WalkFunc imported
type WalkFunc = os.WalkFunc

// WalkGitIgnore imported
var WalkGitIgnore = os.WalkGitIgnore

// WalkIgnoreHidden imported
var WalkIgnoreHidden = os.WalkIgnoreHidden

// Exec imported
var Exec = run.Exec

// ExecContext imported
type ExecContext = run.ExecContext

// Guard imported
var Guard = run.Guard

// GuardContext imported
type GuardContext = run.GuardContext

// GuardDefaultPatterns imported
var GuardDefaultPatterns = run.GuardDefaultPatterns

// KillTree imported
var KillTree = run.KillTree

// MustGoTool imported
var MustGoTool = run.MustGoTool

// Task imported
var Task = run.Task

// TaskCmd imported
type TaskCmd = run.TaskCmd

// TaskContext imported
type TaskContext = run.TaskContext

// Tasks imported
var Tasks = run.Tasks

// TasksContext imported
type TasksContext = run.TasksContext

// TasksNew imported
var TasksNew = run.TasksNew
