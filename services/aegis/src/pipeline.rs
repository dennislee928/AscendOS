use serde::Serialize;
use std::collections::BTreeSet;

#[derive(Debug, Serialize)]
pub struct TextAnalysisResult {
    pub token_count: usize,
    pub score: u8,
    pub flags: Vec<String>,
    pub indicators: Vec<ManipulationIndicator>,
    pub summary: String,
}

#[derive(Debug, Serialize)]
pub struct ManipulationIndicator {
    pub name: &'static str,
    pub hits: usize,
    pub weight: u8,
    pub evidence: Vec<String>,
}

struct Heuristic {
    name: &'static str,
    flag: &'static str,
    weight: u8,
    phrases: &'static [&'static str],
}

const HEURISTICS: &[Heuristic] = &[
    Heuristic {
        name: "urgency pressure",
        flag: "urgency",
        weight: 18,
        phrases: &["urgent", "asap", "immediately", "right now", "act now", "last chance"],
    },
    Heuristic {
        name: "secrecy framing",
        flag: "secrecy",
        weight: 16,
        phrases: &[
            "don't tell",
            "keep this between us",
            "just between us",
            "private",
            "confidential",
        ],
    },
    Heuristic {
        name: "guilt pressure",
        flag: "guilt_pressure",
        weight: 14,
        phrases: &[
            "after all i've done",
            "if you cared",
            "if you really cared",
            "you owe me",
            "disappointed",
            "prove",
        ],
    },
    Heuristic {
        name: "certainty pressure",
        flag: "pressure",
        weight: 12,
        phrases: &["trust me", "must", "have to", "need to", "no choice", "everyone knows"],
    },
    Heuristic {
        name: "scarcity framing",
        flag: "scarcity",
        weight: 10,
        phrases: &["limited time", "last one", "only chance", "exclusive", "now or never"],
    },
];

pub async fn analyze_text(input: &str) -> TextAnalysisResult {
    let normalized = input.trim();
    let token_count = normalized.split_whitespace().count();
    let lower = normalized.to_lowercase();

    let mut indicators = Vec::new();
    let mut flags = BTreeSet::new();
    let mut score = 0u32;

    for heuristic in HEURISTICS {
        let mut hits = 0usize;
        let mut evidence = Vec::new();

        for phrase in heuristic.phrases {
            let count = count_occurrences(&lower, phrase);
            if count > 0 {
                hits += count;
                evidence.push((*phrase).to_string());
            }
        }

        if hits > 0 {
            score += heuristic.weight as u32 * hits as u32;
            flags.insert(heuristic.flag.to_string());
            indicators.push(ManipulationIndicator {
                name: heuristic.name,
                hits,
                weight: heuristic.weight,
                evidence,
            });
        }
    }

    let shouting_hits = normalized
        .split_whitespace()
        .filter(|word| is_shouting_word(word))
        .count();
    if shouting_hits > 0 {
        score += 6 * shouting_hits as u32;
        flags.insert("shouting".to_string());
        indicators.push(ManipulationIndicator {
            name: "shouting",
            hits: shouting_hits,
            weight: 6,
            evidence: normalized
                .split_whitespace()
                .filter(|word| is_shouting_word(word))
                .map(|word| word.trim_matches(|c: char| !c.is_alphanumeric()).to_string())
                .collect(),
        });
    }

    let exclamation_marks = normalized.matches('!').count();
    if exclamation_marks > 1 {
        score += 3 * exclamation_marks as u32;
        flags.insert("exclamation_pressure".to_string());
        indicators.push(ManipulationIndicator {
            name: "exclamation pressure",
            hits: exclamation_marks,
            weight: 3,
            evidence: vec!["!".repeat(exclamation_marks.min(3))],
        });
    }

    let score = score.min(100) as u8;
    let flags = flags.into_iter().collect::<Vec<_>>();
    let summary = if flags.is_empty() {
        format!("No strong manipulation cues detected across {} tokens.", token_count)
    } else {
        format!(
            "Score {} with {} flag(s): {}.",
            score,
            flags.len(),
            flags.join(", ")
        )
    };

    TextAnalysisResult {
        token_count,
        score,
        flags,
        indicators,
        summary,
    }
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

fn is_shouting_word(word: &str) -> bool {
    let letters: String = word.chars().filter(|ch| ch.is_alphabetic()).collect();
    letters.len() >= 4 && letters.chars().all(|ch| ch.is_uppercase())
}

#[cfg(test)]
mod tests {
    use super::analyze_text;

    #[tokio::test]
    async fn scores_text_with_manipulation_cues() {
        let result = analyze_text("This is urgent. Don't tell anyone. Trust me.").await;

        assert_eq!(result.token_count, 8);
        assert_eq!(result.score, 46);
        assert_eq!(
            result.flags,
            vec![
                "pressure".to_string(),
                "secrecy".to_string(),
                "urgency".to_string()
            ]
        );
        assert_eq!(result.indicators.len(), 3);
        assert!(result.summary.contains("Score 46"));
    }

    #[tokio::test]
    async fn returns_low_signal_for_neutral_text() {
        let result = analyze_text("Let's review the project timeline and adjust the agenda.").await;

        assert_eq!(result.token_count, 9);
        assert_eq!(result.score, 0);
        assert!(result.flags.is_empty());
        assert!(result.indicators.is_empty());
        assert_eq!(
            result.summary,
            "No strong manipulation cues detected across 9 tokens."
        );
    }

    #[tokio::test]
    async fn captures_shouting_pressure() {
        let result = analyze_text("THIS IS URGENT!").await;

        assert_eq!(result.score, 30);
        assert_eq!(result.flags, vec!["shouting".to_string(), "urgency".to_string()]);
        assert_eq!(result.indicators.len(), 2);
    }
}
