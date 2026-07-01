// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Rendering diagnostics for display.

use crate::{diagnostic::Diagnostic, file::SourceFile};

/// Renders a diagnostic as plain text.
///
/// The rendered output includes the diagnostic severity, message, source location, and the corresponding line of source
/// text.
pub fn render(diagnostic: &Diagnostic, source_file: &SourceFile) -> String {
    debug_assert_eq!(diagnostic.file_id(), source_file.id());

    let location = source_file
        .line_map()
        .line_column(diagnostic.span().start());

    let line_text = source_file
        .text()
        .lines()
        .nth(location.line - 1)
        .unwrap_or("");

    format!(
        "{severity}: {message}\n --> {path}:{line}:{column}\n  |\n{line_number} | {line_text}\n",
        severity = diagnostic.severity(),
        message = diagnostic.message(),
        path = source_file.path().display(),
        line = location.line,
        column = location.column,
        line_number = location.line,
        line_text = line_text,
    )
}

#[cfg(test)]
mod tests {
    use std::path::PathBuf;

    use super::*;

    use crate::diagnostic::Severity;
    use crate::file::FileId;
    use crate::location::{Position, Span};

    #[test]
    fn render_includes_severity_and_message() {
        let source_file =
            SourceFile::new(FileId(1), PathBuf::from("main.rSpot"), String::from("abc"));

        let diagnostic = Diagnostic::new(
            FileId(1),
            Span::new(Position(0), Position(1)),
            Severity::Error,
            "unexpected token",
        );

        let output = render(&diagnostic, &source_file);

        assert!(output.contains("error: unexpected token"));
    }

    #[test]
    fn render_includes_line_and_column() {
        let source_file = SourceFile::new(
            FileId(1),
            PathBuf::from("main.rSpot"),
            String::from("abc\ndef"),
        );

        let diagnostic = Diagnostic::new(
            FileId(1),
            Span::new(Position(4), Position(7)),
            Severity::Error,
            "unexpected token",
        );

        let output = render(&diagnostic, &source_file);

        assert!(output.contains("main.rSpot:2:1"));
    }

    #[test]
    fn render_includes_source_line() {
        let source_file = SourceFile::new(
            FileId(1),
            PathBuf::from("main.rSpot"),
            String::from("abc\ndef"),
        );

        let diagnostic = Diagnostic::new(
            FileId(1),
            Span::new(Position(4), Position(7)),
            Severity::Error,
            "unexpected token",
        );

        let output = render(&diagnostic, &source_file);

        assert!(output.contains("2 | def"));
    }
}
