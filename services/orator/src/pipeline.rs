use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
pub struct ProsodyInput {
    pub transcript: Option<String>,
    pub sample_rate_hz: Option<u32>,
    pub duration_ms: Option<u32>,
}

impl ProsodyInput {
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

        assert_eq!(result.pacing_wpm, 132);
        assert_eq!(result.energy, "medium");
        assert_eq!(result.note, "Stub prosody analysis complete for 1500 ms of audio at 16000 Hz");
    }

    #[tokio::test]
    async fn analyzes_prosody_without_transcript() {
        let result = analyze_prosody(ProsodyInput {
            transcript: None,
            sample_rate_hz: None,
            duration_ms: None,
        })
        .await;

        assert_eq!(result.note, "Stub prosody analysis complete without transcript");
    }
}
