// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================
//! Byte ranges within source text.
use super::Position;

/// A half-open byte range in source text: `[start, end)`.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Span {
    start: Position,
    end: Position,
}

impl Span {
    /// Creates a half-open byte range.
    ///
    /// # Panics
    ///
    /// Panics when `end` is before `start`.
    pub fn new(start: Position, end: Position) -> Self {
        assert!(start <= end, "span end cannot be before span start");

        Self { start, end }
    }

    /// Returns the first position in the span.
    pub const fn start(self) -> Position {
        self.start
    }

    /// Returns the position immediately after the end of the span.
    pub const fn end(self) -> Position {
        self.end
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn span_preserves_start_and_end_positions() {
        let start = Position(3);
        let end = Position(8);

        let span = Span::new(start, end);

        assert_eq!(span.start(), start);
        assert_eq!(span.end(), end);
    }

    #[test]
    fn span_allows_empty_ranges() {
        let position = Position(3);

        let span = Span::new(position, position);

        assert_eq!(span.start(), position);
        assert_eq!(span.end(), position);
    }

    #[test]
    #[should_panic(expected = "span end cannot be before span start")]
    fn span_rejects_reversed_ranges() {
        Span::new(Position(8), Position(3));
    }
}
