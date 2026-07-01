// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Diagnostic severity levels.
use std::fmt;

/// The severity of a diagnostic message.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Severity {
    /// An error that should prevent successful analysis.
    Error,

    /// A warning that reports suspicious code without failing analysis.
    Warning,

    /// Informational feedback about the analyzed source.
    Info,
}

impl fmt::Display for Severity {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Error => formatter.write_str("error"),
            Self::Warning => formatter.write_str("warning"),
            Self::Info => formatter.write_str("info"),
        }
    }
}
