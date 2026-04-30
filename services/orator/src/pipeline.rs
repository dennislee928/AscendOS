use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct ProsodyInput {
    pub transcript: Option<String>,
    pub sample_rate_hz: Option<u32>,
    pub duration_ms: Option<u32>,
}

#[derive(Debug, Serialize)]
pub struct ProsodyResult {
    pub pacing_wpm: u32,
    pub energy: &'static str,
    pub note: String,
}

pub async fn analyze_prosody(input: ProsodyInput) -> ProsodyResult {
    let pacing_wpm = 132;
    let duration = input.duration_ms.unwrap_or(0);
    let sample_rate = input.sample_rate_hz.unwrap_or(0);
    let has_transcript = input
        .transcript
        .as_ref()
        .map(|s| !s.trim().is_empty())
        .unwrap_or(false);

    let note = if has_transcript {
        format!(
            "Stub prosody analysis complete for {} ms of audio at {} Hz",
            duration, sample_rate
        )
    } else {
        "Stub prosody analysis complete without transcript".to_string()
    };

    ProsodyResult {
        pacing_wpm,
        energy: "medium",
        note,
    }
}
