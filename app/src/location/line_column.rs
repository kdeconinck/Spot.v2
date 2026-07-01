// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Line and column locations within source text.

/// A one-based line and column location in source text.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct LineColumn {
    /// The one-based line number.
    pub line: usize,

    /// The one-based column number.
    pub column: usize,
}
