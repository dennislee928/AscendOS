mod pipeline;

use axum::{extract::Json, routing::{get, post}, Router};

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
) -> Json<pipeline::ProsodyResult> {
    let result = pipeline::analyze_prosody(payload).await;
    Json(result)
}
