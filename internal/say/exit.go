package say

// Exit codes. 75 is EX_TEMPFAIL, so xargs and retry loops read a busy mouth
// correctly; 130 is the shell's 128+SIGINT.
const (
	ExitOK          = 0
	ExitFail        = 1
	ExitUsage       = 2
	ExitBusy        = 75
	ExitInterrupted = 130
)
