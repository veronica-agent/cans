/* Subset of qwen3-tts-native v0.1.0 C API (symbols in libqwen3tts.0.dylib). */
#ifndef QWEN3TTS_C_API_H
#define QWEN3TTS_C_API_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct Qwen3Tts Qwen3Tts;

typedef struct Qwen3TtsParams {
	int32_t max_audio_tokens;
	float   temperature;
	float   top_p;
	int32_t top_k;
	int32_t n_threads;
	float   repetition_penalty;
	int32_t language_id;
} Qwen3TtsParams;

typedef struct Qwen3TtsAudio {
	const float *samples;
	int32_t      n_samples;
	int32_t      sample_rate;
} Qwen3TtsAudio;

void qwen3_tts_default_params(Qwen3TtsParams *params);
Qwen3Tts *qwen3_tts_create(const char *model_dir, int32_t n_threads);
int qwen3_tts_is_loaded(const Qwen3Tts *tts);
Qwen3TtsAudio *qwen3_tts_synthesize(Qwen3Tts *tts, const char *text, const Qwen3TtsParams *params);
void qwen3_tts_free_audio(Qwen3TtsAudio *audio);
void qwen3_tts_destroy(Qwen3Tts *tts);
Qwen3TtsAudio *qwen3_tts_synthesize_with_voice_file(Qwen3Tts *tts, const char *text, const char *reference_audio_path, const Qwen3TtsParams *params);
Qwen3TtsAudio *qwen3_tts_synthesize_with_embedding(Qwen3Tts *tts, const char *text, const float *embedding, int32_t embedding_size, const Qwen3TtsParams *params);
const char *qwen3_tts_get_error(const Qwen3Tts *tts);

#ifdef __cplusplus
}
#endif

#endif
