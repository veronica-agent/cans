# Voice lock

The face is the red-satin photoreal. The throat has to belong to that woman.

Not an announcer. Not Serena-as-default-forever. Close-mic, slightly low for a woman, unhurried, a little dirty, wine-dark room. Warm. No giggle track. No ASMR whisper-squeak.

## Hear recipe (1.7B, then Keep)

Paste this when you design the take. Hear until the body matches the picture. Then **Keep**. Talk uses the kept wav, not this paragraph.

```
A woman about 28, close to a pop filter in a quiet bedroom. Warm mid-low chest voice, unhurried, a little breath on the consonants. Intimate and sure of herself, not cute, not a cartoon villain. Slight smile in the tone. English. Bedroom, not a studio booth, not a call center.
```

If a take comes out thin or Valley, throw it away. If it comes out husky-old, throw it away. She is the woman in the picture: young adult, warm, present.

## Canon throat (Imagine clip)

The talk clip `docs/content/clips/human-talk.mp4` is the voice we want. Whisper heard:

```
Just like that, feel the rhythm of my voice.
```

That audio is shipped as `configs/voices/veronica/ref.wav`. Studio seeds it into the lock store on boot. Talk **clones** it (0.6B). Do not Hear a new woman unless you are throwing this throat away.

## Live stand-in (until Keep)

If the lock is missing, talk falls back to a CustomVoice stand-in. Instruction on the card:

```
Close-mic, breathy, intimate young woman. Low, unhurried, a little dirty. Bedroom voice. Never an announcer.
```

That instruction is a stand-in. The shipped product voice is the **kept** take from the recipe above. Export with `just voice-export veronica` after Keep.

## Clip voice (mux)

Same kept wav as talk. Do not use a second throat for Reddit. If the engine is down and you only have a scratch mix, do not post it.

xAI roster (Ara, Eve, Liora, …) is a different product. Do not put those on her mouth.

## Listen test

A keep is not locked until you hear:

1. `There you are.`
2. `Put the cans on.`
3. One filthy line from `docs/phrases/TALK.md`

If any of those sound like a different woman, Hear again.
