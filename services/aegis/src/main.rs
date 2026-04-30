mod pipeline;

use axum::{extract::Json, routing::{get, post}, Router};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct AnalyzeTextRequest {
    text: String,
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

async fn analyze_text(Json(payload): Json<AnalyzeTextRequest>) -> Json<pipeline::TextAnalysisResult> {
    let result = pipeline::analyze_text(&payload.text).await;
    Json(result)
}
