use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct TextAnalysisResult {
    pub token_count: usize,
    pub sentiment: &'static str,
    pub summary: String,
}

pub async fn analyze_text(input: &str) -> TextAnalysisResult {
    let token_count = input.split_whitespace().count();
    let summary = if input.is_empty() {
        "No text provided".to_string()
    } else {
        format!("Stub analysis complete for {} characters", input.chars().count())
    };

    TextAnalysisResult {
        token_count,
        sentiment: "neutral",
        summary,
    }
}
