// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Diagnostic messages produced during analysis.
use super::Severity;

use crate::file::FileId;
use crate::location::Span;

/// A diagnostic message attached to a source span.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Diagnostic {
    file_id: FileId,
    span: Span,
    severity: Severity,
    message: String,
}

impl Diagnostic {
    /// Creates a diagnostic message.
    pub fn new(
        file_id: FileId,
        span: Span,
        severity: Severity,
        message: impl Into<String>,
    ) -> Self {
        Self {
            file_id,
            span,
            severity,
            message: message.into(),
        }
    }

    /// Returns the file that contains the diagnostic span.
    pub const fn file_id(&self) -> FileId {
        self.file_id
    }

    /// Returns the source span associated with this diagnostic.
    pub const fn span(&self) -> Span {
        self.span
    }

    /// Returns the diagnostic severity.
    pub const fn severity(&self) -> Severity {
        self.severity
    }

    /// Returns the diagnostic message.
    pub fn message(&self) -> &str {
        &self.message
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    use crate::file::FileId;
    use crate::location::{Position, Span};

    #[test]
    fn diagnostic_preserves_file_span_severity_and_message() {
        let file_id = FileId(7);
        let span = Span::new(Position(3), Position(8));

        let diagnostic = Diagnostic::new(file_id, span, Severity::Warning, "unexpected whitespace");

        assert_eq!(diagnostic.file_id(), file_id);
        assert_eq!(diagnostic.span(), span);
        assert_eq!(diagnostic.severity(), Severity::Warning);
        assert_eq!(diagnostic.message(), "unexpected whitespace");
    }
}
