# Booth GIF for the README.
# just vhs booth  (builds bin/cans first)
# Setup stays hidden. The tape is what a user types: cans.

Output docs/booth.gif

Require bash

Set Shell "bash"
Set FontSize 16
Set Width 680
Set Height 300
Set Padding 12
Set TypingSpeed 55ms
Set PlaybackSpeed 0.9
Set CursorBlink false
Set Theme {"name":"cans","background":"#1a1014","foreground":"#F8EDE8","cursor":"#C45C7A","selection":"#3a2028","black":"#1a1014","red":"#C45C7A","green":"#50FA7B","yellow":"#E8C4A8","blue":"#8B6B8A","magenta":"#C45C7A","cyan":"#A890A0","white":"#F8EDE8","brightBlack":"#6a5560","brightRed":"#E07070","brightGreen":"#50FA7B","brightYellow":"#E8C4A8","brightBlue":"#8B6B8A","brightMagenta":"#E8A0B4","brightCyan":"#C4B0B8","brightWhite":"#FFFFFF"}

Env TERM "xterm-256color"
Env COLORTERM "truecolor"
Env CLICOLOR_FORCE "1"
Env CANS_NOPLAY "1"
Env CANS_SAY_BIN "tapes/fake-say"

Hide
Type `export PATH="$PWD/bin:$PATH" PS1="> "`
Enter
Ctrl+L
Type "cans"
Enter
Sleep 1400ms
Show

Sleep 500ms
Type "Put the cans on."
Sleep 350ms
Enter
Sleep 2400ms
