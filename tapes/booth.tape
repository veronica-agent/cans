# Booth GIF for the README.
# just vhs booth  (builds bin/cans first)

Output docs/booth.gif

Set Shell "bash"
Set FontSize 16
Set Width 780
Set Height 420
Set Padding 18
Set TypingSpeed 65ms
Set PlaybackSpeed 0.9
Set Theme {"name":"cans","background":"#1a1014","foreground":"#F8EDE8","cursor":"#C45C7A","black":"#1a1014","red":"#C45C7A","green":"#50FA7B","yellow":"#E8C4A8","blue":"#8B6B8A","magenta":"#C45C7A","cyan":"#A890A0","white":"#F8EDE8"}

Hide
Type `export CANS_ROOT="$PWD" CANS_NOPLAY=1 CANS_SAY_BIN="$PWD/tapes/fake-say"`
Enter
Type "chmod +x tapes/fake-say"
Enter
Type "clear"
Enter
Show
Sleep 400ms
Type `CANS_ROOT="$PWD" CANS_NOPLAY=1 CANS_SAY_BIN="$PWD/tapes/fake-say" ./bin/cans`
Enter
Sleep 1200ms
Type "Put the cans on."
Sleep 400ms
Enter
Sleep 2200ms
