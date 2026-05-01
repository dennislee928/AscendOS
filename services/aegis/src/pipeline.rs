use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct TextAnalysisResult {
    pub token_count: usize,
    pub sentiment: &'static str,
    pub summary: String,
}

pub async fn analyze_text(input: &str) -> TextAnalysisResult {
    let normalized = input.trim();
    let token_count = normalized.split_whitespace().count();
    let summary = if normalized.is_empty() {
        "No text provided".to_string()
    } else {
        format!(
            "Stub analysis complete for {} characters",
            normalized.chars().count()
        )
    };

    TextAnalysisResult {
        token_count,
        sentiment: "neutral",
        summary,
    }
}

#[cfg(test)]
mod tests {
    use super::analyze_text;

    #[tokio::test]
    async fn summarizes_non_empty_text() {
        let result = analyze_text("hello there").await;

        assert_eq!(result.token_count, 2);
        assert_eq!(result.sentiment, "neutral");
        assert_eq!(result.summary, "Stub analysis complete for 11 characters");
    }

    #[tokio::test]
    async fn trims_text_before_analysis() {
        let result = analyze_text("  hello there  ").await;

        assert_eq!(result.token_count, 2);
        assert_eq!(result.summary, "Stub analysis complete for 11 characters");
    }

    #[tokio::test]
    async fn summarizes_blank_text() {
        let result = analyze_text("   ").await;

        assert_eq!(result.token_count, 0);
        assert_eq!(result.summary, "No text provided");
    }
}
