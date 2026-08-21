# Pipe GIF for the README.
# just vhs pipe  (builds bin/cans first)
# Real mouth: no CANS_SAY_BIN, no fake worker. Three lines in, three wavs out.
# The setup stays hidden; the tape is what a user types.

Output docs/pipe.gif

Require bash

Set Shell "bash"
Set FontSize 16
Set Width 680
Set Height 300
Set Padding 12
Set TypingSpeed 55ms
Set PlaybackSpeed 1.0
Set Framerate 12
Set CursorBlink false
Set Theme {"name":"cans","background":"#1a1014","foreground":"#F8EDE8","cursor":"#C45C7A","selection":"#3a2028","black":"#1a1014","red":"#C45C7A","green":"#50FA7B","yellow":"#E8C4A8","blue":"#8B6B8A","magenta":"#C45C7A","cyan":"#A890A0","white":"#F8EDE8","brightBlack":"#6a5560","brightRed":"#E07070","brightGreen":"#50FA7B","brightYellow":"#E8C4A8","brightBlue":"#8B6B8A","brightMagenta":"#E8A0B4","brightCyan":"#C4B0B8","brightWhite":"#FFFFFF"}

Env TERM "xterm-256color"
Env COLORTERM "truecolor"
Env CLICOLOR_FORCE "1"

Hide
Type `export PATH="$PWD/bin:$PATH" PS1="> "`
Enter
Type `cd "$(mktemp -d)"`
Enter
Sleep 500ms
Ctrl+L
Show

Sleep 500ms
Type `printf 'Put the cans on.\nOne worker, one model load.\nFiles land where the script points.\n' > lines.txt`
Sleep 350ms
Enter
Sleep 900ms
Type `cat lines.txt | cans say --stream -o 'out/%03d.wav' --json`
Sleep 350ms
Enter
Wait+Screen@300s /003\.wav/
Sleep 900ms
Type `ls out`
Sleep 350ms
Enter
Wait+Screen@30s /(?m)^001\.wav/
Sleep 2500ms
