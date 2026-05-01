use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct ProsodyInput {
    pub transcript: Option<String>,
    pub sample_rate_hz: Option<u32>,
    pub duration_ms: Option<u32>,
}

impl ProsodyInput {
    fn transcript(&self) -> Option<&str> {
        self.transcript
            .as_deref()
            .map(str::trim)
            .filter(|text| !text.is_empty())
    }

    pub fn validate(&self) -> Result<(), &'static str> {
        if matches!(self.transcript.as_deref(), Some(text) if text.trim().is_empty()) {
            return Err("transcript must not be empty when provided");
        }

        if matches!(self.sample_rate_hz, Some(0)) {
            return Err("sample_rate_hz must be greater than zero when provided");
        }

        if matches!(self.duration_ms, Some(0)) {
            return Err("duration_ms must be greater than zero when provided");
        }

        Ok(())
    }
}

#[derive(Debug, Serialize)]
pub struct ProsodyResult {
    pub word_count: usize,
    pub speech_rate_wpm: Option<u32>,
    pub estimated_pause_ms: u32,
    pub pause_ratio: Option<f32>,
    pub pacing: &'static str,
    pub note: String,
}

pub async fn analyze_prosody(input: ProsodyInput) -> ProsodyResult {
    let transcript = input.transcript();
    let word_count = transcript
        .map(|text| text.split_whitespace().count())
        .unwrap_or(0);
    let duration_ms = input.duration_ms;
    let speech_rate_wpm = match (word_count, duration_ms) {
        (0, _) | (_, None) => None,
        (_, Some(0)) => None,
        (words, Some(duration)) => Some(
            ((words as u64 * 60_000) / duration as u64)
                .max(1)
                .min(u32::MAX as u64) as u32,
        ),
    };

    let estimated_pause_ms = transcript
        .map(estimate_pause_ms)
        .unwrap_or(0)
        .min(duration_ms.unwrap_or(u32::MAX));

    let pause_ratio = match (estimated_pause_ms, duration_ms) {
        (_, None) | (_, Some(0)) => None,
        (pause_ms, Some(duration)) => Some(pause_ms as f32 / duration as f32),
    };

    let pacing = classify_pacing(speech_rate_wpm, pause_ratio, transcript.is_some());
    let note = build_note(
        word_count,
        speech_rate_wpm,
        estimated_pause_ms,
        pause_ratio,
        input.sample_rate_hz,
        transcript.is_some(),
    );

    ProsodyResult {
        word_count,
        speech_rate_wpm,
        estimated_pause_ms,
        pause_ratio,
        pacing,
        note,
    }
}

fn classify_pacing(
    speech_rate_wpm: Option<u32>,
    pause_ratio: Option<f32>,
    has_transcript: bool,
) -> &'static str {
    match speech_rate_wpm {
        None if has_transcript => "insufficient-data",
        None => "insufficient-data",
        Some(rate) => match pause_ratio {
            Some(ratio) if rate < 110 || ratio > 0.18 => "slow",
            Some(ratio) if rate > 165 && ratio < 0.08 => "fast",
            Some(_) => "steady",
            None if rate < 110 => "slow",
            None if rate > 165 => "fast",
            None => "steady",
        },
    }
}

fn build_note(
    word_count: usize,
    speech_rate_wpm: Option<u32>,
    estimated_pause_ms: u32,
    pause_ratio: Option<f32>,
    sample_rate_hz: Option<u32>,
    has_transcript: bool,
) -> String {
    if !has_transcript {
        return match sample_rate_hz {
            Some(hz) => format!("No transcript supplied; audio sampled at {} Hz.", hz),
            None => "No transcript supplied; prosody outputs are insufficient-data.".to_string(),
        };
    }

    let mut fragments = vec![format!("{} words", word_count)];
    if let Some(rate) = speech_rate_wpm {
        fragments.push(format!("{} wpm speech rate", rate));
    } else {
        fragments.push("speech rate unavailable".to_string());
    }
    fragments.push(format!("{} ms estimated pauses", estimated_pause_ms));

    if let Some(ratio) = pause_ratio {
        fragments.push(format!("{:.1}% pause ratio", ratio * 100.0));
    }

    if let Some(hz) = sample_rate_hz {
        fragments.push(format!("audio sampled at {} Hz", hz));
    }

    format!("Derived from {}.", fragments.join(", "))
}

fn estimate_pause_ms(transcript: &str) -> u32 {
    let normalized = transcript.trim();
    let lower = normalized.to_lowercase();
    let filler_hits = count_exact_words(&lower, &["um", "uh", "er", "erm", "like"])
        + count_occurrences(&lower, "you know");

    let ellipsis_count = count_occurrences(normalized, "...");
    let terminal_hits = normalized.matches('!').count()
        + normalized.matches('?').count()
        + normalized.matches('.').count()
        - (ellipsis_count * 3);

    let mut pause_ms = 0u32;
    pause_ms += (normalized.matches(',').count() as u32) * 180;
    pause_ms += (normalized.matches(';').count() as u32) * 160;
    pause_ms += (normalized.matches(':').count() as u32) * 120;
    pause_ms += (ellipsis_count as u32) * 300;
    pause_ms += (terminal_hits as u32) * 220;
    pause_ms += (filler_hits as u32) * 120;

    pause_ms
}

fn count_exact_words(haystack: &str, words: &[&str]) -> usize {
    haystack
        .split_whitespace()
        .filter(|token| {
            let token = token.trim_matches(|ch: char| !ch.is_alphabetic());
            words.iter().any(|word| token == *word)
        })
        .count()
}

fn count_occurrences(haystack: &str, needle: &str) -> usize {
    if needle.is_empty() {
        return 0;
    }

    let mut count = 0;
    let mut start = 0;
    while let Some(index) = haystack[start..].find(needle) {
        count += 1;
        start += index + needle.len();
    }
    count
}

#[cfg(test)]
mod tests {
    use super::{analyze_prosody, ProsodyInput};

    #[tokio::test]
    async fn analyzes_prosody_with_transcript() {
        let result = analyze_prosody(ProsodyInput {
            transcript: Some("hello world".to_string()),
            sample_rate_hz: Some(16_000),
            duration_ms: Some(1500),
        })
        .await;

        assert_eq!(result.word_count, 2);
        assert_eq!(result.speech_rate_wpm, Some(80));
        assert_eq!(result.estimated_pause_ms, 0);
        assert_eq!(result.pause_ratio, Some(0.0));
        assert_eq!(result.pacing, "slow");
        assert_eq!(
            result.note,
            "Derived from 2 words, 80 wpm speech rate, 0 ms estimated pauses, 0.0% pause ratio, audio sampled at 16000 Hz."
        );
    }

    #[tokio::test]
    async fn analyzes_prosody_without_transcript() {
        let result = analyze_prosody(ProsodyInput {
            transcript: None,
            sample_rate_hz: None,
            duration_ms: None,
        })
        .await;

        assert_eq!(result.word_count, 0);
        assert_eq!(result.speech_rate_wpm, None);
        assert_eq!(result.estimated_pause_ms, 0);
        assert_eq!(result.pause_ratio, None);
        assert_eq!(result.pacing, "insufficient-data");
        assert_eq!(
            result.note,
            "No transcript supplied; prosody outputs are insufficient-data."
        );
    }

    #[tokio::test]
    async fn infers_pauses_from_transcript_punctuation() {
        let result = analyze_prosody(ProsodyInput {
            transcript: Some("I think, um we should go now...".to_string()),
            sample_rate_hz: Some(44_100),
            duration_ms: Some(4_000),
        })
        .await;

        assert_eq!(result.word_count, 7);
        assert_eq!(result.speech_rate_wpm, Some(105));
        assert_eq!(result.estimated_pause_ms, 600);
        assert!(matches!(result.pause_ratio, Some(ratio) if (ratio - 0.15).abs() < 0.0001));
        assert_eq!(result.pacing, "slow");
        assert_eq!(
            result.note,
            "Derived from 7 words, 105 wpm speech rate, 600 ms estimated pauses, 15.0% pause ratio, audio sampled at 44100 Hz."
        );
    }
}
