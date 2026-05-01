mod pipeline;

use axum::{
    extract::Json,
    http::StatusCode,
    routing::{get, post},
    Router,
};

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/healthz", get(healthz))
        .route("/v1/audio/prosody", post(analyze_prosody));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8081")
        .await
        .expect("failed to bind orator listener");
    axum::serve(listener, app)
        .await
        .expect("orator server error");
}

async fn healthz() -> &'static str {
    "ok"
}

async fn analyze_prosody(
    Json(payload): Json<pipeline::ProsodyInput>,
) -> Result<Json<pipeline::ProsodyResult>, (StatusCode, &'static str)> {
    payload
        .validate()
        .map_err(|message| (StatusCode::BAD_REQUEST, message))?;
    let result = pipeline::analyze_prosody(payload).await;
    Ok(Json(result))
}

#[cfg(test)]
mod tests {
    use super::pipeline::ProsodyInput;

    #[test]
    fn rejects_blank_transcript() {
        let input = ProsodyInput {
            transcript: Some("  ".to_string()),
            sample_rate_hz: Some(16_000),
            duration_ms: Some(2500),
        };

        assert!(input.validate().is_err());
    }

    #[test]
    fn accepts_minimal_payload() {
        let input = ProsodyInput {
            transcript: None,
            sample_rate_hz: None,
            duration_ms: None,
        };

        assert!(input.validate().is_ok());
    }

    #[test]
    fn rejects_zero_sample_rate() {
        let input = ProsodyInput {
            transcript: Some("hello".to_string()),
            sample_rate_hz: Some(0),
            duration_ms: Some(1000),
        };

        assert_eq!(
            input.validate(),
            Err("sample_rate_hz must be greater than zero when provided")
        );
    }

    #[test]
    fn rejects_zero_duration() {
        let input = ProsodyInput {
            transcript: Some("hello".to_string()),
            sample_rate_hz: Some(16_000),
            duration_ms: Some(0),
        };

        assert_eq!(
            input.validate(),
            Err("duration_ms must be greater than zero when provided")
        );
    }
}
