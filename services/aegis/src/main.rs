mod pipeline;

use axum::{
    extract::Json,
    http::StatusCode,
    routing::{get, post},
    Router,
};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct AnalyzeTextRequest {
    text: String,
}

impl AnalyzeTextRequest {
    fn validate(&self) -> Result<(), &'static str> {
        if self.text.trim().is_empty() {
            Err("text must not be empty")
        } else {
            Ok(())
        }
    }
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/healthz", get(healthz))
        .route("/v1/text/analyze", post(analyze_text));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080")
        .await
        .expect("failed to bind aegis listener");
    axum::serve(listener, app)
        .await
        .expect("aegis server error");
}

async fn healthz() -> &'static str {
    "ok"
}

async fn analyze_text(
    Json(payload): Json<AnalyzeTextRequest>,
) -> Result<Json<pipeline::TextAnalysisResult>, (StatusCode, &'static str)> {
    payload
        .validate()
        .map_err(|message| (StatusCode::BAD_REQUEST, message))?;
    let result = pipeline::analyze_text(&payload.text).await;
    Ok(Json(result))
}

#[cfg(test)]
mod tests {
    use super::AnalyzeTextRequest;

    #[test]
    fn rejects_blank_text() {
        let request = AnalyzeTextRequest {
            text: "   ".to_string(),
        };

        assert!(request.validate().is_err());
    }

    #[test]
    fn accepts_non_empty_text() {
        let request = AnalyzeTextRequest {
            text: "hello world".to_string(),
        };

        assert!(request.validate().is_ok());
    }
}
